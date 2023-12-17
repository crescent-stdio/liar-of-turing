package handlers

import (
	"fmt"
	"liar-of-turing/models"
	"liar-of-turing/services"
	"liar-of-turing/utils"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleWebSocketRequest(w http.ResponseWriter, r *http.Request, userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	ws, err := webSocketService.UpgradeConfig.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	conn := models.WebSocketConnection{Conn: ws}
	log.Println("Client connected to endpoint")

	go ListenWebSocketConnections(webSocketService, &conn)
}

func ListenWebSocketConnections(webSocketService *services.WebSocketService, conn *models.WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in ListenWebSocketConnections:", r)
		}
		webSocketService.RemoveClient(*conn) // Safely close the connection
	}()

	var payload models.WsPayload

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
		webSocketService.Broadcast <- payload
	}
}

func ListenToWebSocketChannel(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	clients := webSocketService.GetClients()
	adminUser := userManager.GetAdminUser()
	GPTUsers := userManager.GetGPTUsers()

	for {

		e := <-webSocketService.Broadcast

		// Consolidated logging
		log.Printf("Action: %s, User: %v\n", e.Action, e.User)

		switch e.Action {
		case "broadcast", "new_message_admin":
			BroadcastWebSocketMessage(userManager, webSocketService, gameState, e)

		case "new_message":
			BroadcastWebSocketMessage(userManager, webSocketService, gameState, e)
			gameState.SetNextTurnInfo()
			// ProcessNextTurn(userManager, webSocketService, gameState) //, e.Timestamp)
		case "enter_human":
			HandleHumanUserEntry(userManager, webSocketService, gameState, e)
			// broadcastToAll(clients, response)
		case "left_user":
			HandleUserLeave(userManager, webSocketService, gameState, e)

		case "user_is_ready":
			HandleUserReady(userManager, webSocketService, gameState, e)

		case "choose_ai":
			HandleAISelection(gameState, e)

		case "set_max_player":
			gameState.SetMaxPlayer(e.MaxPlayer)
			userManager.SetAllUsersAsWatchers()
			for _, gpt := range GPTUsers {
				gpt.PlayerType = "watcher"
			}
			userManager.SetGPTUsers(GPTUsers)

			message := utils.CreateMessageFromUser(userManager, e.User, e.Timestamp)
			message.Message = fmt.Sprintf("최대 인원이 %d명으로 설정되었습니다.", e.MaxPlayer)
			message.MessageType = "info"

			response := utils.CreateResponseUsingPayload(userManager, gameState, e)
			response.Action = "update_state"
			response.Message = message.Message
			response.MessageType = "info"

			broadcastToAll(clients, response)

			userManager.AddPrevMessagesFromMessages()
			userManager.ClearMessages()
			userManager.AddMessage(message)

		case "set_game_round":
			gameState.SetGameRoundNum(e.GameRoundNum)
			userManager.AddPrevMessagesFromMessages()
			userManager.ClearMessages()

			message := utils.CreateMessageFromUser(userManager, e.User, e.Timestamp)
			message.Message = fmt.Sprintf("라운드 수가 %d개로 설정되었습니다.", e.GameRoundNum)
			message.MessageType = "info"

			response := utils.CreateResponseUsingPayload(userManager, gameState, e)
			response.Action = "update_state"
			response.Message = message.Message
			response.MessageType = "info"
			response.User = adminUser

			broadcastToAll(clients, response)

		case "set_game_turn":
			gameState.SetGameTurnNum(e.GameTurnNum)
			message := utils.CreateMessageFromUser(userManager, adminUser, e.Timestamp)
			message.Message = fmt.Sprintf("턴 수가 %d개로 설정되었습니다.", e.GameTurnNum)
			message.MessageType = "info"

			response := utils.CreateResponseUsingPayload(userManager, gameState, e)
			response.Action = "update_state"
			response.Message = message.Message
			response.MessageType = "info"
			response.User = adminUser

			broadcastToAll(clients, response)

			userManager.AddMessage(message)
		case "clear_messages":
			userManager.AddPrevMessagesFromMessages()
			userManager.ClearMessages()

			// Broadcast to all
			response := utils.CreateResponseUsingPayload(userManager, gameState, e)
			response.Action = "update_state"
			response.Message = "메시지가 초기화되었습니다."
			response.MessageType = "info"
			response.User = adminUser

			broadcastToAll(clients, response)

		case "restart_game":
			HandleRestartGameEvent(userManager, webSocketService, gameState, e)

		case "restart_round":
			gameRound := gameState.GetNowGameInfo().Round
			gameState.SetIfResetRound(userManager)
			message := utils.CreateMessageFromUser(userManager, adminUser, e.Timestamp)
			message.Message = fmt.Sprintf("%d라운드가 초기화되었습니다.", gameRound)
			message.MessageType = "info"

			userManager.AddPrevMessagesFromMessages()
			userManager.ClearMessages()
			userManager.AddMessage(message)

			response := utils.CreateResponseUsingPayload(userManager, gameState, e)
			response.Action = "restart_round"
			response.Message = message.Message
			response.MessageType = "info"
			response.User = adminUser
			broadcastToAll(clients, response)

		case "get_game_Info":
			response := utils.CreateResponseUsingPayload(userManager, gameState, e)
			response.Action = "get_game_Info"
			response.Message = "게임 정보를 가져왔습니다."
			gameInfo := gameState.GetAllGameInfo()
			response.MessageLogList = CreateMessagesFromGameStatus(gameInfo)
			response.MessageType = "info"
			response.User = adminUser
			broadcastToSomeone(clients, e.Conn, response)

		default:
			log.Printf("Unknown action: %s\n", e.Action)
		}
		ProcessGPTEntering(userManager, webSocketService, gameState)
		log.Println("ProcessGPTReady")
		ProcessGPTReady(userManager, webSocketService, gameState)
		log.Println("ProcessAllPlayersReady")
		HandleAllHumanUserReady(userManager, webSocketService, gameState)
		ProcessAllPlayersReady(userManager, webSocketService, gameState)
		log.Println("ProcessAllPlayersVoted")
		ProcessNextTurn(userManager, webSocketService, gameState) //, e.Timestamp)
		log.Println("ProcessAllPlayersVoted")
		ProcessAllPlayersVoted(userManager, webSocketService, gameState)
		log.Println("HandleRoundIsOver")
		HandleRoundIsOver(userManager, webSocketService, gameState)
		log.Println("GPTEnterNums:", gameState.GetGPTEntryNums())
		log.Println("GPTReadyNums:", gameState.GetGPTReadyNums())
		log.Println("userSelections:", gameState.GetNowUserSelections())
		log.Println("len(userSelections):", len(gameState.GetNowUserSelections()))
		log.Println("gameState.CheckAllUserVoted(userManager):", gameState.CheckAllUserVoted(userManager))
		log.Println("=================================================")
	}
}

