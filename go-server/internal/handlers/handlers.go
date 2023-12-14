package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"liar-of-turing/global"
	"liar-of-turing/models"
	"liar-of-turing/utils"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"

	. "liar-of-turing/models"

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
	GPTEnterNum, GPTReadyNum = SelectRandomReadyAndEnteringUser(MaxPlayer)
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
var GPTNum = 1
var GPTEnterNum = 1
var GPTReadyNum = 0
var Broadcast = make(chan WsPayload)
var GPTBroadcast = make(chan GPTWsPayload)

var clients = make(map[WebSocketConnection]User)

var players = make(map[string]User)
var sorted_players = make([]User, 0)

var messages = make([]Message, 0)
var prevMessages = make([]Message, 0)
var MaxPlayer = 6

var gameRoundNum = 1
var gameTurnNum = 2

var isGameStarted = false
var isGameOver = false
var gameInfo = make([]Game, 0)
var gameTurnsLeft = gameTurnNum * MaxPlayer
var gameRound = 3
var userSelections = make([]UserSelection, 0)

// upgradeConnection is the websocket upgrader from gorilla/websockets
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
	HandshakeTimeout: 1024,
}

// HandleWebSocketRequest upgrades connection to websocket

func processEnterHuman(e WsPayload) {
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

	response := WsJsonResponse{
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
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	broadcastToAll(clients, response)
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
			UUID:       GPTUser.UUID,
			IsOnline:   true,
			PlayerType: "watcher",
		}
		players[GPTUser.UUID] = GPTUser
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
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
		broadcastToAll(clients, GPTresponse)
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
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
			}
			broadcastToAll(clients, resultResponse)
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
					IsGameStarted:  isGameStarted,
					IsGameOver:     isGameOver,
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
					IsGameStarted:  isGameStarted,
					IsGameOver:     isGameOver,
				}
			}
			broadCastToSomeone(clients, e.Conn, waitResponse)
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
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
			}
			nextConn, _ := GetWebSocketClientByUUID(nextUser.UUID)
			// if exists {
			broadCastToSomeone(clients, nextConn, someoneResponse)
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
		e.Conn.CloseWebSocketConnection()

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
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
	} else {
		log.Println("User not found in clients map")
		// *response = WsJsonResponse{} // Empty response
		response := WsJsonResponse{
			Action:    "update_state",
			Timestamp: e.Timestamp,
			// Action:         "user_list",
			MessageLogList: messages,
			OnlineUserList: getUserList(),
			MaxPlayer:      MaxPlayer,
			PlayerList:     getPlayerList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
		broadcastToAll(clients, response)
	}
	mutex.Unlock()
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

func GetWebSocketClientByUUID(uuid string) (WebSocketConnection, bool) {
	for client := range clients {
		if client.Conn != nil && clients[client].UUID == uuid {
			return client, true
		}
	}
	return WebSocketConnection{}, false
}

func BroadcastMessageToWebSockets(e WsPayload) {
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
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	broadcastToAll(clients, response)
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
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	broadcastToAll(clients, response)

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
	if len(userSelections) >= MaxPlayer-GPTNum {
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
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
		broadcastToAll(clients, resultResponse)
		mutex.Unlock()

		sendUserSelectionsForConsole(e.Timestamp)
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
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
			}
			broadcastToAll(clients, finalResponse)
			mutex.Unlock()

			return

		}

		sorted_players = ShuffleUsersRandomly(sorted_players)
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
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
		broadcastToAll(clients, nextResponse)
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

		// TODO: for range for GPT
		selection := UserSelection{
			User:      GPTUser,
			Selection: eliminatedPlayer.UserName,
			Reason:    "AI같다.",
		}
		userSelections = append(userSelections, selection)
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
	nextConn, _ := GetWebSocketClientByUUID(nextUser.UUID)
	// if exists {
	// 	broadCastToSomeone(someoneResponse, nextConn)
	// 	mutex.Unlock()
	// } else {
	// 	mutex.Unlock()
	// 	SendGPTMessage(timestamp)
	// }
	broadCastToSomeone(clients, nextConn, someoneResponse)
	mutex.Unlock()
}

func makeMessagesFromGameInfo(gameInfo []Game) []Message {
	messages := make([]Message, 0)
	for _, v := range gameInfo {
		messages = append(messages, v.Messages...)
	}
	return messages

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
	gameTurnsLeft = utils.Max(gameTurnsLeft-1, 0)
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
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	broadcastToAll(clients, chatResponse)
	mutex.Unlock()

	if gameTurnsLeft == 0 {
		processGameOver(timestamp)
		return
	}
	processNextTurn(utils.GetCurrentTimestamp())
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

func SetAllUserToWatcher(usermap map[string]User) map[string]User {
	playerMap := make(map[string]User)
	for UUID, v := range usermap {
		player := User{
			UserId:     v.UserId,
			UserName:   v.UserName,
			NicknameId: v.NicknameId,
			Role:       v.Role,
			UUID:       v.UUID,
			IsOnline:   v.IsOnline,
			PlayerType: "watcher",
		}
		playerMap[UUID] = player
	}
	return playerMap
}
