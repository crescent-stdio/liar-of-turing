package handlers

import (
	"fmt"
	"liar-of-turing/common"
	"liar-of-turing/models"
	"liar-of-turing/services"
	"liar-of-turing/utils"
	"log"
)

func HandleHumanUserEntry(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, e models.WsPayload) {
	log.Println("HandleHumanUserEntry")
	nowUser, exists := userManager.GetPlayerByUUID(e.User.UUID)
	clients := webSocketService.GetClients()

	isGameStarted := gameState.GetStatus().IsStarted
	isGameOver := gameState.GetStatus().IsOver
	log.Println("isGameStarted:", isGameStarted, "isGameOver:", isGameOver)

	// If user is not in players map, create new user
	if !exists {
		nicknameId, userName := userManager.GenerateRandomUsername()
		nowUser = common.User{
			UserId:     int64(len(clients)),
			UserName:   userName,
			NicknameId: nicknameId,
			Role:       "human",
			UUID:       e.User.UUID,
			IsOnline:   true,
			PlayerType: "watcher",
		}
	}
	nowUser.IsOnline = true
	// Add Active Players and Clients
	userManager.AddPlayerByUser(nowUser)
	webSocketService.AddClient(e.Conn, nowUser)

	response := utils.CreateResponseUsingPayload(userManager, gameState, e)
	response.Action = "update_state"
	broadcastToAll(clients, response)

	// if isGameStarted && !isGameOver {
	// 	adminUser := userManager.GetAdminUser()
	// 	nextUser, nextExists := gameState.GetNextTurnPlayer()

	// 	// Send message to next user
	// 	if nextExists && nextUser.UUID == nowUser.UUID {
	// 		// nowUser.PlayerType = "player"
	// 		// someoneMessage := utils.CreateMessageFromUser(userManager, adminUser, e.Timestamp)
	// 		// someoneMessage.Message = fmt.Sprintf("%s님의 차례입니다.", nextUser.UserName)
	// 		// someoneMessage.MessageType = "info"

	// 		// someoneResponse := utils.CreateResponseUsingTimestamp(userManager, gameState, e.Timestamp)
	// 		// someoneResponse.Action = "your_turn"
	// 		// someoneResponse.MessageType = "info"
	// 		// someoneResponse.Message = someoneMessage.Message
	// 		// someoneResponse.User = adminUser

	// 		// nextConn, _ := webSocketService.RetrieveClientByUUID(nextUser.UUID)
	// 		// broadcastToSomeone(clients, nextConn, someoneResponse)
	// 	} else {
	// 		// Send message to now user
	// 		response := utils.CreateResponseUsingPayload(userManager, gameState, e)
	// 		response.Action = "update_state"
	// 		response.User = adminUser

	// 		broadcastToAll(clients, response)

	// 	}
	// } else {
	// broadcast to all
	adminUser := userManager.GetAdminUser()
	message := utils.CreateMessageFromUser(userManager, adminUser, e.Timestamp)
	message.Message = fmt.Sprintf("%s님이 입장했습니다.", nowUser.UserName)
	message.MessageType = "system"

	response = utils.CreateResponseUsingPayload(userManager, gameState, e)
	response.Action = "human_info"
	response.User = adminUser
	response.MessageType = "system"
	response.Message = message.Message
	log.Println("response:", response)

	clients = webSocketService.GetClients()
	broadcastToAll(clients, response)
	// }

	// Round is over And enter the web application(previous User was player)
	if gameState.CheckIsRoundOver() && exists && nowUser.PlayerType == "player" {
		var waitResponse models.WsJsonResponse
		_, exists := gameState.SearchUserInUserSelections(nowUser)
		if exists {
			waitResponse = utils.CreateResponseUsingTimestamp(userManager, gameState, e.Timestamp)
			waitResponse.Action = "wait_for_players"
			waitResponse.MessageType = "info"
			waitResponse.Message = "라운드가 종료되었습니다. 다음 플레이어의 선택을 기다리세요."
		} else {
			waitResponse = utils.CreateResponseUsingTimestamp(userManager, gameState, e.Timestamp)
			waitResponse.Action = "submit_ai"
			waitResponse.MessageType = "info"
			waitResponse.Message = "라운드가 종료되었습니다. 다음 라운드를 위해 AI를 선택하세요."
		}
		broadcastToSomeone(clients, e.Conn, waitResponse)

		return
	}

}

func HandleRestartGameEvent(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, e models.WsPayload) {
	log.Println("HandleRestartGameEvent")
	clients := webSocketService.GetClients()
	adminUser := userManager.GetAdminUser()
	GPTUsers := userManager.GetGPTUsers()

	// Reset Game Status

	// Set all users as watchers
	gameState.SetIfGameTotallyReset(userManager)
	userManager.SetAllUsersAsWatchers()
	for _, gpt := range GPTUsers {
		gpt.PlayerType = "watcher"
	}
	userManager.SetGPTUsers(GPTUsers)

	// Send message about game reset
	message := utils.CreateMessageFromUser(userManager, adminUser, e.Timestamp)
	message.Message = "게임이 초기화되었습니다."
	message.MessageType = "info"

	// Reset Messages
	userManager.AddPrevMessagesFromMessages()
	userManager.ClearMessages()
	userManager.AddMessage(message)

	response := utils.CreateResponseUsingPayload(userManager, gameState, e)
	response.Action = "restart_game"
	response.Message = message.Message
	response.MessageType = message.MessageType
	response.User = message.User
	broadcastToAll(clients, response)

}
