package handlers

import (
	"fmt"
	"log"
)

func initGameInfo() {

	gameInfo = make([]Game, 0)
	for i := 0; i < gameRoundNum; i++ {
		game := Game{
			NowUserIndex:   0,
			MaxPlayer:      MaxPlayer,
			OnlineUserList: getUserList(),
			PlayerList:     getPlayerList(),
			TurnsLeft:      gameTurnNum * MaxPlayer,
			UserSelections: make([]UserSelection, 0),
		}
		gameInfo = append(gameInfo, game)
	}

}

func processGameOver(timestamp int64) {
	// time.Sleep(time.Second * 5)
	message := Message{
		Timestamp:   timestamp,
		MessageId:   int64(len(messages)),
		User:        adminUser,
		Message:     "게임이 종료되었습니다. 누가 AI였는지 확인해보세요.",
		MessageType: "alert",
	}

	response := WsJsonResponse{
		Timestamp:      timestamp,
		Action:         "choose_ai",
		MessageType:    "alert",
		Message:        message.Message,
		MaxPlayer:      MaxPlayer,
		User:           adminUser,
		MessageLogList: messages,
		OnlineUserList: getUserList(),
		PlayerList:     getPlayerList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      gameRound,
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	broadcastToAll(clients, response)

	// gameInfo[gameRound-1].UserSelections = make([]UserSelection, 0)
	gameInfo[gameRound-1].Messages = messages

}
func sendUserSelectionsForConsole(timestamp int64) {
	log.Println("=================================================")
	log.Println("sendUserSelectionsForConsole")

	conn, exists := GetWebSocketClientByUUID(adminUser.UUID)
	if !exists {
		log.Println("admin is not connected")
		return
	}

	mutex.Lock()
	chatResponse := WsJsonResponse{
		Timestamp:      timestamp,
		Action:         "send_result",
		MessageType:    "message",
		Message:        "send_result",
		MaxPlayer:      MaxPlayer,
		User:           GPTUser,
		MessageLogList: prevMessages,
		OnlineUserList: getUserList(),
		PlayerList:     getPlayerList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      gameRound,
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
		UserSelection:  userSelections,
	}
	broadCastToSomeone(clients, conn, chatResponse)
	mutex.Unlock()

}

func processChooseAI(e WsPayload) {
	selection := UserSelection{
		User:      e.User,
		Selection: e.UserSelection.Selection,
		Reason:    e.UserSelection.Reason,
	}
	userSelections = append(userSelections, selection)

	// next round
	if len(userSelections) >= MaxPlayer-GPTNum {
		gameInfo[gameRound-1].PlayerList = sorted_players
		vote, eliminatedPlayer, remainingPlayerList := removePlayerListFromUserSelection(sorted_players)
		sorted_players = remainingPlayerList
		mutex.Lock()
		resultMessage := Message{
			Timestamp:   e.Timestamp,
			MessageId:   int64(len(messages)),
			User:        adminUser,
			Message:     fmt.Sprintf("%d라운드가 종료되었습니다. 탈락자는 %d표를 받은 [%s]입니다.", gameRound, vote, eliminatedPlayer.UserName),
			MessageType: "alert",
		}
		messages = []Message{resultMessage}
		resultResponse := WsJsonResponse{
			Timestamp:      e.Timestamp,
			Action:         "next_round",
			MessageType:    "alert",
			Message:        resultMessage.Message,
			MaxPlayer:      MaxPlayer,
			User:           User{},
			MessageLogList: messages,
			OnlineUserList: getUserList(),
			PlayerList:     getPlayerList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
		broadcastToAll(clients, resultResponse)
		mutex.Unlock()

		sendUserSelectionsForConsole(e.Timestamp)
		gameRound++

		if gameRound >= gameRoundNum {
			isGameOver = true
			mutex.Lock()
			alivePlayerString := ""
			for _, v := range remainingPlayerList {
				alivePlayerString += fmt.Sprintf("[%s] ", v.UserName)
			}
			finalMessage := Message{
				Timestamp:   e.Timestamp,
				MessageId:   int64(len(messages)),
				User:        adminUser,
				Message:     fmt.Sprintf("게임이 종료되었습니다. 최종 생존자는 %s입니다.", alivePlayerString),
				MessageType: "alert",
			}
			messages = []Message{finalMessage}
			finalResponse := WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "game_over",
				MessageType:    "alert",
				Message:        finalMessage.Message,
				MaxPlayer:      MaxPlayer,
				User:           User{},
				MessageLogList: messages,
				OnlineUserList: getUserList(),
				PlayerList:     getPlayerList(),
				GameTurnsLeft:  gameTurnsLeft,
				GameRound:      gameRound,
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
			}
			broadcastToAll(clients, finalResponse)
			mutex.Unlock()

			return

		}

		sorted_players = ShuffleUsersRandomly(sorted_players)
		gameTurnsLeft = gameTurnNum * len(sorted_players)
		gameInfo[gameRound-1].TurnsLeft = gameTurnsLeft
		gameInfo[gameRound-1].PlayerList = sorted_players
		gameInfo[gameRound-1].NowUserIndex = 0
		userSelections = make([]UserSelection, 0)

		mutex.Lock()
		nextMessage := Message{
			Timestamp:   e.Timestamp,
			MessageId:   int64(len(messages)),
			User:        adminUser,
			Message:     fmt.Sprintf("%d라운드가 시작되었습니다.", gameRound),
			MessageType: "alert",
		}
		messages = append(messages, nextMessage)
		nextResponse := WsJsonResponse{
			Timestamp:      e.Timestamp,
			Action:         "update_state",
			MessageType:    "alert",
			Message:        nextMessage.Message,
			MaxPlayer:      MaxPlayer,
			User:           User{},
			MessageLogList: messages,
			OnlineUserList: getUserList(),
			PlayerList:     getPlayerList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
		broadcastToAll(clients, nextResponse)
		mutex.Unlock()

		processNextTurn(e.Timestamp)

	}
}
func processNextTurn(timestamp int64) {

	mutex.Lock()
	nextUser := gameInfo[gameRound-1].PlayerList[gameInfo[gameRound-1].NowUserIndex]
	log.Println("nextUser:", nextUser)
	// if nextUser.UUID == GPTUser.UUID {
	// 	SendGPTMessage(timestamp)
	// 	return
	// }
	someoneMessage := Message{
		Timestamp:   timestamp,
		MessageId:   int64(len(messages)),
		User:        adminUser,
		Message:     fmt.Sprintf("%s님의 차례입니다.", nextUser.UserName),
		MessageType: "alert",
	}
	// messages = append(messages, someoneMessage)
	someoneResponse := WsJsonResponse{
		Timestamp:   timestamp,
		Action:      "your_turn",
		MessageType: "alert",
		Message:     someoneMessage.Message,
		MaxPlayer:   MaxPlayer,
		// User:           nextUser,
		User:           adminUser,
		MessageLogList: messages,
		PlayerList:     getPlayerList(),
		OnlineUserList: getUserList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      1,
	}
	nextConn, _ := GetWebSocketClientByUUID(nextUser.UUID)
	// if exists {
	// 	broadCastToSomeone(someoneResponse, nextConn)
	// 	mutex.Unlock()
	// } else {
	// 	mutex.Unlock()
	// 	SendGPTMessage(timestamp)
	// }
	broadCastToSomeone(clients, nextConn, someoneResponse)
	mutex.Unlock()
}

func makeMessagesFromGameInfo(gameInfo []Game) []Message {
	messages := make([]Message, 0)
	for _, v := range gameInfo {
		messages = append(messages, v.Messages...)
	}
	return messages

}
