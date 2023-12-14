package models

// GPTWsJsonResponse represents the JSON response for a GPT-WebSocket message.
type GPTWsJsonResponse struct {
	UserUUID         string `json:"user_uuid"`
	MessageLogString string `json:"message_log_string"`
	MessageType      string `json:"message_type"`
}

// GPTWsPayload represents the payload of a GPT-WebSocket message.
type GPTWsPayload struct {
	UserUUID string              `json:"user_uuid"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

type MessageData struct {
	UserUUID string `json:"user_UUID"`
	Message  string `json:"message"`
}
