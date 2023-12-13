package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

var FastAPIURL = ""
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
	GPTEnterNum, GPTReadyNum = randomChooseReadyUserAndEnterUser(MaxPlayer)
	//env

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

var GPTUser = User{
	UUID:       "999",
	UserId:     999,
	UserName:   "",
	NicknameId: 999,
	Role:       "player",
	IsOnline:   false,
	PlayerType: "watcher",
}
var GPTEnterNum = 1
var GPTReadyNum = 0
var Broadcast = make(chan WsPayload)
var GPTBroadcast = make(chan GPTWsPayload)

var clients = make(map[WebSocketConnection]User)

var players = make(map[string]User)
var sorted_players = make([]User, 0)

var messages = make([]Message, 0)
var MaxPlayer = 5

var gameRoundNum = 2
var gameTurnNum = 1

var isGameStarted = false
var isGameOver = false
var gameInfo = make([]Game, 0)
var gameTurnsLeft = gameTurnNum * MaxPlayer
var gameRound = 1
var userSelections = make([]UserSelection, 0)

type Game struct {
	NowUserIndex   int             `json:"now_user_index"`
	MaxPlayer      int             `json:"max_player"`
	OnlineUserList []User          `json:"online_user_list"`
	PlayerList     []User          `json:"player_list"`
	TurnsLeft      int             `json:"turns_left"`
	UserSelections []UserSelection `json:"user_round_selections"`
	Messages       []Message       `json:"messages"`
}

