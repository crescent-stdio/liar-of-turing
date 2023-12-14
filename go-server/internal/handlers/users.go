package handlers

import (
	"fmt"
	. "liar-of-turing/models"
	"log"
	"math/rand"
	"sort"
	"time"
)

func getUserList() []User {
	var userList []User
	for _, v := range players { //TODO:  or clients?
		if v.UserName != "admin" && v.IsOnline {
			userList = append(userList, v)
		}
	}
	// sortUserList(userList)
	return userList
}

func getPlayerList() []User {
	var playerList []User
	for _, v := range sorted_players {
		// playerList { //TODO:  or clients?
		if v.UserName != "admin" && v.IsOnline && v.PlayerType == "player" {
			playerList = append(playerList, v)
		}
	}
	return playerList
}

func getRandomUsername() (int, string) {
	//random suffle nicknames and return not used nickname
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := rand.Perm(len(nicknames))
	for _, v := range perm {
		if !nicknames[v].IsUsed {
			nicknames[v].IsUsed = true
			log.Println("nickname,", nicknames[v])
			return nicknames[v].Id, nicknames[v].Nickname
		}
	}
	return -1, ""
}

func findUserInUserSelection(userSelections []UserSelection, user User) (UserSelection, bool) {
	for _, v := range userSelections {
		if v.User.UUID == user.UUID {
			return v, true
		}
	}
	return UserSelection{}, false
}

func SortUsersByUserName(users []User) []User {
	sort.Slice(users, func(i, j int) bool {
		return users[i].UserName < users[j].UserName
	})
	return users
}

func ShuffleUsersRandomly(users []User) []User {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rand := rand.New(src)

	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
	readyNum := 1
	if len(users) < 3 {
		readyNum = 1
	} else {
		readyNum = rand.Intn(MaxPlayer-2) + 2
	}
	//swap users
	//find GPTUser
	for i, user := range users {
		if user.UUID == GPTUser.UUID {
			users[i] = users[readyNum]
			users[readyNum] = user
			break
		}
	}

	//getRandomUsername
	for i, user := range users {
		nicknameId, userName := getRandomUsername()
		users[i] = User{
			UserId:     user.UserId,
			UserName:   userName,
			NicknameId: nicknameId,
			Role:       user.Role,
			UUID:       user.UUID,
			IsOnline:   user.IsOnline,
			PlayerType: user.PlayerType,
		}
		for conn, client := range clients {
			if client.UUID == user.UUID {
				clients[conn] = users[i]
				break
			}
		}
		if players[user.UUID].UUID == user.UUID {
			players[user.UUID] = users[i]
		}
		if user.UUID == GPTUser.UUID {
			GPTUser = users[i]
		}

	}
	return users
}

func SelectRandomReadyAndEnteringUser(max_player int) (int, int) {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rand := rand.New(src)

	//1 to max_player-2
	if max_player < 3 {
		return 0, 1
	}
	enterNum := rand.Intn(max_player-2) + 1
	// 2 to max_player-1
	readyNum := rand.Intn(max_player-2) + 2
	return enterNum, readyNum
}

func SetAllUserToWatcher(usermap map[string]User) map[string]User {
	playerMap := make(map[string]User)
	for UUID, v := range usermap {
		player := User{
			UserId:     v.UserId,
			UserName:   v.UserName,
			NicknameId: v.NicknameId,
			Role:       v.Role,
			UUID:       v.UUID,
			IsOnline:   v.IsOnline,
			PlayerType: "watcher",
		}
		playerMap[UUID] = player
	}
	return playerMap
}

