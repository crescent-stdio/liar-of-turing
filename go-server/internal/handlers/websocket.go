package handlers

import (
	"fmt"
	"liar-of-turing/utils"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// GetWebSocketClientByUUID returns the WebSocket client with the given UUID.
func GetWebSocketClientByUUID(uuid string) (WebSocketConnection, bool) {
	for client := range clients {
		if client.Conn != nil && clients[client].UUID == uuid {
			return client, true
		}
	}
	return WebSocketConnection{}, false
}

func HandleWebSocketRequest(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	conn := WebSocketConnection{Conn: ws}
	log.Println("Client connected to endpoint")

	go ListenWebSocketConnections(&conn)
}

func ListenWebSocketConnections(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in ListenForWs:", r)
		}
		conn.CloseWebSocketConnection() // Safely close the connection
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

func ListenToWebSocketChannel() {
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
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
			}
			broadcastToAll(clients, response)
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
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
			}
			broadcastToAll(clients, response)
			mutex.Unlock()
		case "new_message":
			HandleNewWebSocketMessage(e, &response)
			// broadcastToAll(clients, response)
		case "enter_human":
			processEnterHuman(e)
			// broadcastToAll(clients, response)

		case "left_user":
			processLeftUser(e, &response)

		case "user_is_ready":
			processReadyUser(e)

		case "choose_ai":
			processChooseAI(e)

		case "set_max_player":
			mutex.Lock()
			MaxPlayer = e.MaxPlayer
			// randomChooseReadyUserAndEnterUser
			GPTEnterNum, GPTReadyNum = SelectRandomReadyAndEnteringUser(MaxPlayer)
			log.Println("e.MaxPlayer:", MaxPlayer)
			sorted_players = []User{}
			players = SetAllUserToWatcher(players)
			GPTUser.PlayerType = "watcher"

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
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
			}
			broadcastToAll(clients, response)
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
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
			}
			broadcastToAll(clients, response)
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
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
			}
			broadcastToAll(clients, response)
			mutex.Unlock()

		case "clear_messages":
			mutex.Lock()
			prevMessages = messages
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
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
			}
			broadcastToAll(clients, response)
			mutex.Unlock()
		case "restart_game":
			mutex.Lock()
			isGameStarted = false
			isGameOver = false
			gameInfo = make([]Game, 0)
			gameTurnsLeft = gameTurnNum * MaxPlayer
			gameRound = 1
			userSelections = make([]UserSelection, 0)
			sorted_players = []User{}
			players = SetAllUserToWatcher(players)
			GPTUser.PlayerType = "watcher"
			message := Message{
				Timestamp:   e.Timestamp,
				MessageId:   int64(len(messages)),
				User:        adminUser,
				Message:     "게임이 초기화되었습니다.",
				MessageType: "alert",
			}
			prevMessages = messages
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
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
			}
			broadcastToAll(clients, response)
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
			prevMessages = messages
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
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
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
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
				UserSelection:  userSelections,
			}
			broadCastToSomeone(clients, e.Conn, response)
			mutex.Unlock()

		default:
			log.Printf("Unknown action: %s\n", e.Action)
		}
		// log.Println("users:", getUserList())
		// if response.Action != "" {
		// 	broadcastToAll(clients, response)
		// }
	}
}

func HandleNewWebSocketMessage(e WsPayload, response *WsJsonResponse) {
	log.Println("processNewMessage")
	// log.Println("message", e.Message)
	BroadcastMessageToWebSockets(e)

	// GameRound := e.GameRound
	gameInfo[gameRound-1].NowUserIndex = (gameInfo[gameRound-1].NowUserIndex + 1) % gameInfo[gameRound-1].MaxPlayer
	gameTurnsLeft = utils.Max(gameTurnsLeft-1, 0)
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
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	nextConn, _ := GetWebSocketClientByUUID(nextUser.UUID)
	// if exists {
	// 	broadCastToSomeone(someoneResponse, nextConn)
	// 	mutex.Unlock()
	// } else {
	// 	mutex.Unlock()

	// }
	broadCastToSomeone(clients, nextConn, someoneResponse)
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
		sorted_players = ShuffleUsersRandomly(sorted_players)
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
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
		broadcastToAll(clients, response)
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
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	broadcastToAll(clients, response)
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

	timestamp := utils.GetCurrentTimestamp()

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
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	broadcastToAll(clients, response)
	mutex.Unlock()

	if MaxPlayer == len(getPlayerList()) {
		sorted_players = ShuffleUsersRandomly(sorted_players)
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
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
		broadcastToAll(clients, response)
		mutex.Unlock()

		processNextTurn(timestamp)
		isGameStarted = true
		return
	}

}

func HandleNewWebSocketMessage(e WsPayload, response *WsJsonResponse) {
	log.Println("processNewMessage")
	// log.Println("message", e.Message)
	BroadcastMessageToWebSockets(e)

	// GameRound := e.GameRound
	gameInfo[gameRound-1].NowUserIndex = (gameInfo[gameRound-1].NowUserIndex + 1) % gameInfo[gameRound-1].MaxPlayer
	gameTurnsLeft = utils.Max(gameTurnsLeft-1, 0)
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
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	nextConn, _ := GetWebSocketClientByUUID(nextUser.UUID)
	// if exists {
	// 	broadCastToSomeone(someoneResponse, nextConn)
	// 	mutex.Unlock()
	// } else {
	// 	mutex.Unlock()

	// }
	broadCastToSomeone(clients, nextConn, someoneResponse)
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
		sorted_players = ShuffleUsersRandomly(sorted_players)
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
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
		broadcastToAll(clients, response)
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
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	broadcastToAll(clients, response)
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

	timestamp := utils.GetCurrentTimestamp()

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
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	broadcastToAll(clients, response)
	mutex.Unlock()

	if MaxPlayer == len(getPlayerList()) {
		sorted_players = ShuffleUsersRandomly(sorted_players)
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
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
		broadcastToAll(clients, response)
		mutex.Unlock()

		processNextTurn(timestamp)
		isGameStarted = true
		return
	}

}
