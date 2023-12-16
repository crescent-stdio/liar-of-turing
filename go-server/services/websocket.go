package services

import (
	"liar-of-turing/common"
	"liar-of-turing/models"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketService is a service that handles websocket connections
type WebSocketService struct {
	mutex         *sync.Mutex
	Clients       map[models.WebSocketConnection]common.User
	Broadcast     chan models.WsPayload
	GPTBroadcast  chan models.GPTWsPayload
	UpgradeConfig websocket.Upgrader
}

// NewWebSocketService creates a new WebSocketService
func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		mutex:   &sync.Mutex{},
		Clients: make(map[models.WebSocketConnection]common.User),

		Broadcast:    make(chan models.WsPayload),
		GPTBroadcast: make(chan models.GPTWsPayload),
		UpgradeConfig: websocket.Upgrader{
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
			CheckOrigin:      func(r *http.Request) bool { return true },
			HandshakeTimeout: 1024,
		},
	}
}

// GetUpgradeConfig
func (ws *WebSocketService) GetUpgradeConfig() websocket.Upgrader {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	return ws.UpgradeConfig
}

//

// AddClient: add client to Clients
func (ws *WebSocketService) AddClient(conn models.WebSocketConnection, user common.User) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	ws.Clients[conn] = user
}

// RemoveClient: remove client from Clients
func (ws *WebSocketService) RemoveClient(conn models.WebSocketConnection) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	delete(ws.Clients, conn)
	conn.CloseWebSocket()
}

// GetClientByConn: get client from Clients
func (ws *WebSocketService) GetClientByConn(conn models.WebSocketConnection) (common.User, bool) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	user, exists := ws.Clients[conn]
	return user, exists
}

// SetClientByUUID: set client by user.uuid
func (ws *WebSocketService) SetClientByUUID(uuid string, user common.User) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	for conn, client := range ws.Clients {
		if client.UUID == uuid {
			ws.Clients[conn] = user
			break
		}
	}
}

// RetrieveClientByUUID: retrieve client by user.uuid
func (ws *WebSocketService) RetrieveClientByUUID(uuid string) (models.WebSocketConnection, bool) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	for conn, client := range ws.Clients {
		if client.UUID == uuid {
			return conn, true
		}
	}
	return models.WebSocketConnection{}, false
}

// GetClients gets clients from players
func (ws *WebSocketService) GetClients() map[models.WebSocketConnection]common.User {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	clients := make(map[models.WebSocketConnection]common.User)
	for k, v := range ws.Clients {
		clients[k] = v
	}
	return clients
}

// SetClientsByUserUUID sets clients by user.uuid
func (ws *WebSocketService) SetClientsByUserUUID(user common.User) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	for conn, client := range ws.Clients {
		if client.UUID == user.UUID {
			ws.Clients[conn] = user
			break
		}
	}
}
func (ws *WebSocketService) GetTotalUserNum() int {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	return len(ws.GetClients())
}

// GetOnlineUserList: Get online user list
func (ws *WebSocketService) GetOnlineUserList() []common.User {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	clients := ws.GetClients()
	users := make([]common.User, 0)
	for _, v := range clients {
		if v.IsOnline {
			users = append(users, v)
		}
	}
	return users
}
