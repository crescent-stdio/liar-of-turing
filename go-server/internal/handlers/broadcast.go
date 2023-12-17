package handlers

import (
	"liar-of-turing/common"
	"liar-of-turing/models"
	"liar-of-turing/services"
	"liar-of-turing/utils"
	"log"
)

// WebScok
func broadcastToAll(clients map[models.WebSocketConnection]common.User, response models.WsJsonResponse) {
	common.GlobalMutex.Lock()
	defer common.GlobalMutex.Unlock()
	for client := range clients {
		if err := client.WriteJSON(response); err != nil {
			log.Println("[broadcastToAll] Websocket error:", err)
			client.CloseWebSocket()
			delete(clients, client)
		}
	}
	log.Println("Broadcasted message")
}

func broadcastToSomeone(clients map[models.WebSocketConnection]common.User, client models.WebSocketConnection, response models.WsJsonResponse) {
	common.GlobalMutex.Lock()
	defer common.GlobalMutex.Unlock()

	if err := client.WriteJSON(response); err != nil {
		log.Println("[braodCastToSomeone] Websocket error:", err)
		client.CloseWebSocket()
		delete(clients, client)
	}
	log.Println("Broadcasted message")
}

func broadcastChooseAIToAll(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) {
	isUsersVoting := gameState.GetStatus().IsUsersVoting
	if isUsersVoting {
		return
	}
	gameState.SetIsUsersVotingTrue()

	adminUser := userManager.GetAdminUser()
	message := utils.CreateMessageWithAutoTimestamp(userManager, adminUser)
	message.Message = "AI를 선택해주세요."
	message.MessageType = "info"

	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.Action = "choose_ai"
	response.Message = message.Message
	response.MessageType = message.MessageType
	response.User = message.User

	clients := webSocketService.GetClients()
	// broadcastToAll(clients, response)
	players := gameState.GetNowGameInfo().PlayerList
	for _, player := range players {
		conn, exists := webSocketService.RetrieveClientByUUID(player.UUID)
		_, isVoted := gameState.SearchUserInUserSelections(player)
		if player.Role == "human" && exists && !isVoted {
			broadcastToSomeone(clients, conn, response)
		}
	}
}

// func HandleUserJoin(userManager *services.UserManager, e models.WsPayload, clients map[models.WebSocketConnection]common.User) {
// 	message := models.Message{
// 		Timestamp:   e.Timestamp,
// 		MessageId:   int64(len(messages)),
// 		User:        e.User,
// 		Message:     e.Message,
// 		MessageType: "message",
// 	}
// 	messages = append(messages, message)

// 	response := models.WsJsonResponse{
// 		Timestamp:      e.Timestamp,
// 		Action:         e.Action,
// 		User:           e.User,
// 		Message:        e.Message,
// 		MessageLogList: messages,
// 		MessageType:    "message",
// 		MaxPlayer:      MaxPlayer,

// 		OnlineUserList: utils.RetrieveUserList(userManager),
// 		PlayerList:     utils.RetrieveReadyUserList(userManager),
// 		GameTurnsLeft:  gameTurnsLeft,
// 		GameRound:      gameRound,
// 		IsGameStarted:  isGameStarted,
// 		IsGameOver:     isGameOver,
// 	}
// 	broadcastToAll(clients, response)
// }
