package handlers

import (
	"liar-of-turing/models"
	"liar-of-turing/services"
	"liar-of-turing/utils"
	"log"
)

func HandleUserLeave(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, e models.WsPayload) {
	clients := webSocketService.GetClients()
	if leftUser, ok := clients[e.Conn]; ok {
		leftUser.IsOnline = false
		userManager.AddPlayerByUser(leftUser)
		webSocketService.RemoveClient(e.Conn)

		response := utils.CreateResponseUsingPayload(userManager, gameState, e)
		response.Action = "update_state"
		response.MessageType = "system"
		response.Message = leftUser.UserName + "님이 퇴장했습니다."
		broadcastToAll(clients, response)

	} else {
		log.Println("User not found in clients map")
		response := utils.CreateResponseUsingPayload(userManager, gameState, e)
		response.Action = "update_state"
		response.MessageType = "system"
		response.Message = "User not found in clients map"
		broadcastToAll(clients, response)
	}
}
