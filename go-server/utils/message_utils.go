package utils

import (
	"liar-of-turing/common"
	"liar-of-turing/models"
	"liar-of-turing/services"
)

func CreateResponseUsingPayload(userManager *services.UserManager, gameState *services.GameState, e models.WsPayload) models.WsJsonResponse {
	return models.WsJsonResponse{
		Timestamp:      e.Timestamp,
		Action:         e.Action,
		MessageType:    "message",
		Message:        e.Message,
		MaxPlayer:      userManager.GetMaxPlayer(),
		User:           e.User,
		MessageLogList: userManager.GetMessages(),
		OnlineUserList: RetrieveUserList(userManager),
		PlayerList:     RetrieveReadyUserList(userManager),
		GameTurnsLeft:  gameState.GetNowGameInfo().TurnsLeft,
		GameRound:      gameState.GetNowGameInfo().Round,
		IsGameStarted:  gameState.GetStatus().IsStarted,
		IsGameOver:     gameState.GetStatus().IsOver,
	}
}

func CreateResponseUsingTimestamp(userManager *services.UserManager, gameState *services.GameState, timestamp int64) models.WsJsonResponse {
	return models.WsJsonResponse{
		Timestamp:      timestamp,
		Action:         "update_state",
		MessageType:    "message",
		Message:        "",
		MaxPlayer:      userManager.GetMaxPlayer(),
		User:           common.User{},
		MessageLogList: userManager.GetMessages(),
		OnlineUserList: RetrieveUserList(userManager),
		PlayerList:     RetrieveReadyUserList(userManager),
		GameTurnsLeft:  gameState.GetNowGameInfo().TurnsLeft,
		GameRound:      gameState.GetNowGameInfo().Round,
		IsGameStarted:  gameState.GetStatus().IsStarted,
		IsGameOver:     gameState.GetStatus().IsOver,
	}
}

func CreateInitalizeResponse(userManager *services.UserManager, gameState *services.GameState) models.WsJsonResponse {
	return models.WsJsonResponse{
		Timestamp:      GetCurrentTimestamp(),
		Action:         "update_state",
		MessageType:    "message",
		Message:        "",
		MaxPlayer:      userManager.GetMaxPlayer(),
		User:           common.User{},
		MessageLogList: userManager.GetMessages(),
		OnlineUserList: RetrieveUserList(userManager),
		PlayerList:     RetrieveReadyUserList(userManager),
		GameTurnsLeft:  gameState.GetNowGameInfo().TurnsLeft,
		GameRound:      gameState.GetNowGameInfo().Round,
		IsGameStarted:  gameState.GetStatus().IsStarted,
		IsGameOver:     gameState.GetStatus().IsOver,
	}
}
func CreateMessageFromUser(userManager *services.UserManager, user common.User, timestamp int64) models.Message {
	return models.Message{
		Timestamp:   timestamp,
		MessageId:   int64(len(userManager.GetMessages())),
		User:        user,
		Message:     "",
		MessageType: "message",
	}
}

func CreateMessageWithAutoTimestamp(userManager *services.UserManager, user common.User) models.Message {
	return models.Message{
		Timestamp:   GetCurrentTimestamp(),
		MessageId:   int64(len(userManager.GetMessages())),
		User:        user,
		Message:     "",
		MessageType: "message",
	}
}
