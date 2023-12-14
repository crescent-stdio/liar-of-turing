package handlers

import (
	. "liar-of-turing/models"
	"log"
)

// WebScok
func broadcastToAll(clients map[WebSocketConnection]User, response WsJsonResponse) {

	for client := range clients {
		if err := client.WriteJSON(response); err != nil {
			log.Println("[broadcastToAll] Websocket error:", err)
			client.CloseWebSocketConnection()
			delete(clients, client)
		}
	}
	log.Println("Broadcasted message")
}

func broadCastToSomeone(clients map[WebSocketConnection]User, client WebSocketConnection, response WsJsonResponse) {

	if err := client.WriteJSON(response); err != nil {
		log.Println("[braodCastToSomeone] Websocket error:", err)
		client.CloseWebSocketConnection()
		delete(clients, client)
	}
	log.Println("Broadcasted message")
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
