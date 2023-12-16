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

	// broadcast to all

	response := utils.CreateResponseUsingPayload(userManager, gameState, e)
	response.Action = "human_info"
	response.MessageType = "system"
	response.Message = fmt.Sprintf("%s님이 입장했습니다.", nowUser.UserName)
	broadcastToAll(clients, response)

	// Round is over And enter the web application(previous User was player)
	if gameState.CheckIsRoundOver() && exists && nowUser.PlayerType == "player" {
		var waitResponse models.WsJsonResponse
		gameRound := gameState.GetNowGameInfo().Round
		_, exists := gameState.SearchUserInUserSelections(gameRound-1, nowUser)
		if exists {
			waitResponse = utils.CreateResponseUsingTimestamp(userManager, gameState, e.Timestamp)
			waitResponse.Action = "wait_for_players"
			waitResponse.MessageType = "alert"
			waitResponse.Message = "라운드가 종료되었습니다. 다음 플레이어의 선택을 기다리세요."
		} else {
			waitResponse = utils.CreateResponseUsingTimestamp(userManager, gameState, e.Timestamp)
			waitResponse.Action = "submit_ai"
			waitResponse.MessageType = "alert"
			waitResponse.Message = "라운드가 종료되었습니다. 다음 라운드를 위해 AI를 선택하세요."
		}
		broadCastToSomeone(clients, e.Conn, waitResponse)

		return
	}

	// If game is started, send message to next user
	if isGameStarted && !isGameOver {
		adminUser := userManager.GetAdminUser()
		nextUser, exists := gameState.GetNextTurnPlayer()

		// Send message to next user
		if exists && nextUser.UUID == e.User.UUID {
			nowUser.PlayerType = "player"
			someoneMessage := utils.CreateMessageFromUser(userManager, adminUser, e.Timestamp)
			someoneMessage.Message = fmt.Sprintf("%s님의 차례입니다.", nextUser.UserName)
			someoneMessage.MessageType = "alert"

			someoneResponse := utils.CreateResponseUsingTimestamp(userManager, gameState, e.Timestamp)
			someoneResponse.Action = "your_turn"
			someoneResponse.MessageType = "alert"
			someoneResponse.Message = someoneMessage.Message
			someoneResponse.User = adminUser

			nextConn, _ := webSocketService.RetrieveClientByUUID(nextUser.UUID)
			broadCastToSomeone(clients, nextConn, someoneResponse)
		}
	}

}

func HandleRestartGameEvent(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, e models.WsPayload) {
	log.Println("HandleRestartGameEvent")
	clients := webSocketService.GetClients()
	adminUser := userManager.GetAdminUser()
	GPTUsers := userManager.GetGPTUsers()

	// Reset Game Status
	gameState.SetIfGameTotallyReset(userManager)

	// Set all users as watchers
	userManager.SetAllUsersAsWatchers()
	for _, gpt := range GPTUsers {
		gpt.PlayerType = "watcher"
	}
	userManager.SetGPTUsers(GPTUsers)

	// Send message about game reset
	message := utils.CreateMessageFromUser(userManager, adminUser, e.Timestamp)
	message.Message = "게임이 초기화되었습니다."
	message.MessageType = "alert"

	response := utils.CreateResponseUsingPayload(userManager, gameState, e)
	response.Action = "restart_game"
	response.Message = message.Message
	response.MessageType = message.MessageType
	response.User = message.User
	broadcastToAll(clients, response)

	// Reset Messages
	userManager.AddPrevMessagesFromMessages()
	userManager.ClearMessages()
	userManager.AddMessage(message)
}
