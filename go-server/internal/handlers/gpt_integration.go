package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"liar-of-turing/utils"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func ListenGPTWebSocketConnections(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in ListenForGPTWs:", r)
		}
		conn.CloseWebSocketConnection() // Safely close the connection
		mutex.Lock()
		delete(clients, *conn) // Remove the client from the map
		mutex.Unlock()
	}()

	var payload GPTWsPayload

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
		GPTBroadcast <- payload
	}
}

func ListenToGPTWebSocketChannel(fastAPIURL string) {
	FastAPIURL = fastAPIURL
	var response WsJsonResponse
	for {
		e := <-GPTBroadcast
		log.Println("GPTBroadcast:", e)
		mutex.Lock()
		message := Message{
			Timestamp:   utils.GetCurrentTimestamp(),
			MessageId:   int64(len(messages)),
			User:        GPTUser,
			Message:     e.Message,
			MessageType: "message",
		}
		messages = append(messages, message)
		response = WsJsonResponse{
			Timestamp:      utils.GetCurrentTimestamp(),
			Action:         "new_message",
			User:           GPTUser,
			Message:        e.Message,
			MessageLogList: messages,
			MessageType:    "message",
			MaxPlayer:      MaxPlayer,
		}
		broadcastToAll(clients, response)
	}
}

func HandleGPTWebSocketRequest(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	conn := WebSocketConnection{Conn: ws}
	log.Println("Client connected to endpoint")

	go ListenGPTWebSocketConnections(&conn)
}

func sendGPTMessageToFastAPI(message string) MessageData {
	url := FastAPIURL + "/useGPT"
	msgData := MessageData{
		UserUUID: GPTUser.UUID,
		Message:  message,
	}
	// msgData := map[string]string{"message": message}
	msgBytes, err := json.Marshal(msgData)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		log.Fatalf("Error occurred during sending request to FastAPI. Error: %s", err.Error())
	}
	defer resp.Body.Close()

	log.Printf("Message sent to FastAPI, received response status: %s", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error occurred during reading response body. Error: %s", err.Error())
	}
	log.Printf("Response body: %s", string(body))

	var gptResponse MessageData
	err = json.Unmarshal(body, &gptResponse)
	if err != nil {
		log.Fatalf("Error occurred during unmarshaling. Error: %s", err.Error())
	}
	log.Printf("gptResponse: %s", gptResponse)
	return gptResponse
}

func SendGPTMessage(timestamp int64) {
	log.Println("=================================================")
	log.Println("SendGPTMessage")

	// var gptResponse WsJsonResponse
	MessageLogString := ""
	for _, v := range messages {
		if v.User.UserName != "admin" {
			MessageLogString += fmt.Sprintf("%s: %s\n", v.User.UserName, v.Message)
		}
	}
	log.Println("MessageLogString:", MessageLogString)
	gptResponse := sendGPTMessageToFastAPI(MessageLogString)
	log.Println("gptResponse:", gptResponse)

	gameInfo[gameRound-1].NowUserIndex = (gameInfo[gameRound-1].NowUserIndex + 1) % gameInfo[gameRound-1].MaxPlayer
	gameTurnsLeft = utils.Max(gameTurnsLeft-1, 0)
	gameInfo[gameRound-1].TurnsLeft = gameTurnsLeft

	mutex.Lock()

	message := Message{
		Timestamp:   timestamp,
		MessageId:   int64(len(messages)),
		User:        GPTUser,
		Message:     gptResponse.Message,
		MessageType: "message",
	}
	messages = append(messages, message)
	chatResponse := WsJsonResponse{
		Timestamp:      timestamp,
		Action:         "new_message",
		MessageType:    "message",
		Message:        message.Message,
		MaxPlayer:      MaxPlayer,
		User:           GPTUser,
		MessageLogList: messages,
		OnlineUserList: getUserList(),
		PlayerList:     getPlayerList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      gameRound,
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	broadcastToAll(clients, chatResponse)
	mutex.Unlock()

	if gameTurnsLeft == 0 {
		processGameOver(timestamp)
		return
	}
	processNextTurn(utils.GetCurrentTimestamp())
}