type UserSelection struct {
	User      User   `json:"user"`
	Selection string `json:"selection"`
	Reason    string `json:"reason"`
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
	Timestamp      int64           `json:"timestamp"`
	MaxPlayer      int             `json:"max_player"`
	Action         string          `json:"action"`
	User           User            `json:"user"`
	Message        string          `json:"message"`
	MessageType    string          `json:"message_type"`
	MessageLogList []Message       `json:"message_log_list"`
	OnlineUserList []User          `json:"online_user_list"`
	PlayerList     []User          `json:"player_list"`
	GameTurnsLeft  int             `json:"game_turns_left"`
	GameRound      int             `json:"game_round"`
	GameTurnNum    int             `json:"game_turn_num"`
	GameRoundNum   int             `json:"game_round_num"`
	GamsSelections []UserSelection `json:"game_selections"`
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

func randomSuffleUserList(users []User) []User {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rand := rand.New(src)

	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
	readyNum := rand.Intn(MaxPlayer-2) + 2

	//swap users
	//find GPTUser
	for i, user := range users {
		if user.UUID == GPTUser.UUID {
			users[i] = users[readyNum]
			users[readyNum] = user
			break
		}
	}

	//getRandomUsername
	for i, user := range users {
		nicknameId, userName := getRandomUsername()
		users[i] = User{
			UserId:     user.UserId,
			UserName:   userName,
			NicknameId: nicknameId,
			Role:       user.Role,
			UUID:       user.UUID,
			IsOnline:   user.IsOnline,
			PlayerType: user.PlayerType,
		}
		for conn, client := range clients {
			if client.UUID == user.UUID {
				clients[conn] = users[i]
				break
			}
		}
		if players[user.UUID].UUID == user.UUID {
			players[user.UUID] = users[i]
		}
		if user.UUID == GPTUser.UUID {
			GPTUser = users[i]
		}

	}
	return users
}

func randomChooseReadyUserAndEnterUser(max_player int) (int, int) {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rand := rand.New(src)

	//1 to max_player-2
	enterNum := rand.Intn(max_player-2) + 1
	// 2 to max_player-1
	readyNum := rand.Intn(max_player-2) + 2
	return enterNum, readyNum
}

// WsPayload defines the websocket request from the client
type WsPayload struct {
	Action string `json:"action"`
	// RoomId    int64               `json:"room_id"`
	MaxPlayer     int                 `json:"max_player"`
	User          User                `json:"user"`
	Timestamp     int64               `json:"timestamp"`
	Message       string              `json:"message"`
	GameTurnsLeft int                 `json:"game_turns_left"`
	GameRound     int                 `json:"game_round"`
	GameTurnNum   int                 `json:"game_turn_num"`
	GameRoundNum  int                 `json:"game_round_num"`
	UserSelection UserSelection       `json:"user_selection"`
	Conn          WebSocketConnection `json:"-"` // ignore this field
}

type GPTWsPayload struct {
	UserUUID string              `json:"user_uuid"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

type GPTWsJsonResponse struct {
	UserUUID         string `json:"user_uuid"`
	MessageLogString string `json:"message_log_string"`
	MessageType      string `json:"message_type"`
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

func WithGPTWsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	conn := WebSocketConnection{Conn: ws}
	log.Println("Client connected to endpoint")

	go ListenForGPTWs(&conn)
}

func ListenForGPTWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in ListenForGPTWs:", r)
		}
		conn.SafeClose() // Safely close the connection
		mutex.Lock()
		delete(clients, *conn) // Remove the client from the map
		mutex.Unlock()
	}()

	var payload GPTWsPayload

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
		GPTBroadcast <- payload
	}
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
func ListenToGPTWsChannel(fastAPIURL string) {
	FastAPIURL = fastAPIURL
	var response WsJsonResponse
	for {
		e := <-GPTBroadcast
		log.Println("GPTBroadcast:", e)
		mutex.Lock()
		message := Message{
			Timestamp:   getTimeStamp(),
			MessageId:   int64(len(messages)),
			User:        GPTUser,
			Message:     e.Message,
			MessageType: "message",
		}
		messages = append(messages, message)
		response = WsJsonResponse{
			Timestamp:      getTimeStamp(),
			Action:         "new_message",
			User:           GPTUser,
			Message:        e.Message,
			MessageLogList: messages,
			MessageType:    "message",
			MaxPlayer:      MaxPlayer,
		}
		broadcastToAll(response)
	}
}

func ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-Broadcast

		// Consolidated logging
		log.Printf("Action: %s, User: %v\n", e.Action, e.User)

		// MaxPlayer = e.MaxPlayer
		log.Println("MaxPlayer:", MaxPlayer)
		log.Println("messages:", messages)
		log.Println("userSelections:", userSelections)
		switch e.Action {
		case "broadcast", "new_message_admin":
			mutex.Lock()
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
			mutex.Unlock()

		case "list_users":
			mutex.Lock()
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
			mutex.Unlock()
		case "new_message":
			processNewMessage(e, &response)
			// broadcastToAll(response)
		case "enter_human":
			processEnterHuman(e, &response)
			// broadcastToAll(response)

		case "left_user":
			processLeftUser(e, &response)

		case "user_is_ready":
			processReadyUser(e)

		case "choose_ai":
			processChooseAI(e)

		case "set_max_player":
			mutex.Lock()
			MaxPlayer = e.MaxPlayer
			log.Println("e.MaxPlayer:", MaxPlayer)
			sorted_players = []User{}
			message := Message{
				Timestamp:   e.Timestamp,
				MessageId:   int64(len(messages)),
				User:        e.User,
				Message:     fmt.Sprintf("최대 인원이 %d명으로 설정되었습니다.", MaxPlayer),
				MessageType: "info", // "alert",
			}
			messages = append(messages, message)
			response = WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "update_state",
				MessageType:    "info",
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
			mutex.Unlock()
		case "set_game_round":
			mutex.Lock()
			gameRoundNum = e.GameRound
			log.Println("e.GameRound:", gameRoundNum)
			message := Message{
				Timestamp:   e.Timestamp,
				MessageId:   int64(len(messages)),
				User:        e.User,
				Message:     fmt.Sprintf("라운드 수가 %d개로 설정되었습니다.", gameRoundNum),
				MessageType: "info", // "alert",
			}
			messages = append(messages, message)
			response = WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "update_state",
				MessageType:    "info",
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
			mutex.Unlock()
		case "set_game_turn":
			mutex.Lock()
			gameTurnNum = e.GameTurnNum
			log.Println("e.GameTurn:", gameTurnNum)
			message := Message{
				Timestamp:   e.Timestamp,
				MessageId:   int64(len(messages)),
				User:        e.User,
				Message:     fmt.Sprintf("턴 수가 %d개로 설정되었습니다.", gameTurnNum),
				MessageType: "info", // "alert",
			}
			messages = append(messages, message)
			response = WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "update_state",
				MessageType:    "info",
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
			mutex.Unlock()

		case "clear_messages":
			mutex.Lock()
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
			mutex.Unlock()
		case "restart_game":
			mutex.Lock()
			isGameStarted = false
			isGameOver = false
			gameInfo = make([]Game, 0)
			gameTurnsLeft = gameTurnNum * MaxPlayer
			gameRound = 1
			userSelections = make([]UserSelection, 0)
			sorted_players = make([]User, 0)
			messages = make([]Message, 0)
			message := Message{
				Timestamp:   e.Timestamp,
				MessageId:   int64(len(messages)),
				User:        adminUser,
				Message:     "게임이 초기화되었습니다.",
				MessageType: "alert",
			}
			messages = []Message{message}
			response = WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "restart_game",
				MessageType:    "alert",
				Message:        "게임이 초기화되었습니다.",
				MaxPlayer:      MaxPlayer,
				User:           User{},
				MessageLogList: messages,
				OnlineUserList: getUserList(),
				PlayerList:     getPlayerList(),
				GameTurnsLeft:  gameTurnsLeft,
				GameRound:      gameRound,
			}
			broadcastToAll(response)
			mutex.Unlock()
		case "restart_round":
			mutex.Lock()
			gameInfo[gameRound-1].NowUserIndex = 0
			gameTurnsLeft = gameTurnNum * len(gameInfo[gameRound-1].PlayerList)
			gameInfo[gameRound-1].TurnsLeft = gameTurnNum * len(gameInfo[gameRound-1].PlayerList)
			userSelections = make([]UserSelection, 0)
			message := Message{
				Timestamp:   e.Timestamp,
				MessageId:   int64(len(messages)),
				User:        e.User,
				Message:     fmt.Sprintf("%d라운드가 초기화되었습니다.", gameRound),
				MessageType: "alert",
			}
			messages = []Message{message}
			response = WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "restart_round",
				MessageType:    "alert",
				Message:        "라운드가 초기화되었습니다.",
				MaxPlayer:      MaxPlayer,
				User:           User{},
				MessageLogList: messages,
				OnlineUserList: getUserList(),
				PlayerList:     getPlayerList(),
				GameTurnsLeft:  gameTurnsLeft,
				GameRound:      gameRound,
			}

		case "get_game_Info":
			mutex.Lock()
			response = WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "get_game_Info",
				MessageType:    "info",
				Message:        "게임 정보를 가져왔습니다.",
				MaxPlayer:      MaxPlayer,
				User:           User{},
				MessageLogList: makeMessagesFromGameInfo(gameInfo),
				OnlineUserList: getUserList(),
				PlayerList:     getPlayerList(),
				GameTurnsLeft:  gameTurnsLeft,
				GameRound:      gameRound,
				GamsSelections: userSelections,
			}
			broadCastToSomeone(response, e.Conn)
			mutex.Unlock()

		default:
			log.Printf("Unknown action: %s\n", e.Action)
		}
		// log.Println("users:", getUserList())
		// if response.Action != "" {
		// 	broadcastToAll(response)
		// }
	}
}

func processEnterHuman(e WsPayload, response *WsJsonResponse) {
	nowUser, exists := players[e.User.UUID]
	log.Println("processEnterHuman")
	log.Println("isGameStarted:", isGameStarted)
	log.Println("gameTurnsLeft:", gameTurnsLeft)

	mutex.Lock()
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
	mutex.Unlock()

	log.Println("GPTReadyNum:", GPTReadyNum)
	log.Println("len(getPlayerList()):", len(getPlayerList()))
	log.Println("GPTEnterNum:", GPTEnterNum)

	if GPTEnterNum <= len(getUserList()) && !GPTUser.IsOnline {
		time.Sleep(time.Second * 1)
		log.Println("GPTEnterNum:", GPTEnterNum)
		log.Println("len(getUserList()):", len(getUserList()))
		log.Println("GPTUser:", GPTUser)
		mutex.Lock()
		nicknameId, userName := getRandomUsername()
		GPTUser = User{
			UserId:     int64(len(clients)),
			UserName:   userName,
			NicknameId: nicknameId,
			Role:       "player",
			UUID:       "999",
			IsOnline:   true,
			PlayerType: "watcher",
		}
		players["999"] = GPTUser
		// clients[WebSocketConnection{Conn: nil}] = GPTUser

		GPTresponse := WsJsonResponse{
			MaxPlayer:      MaxPlayer,
			Timestamp:      e.Timestamp,
			Action:         "human_info",
			MessageType:    "system",
			Message:        fmt.Sprintf("%s님이 입장했습니다.", GPTUser.UserName),
			MessageLogList: messages,
			User:           GPTUser,
			OnlineUserList: getUserList(),
			PlayerList:     getPlayerList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
		}
		broadcastToAll(GPTresponse)
		mutex.Unlock()
	} else if GPTReadyNum >= len(getPlayerList()) && GPTUser.PlayerType == "watcher" {
		log.Println("GPTReadyNum:", GPTReadyNum)
		log.Println("len(getPlayerList()):", len(getPlayerList()))
		log.Println("GPTUser:", GPTUser)

		time.Sleep(time.Second * 1)
		processReadyGPT()
	}

	if gameTurnsLeft == 0 {
		if isGameOver {
			processGameOver(e.Timestamp)
			mutex.Lock()
			resultResponse := WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "game_over",
				MessageType:    "alert",
				Message:        "게임이 종료되었습니다.",
				MaxPlayer:      MaxPlayer,
				User:           adminUser,
				MessageLogList: messages,
				OnlineUserList: getUserList(),
				PlayerList:     getPlayerList(),
				GameTurnsLeft:  gameTurnsLeft,
				GameRound:      gameRound,
			}
			broadcastToAll(resultResponse)
			mutex.Unlock()
			return
		}
		// messages = append(messages, waitMessage)
		// in userSelection
		if exists && nowUser.PlayerType == "player" {
			mutex.Lock()
			var waitResponse WsJsonResponse
			_, exists := findUserInUserSelection(userSelections, nowUser)
			log.Println("exists:", exists)
			log.Println(userSelections)
			log.Println(nowUser)
			if exists {
				waitResponse = WsJsonResponse{
					Timestamp:      e.Timestamp,
					Action:         "wait_for_players",
					MessageType:    "alert",
					Message:        "라운드가 종료되었습니다. 다음 플레이어의 선택을 기다리세요.",
					MaxPlayer:      MaxPlayer,
					User:           User{},
					MessageLogList: messages,
					OnlineUserList: getUserList(),
					PlayerList:     getPlayerList(),
					GameTurnsLeft:  gameTurnsLeft,
					GameRound:      gameRound,
				}
			} else {
				waitResponse = WsJsonResponse{
					Timestamp:      e.Timestamp,
					Action:         "submit_ai",
					MessageType:    "alert",
					Message:        "라운드가 종료되었습니다. 다음 라운드를 위해 AI를 선택하세요.",
					MaxPlayer:      MaxPlayer,
					User:           User{},
					MessageLogList: messages,
					OnlineUserList: getUserList(),
					PlayerList:     getPlayerList(),
					GameTurnsLeft:  gameTurnsLeft,
					GameRound:      gameRound,
				}
			}
			broadCastToSomeone(waitResponse, e.Conn)
			mutex.Unlock()
		}
		return
	}

	if isGameStarted {
		nextUser := gameInfo[gameRound-1].PlayerList[gameInfo[gameRound-1].NowUserIndex]
		if nextUser.UUID == e.User.UUID {
			mutex.Lock()
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
			nextConn, _ := getUserClientByUUID(nextUser.UUID)
			// if exists {
			broadCastToSomeone(someoneResponse, nextConn)
			mutex.Unlock()
			// } else {
			// 	mutex.Unlock()
		}
	}

}

func processLeftUser(e WsPayload, response *WsJsonResponse) {
	mutex.Lock()
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
	broadcastToAll(*response)
	mutex.Unlock()
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
		processGameOver(e.Timestamp)
	}

	nextUser := gameInfo[gameRound-1].PlayerList[gameInfo[gameRound-1].NowUserIndex]
	if nextUser.UUID == GPTUser.UUID {
		SendGPTMessage(e.Timestamp)
		return
	}
	mutex.Lock()
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
	nextConn, _ := getUserClientByUUID(nextUser.UUID)
	// if exists {
	// 	broadCastToSomeone(someoneResponse, nextConn)
	// 	mutex.Unlock()
	// } else {
	// 	mutex.Unlock()

	// }
	broadCastToSomeone(someoneResponse, nextConn)
	mutex.Unlock()

}

func processReadyUser(e WsPayload) {
	nowUser := players[e.User.UUID]
	nowUser.PlayerType = "player"
	log.Println(nowUser.NicknameId)
	players[e.User.UUID] = nowUser
	clients[e.Conn] = nowUser
	sorted_players = append(sorted_players, nowUser)

	if MaxPlayer == len(getPlayerList()) {
		sorted_players = randomSuffleUserList(sorted_players)
		gameRound = 1
		gameTurnsLeft = gameTurnNum * MaxPlayer
		initGameInfo()

		mutex.Lock()
		message := Message{
			Timestamp:   e.Timestamp,
			MessageId:   int64(len(messages)),
			User:        adminUser,
			Message:     "게임이 시작되었습니다.",
			MessageType: "alert",
		}
		messages = []Message{message}

		gameInfo[0].TurnsLeft = gameTurnsLeft
		gameInfo[0].PlayerList = sorted_players
		response := WsJsonResponse{
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
		broadcastToAll(response)
		mutex.Unlock()
		isGameStarted = true

		processNextTurn(e.Timestamp)
		return
	}

	mutex.Lock()
	response := WsJsonResponse{
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
	broadcastToAll(response)
	mutex.Unlock()

	if GPTReadyNum >= len(getPlayerList()) && GPTUser.PlayerType == "watcher" {
		log.Println("GPTReadyNum:", GPTReadyNum)
		log.Println("len(getPlayerList()):", len(getPlayerList()))
		log.Println("GPTUser:", GPTUser)

		time.Sleep(time.Second * 1)
		processReadyGPT()
	}

}

func processReadyGPT() {
	mutex.Lock()
	GPTUser.PlayerType = "player"
	nowUser := GPTUser

	sorted_players = append(sorted_players, nowUser)

	timestamp := getTimeStamp()

	response := WsJsonResponse{
		Timestamp:      timestamp,
		Action:         "update_state",
		MessageType:    "alert",
		MaxPlayer:      MaxPlayer,
		User:           nowUser,
		Message:        fmt.Sprintf("%s님이 게임에 참여했습니다.", nowUser.UserName),
		MessageLogList: messages,
		PlayerList:     getPlayerList(),
		OnlineUserList: getUserList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      gameRound,
	}
	broadcastToAll(response)
	mutex.Unlock()

	if MaxPlayer == len(getPlayerList()) {
		sorted_players = randomSuffleUserList(sorted_players)
		gameRound = 1
		gameTurnsLeft = gameTurnNum * MaxPlayer
		initGameInfo()

		mutex.Lock()
		message := Message{
			Timestamp:   timestamp,
			MessageId:   int64(len(messages)),
			User:        adminUser,
			Message:     "게임이 시작되었습니다.",
			MessageType: "alert",
		}
		messages = []Message{message}

		gameInfo[0].TurnsLeft = gameTurnsLeft
		gameInfo[0].PlayerList = sorted_players
		response := WsJsonResponse{
			Timestamp:      timestamp,
			Action:         "update_state",
			MessageType:    "alert",
			Message:        message.Message,
			MaxPlayer:      MaxPlayer,
			User:           GPTUser,
			MessageLogList: messages,
			PlayerList:     getPlayerList(),
			OnlineUserList: getUserList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
		}
		broadcastToAll(response)
		mutex.Unlock()

		processNextTurn(timestamp)
		isGameStarted = true
		return
	}

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

func broadCastToSomeone(response WsJsonResponse, client WebSocketConnection) {
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

func initGameInfo() {

	gameInfo = make([]Game, 0)
	for i := 0; i < gameRoundNum; i++ {
		game := Game{
			NowUserIndex:   0,
			MaxPlayer:      MaxPlayer,
			OnlineUserList: getUserList(),
			PlayerList:     getPlayerList(),
			TurnsLeft:      gameTurnNum * MaxPlayer,
			UserSelections: make([]UserSelection, 0),
		}
		gameInfo = append(gameInfo, game)
	}

}

func getUserClientByUUID(uuid string) (WebSocketConnection, bool) {
	for client := range clients {
		if client.Conn != nil && clients[client].UUID == uuid {
			return client, true
		}
	}
	return WebSocketConnection{}, false
}

func broadcastNewMessage(e WsPayload) {
	mutex.Lock()
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
	mutex.Unlock()
}

func processGameOver(timestamp int64) {
	// time.Sleep(time.Second * 5)
	message := Message{
		Timestamp:   timestamp,
		MessageId:   int64(len(messages)),
		User:        adminUser,
		Message:     "게임이 종료되었습니다. 누가 AI였는지 확인해보세요.",
		MessageType: "alert",
	}

	response := WsJsonResponse{
		Timestamp:      timestamp,
		Action:         "choose_ai",
		MessageType:    "alert",
		Message:        message.Message,
		MaxPlayer:      MaxPlayer,
		User:           adminUser,
		MessageLogList: messages,
		OnlineUserList: getUserList(),
		PlayerList:     getPlayerList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      gameRound,
	}
	broadcastToAll(response)

	// gameInfo[gameRound-1].UserSelections = make([]UserSelection, 0)
	gameInfo[gameRound-1].Messages = messages

}

func findUserInUserSelection(userSelections []UserSelection, user User) (UserSelection, bool) {
	for _, v := range userSelections {
		if v.User.UUID == user.UUID {
			return v, true
		}
	}
	return UserSelection{}, false
}

func processChooseAI(e WsPayload) {
	selection := UserSelection{
		User:      e.User,
		Selection: e.UserSelection.Selection,
		Reason:    e.UserSelection.Reason,
	}
	userSelections = append(userSelections, selection)

	// next round
	if len(userSelections) == MaxPlayer {
		gameInfo[gameRound-1].PlayerList = sorted_players
		vote, eliminatedPlayer, remainingPlayerList := removePlayerListFromUserSelection(sorted_players)
		sorted_players = remainingPlayerList
		mutex.Lock()
		resultMessage := Message{
			Timestamp:   e.Timestamp,
			MessageId:   int64(len(messages)),
			User:        adminUser,
			Message:     fmt.Sprintf("%d라운드가 종료되었습니다. 탈락자는 %d표를 받은 [%s]입니다.", gameRound, vote, eliminatedPlayer.UserName),
			MessageType: "alert",
		}
		messages = []Message{resultMessage}
		resultResponse := WsJsonResponse{
			Timestamp:      e.Timestamp,
			Action:         "next_round",
			MessageType:    "alert",
			Message:        resultMessage.Message,
			MaxPlayer:      MaxPlayer,
			User:           User{},
			MessageLogList: messages,
			OnlineUserList: getUserList(),
			PlayerList:     getPlayerList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
		}
		broadcastToAll(resultResponse)
		mutex.Unlock()

		gameRound++

		if gameRound >= gameRoundNum {
			isGameOver = true
			mutex.Lock()
			alivePlayerString := ""
			for _, v := range remainingPlayerList {
				alivePlayerString += fmt.Sprintf("[%s] ", v.UserName)
			}
			finalMessage := Message{
				Timestamp:   e.Timestamp,
				MessageId:   int64(len(messages)),
				User:        adminUser,
				Message:     fmt.Sprintf("게임이 종료되었습니다. 최종 생존자는 %s입니다.", alivePlayerString),
				MessageType: "alert",
			}
			messages = []Message{finalMessage}
			finalResponse := WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "game_over",
				MessageType:    "alert",
				Message:        finalMessage.Message,
				MaxPlayer:      MaxPlayer,
				User:           User{},
				MessageLogList: messages,
				OnlineUserList: getUserList(),
				PlayerList:     getPlayerList(),
				GameTurnsLeft:  gameTurnsLeft,
				GameRound:      gameRound,
			}
			broadcastToAll(finalResponse)
			mutex.Unlock()

			return

		}

		sorted_players = randomSuffleUserList(sorted_players)
		gameTurnsLeft = gameTurnNum * len(sorted_players)
		gameInfo[gameRound-1].TurnsLeft = gameTurnsLeft
		gameInfo[gameRound-1].PlayerList = sorted_players
		gameInfo[gameRound-1].NowUserIndex = 0
		userSelections = make([]UserSelection, 0)

		mutex.Lock()
		nextMessage := Message{
			Timestamp:   e.Timestamp,
			MessageId:   int64(len(messages)),
			User:        adminUser,
			Message:     fmt.Sprintf("%d라운드가 시작되었습니다.", gameRound),
			MessageType: "alert",
		}
		messages = append(messages, nextMessage)
		nextResponse := WsJsonResponse{
			Timestamp:      e.Timestamp,
			Action:         "update_state",
			MessageType:    "alert",
			Message:        nextMessage.Message,
			MaxPlayer:      MaxPlayer,
			User:           User{},
			MessageLogList: messages,
			OnlineUserList: getUserList(),
			PlayerList:     getPlayerList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
		}
		broadcastToAll(nextResponse)
		mutex.Unlock()

		processNextTurn(e.Timestamp)

	}
}
func removePlayerListFromUserSelection(playerList []User) (vote int, eliminatedPlayer User, remainingPlayerList []User) {
	for _, v := range userSelections {
		log.Println("userSelection:", v.Selection)
	}
	votes := make(map[string]int)
	for _, v := range userSelections {
		votes[v.Selection]++
	}
	log.Println("votes:", votes)

	//sort decreasing order by vote
	sort.Slice(playerList, func(i, j int) bool {
		return votes[playerList[i].UserName] > votes[playerList[j].UserName]
	})

	if len(playerList) > 0 {
		eliminatedPlayer = playerList[0]
		remainingPlayerList = playerList[1:]
	} else {
		remainingPlayerList = []User{}
	}

	log.Println("remainingPlayerList:", remainingPlayerList)
	log.Println("eliminatedPlayer:", eliminatedPlayer)

	return votes[eliminatedPlayer.UserName], eliminatedPlayer, remainingPlayerList
}

func processNextTurn(timestamp int64) {

	mutex.Lock()
	nextUser := gameInfo[gameRound-1].PlayerList[gameInfo[gameRound-1].NowUserIndex]
	log.Println("nextUser:", nextUser)
	// if nextUser.UUID == GPTUser.UUID {
	// 	SendGPTMessage(timestamp)
	// 	return
	// }
	someoneMessage := Message{
		Timestamp:   timestamp,
		MessageId:   int64(len(messages)),
		User:        adminUser,
		Message:     fmt.Sprintf("%s님의 차례입니다.", nextUser.UserName),
		MessageType: "alert",
	}
	// messages = append(messages, someoneMessage)
	someoneResponse := WsJsonResponse{
		Timestamp:   timestamp,
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
	nextConn, _ := getUserClientByUUID(nextUser.UUID)
	// if exists {
	// 	broadCastToSomeone(someoneResponse, nextConn)
	// 	mutex.Unlock()
	// } else {
	// 	mutex.Unlock()
	// 	SendGPTMessage(timestamp)
	// }
	broadCastToSomeone(someoneResponse, nextConn)
	mutex.Unlock()
}

func makeMessagesFromGameInfo(gameInfo []Game) []Message {
	messages := make([]Message, 0)
	for _, v := range gameInfo {
		messages = append(messages, v.Messages...)
	}
	return messages

}

func getTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type MessageData struct {
	UserUUID string `json:"user_UUID"`
	Message  string `json:"message"`
}

func SendGPTMessage(timestamp int64) {
	log.Println("=================================================")
	log.Println("SendGPTMessage")

	// var gptResponse WsJsonResponse
	MessageLogString := ""
	for _, v := range messages {
		if v.User.UserName != "admin" {
			MessageLogString += fmt.Sprintf("%s: %s\n", v.User.UserName, v.Message)
		}
	}
	log.Println("MessageLogString:", MessageLogString)
	gptResponse := sendGPTMessageToFastAPI(MessageLogString)
	log.Println("gptResponse:", gptResponse)

	gameInfo[gameRound-1].NowUserIndex = (gameInfo[gameRound-1].NowUserIndex + 1) % gameInfo[gameRound-1].MaxPlayer
	gameTurnsLeft = Max(gameTurnsLeft-1, 0)
	gameInfo[gameRound-1].TurnsLeft = gameTurnsLeft

	mutex.Lock()

	message := Message{
		Timestamp:   timestamp,
		MessageId:   int64(len(messages)),
		User:        GPTUser,
		Message:     gptResponse.Message,
		MessageType: "message",
	}
	messages = append(messages, message)
	chatResponse := WsJsonResponse{
		Timestamp:      timestamp,
		Action:         "new_message",
		MessageType:    "message",
		Message:        message.Message,
		MaxPlayer:      MaxPlayer,
		User:           GPTUser,
		MessageLogList: messages,
		OnlineUserList: getUserList(),
		PlayerList:     getPlayerList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      gameRound,
	}
	broadcastToAll(chatResponse)
	mutex.Unlock()

	if gameTurnsLeft == 0 {
		processGameOver(timestamp)
		return
	}
	processNextTurn(getTimeStamp())
}

func sendGPTMessageToFastAPI(message string) MessageData {
	url := FastAPIURL + "/useGPT"
	msgData := MessageData{
		UserUUID: GPTUser.UUID,
		Message:  message,
	}
	// msgData := map[string]string{"message": message}
	msgBytes, err := json.Marshal(msgData)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		log.Fatalf("Error occurred during sending request to FastAPI. Error: %s", err.Error())
	}
	defer resp.Body.Close()

	log.Printf("Message sent to FastAPI, received response status: %s", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error occurred during reading response body. Error: %s", err.Error())
	}
	log.Printf("Response body: %s", string(body))

	var gptResponse MessageData
	err = json.Unmarshal(body, &gptResponse)
	if err != nil {
		log.Fatalf("Error occurred during unmarshaling. Error: %s", err.Error())
	}
	log.Printf("gptResponse: %s", gptResponse)
	return gptResponse
}
