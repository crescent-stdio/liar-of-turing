package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"liar-of-turing/common"
	"liar-of-turing/models"
	"liar-of-turing/services"
	"liar-of-turing/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func HandleGPTWebSocketRequest(w http.ResponseWriter, r *http.Request, webSocketService *services.WebSocketService) {
	ws, err := webSocketService.UpgradeConfig.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	conn := models.WebSocketConnection{Conn: ws}
	log.Println("Client connected to endpoint")

	go ListenToGPTWebSocketConnections(webSocketService, &conn)
}

func ListenToGPTWebSocketConnections(webSocketService *services.WebSocketService, conn *models.WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in ListenToGPTWebSocketConnections:", r)
		}
		webSocketService.RemoveClient(*conn) // Safely close the connection
	}()

	var payload models.GPTWsPayload

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
		webSocketService.GPTBroadcast <- payload
	}
}

func ForwardMessageToGPTAPI(GPTUser common.User, message string) models.MessageData {
	FastAPIURL := common.GetFastAPIURL()
	url := FastAPIURL + "/useGPT"
	msgData := models.MessageData{
		UserName: GPTUser.UserName,
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

	var gptResponse models.MessageData
	err = json.Unmarshal(body, &gptResponse)
	if err != nil {
		log.Fatalf("Error occurred during unmarshaling. Error: %s", err.Error())
	}
	log.Printf("gptResponse: %s", gptResponse)
	return gptResponse
}

func ProcessGPTSendMessage(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState, gptIndex int) {
	log.Println("=================================================")
	log.Println("SendGPTWsMessage")
	clients := webSocketService.GetClients()
	GPTUser := userManager.GetGPTUsers()[gptIndex]

	// var gptResponse WsJsonResponse
	messages := userManager.GetMessages()

	MessageLogString := ""

	// Add Players
	players := userManager.GetSortedPlayers()
	MessageLogString += "[단체 대화방 참여자]: "
	for _, v := range players {
		MessageLogString += fmt.Sprintf("%s, ", v.UserName)
	}
	MessageLogString += "\n"

	for _, v := range messages {
		if v.User.UserName == "server" && v.MessageType == "alert" {
			question := strings.Split(strings.Split(v.Message, "'")[1], "'")[0]
			MessageLogString += fmt.Sprintf("[주제]: %s\n", question)
		} else if v.User.UserName != "server" {
			MessageLogString += fmt.Sprintf("%s: %s\n", v.User.UserName, v.Message)
		}
	}
	MessageLogString += fmt.Sprintf("%s: ", GPTUser.UserName)
	log.Println("MessageLogString:", MessageLogString)
	gptResponse := ForwardMessageToGPTAPI(GPTUser, MessageLogString)
	log.Println("gptResponse:", gptResponse)

	gameState.SetNextTurnInfo()

	message := utils.CreateMessageWithAutoTimestamp(userManager, GPTUser)
	message.Message = gptResponse.Message

	userManager.AddMessage(message)

	response := utils.CreateInitalizeResponse(userManager, gameState)
	response.Action = "new_message"
	response.Message = message.Message
	response.User = GPTUser
	broadcastToAll(clients, response)

	ProcessNextTurn(userManager, webSocketService, gameState) //, message.Timestamp)
}

func CreateMessagesFromGameStatus(gameInfo []models.Game) []models.Message {
	messages := make([]models.Message, 0)
	for _, v := range gameInfo {
		messages = append(messages, v.Messages...)
	}
	return messages

}
