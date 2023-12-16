package handlers

import (
	"liar-of-turing/common"
	"liar-of-turing/models"
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

func broadCastToSomeone(clients map[models.WebSocketConnection]common.User, client models.WebSocketConnection, response models.WsJsonResponse) {
	common.GlobalMutex.Lock()
	defer common.GlobalMutex.Unlock()

	if err := client.WriteJSON(response); err != nil {
		log.Println("[braodCastToSomeone] Websocket error:", err)
		client.CloseWebSocket()
		delete(clients, client)
	}
	log.Println("Broadcasted message")
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
// 		PlayerList:     utils.RetrievePlayerList(userManager),
// 		GameTurnsLeft:  gameTurnsLeft,
// 		GameRound:      gameRound,
// 		IsGameStarted:  isGameStarted,
// 		IsGameOver:     isGameOver,
// 	}
// 	broadcastToAll(clients, response)
// }
