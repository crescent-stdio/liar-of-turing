package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Broadcast = make(chan WsPayload)

var clients = make(map[WebSocketConnection]string)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

// upgradeConnection is the websocket upgrader from gorilla/websockets
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebSocketConnection struct {
	*websocket.Conn
}

// WsJsonResponse defines the response sent back from websocket
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
}

// WsPayload defines the websocket request from the client
type WsPayload struct {
	Action    string              `json:"action"`
	RoomId    string              `json:"roomId"`
	Username  string              `json:"username"`
	UserId    string              `json:"userId"`
	Role      string              `json:"role"`
	Timestamp int64               `json:"timestamp"`
	Message   string              `json:"message"`
	Conn      WebSocketConnection `json:"-"`
}

// WsEndpoint upgrades connection to websocket
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	var response WsJsonResponse
	response.Action = "Connected"
	response.Message = "Connected to server"
	response.MessageType = "info"

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = "Anonymous"

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}
	go ListenForWs(&conn)
}

func ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-Broadcast
		// switch v := e.(role) {
		// case Message:
		// 	response.Action = "message"
		// 	response.Message = v.Message
		// 	response.MessageType = v.Type
		// case User:
		// 	response.Action = "user"
		// 	response.Message = v.UserName
		// 	response.MessageType = v.Role
		// case Room:
		// 	response.Action = "room"
		// 	response.Message = v.RoomID
		// 	response.MessageType = v.Users
		// case RoomList:
		// 	response.Action = "roomList"
		// 	response.Message = v.Rooms
		// 	response.MessageType = "roomList"
		// case RoomInfo:
		// 	response.Action = "roomInfo"
		// 	response.Message = v.RoomID
		// 	response.MessageType = v.Users
		// }

		response.Action = "Got your message"
		response.Message = fmt.Sprintf("Someone sent a message and Action is %s", e.Action)
		response.MessageType = "info"
		broadcastToAll(response)
	}
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Websocket error", err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", r)
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			Broadcast <- payload
		}
	}
}