// HandleUserReady: If user is ready, then set user to "PLAYER"
func HandleUserReady(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, e models.WsPayload) {
	clients := webSocketService.GetClients()
	adminUser := userManager.GetAdminUser()

	// Set user's type to "PLAYER"
	nowUser, exists := userManager.GetPlayerByUUID(e.User.UUID)
	if !exists {
		log.Println("User is not a player")
		return
	}
	nowUser.PlayerType = "player"
	userManager.AddPlayerByUser(nowUser)
	userManager.AddSortedPlayerByUser(nowUser)
	webSocketService.AddClient(e.Conn, nowUser)

	message := utils.CreateMessageFromUser(userManager, adminUser, e.Timestamp)
	message.User = adminUser
	message.MessageType = "info"
	message.Message = fmt.Sprintf("%s님이 게임에 참여했습니다.", nowUser.UserName)
	userManager.AddMessage(message)
	response := utils.CreateResponseUsingPayload(userManager, gameState, e)
	response.Action = "human_info"
	response.MessageType = "system"
	response.User = nowUser
	response.Message = message.Message
	broadcastToAll(clients, response)

}

func BroadcastWebSocketMessage(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, e models.WsPayload) {
	clients := webSocketService.GetClients()
	message := utils.CreateMessageFromUser(userManager, e.User, e.Timestamp)
	message.Message = e.Message
	message.MessageType = "message"
	userManager.AddMessage(message)
	log.Println(userManager.GetMessages())
	response := utils.CreateResponseUsingPayload(userManager, gameState, e)
	response.Action = "new_message"
	response.Message = message.Message
	response.User = e.User

	broadcastToAll(clients, response)
}

