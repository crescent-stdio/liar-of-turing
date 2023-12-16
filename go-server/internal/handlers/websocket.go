package handlers

import (
	"fmt"
	"liar-of-turing/models"
	"liar-of-turing/services"
	"liar-of-turing/utils"
	"log"
	"net/http"
	"time"

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
		case "broadcast", "new_message_admin", "new_message":
			BroadcastWebSocketMessage(userManager, webSocketService, gameState, e)

		// case "new_message":
		// 	HandleNewWebSocketMessage(userManager, webSocketService, gameState, e)
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
			message.MessageType = "alert"

			response := utils.CreateResponseUsingPayload(userManager, gameState, e)
			response.Action = "restart_round"
			response.Message = message.Message
			response.MessageType = "alert"
			response.User = adminUser

			broadcastToAll(clients, response)
			userManager.AddPrevMessagesFromMessages()
			userManager.ClearMessages()
			userManager.AddMessage(message)
		case "get_game_Info":
			response := utils.CreateResponseUsingPayload(userManager, gameState, e)
			response.Action = "get_game_Info"
			response.Message = "게임 정보를 가져왔습니다."
			gameInfo := gameState.GetAllGameInfo()
			response.MessageLogList = CreateMessagesFromGameStatus(gameInfo)
			response.MessageType = "info"
			response.User = adminUser
			broadCastToSomeone(clients, e.Conn, response)

		default:
			log.Printf("Unknown action: %s\n", e.Action)
		}
		ProcessGPTEntering(userManager, webSocketService, gameState)
		ProcessGPTReady(userManager, webSocketService, gameState)
		ProcessAllPlayersReady(userManager, webSocketService, gameState)
		ProcessNextTurn(userManager, webSocketService, gameState, e.Timestamp)
		ProcessAllPlayersVoted(userManager, webSocketService, gameState)
		// next game or next round
		HandleRoundIsOver(userManager, webSocketService, gameState, e)
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
	userManager.AddSortedPlayer(nowUser)
	webSocketService.AddClient(e.Conn, nowUser)

	message := utils.CreateMessageFromUser(userManager, adminUser, e.Timestamp)
	message.MessageType = "info"
	message.Message = fmt.Sprintf("%s님이 게임에 참여했습니다.", nowUser.UserName)
	userManager.AddMessage(message)
	response := utils.CreateResponseUsingPayload(userManager, gameState, e)
	response.Action = "human_info"
	response.MessageType = "system"
	response.Message = message.Message
	broadcastToAll(clients, response)

}

func BroadcastWebSocketMessage(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, e models.WsPayload) {
	clients := webSocketService.GetClients()
	message := utils.CreateMessageFromUser(userManager, e.User, e.Timestamp)
	userManager.AddMessage(message)
	response := utils.CreateResponseUsingPayload(userManager, gameState, e)
	broadcastToAll(clients, response)
}

func ProcessAllPlayersReady(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	if !gameState.CheckAllUserReady(userManager) {
		return
	}
	userManager.SetRandomlyShuffledPlayers(webSocketService, gameState)
	gameState.InitializeRoundInfo(userManager, webSocketService)

	clients := webSocketService.GetClients()
	adminUser := userManager.GetAdminUser()

	message := utils.CreateMessageWithAutoTimestamp(userManager, adminUser)
	message.MessageType = "alert"
	message.Message = "게임이 시작되었습니다."

	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.MessageType = "alert"
	response.Message = message.Message
	broadcastToAll(clients, response)

	userManager.AddPrevMessagesFromMessages()
	userManager.ClearMessages()
	userManager.AddMessage(message)

	// HandleRoundIsOver(e.Timestamp)
}

func ProcessGPTEntering(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	OnlineUserNum := len(utils.RetrieveUserList(userManager))

	GPTUsers := userManager.GetGPTUsers()
	GPTEntryNums := gameState.GetGPTEntryNums()
	for idx, GPTEntryNum := range GPTEntryNums {
		GPTUser := GPTUsers[idx]
		// If GPT's Entring time cames...(GPTUser isn't ONLINE)
		if !GPTUser.IsOnline && GPTEntryNum >= OnlineUserNum {
			time.Sleep(time.Second * 1)
			HandelGPTEntry(userManager, webSocketService, gameState, idx)
			time.Sleep(time.Second * 1)
		}
	}
}

