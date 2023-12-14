package models

import "github.com/gorilla/websocket"

// WebSocketConnection represents a WebSocket connection.
type WebSocketConnection struct {
	*websocket.Conn
}

// CloseWebSocketConnection closes the WebSocket connection.
func (conn *WebSocketConnection) CloseWebSocketConnection() {
	if conn != nil && conn.Conn != nil {
		conn.Close()
	}
}

// WsPayload represents the payload of a WebSocket message.
type WsPayload struct {
	Action string `json:"action"`
	// RoomId    int64               `json:"room_id"`
	MaxPlayer     int                 `json:"max_player"`
	User          User                `json:"user"`
	Timestamp     int64               `json:"timestamp"`
	Message       string              `json:"message"`
	GameTurnsLeft int                 `json:"game_turns_left"`
	GameRound     int                 `json:"game_round"`
	GameTurnNum   int                 `json:"game_turn_num"`
	GameRoundNum  int                 `json:"game_round_num"`
	UserSelection UserSelection       `json:"user_selection"`
	Conn          WebSocketConnection `json:"-"` // ignore this field
}

// WsJsonResponse represents the JSON response for a WebSocket message.
type WsJsonResponse struct {
	Timestamp      int64           `json:"timestamp"`
	MaxPlayer      int             `json:"max_player"`
	Action         string          `json:"action"`
	User           User            `json:"user"`
	Message        string          `json:"message"`
	MessageType    string          `json:"message_type"`
	MessageLogList []Message       `json:"message_log_list"`
	OnlineUserList []User          `json:"online_user_list"`
	PlayerList     []User          `json:"player_list"`
	GameTurnsLeft  int             `json:"game_turns_left"`
	GameRound      int             `json:"game_round"`
	GameTurnNum    int             `json:"game_turn_num"`
	GameRoundNum   int             `json:"game_round_num"`
	IsGameStarted  bool            `json:"is_game_started"`
	IsGameOver     bool            `json:"is_game_over"`
	UserSelection  []UserSelection `json:"user_selections_list"`
}