func ProcessAllPlayersReady(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	isGameStarted := gameState.GetStatus().IsStarted
	if !gameState.CheckAllUserReady(userManager) || isGameStarted || gameState.CheckIsRoundOver() || gameState.CheckIsGameOver() || gameState.CheckAllUserVoted(userManager) {
		return
	}
	log.Println("ProcessAllPlayersReady is called")
	userManager.SetPlayersRandomlyShuffled(webSocketService, gameState)
	gameState.SetQuestionsRandomly()
	gameState.InitializeRoundInfo(userManager, webSocketService)

	clients := webSocketService.GetClients()
	adminUser := userManager.GetAdminUser()

	// Make Message & Response
	message := utils.CreateMessageWithAutoTimestamp(userManager, adminUser)
	message.User = adminUser
	message.MessageType = "info"
	message.Message = "게임이 시작되었습니다."

	userManager.AddPrevMessagesFromMessages()
	userManager.ClearMessages()
	userManager.AddMessage(message)

	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.MessageType = "info"
	response.Message = message.Message
	broadcastToAll(clients, response)

	// Broadcast Question to all
	question := gameState.GetQuestion()
	QMessage := utils.CreateMessageWithAutoTimestamp(userManager, adminUser)
	QMessage.User = adminUser
	QMessage.MessageType = "alert"
	QMessage.Message = fmt.Sprintf("조건: '%s'", question)
	userManager.AddMessage(QMessage)

	QResponse := utils.CreateInitalizeResponse(userManager, gameState)
	QResponse.MessageType = "alert"
	QResponse.Message = QMessage.Message
	broadcastToAll(clients, QResponse)

	// HandleRoundIsOver(e.Timestamp)
}

func ProcessGPTEntering(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	OnlineUserNum := len(utils.RetrieveUserList(userManager))

	GPTUsers := userManager.GetGPTUsers()
	GPTEntryNums := gameState.GetGPTEntryNums()
	for idx, GPTEntryNum := range GPTEntryNums {
		GPTUser := GPTUsers[idx]
		// If GPT's Entring time cames...(GPTUser isn't ONLINE)
		if !GPTUser.IsOnline && GPTEntryNum <= OnlineUserNum+1 {
			utils.RandomTimeSleep()
			HandleGPTEntry(userManager, webSocketService, gameState, idx)
		}
	}
}

func HandleGPTEntry(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, index int) {
	clients := webSocketService.GetClients()
	adminUser := userManager.GetAdminUser()
	GPTUser := userManager.GetGPTUsers()[index]
	nicknameId, userName := userManager.GenerateRandomUsername()
	GPTUser.NicknameId = nicknameId
	GPTUser.UserName = userName
	GPTUser.PlayerType = "watcher"
	GPTUser.IsOnline = true
	userManager.SetGPTUser(index, GPTUser)
	log.Println("GPTUser:", GPTUser)

	// set GPTUser to true
	GPTUser.IsOnline = true
	userManager.SetGPTUser(index, GPTUser)
	userManager.AddPlayerByUser(GPTUser)

	// Make Message & Response
	message := utils.CreateMessageWithAutoTimestamp(userManager, adminUser)
	message.User = adminUser
	message.MessageType = "info"
	message.Message = fmt.Sprintf("%s님이 채팅방에 입장했습니다.", GPTUser.UserName)
	userManager.AddMessage(message)

	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.Action = "human_info"
	response.MessageType = "system"
	response.User = GPTUser
	response.Message = message.Message

	broadcastToAll(clients, response)
}
func ProcessGPTReady(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	ReadyPlayerNum := len(utils.RetrieveReadyUserList(userManager))
	log.Println("ReadyPlayerNum:", ReadyPlayerNum)

	GPTUsers := userManager.GetGPTUsers()
	GPTReadyNums := gameState.GetGPTReadyNums()
	for idx, GPTReadyNum := range GPTReadyNums {
		GPTUser := GPTUsers[idx]
		// If GPT's Ready time cames...(GPTUser isn't WATCHER)
		if GPTUser.PlayerType == "watcher" && GPTReadyNum <= ReadyPlayerNum+1 {
			log.Println("GPTReadyNum:", GPTReadyNum)
			utils.RandomTimeSleep()
			HandleGPTReady(userManager, webSocketService, gameState, idx)
		}
	}
}
func HandleAllHumanUserReady(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	if !gameState.CheckAllHumanPlayerReady(userManager) {
		return
	}
	GPTUsers := userManager.GetGPTUsers()
	GPTReadyNums := gameState.GetGPTReadyNums()
	log.Println(GPTUsers)
	for idx, GPTReadyNum := range GPTReadyNums {
		GPTUser := GPTUsers[idx]
		if GPTUser.PlayerType == "watcher" {
			log.Println("GPTReadyNum:", GPTReadyNum)
			utils.RandomTimeSleep()
			HandleGPTReady(userManager, webSocketService, gameState, idx)
		}
	}
}

