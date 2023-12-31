package handlers

import (
	"fmt"
	"liar-of-turing/models"
	"liar-of-turing/services"
	"liar-of-turing/utils"
	"log"
)

func BroadcastConsoleSelections(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, timestamp int64) {
	log.Println("BroadcastConsoleSelections")

	clients := webSocketService.GetClients()
	adminUser := userManager.GetAdminUser()
	conn, exists := webSocketService.RetrieveClientByUUID(adminUser.UUID)
	if !exists {
		log.Println("admin is not connected")
		return
	}

	chatResponse := utils.CreateResponseUsingTimestamp(userManager, gameState, timestamp)
	chatResponse.Action = "send_result"
	chatResponse.Message = "send_result"
	chatResponse.UserSelection = gameState.GetNowUserSelections()

	broadcastToSomeone(clients, conn, chatResponse)

}

func HandleAISelection(gameState *services.GameState, e models.WsPayload) {
	selection := models.UserSelection{
		Timestamp: e.Timestamp,
		User:      e.User,
		Selection: e.UserSelection.Selection,
		Reason:    e.UserSelection.Reason,
	}
	log.Println("HandleAISelection:", selection)
	gameState.AddUserSelection(selection)
}
func ProcessNextTurn(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	// , timestamp int64) {
	if !gameState.CheckAllUserReady(userManager) {
		return
	}
	if gameState.CheckIsRoundOver() {
		gameState.SetIfRoundIsOver()
		// ProcessAllPlayersVoted(userManager, webSocketService, gameState)
		// broadcastChooseAIToAll(userManager, webSocketService, gameState)
		return
	}
	log.Println("ProcessNextTurn")

	adminUser := userManager.GetAdminUser()
	GPTUsers := userManager.GetGPTUsers()

	nextUser, exists := gameState.GetNextTurnPlayer()

	if !exists {
		log.Println("nextUser is not exists")
		return
	}
	log.Println("nextUser:", nextUser)

	// GPT send message
	for idx, GPTUser := range GPTUsers {
		if nextUser.UUID == GPTUser.UUID {
			ProcessGPTSendMessage(userManager, webSocketService, gameState, idx)
			return
		}
	}

	// message := utils.CreateMessageFromUser(userManager, adminUser, timestamp)
	message := utils.CreateMessageWithAutoTimestamp(userManager, adminUser)
	message.Message = fmt.Sprintf("%s님의 차례입니다.", nextUser.UserName)
	message.MessageType = "info"

	// response := utils.CreateResponseUsingTimestamp(userManager, gameState, timestamp)
	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.Action = "your_turn"
	response.MessageType = "info"
	response.Message = message.Message
	response.User = adminUser

	nextConn, exists := webSocketService.RetrieveClientByUUID(nextUser.UUID)
	if !exists {
		log.Println("nextUser is not connected")
		return
	}

	clients := webSocketService.GetClients()
	broadcastToSomeone(clients, nextConn, response)

}

func HandleGameOverEvent(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	// time.Sleep(5 * time.Second)
	clients := webSocketService.GetClients()
	adminUser := userManager.GetAdminUser()
	message := utils.CreateMessageWithAutoTimestamp(userManager, adminUser)
	message.Message = "게임이 종료되었습니다. 자동으로 재시작됩니다."
	message.MessageType = "info"

	userManager.AddMessage(message)
	// userManager.AddPrevMessagesFromMessages()
	// userManager.ClearMessages()

	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.Action = "update_state"
	response.Message = message.Message
	response.MessageType = message.MessageType
	response.User = message.User

	broadcastToAll(clients, response)
	// time.Sleep(10 * time.Second)

}