func HandelGPTEntry(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, index int) {
	clients := webSocketService.GetClients()
	adminUser := userManager.GetAdminUser()
	GPTUser := userManager.GetGPTUsers()[index]

	// set GPTUser to true
	GPTUser.IsOnline = true
	userManager.SetGPTUser(index, GPTUser)
	userManager.AddPlayerByUser(GPTUser)

	// Make Message & Response
	message := utils.CreateMessageWithAutoTimestamp(userManager, adminUser)
	message.MessageType = "info"
	message.Message = fmt.Sprintf("%s님이 채팅방에 입장했습니다.", GPTUser.UserName)
	userManager.AddMessage(message)

	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.Action = "human_info"
	response.MessageType = "system"
	response.Message = message.Message

	broadcastToAll(clients, response)
}
func ProcessGPTReady(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	ReadyPlayerNum := len(utils.RetrievePlayerList(userManager))

	GPTUsers := userManager.GetGPTUsers()
	GPTReadyNums := gameState.GetGPTReadyNums()
	for idx, GPTReadyNum := range GPTReadyNums {
		GPTUser := GPTUsers[idx]
		// If GPT's Ready time cames...(GPTUser isn't WATCHER)
		if GPTUser.PlayerType == "watcher" && GPTReadyNum >= ReadyPlayerNum {
			time.Sleep(time.Second * 1)
			HandleGPTReady(userManager, webSocketService, gameState, idx)
			time.Sleep(time.Second * 1)
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
	userManager.AddSortedPlayer(GPTUser)
	userManager.SetGPTUser(index, GPTUser)

	// Make Message & Response
	message := utils.CreateMessageWithAutoTimestamp(userManager, adminUser)
	message.MessageType = "info"
	message.Message = fmt.Sprintf("%s님이 게임에 참여했습니다.", GPTUser.UserName)
	userManager.AddMessage(message)

	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.Action = "human_info"
	response.MessageType = "system"
	response.Message = message.Message

	broadcastToAll(clients, response)
}

func ProcessAllPlayersVoted(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	clients := webSocketService.GetClients()
	gameInfo := gameState.GetNowGameInfo()
	gameRound := gameInfo.Round
	userSelections := gameState.GetNowUserSelections()
	AdmminUser := userManager.GetAdminUser()
	GPTUsers := userManager.GetGPTUsers()

	if len(userSelections) == gameInfo.MaxPlayer-len(GPTUsers) {
		gameState.SetIfRoundIsOver()
		voteNum, eliminatedPlayer, remainingPlayerList := userManager.ExcludePlayersFromSelections(webSocketService, gameState)

		userManager.SetSortedPlayers(remainingPlayerList)

		messages := utils.CreateMessageWithAutoTimestamp(userManager, AdmminUser)
		messages.MessageType = "alert"
		messages.Message = fmt.Sprintf("%d라운드가 종료되었습니다. 탈락자는 %d표를 받은 [%s]입니다.", gameRound, voteNum, eliminatedPlayer.UserName)

		response := utils.CreateInitalizeResponse(userManager, gameState)
		response.Action = ""
		response.MessageType = "alert"
		response.Message = messages.Message

		broadcastToAll(clients, response)

		// Broadcast vote result to console
		broadcastSelectionResultToAll(userManager, webSocketService, gameState)

		time.Sleep(time.Second * 10)
	}
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
		message.MessageType = "alert"

		response := utils.CreateResponseUsingTimestamp(userManager, gameState, Timestamp)
		response.Action = "new_message_admin"
		response.MessageType = "alert"
		response.Message = message.Message
		response.User = VotedUser

		broadcastToAll(clients, response)
		userManager.AddMessage(message)
	}
}

func HandleRoundIsOver(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, e models.WsPayload) {
	isGameStarted := gameState.GetStatus().IsStarted
	gameTurnsLeft := gameState.GetNowGameInfo().TurnsLeft
	gameRound := gameState.GetNowGameInfo().Round
	roundNum := gameState.GetStatus().RoundNum

	// Game is not started
	if !isGameStarted || gameTurnsLeft != 0 {
		return
	}

	if roundNum == gameRound {
		HandleGameOverEvent(userManager, webSocketService, gameState)
	} else {
		HandleRoundOverEvent(userManager, webSocketService, gameState)
	}
}

func HandleRoundOverEvent(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	gameRound := gameState.GetNowGameInfo().Round
	gameState.SetIfResetRound(userManager)
	message := utils.CreateMessageWithAutoTimestamp(userManager, userManager.GetAdminUser())
	message.MessageType = "alert"
	message.Message = fmt.Sprintf("%d라운드가 종료되었습니다.", gameRound)
	userManager.AddMessage(message)

	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.Action = "restart_round"
	response.MessageType = "alert"
	response.Message = message.Message
	response.User = userManager.GetAdminUser()

	clients := webSocketService.GetClients()
	broadcastToAll(clients, response)

	userManager.AddPrevMessagesFromMessages()
	userManager.ClearMessages()
	userManager.AddMessage(message)
}
