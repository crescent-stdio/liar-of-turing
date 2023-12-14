package models

type User struct {
	UUID   string `json:"uuid"`
	UserId int64  `json:"user_id"`
	// RoomId     int64  `json:"room_id"`
	NicknameId int    `json:"nickname_id"`
	UserName   string `json:"username"`
	Role       string `json:"role"`
	IsOnline   bool   `json:"is_online"`
	PlayerType string `json:"player_type"`
}

type Message struct {
	Timestamp   int64  `json:"timestamp"`
	MessageId   int64  `json:"message_id"`
	User        User   `json:"user"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}
