package handlers

import (
	"fmt"
	"liarOfTuring/global"
	"liarOfTuring/models"
	"liarOfTuring/utils"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var mutex = &sync.Mutex{}
var nicknames []models.Nickname

func init() {
	// 애플리케이션 시작 시 전역 변수 초기화
	if err := utils.LoadNicknames(); err != nil {
		fmt.Println("Error loading nicknames:", err)
	}
	nicknames = global.GetGlobalNicknames()
}

var Broadcast = make(chan WsPayload)

var clients = make(map[WebSocketConnection]User)

var players = make(map[string]User)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

// upgradeConnection is the websocket upgrader from gorilla/websockets
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
	HandshakeTimeout: 1024,
}

type WebSocketConnection struct {
	*websocket.Conn
}

// WsJsonResponse defines the response sent back from websocket
type WsJsonResponse struct {
	Timestamp      int64  `json:"timestamp"`
	Action         string `json:"action"`
	User           User   `json:"user"`
	Message        string `json:"message"`
	MessageType    string `json:"message_type"`
	ConnectedUsers []User `json:"connected_users"`
}

type User struct {
	UUID       string `json:"uuid"`
	UserId     int64  `json:"user_id"`
	NicknameId int    `json:"nickname_id"`
	UserName   string `json:"username"`
	Role       string `json:"role"`
	isOnline   bool   `json:"is_online"`
}

func sortUserList(users []User) []User {
	sort.Slice(users, func(i, j int) bool {
		return users[i].UserName < users[j].UserName
	})
	return users
}

// WsPayload defines the websocket request from the client
type WsPayload struct {
	Action    string              `json:"action"`
	RoomId    int64               `json:"room_id"`
	User      User                `json:"user"`
	Timestamp int64               `json:"timestamp"`
	Message   string              `json:"message"`
	Conn      WebSocketConnection `json:"-"` // ignore this field
}

// WsEndpoint upgrades connection to websocket
func WsEndpoint(w http.ResponseWriter, r *http.Request) {

	// nicknames, err := services.LoadNicknames()
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	// var response WsJsonResponse
	// response.Action = "Connected"
	// response.Message = "Connected to server"
	// response.MessageType = "server"

	conn := WebSocketConnection{Conn: ws}
	// clients[conn] = User{UserName: "server", Role: "server", UserId: 0}
	// log.Println(clients)

	// err = ws.WriteJSON(response)
	// if err != nil {
	// 	log.Println(err)
	// }
	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", r)
		}
		conn.Close()           // Ensure the connection is closed
		delete(clients, *conn) // Remove the client from the ma
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		log.Println("ListenForWs")
		log.Println(payload)
		log.Println("err:", err)
		if err != nil {
			log.Println("Error reading json:", err)
			break // Exit loop on error
		} else {
			payload.Conn = *conn
			Broadcast <- payload
		}
	}

}

func ListenToWsChannel() {
	var response WsJsonResponse
	for {
		log.Println("Got here!!")
		e := <-Broadcast
		log.Println("Got here")
		log.Println("e.User:", e.User)
		mutex.Lock()
		switch e.Action {
		case "new_message":
			response.Timestamp = e.Timestamp
			response.Action = "message"
			response.User = e.User
			response.Message = e.Message
			log.Println(response.Message)
			response.MessageType = "message"
			broadcastToAll(response)
		case "list_users":
			response.Timestamp = e.Timestamp
			response.Action = "user_list"
			response.MessageType = "info"
			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "enter_human":
			// clients[e.Conn] = e.User
			nowUser := findUser(e.User.UUID)
			if nowUser == (User{}) {
				nicknameId, userName := getRandomUsername()
				nowUser = User{
					UserId:     int64(len(clients)),
					UserName:   userName,
					NicknameId: nicknameId,
					Role:       "human",
					UUID:       e.User.UUID,
					isOnline:   true,
				}
				players[e.User.UUID] = nowUser
			}
			log.Println("enter_human")
			clients[e.Conn] = nowUser
			log.Println(clients)
			log.Println("e.User:", e.User)
			log.Println("nowUser:", nowUser)
			// e.User = nowUser
			response.Action = "human_info"
			response.Message = fmt.Sprintf("%s님이 입장했습니다.", nowUser.UserName)
			response.MessageType = "info"
			response.User = nowUser
			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(response)

		case "left_user":
			log.Println("left_user")
			// response.Action = "update_state"

			leftUser := clients[e.Conn]
			leftUser.isOnline = false
			players[e.User.UUID] = leftUser

			delete(clients, e.Conn)
			// response.Message = fmt.Sprintf("%s님이 퇴장했습니다.", leftUser.UserName)
			response.Action = "user_list"
			users := getUserList()
			log.Println(users)
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("%s: %s", e.User.UserName, e.Message)
			response.MessageType = "message"
			broadcastToAll(response)

		}
		// clients[e.Conn] = e.Username

		// response.Action = "Got your message"
		// response.Message = fmt.Sprintf("Someone sent a message and Action is %s", e.Action)
		// response.MessageType = "info"
		// broadcastToAll(response)
		mutex.Unlock()
	}

}

func getUserList() []User {
	var userList []User
	for _, v := range clients {
		if v.UserName != "server" && v.isOnline {
			userList = append(userList, v)
		}
	}
	sortUserList(userList)
	return userList
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

func getRandomUsername() (int, string) {
	//random suffle nicknames and return not used nickname
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(nicknames))
	for idx, v := range perm {
		if !nicknames[v].IsUsed {
			nicknames[v].IsUsed = true
			return idx, nicknames[v].Nickname
		}
	}
	return -1, ""
}

func findUser(uuid string) User {
	for _, v := range players {
		log.Println("v.UUID:", v.UUID)
		if v.UUID == uuid {
			log.Println("dddddddddd")
			log.Println("v:", v)
			v.isOnline = true
			return v
		}
	}
	return User{}
}
