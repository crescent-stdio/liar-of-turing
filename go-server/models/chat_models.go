package models

import "liar-of-turing/common"

type Message struct {
	Timestamp   int64       `json:"timestamp"`
	MessageId   int64       `json:"message_id"`
	User        common.User `json:"user"`
	Message     string      `json:"message"`
	MessageType string      `json:"message_type"`
}
