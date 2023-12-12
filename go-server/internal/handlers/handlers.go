package handlers

import (
	"fmt"
	"liarOfTuring/global"
	"liarOfTuring/models"
	"liarOfTuring/utils"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var mutex = &sync.Mutex{}
var nicknames []models.Nickname

func init() {
	// 애플리케이션 시작 시 전역 변수 초기화
	if err := utils.LoadNicknames(); err != nil {
		fmt.Println("Error loading nicknames:", err)
		return
	}
	nicknames = global.GetGlobalNicknames()

	// clients[WebSocketConnection{Conn: nil}] = adminUser
	players["0"] = adminUser
}

var adminUser = User{
	UUID:       "0",
	UserId:     0,
	UserName:   "server",
	NicknameId: 999,
	Role:       "admin",
	IsOnline:   false,
	PlayerType: "admin",
}

var Broadcast = make(chan WsPayload)

var clients = make(map[WebSocketConnection]User)

var players = make(map[string]User)
var sorted_players = make([]User, 0)

var messages = make([]Message, 0)
var MaxPlayer = 2

var gameRoundNum = 2
var gameTurnNum = 1

var isGameStarted = false
var gameInfo = make([]Game, 0)
var gameTurnsLeft = gameTurnNum * MaxPlayer
var gameRound = 1

type Game struct {
	NowUserIndex        int                  `json:"now_user_index"`
	MaxPlayer           int                  `json:"max_player"`
	OnlineUserList      []User               `json:"online_user_list"`
	PlayerList          []User               `json:"player_list"`
	TurnsLeft           int                  `json:"turns_left"`
	UserRoundSelections []UserRoundSelection `json:"user_round_selections"`
}

type UserRoundSelection struct {
	User          User   `json:"user"`
	WhoIsAI       bool   `json:"who_is_ai"`
	ReasonMessage string `json:"reason_message"`
}

// upgradeConnection is the websocket upgrader from gorilla/websockets
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
	HandshakeTimeout: 1024,
}

type WebSocketConnection struct {
	*websocket.Conn
}

// SafeClose closes the websocket connection safely
func (conn *WebSocketConnection) SafeClose() {
	if conn != nil && conn.Conn != nil {
		conn.Close()
	}
}

// WsJsonResponse defines the response sent back from websocket
type WsJsonResponse struct {
	Timestamp      int64     `json:"timestamp"`
	MaxPlayer      int       `json:"max_player"`
	Action         string    `json:"action"`
	User           User      `json:"user"`
	Message        string    `json:"message"`
	MessageType    string    `json:"message_type"`
	MessageLogList []Message `json:"message_log_list"`
	OnlineUserList []User    `json:"online_user_list"`
	PlayerList     []User    `json:"player_list"`
	GameTurnsLeft  int       `json:"game_turns_left"`
	GameRound      int       `json:"game_round"`
}

type User struct {
	UUID   string `json:"uuid"`
	UserId int64  `json:"user_id"`
	// RoomId     int64  `json:"room_id"`
	NicknameId int    `json:"nickname_id"`
	UserName   string `json:"username"`
	Role       string `json:"role"`
	IsOnline   bool   `json:"is_online"`
	PlayerType string `json:"player_type"`
}