func removePlayerListFromUserSelection(playerList []User) (vote int, eliminatedPlayer User, remainingPlayerList []User) {

	for _, v := range userSelections {
		log.Println("userSelection:", v.Selection)
	}
	votes := make(map[string]int)
	for _, v := range userSelections {
		votes[v.Selection]++
	}
	log.Println("votes:", votes)

	//sort decreasing order by vote
	sort.Slice(playerList, func(i, j int) bool {
		return votes[playerList[i].UserName] > votes[playerList[j].UserName]
	})

	if len(playerList) > 0 {
		eliminatedPlayer = playerList[0]
		remainingPlayerList = playerList[1:]

		// TODO: for range for GPT
		selection := UserSelection{
			User:      GPTUser,
			Selection: eliminatedPlayer.UserName,
			Reason:    "AI같다.",
		}
		userSelections = append(userSelections, selection)
	} else {
		remainingPlayerList = []User{}
	}

	log.Println("remainingPlayerList:", remainingPlayerList)
	log.Println("eliminatedPlayer:", eliminatedPlayer)

	return votes[eliminatedPlayer.UserName], eliminatedPlayer, remainingPlayerList
}

func processEnterHuman(e WsPayload) {
	nowUser, exists := players[e.User.UUID]
	log.Println("processEnterHuman")
	log.Println("isGameStarted:", isGameStarted)
	log.Println("gameTurnsLeft:", gameTurnsLeft)

	mutex.Lock()
	if !exists {
		nicknameId, userName := getRandomUsername()
		nowUser = User{
			UserId:     int64(len(clients)),
			UserName:   userName,
			NicknameId: nicknameId,
			Role:       "human",
			UUID:       e.User.UUID,
			IsOnline:   true,
			PlayerType: "watcher",
		}
		// players[e.User.UUID] = nowUser
	}
	nowUser.IsOnline = true
	log.Println(nowUser.NicknameId)
	players[e.User.UUID] = nowUser
	clients[e.Conn] = nowUser

	response := WsJsonResponse{
		MaxPlayer:      MaxPlayer,
		Timestamp:      e.Timestamp,
		Action:         "human_info",
		MessageType:    "system",
		Message:        fmt.Sprintf("%s님이 입장했습니다.", nowUser.UserName),
		MessageLogList: messages,
		User:           nowUser,
		OnlineUserList: getUserList(),
		PlayerList:     getPlayerList(),
		GameTurnsLeft:  gameTurnsLeft,
		GameRound:      gameRound,
		IsGameStarted:  isGameStarted,
		IsGameOver:     isGameOver,
	}
	broadcastToAll(clients, response)
	mutex.Unlock()

	log.Println("GPTReadyNum:", GPTReadyNum)
	log.Println("len(getPlayerList()):", len(getPlayerList()))
	log.Println("GPTEnterNum:", GPTEnterNum)

	if GPTEnterNum <= len(getUserList()) && !GPTUser.IsOnline {
		time.Sleep(time.Second * 1)
		log.Println("GPTEnterNum:", GPTEnterNum)
		log.Println("len(getUserList()):", len(getUserList()))
		log.Println("GPTUser:", GPTUser)
		mutex.Lock()
		nicknameId, userName := getRandomUsername()
		GPTUser = User{
			UserId:     int64(len(clients)),
			UserName:   userName,
			NicknameId: nicknameId,
			Role:       "player",
			UUID:       GPTUser.UUID,
			IsOnline:   true,
			PlayerType: "watcher",
		}
		players[GPTUser.UUID] = GPTUser
		// clients[WebSocketConnection{Conn: nil}] = GPTUser

		GPTresponse := WsJsonResponse{
			MaxPlayer:      MaxPlayer,
			Timestamp:      e.Timestamp,
			Action:         "human_info",
			MessageType:    "system",
			Message:        fmt.Sprintf("%s님이 입장했습니다.", GPTUser.UserName),
			MessageLogList: messages,
			User:           GPTUser,
			OnlineUserList: getUserList(),
			PlayerList:     getPlayerList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
		broadcastToAll(clients, GPTresponse)
		mutex.Unlock()
	} else if GPTReadyNum >= len(getPlayerList()) && GPTUser.PlayerType == "watcher" {
		log.Println("GPTReadyNum:", GPTReadyNum)
		log.Println("len(getPlayerList()):", len(getPlayerList()))
		log.Println("GPTUser:", GPTUser)

		time.Sleep(time.Second * 1)
		processReadyGPT()
	}

	if gameTurnsLeft == 0 {
		if isGameOver {
			processGameOver(e.Timestamp)
			mutex.Lock()
			resultResponse := WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "game_over",
				MessageType:    "alert",
				Message:        "게임이 종료되었습니다.",
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
			broadcastToAll(clients, resultResponse)
			mutex.Unlock()
			return
		}
		// messages = append(messages, waitMessage)
		// in userSelection
		if exists && nowUser.PlayerType == "player" {
			mutex.Lock()
			var waitResponse WsJsonResponse
			_, exists := findUserInUserSelection(userSelections, nowUser)
			log.Println("exists:", exists)
			log.Println(userSelections)
			log.Println(nowUser)
			if exists {
				waitResponse = WsJsonResponse{
					Timestamp:      e.Timestamp,
					Action:         "wait_for_players",
					MessageType:    "alert",
					Message:        "라운드가 종료되었습니다. 다음 플레이어의 선택을 기다리세요.",
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
			} else {
				waitResponse = WsJsonResponse{
					Timestamp:      e.Timestamp,
					Action:         "submit_ai",
					MessageType:    "alert",
					Message:        "라운드가 종료되었습니다. 다음 라운드를 위해 AI를 선택하세요.",
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
			}
			broadCastToSomeone(clients, e.Conn, waitResponse)
			mutex.Unlock()
		}
		return
	}

	if isGameStarted {
		nextUser := gameInfo[gameRound-1].PlayerList[gameInfo[gameRound-1].NowUserIndex]
		if nextUser.UUID == e.User.UUID {
			mutex.Lock()
			nowUser.PlayerType = "player"
			someoneMessage := Message{
				Timestamp:   e.Timestamp,
				MessageId:   int64(len(messages)),
				User:        adminUser,
				Message:     fmt.Sprintf("%s님의 차례입니다.", nextUser.UserName),
				MessageType: "alert",
			}
			// messages = append(messages, someoneMessage)
			someoneResponse := WsJsonResponse{
				Timestamp:      e.Timestamp,
				Action:         "your_turn",
				MessageType:    "alert",
				Message:        someoneMessage.Message,
				MaxPlayer:      MaxPlayer,
				User:           adminUser, // nextUser,
				MessageLogList: messages,
				PlayerList:     getPlayerList(),
				OnlineUserList: getUserList(),
				GameTurnsLeft:  gameTurnsLeft,
				GameRound:      gameRound,
				IsGameStarted:  isGameStarted,
				IsGameOver:     isGameOver,
			}
			nextConn, _ := GetWebSocketClientByUUID(nextUser.UUID)
			// if exists {
			broadCastToSomeone(clients, nextConn, someoneResponse)
			mutex.Unlock()
			// } else {
			// 	mutex.Unlock()
		}
	}

}

func processLeftUser(e WsPayload, response *WsJsonResponse) {
	mutex.Lock()
	if leftUser, ok := clients[e.Conn]; ok {
		leftUser.IsOnline = false
		players[e.User.UUID] = leftUser

		delete(clients, e.Conn)
		e.Conn.CloseWebSocketConnection()

		*response = WsJsonResponse{
			Action:    "update_state",
			Timestamp: e.Timestamp,
			MaxPlayer: MaxPlayer,

			// Action:         "user_list",
			MessageLogList: messages,
			OnlineUserList: getUserList(),
			PlayerList:     getPlayerList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
	} else {
		log.Println("User not found in clients map")
		// *response = WsJsonResponse{} // Empty response
		response := WsJsonResponse{
			Action:    "update_state",
			Timestamp: e.Timestamp,
			// Action:         "user_list",
			MessageLogList: messages,
			OnlineUserList: getUserList(),
			MaxPlayer:      MaxPlayer,
			PlayerList:     getPlayerList(),
			GameTurnsLeft:  gameTurnsLeft,
			GameRound:      gameRound,
			IsGameStarted:  isGameStarted,
			IsGameOver:     isGameOver,
		}
		broadcastToAll(clients, response)
	}
	mutex.Unlock()
}