// HandleGPTReady: If GPT is ready, then set user to "PLAYER"
func HandleGPTReady(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, index int) {
	clients := webSocketService.GetClients()
	adminUser := userManager.GetAdminUser()
	GPTUser := userManager.GetGPTUsers()[index]

	// Set GPT's type to "PLAYER"
	GPTUser.PlayerType = "player"
	userManager.AddPlayerByUser(GPTUser)
	userManager.AddSortedPlayerByUser(GPTUser)
	userManager.SetGPTUser(index, GPTUser)

	// Make Message & Response
	message := utils.CreateMessageWithAutoTimestamp(userManager, adminUser)
	message.MessageType = "info"
	message.Message = fmt.Sprintf("%s님이 게임에 참여했습니다.", GPTUser.UserName)
	userManager.AddMessage(message)

	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.Action = "human_info"
	response.MessageType = "system"
	response.User = GPTUser
	response.Message = message.Message

	broadcastToAll(clients, response)
}

func ProcessAllPlayersVoted(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	if !gameState.CheckAllUserVoted(userManager) {
		return
	}
	clients := webSocketService.GetClients()
	gameInfo := gameState.GetNowGameInfo()
	gameRound := gameInfo.Round
	AdmminUser := userManager.GetAdminUser()

	gameState.SetIfRoundIsOver()
	voteNum, eliminatedPlayer, remainingPlayerList := userManager.ExcludePlayersFromSelections(webSocketService, gameState)

	userManager.SetSortedPlayers(remainingPlayerList)

	messages := utils.CreateMessageWithAutoTimestamp(userManager, AdmminUser)
	messages.MessageType = "info"
	messages.Message = fmt.Sprintf("%d라운드가 종료되었습니다. 탈락자는 %d표를 받은 [%s]입니다.", gameRound, voteNum, eliminatedPlayer.UserName)

	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.Action = "show_result"
	response.MessageType = "info"
	response.Message = messages.Message

	broadcastToAll(clients, response)

	// Broadcast vote result to console
	broadcastSelectionResultToAll(userManager, webSocketService, gameState)

	// time.Sleep(time.Second * 10)
}

// Broadcast vote result to console
func broadcastSelectionResultToAll(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	clients := webSocketService.GetClients()

	userSelections := gameState.GetNowUserSelections()

	for _, userSelection := range userSelections {
		VotedUser := userSelection.User
		Selection := userSelection.Selection
		Reason := userSelection.Reason
		Timestamp := userSelection.Timestamp
		message := utils.CreateMessageFromUser(userManager, VotedUser, Timestamp)
		message.Message = fmt.Sprintf("%s님이 [%s]에게 투표했습니다. 사유: %s", VotedUser.UserName, Selection, Reason)
		message.MessageType = "info"

		response := utils.CreateResponseUsingTimestamp(userManager, gameState, Timestamp)
		response.Action = "new_message_admin"
		response.MessageType = "info"
		response.Message = message.Message
		response.User = VotedUser

		broadcastToAll(clients, response)
		userManager.AddMessage(message)
	}
}

func HandleRoundIsOver(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	// Game is not started
	log.Println("HandleRoundIsOver")
	log.Println("gameState.CheckIsRoundOver():", gameState.CheckIsRoundOver())
	log.Println("gameState.CheckAllUserVoted(userManager):", gameState.CheckAllUserVoted(userManager))
	log.Println("gameState.CheckAllUserReady(userManager):", gameState.CheckAllUserReady(userManager))
	log.Println("gameState.CheckIsGameOver():", gameState.CheckIsGameOver())

	if !gameState.CheckIsRoundOver() {
		return
	} else {
		gameState.SetIfRoundIsOver()
	}
	if !gameState.CheckAllUserReady(userManager) {
		return
	}
	if gameState.CheckIsRoundOver() && gameState.CheckAllUserVoted(userManager) { // TODO: CheckIsGameOver
		log.Println("HandleGameOverEvent")
		HandleGameOverEvent(userManager, webSocketService, gameState)
	} else if gameState.CheckIsRoundOver() && !gameState.CheckAllUserVoted(userManager) {
		HandleRoundOverEvent(userManager, webSocketService, gameState)
	}
}

func HandleRoundOverEvent(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	adminUser := userManager.GetAdminUser()
	gameRound := gameState.GetNowGameInfo().Round
	// gameState.SetIfResetRound(userManager)
	prevMessage, exists := userManager.GetPrevMessage()

	message := utils.CreateMessageWithAutoTimestamp(userManager, adminUser)
	message.Message = fmt.Sprintf("%d라운드가 종료되었습니다.", gameRound)
	if exists && prevMessage.Message == message.Message {
		message.MessageType = "hide"
	} else {
		message.MessageType = "info"
		userManager.AddMessage(message)
	}

	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.MessageType = "info"
	response.Message = message.Message
	response.User = adminUser

	clients := webSocketService.GetClients()
	broadcastToAll(clients, response)

	broadcastChooseAIToAll(userManager, webSocketService, gameState)

	// userManager.AddPrevMessagesFromMessages()
	// userManager.ClearMessages()
	// userManager.AddMessage(message)
}