type Message struct {
	Timestamp   int64  `json:"timestamp"`
	MessageId   int64  `json:"message_id"`
	User        User   `json:"user"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

func sortUserList(users []User) []User {
	sort.Slice(users, func(i, j int) bool {
		return users[i].UserName < users[j].UserName
	})
	return users
}

// WsPayload defines the websocket request from the client
type WsPayload struct {
	Action string `json:"action"`
	// RoomId    int64               `json:"room_id"`
	MaxPlayer     int                 `json:"max_player"`
	User          User                `json:"user"`
	Timestamp     int64               `json:"timestamp"`
	Message       string              `json:"message"`
	Conn          WebSocketConnection `json:"-"` // ignore this field
	GameTurnsLeft int                 `json:"game_turns_left"`
	GameRound     int                 `json:"game_round"`
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// WsEndpoint upgrades connection to websocket
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	conn := WebSocketConnection{Conn: ws}
	log.Println("Client connected to endpoint")

	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in ListenForWs:", r)
		}
		conn.SafeClose() // Safely close the connection
		mutex.Lock()
		delete(clients, *conn) // Remove the client from the map
		mutex.Unlock()
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// Check if the error is a normal WebSocket closure
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
				log.Printf("Client disconnected: %v\n", err)
			} else {
				log.Printf("Error reading json: %v\n", err)
			}
			break // Exit loop on client disconnection or error
		}

		payload.Conn = *conn
		Broadcast <- payload
	}
}

func ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-Broadcast

		// Consolidated logging
		log.Printf("Action: %s, User: %v\n", e.Action, e.User)

		mutex.Lock()
		// MaxPlayer = e.MaxPlayer
		log.Println("MaxPlayer:", MaxPlayer)
		log.Println("messages:", messages)
		switch e.Action {
		case "broadcast", "new_message_admin":
			message := Message{
				Timestamp:   e.Timestamp,
				MessageId:   int64(len(messages)),
				User:        e.User,
				Message:     e.Message,
				MessageType: "message",
			}
			messages = append(messages, message)

			response = WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         e.Action,
				User:           e.User,
				Message:        e.Message,
				MessageLogList: messages,
				MessageType:    "message",
				MaxPlayer:      MaxPlayer,

				OnlineUserList: getUserList(),
				PlayerList:     getPlayerList(),
				GameTurnsLeft:  gameTurnsLeft,
				GameRound:      gameRound,
			}
			broadcastToAll(response)

		case "list_users":
			response = WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "update_status",
				MessageType:    "info",
				MessageLogList: messages,
				MaxPlayer:      MaxPlayer,
				OnlineUserList: getUserList(),
				PlayerList:     getPlayerList(),
				GameTurnsLeft:  gameTurnsLeft,
				GameRound:      gameRound,
			}
			broadcastToAll(response)
		case "new_message":
			processNewMessage(e, &response)
			// broadcastToAll(response)
		case "enter_human":
			processEnterHuman(e, &response)
			// broadcastToAll(response)

		case "left_user":
			processLeftUser(e, &response)
			broadcastToAll(response)

		case "user_is_ready":
			processReadyUser(e, &response)

		case "set_max_player":
			MaxPlayer = e.MaxPlayer
			log.Println("e.MaxPlayer:", MaxPlayer)
			message := Message{
				Timestamp:   e.Timestamp,
				MessageId:   int64(len(messages)),
				User:        e.User,
				Message:     fmt.Sprintf("최대 인원이 %d명으로 설정되었습니다.", MaxPlayer),
				MessageType: "alert",
			}
			messages = append(messages, message)
			response = WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "update_state",
				MessageType:    "alert",
				Message:        message.Message,
				MaxPlayer:      MaxPlayer,
				User:           User{},
				MessageLogList: messages,
				OnlineUserList: getUserList(),
				PlayerList:     getPlayerList(),
				GameTurnsLeft:  gameTurnsLeft,
				GameRound:      gameRound,
			}
			broadcastToAll(response)

		case "clear_messages":
			messages = []Message{}
			response = WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "update_state",
				MessageType:    "info",
				Message:        "메시지가 초기화되었습니다.",
				MaxPlayer:      MaxPlayer,
				User:           User{},
				MessageLogList: messages,
				OnlineUserList: getUserList(),
				PlayerList:     getPlayerList(),
				GameTurnsLeft:  gameTurnsLeft,
				GameRound:      gameRound,
			}
			broadcastToAll(response)

		default:
			log.Printf("Unknown action: %s\n", e.Action)
		}
		// log.Println("users:", getUserList())
		// if response.Action != "" {
		// 	broadcastToAll(response)
		// }
		mutex.Unlock()
	}
}

func processEnterHuman(e WsPayload, response *WsJsonResponse) {
	nowUser, exists := players[e.User.UUID]
	log.Println("processEnterHuman")
	log.Println("isGameStarted:", isGameStarted)

	if !exists {
		nicknameId, userName := getRandomUsername()
		nowUser = User{
			UserId:     int64(len(clients)),
			UserName:   userName,
			NicknameId: nicknameId,
			Role:       "human",
			UUID:       e.User.UUID,
			IsOnline:   true,
			PlayerType: "watcher",
		}
		// players[e.User.UUID] = nowUser
	}
	nowUser.IsOnline = true
	log.Println(nowUser.NicknameId)
	players[e.User.UUID] = nowUser
	clients[e.Conn] = nowUser

	*response = WsJsonResponse{
		MaxPlayer:      MaxPlayer,
		Timestamp:      e.Timestamp,
		Action:         "human_info",
		MessageType:    "system",
		Message:        fmt.Sprintf("%s님이 입장했습니다.", nowUser.UserName),
		MessageLogList: messages,
		User:           nowUser,
		OnlineUserList: getUserList(),
		PlayerList:     getPlayerList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      gameRound,
	}
	broadcastToAll(*response)

	if isGameStarted {
		nextUser := gameInfo[gameRound-1].PlayerList[gameInfo[gameRound-1].NowUserIndex]
		if nextUser.UUID == e.User.UUID {
			nowUser.PlayerType = "player"
			someoneMessage := Message{
				Timestamp:   e.Timestamp,
				MessageId:   int64(len(messages)),
				User:        adminUser,
				Message:     fmt.Sprintf("%s님의 차례입니다.", nextUser.UserName),
				MessageType: "alert",
			}
			// messages = append(messages, someoneMessage)
			someoneResponse := WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "your_turn",
				MessageType:    "alert",
				Message:        someoneMessage.Message,
				MaxPlayer:      MaxPlayer,
				User:           adminUser, // nextUser,
				MessageLogList: messages,
				PlayerList:     getPlayerList(),
				OnlineUserList: getUserList(),
				GameTurnsLeft:  gameTurnsLeft,
				GameRound:      gameRound,
			}
			nextConn := getNextUserClientByUUID(nextUser.UUID)
			braodCastToSomeone(someoneResponse, nextConn, nextUser)
		}
	}
	if gameTurnsLeft == 0 {
		processGameOver(e)
	}
}

func processLeftUser(e WsPayload, response *WsJsonResponse) {
	if leftUser, ok := clients[e.Conn]; ok {
		leftUser.IsOnline = false
		players[e.User.UUID] = leftUser

		delete(clients, e.Conn)
		e.Conn.SafeClose()

		*response = WsJsonResponse{
			Action:    "update_state",
			Timestamp: e.Timestamp,
			MaxPlayer: MaxPlayer,

			// Action:         "user_list",
			MessageLogList: messages,
			OnlineUserList: getUserList(),
			PlayerList:     getPlayerList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
		}
	} else {
		log.Println("User not found in clients map")
		// *response = WsJsonResponse{} // Empty response
		*response = WsJsonResponse{
			Action:    "update_state",
			Timestamp: e.Timestamp,
			// Action:         "user_list",
			MessageLogList: messages,
			OnlineUserList: getUserList(),
			MaxPlayer:      MaxPlayer,
			PlayerList:     getPlayerList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
		}
	}
}

func processNewMessage(e WsPayload, response *WsJsonResponse) {
	log.Println("processNewMessage")
	// log.Println("message", e.Message)
	broadcastNewMessage(e)

	// GameRound := e.GameRound
	gameInfo[gameRound-1].NowUserIndex = (gameInfo[gameRound-1].NowUserIndex + 1) % gameInfo[gameRound-1].MaxPlayer
	gameTurnsLeft = Max(gameTurnsLeft-1, 0)
	gameInfo[gameRound-1].TurnsLeft = gameTurnsLeft

	log.Println("gameRound:", gameRound)
	log.Println("gameTurnsLeft:", gameTurnsLeft)

	if gameTurnsLeft == 0 {
		processGameOver(e)
	}

	nextUser := gameInfo[gameRound-1].PlayerList[gameInfo[gameRound-1].NowUserIndex]
	someoneMessage := Message{
		Timestamp:   e.Timestamp,
		MessageId:   int64(len(messages)),
		User:        adminUser,
		Message:     fmt.Sprintf("%s님의 차례입니다.", nextUser.UserName),
		MessageType: "alert",
	}
	// messages = append(messages, someoneMessage)
	someoneResponse := WsJsonResponse{
		Timestamp:      e.Timestamp,
		Action:         "your_turn",
		MessageType:    "alert",
		Message:        someoneMessage.Message,
		MaxPlayer:      MaxPlayer,
		User:           adminUser, // nextUser,
		MessageLogList: messages,
		PlayerList:     getPlayerList(),
		OnlineUserList: getUserList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      gameRound,
	}
	nextConn := getNextUserClientByUUID(nextUser.UUID)
	braodCastToSomeone(someoneResponse, nextConn, nextUser)
}

func processReadyUser(e WsPayload, response *WsJsonResponse) {
	nowUser := players[e.User.UUID]
	nowUser.PlayerType = "player"
	log.Println(nowUser.NicknameId)
	players[e.User.UUID] = nowUser
	clients[e.Conn] = nowUser
	sorted_players = append(sorted_players, nowUser)

	if MaxPlayer == len(getPlayerList()) {
		gameRound = 1
		initGameInfo()
		message := Message{
			Timestamp:   e.Timestamp,
			MessageId:   int64(len(messages)),
			User:        adminUser,
			Message:     "게임이 시작되었습니다.",
			MessageType: "alert",
		}
		messages = []Message{message}

		gameTurnsLeft := gameInfo[0].TurnsLeft
		*response = WsJsonResponse{
			Timestamp:      e.Timestamp,
			Action:         "update_state",
			MessageType:    "alert",
			Message:        message.Message,
			MaxPlayer:      MaxPlayer,
			User:           User{},
			MessageLogList: messages,
			PlayerList:     getPlayerList(),
			OnlineUserList: getUserList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
		}
		broadcastToAll(*response)

		nextUser := gameInfo[gameRound-1].PlayerList[gameInfo[gameRound-1].NowUserIndex]
		log.Println("nextUser:", nextUser)

		someoneMessage := Message{
			Timestamp:   e.Timestamp,
			MessageId:   int64(len(messages)),
			User:        adminUser,
			Message:     fmt.Sprintf("%s님의 차례입니다.", nextUser.UserName),
			MessageType: "alert",
		}
		// messages = append(messages, someoneMessage)
		someoneResponse := WsJsonResponse{
			Timestamp:   e.Timestamp,
			Action:      "your_turn",
			MessageType: "alert",
			Message:     someoneMessage.Message,
			MaxPlayer:   MaxPlayer,
			// User:           nextUser,
			User:           adminUser,
			MessageLogList: messages,
			PlayerList:     getPlayerList(),
			OnlineUserList: getUserList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      1,
		}
		nextConn := getNextUserClientByUUID(nextUser.UUID)
		braodCastToSomeone(someoneResponse, nextConn, nextUser)

		return
	}

	*response = WsJsonResponse{
		Timestamp:      e.Timestamp,
		Action:         "human_info",
		MessageType:    "system",
		Message:        fmt.Sprintf("%s님이 게임에 참여했습니다.", nowUser.UserName),
		User:           nowUser,
		OnlineUserList: getUserList(),
		MessageLogList: messages,
		MaxPlayer:      MaxPlayer,
		PlayerList:     getPlayerList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      gameRound,
	}
	isGameStarted = true
	broadcastToAll(*response)

}

func getUserList() []User {
	var userList []User
	for _, v := range players { //TODO:  or clients?
		if v.UserName != "admin" && v.IsOnline {
			userList = append(userList, v)
		}
	}
	// sortUserList(userList)
	return userList
}

func getPlayerList() []User {
	var playerList []User
	for _, v := range sorted_players {
		// playerList { //TODO:  or clients?
		if v.UserName != "admin" && v.IsOnline && v.PlayerType == "player" {
			playerList = append(playerList, v)
		}
	}
	return playerList
}

func broadcastToAll(response WsJsonResponse) {
	// mutex.Lock()
	// defer mutex.Unlock()

	for client := range clients {
		if err := client.WriteJSON(response); err != nil {
			log.Println("[broadcastToAll] Websocket error:", err)
			client.SafeClose()
			delete(clients, client)
		}
	}
	log.Println("Broadcasted message")
}

func braodCastToSomeone(response WsJsonResponse, client WebSocketConnection, user User) {
	// mutex.Lock()
	// defer mutex.Unlock()
	if err := client.WriteJSON(response); err != nil {
		log.Println("[braodCastToSomeone] Websocket error:", err)
		client.SafeClose()
		delete(clients, client)
	}
	log.Println("Broadcasted message")
}

func getRandomUsername() (int, string) {
	//random suffle nicknames and return not used nickname
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := rand.Perm(len(nicknames))
	for _, v := range perm {
		if !nicknames[v].IsUsed {
			nicknames[v].IsUsed = true
			log.Println("nickname,", nicknames[v])
			return nicknames[v].Id, nicknames[v].Nickname
		}
	}
	return -1, ""
}

// func findUser(uuid string) User {
// 	for _, v := range players {
// 		log.Println("v.UUID:", v.UUID)
// 		if v.UUID == uuid {
// 			log.Println("dddddddddd")
// 			log.Println("v:", v)
// 			v.IsOnline = true
// 			return v
// 		}
// 	}
// 	return User{}
// }

func initGameInfo() {

	gameInfo = make([]Game, 0)
	for i := 0; i < gameRoundNum; i++ {
		game := Game{
			NowUserIndex:        0,
			MaxPlayer:           MaxPlayer,
			OnlineUserList:      getUserList(),
			PlayerList:          getPlayerList(),
			TurnsLeft:           gameTurnNum * MaxPlayer,
			UserRoundSelections: make([]UserRoundSelection, 0),
		}
		gameInfo = append(gameInfo, game)
	}

}

func getNextUserClientByUUID(uuid string) WebSocketConnection {
	for client := range clients {
		if client.Conn != nil && clients[client].UUID == uuid {
			return client
		}
	}
	return WebSocketConnection{}
}

func broadcastNewMessage(e WsPayload) {
	message := Message{
		Timestamp:   e.Timestamp,
		MessageId:   int64(len(messages)),
		User:        e.User,
		Message:     e.Message,
		MessageType: "message",
	}
	messages = append(messages, message)

	response := WsJsonResponse{
		Timestamp:      e.Timestamp,
		Action:         e.Action,
		User:           e.User,
		Message:        e.Message,
		MessageLogList: messages,
		MessageType:    "message",
		MaxPlayer:      MaxPlayer,

		OnlineUserList: getUserList(),
		PlayerList:     getPlayerList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      gameRound,
	}
	broadcastToAll(response)
}

func processGameOver(e WsPayload) {
	message := Message{
		Timestamp:   e.Timestamp,
		MessageId:   int64(len(messages)),
		User:        adminUser,
		Message:     "게임이 종료되었습니다. 누가 AI였는지 확인해보세요.",
		MessageType: "alert",
	}

	response := WsJsonResponse{
		Timestamp:      e.Timestamp,
		Action:         "choose_ai",
		MessageType:    "alert",
		Message:        message.Message,
		MaxPlayer:      MaxPlayer,
		User:           User{},
		MessageLogList: messages,
		OnlineUserList: getUserList(),
		PlayerList:     getPlayerList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      gameRound,
	}
	broadcastToAll(response)

}
