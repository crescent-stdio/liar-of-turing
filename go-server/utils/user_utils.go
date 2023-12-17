package utils

import (
	"liar-of-turing/common"
	"liar-of-turing/services"
	"log"
	"sort"
)

// RetrieveUserList: Retrieve ONLINE user list
func RetrieveUserList(userManager *services.UserManager) []common.User {
	players := userManager.GetPlayers()
	var userList []common.User
	for _, v := range players {
		if v.UserName != "admin" && v.IsOnline {
			log.Println("Online User:", v)
			userList = append(userList, v)
		}
	}
	return userList
}

// RetrieveReadyUserList: Retrieve ONLINE player list
func RetrieveReadyUserList(userManager *services.UserManager) []common.User {
	sorted_players := userManager.GetSortedPlayers()
	var playerList []common.User
	for _, v := range sorted_players {
		if v.UserName != "admin" && v.IsOnline && v.PlayerType == "player" {
			playerList = append(playerList, v)
		}
	}
	return playerList
}

func SortUsersByUserName(users []common.User) []common.User {
	sort.Slice(users, func(i, j int) bool {
		return users[i].UserName < users[j].UserName
	})
	return users
}

func CreateRandomUserData(userManager *services.UserManager, webSocketService *services.WebSocketService, uuid string) common.User {
	nicknameId, userName := userManager.GenerateRandomUsername()
	return common.User{
		UserId:     int64(webSocketService.GetTotalUserNum()),
		UserName:   userName,
		NicknameId: nicknameId,
		Role:       "player",
		UUID:       uuid,
		IsOnline:   true,
		PlayerType: "watcher",
	}
}
