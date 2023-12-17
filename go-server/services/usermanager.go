package services

import (
	"liar-of-turing/common"
	"liar-of-turing/global"
	"liar-of-turing/models"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

// UserManager is a struct that manages users
type UserManager struct {
	mutex         sync.Mutex
	MaxPlayer     int
	TotalUsers    int
	Players       map[string]common.User
	SortedPlayers []common.User
	Messages      []models.Message
	PrevMessages  [][]models.Message
	Nicknames     []common.Nickname
	adminUser     common.User
	gptUsers      []common.User
}

// NewUserManager is a constructor for UserManager
func NewUserManager() *UserManager {
	Nicknames := common.GetNicknames()
	players := make(map[string]common.User)

	// add admin user into players
	adminUser := GenerateAdminUser()
	players[adminUser.UUID] = adminUser

	// add gpt users into players
	gptUsers := GenerateGPTUsers(global.GPTNum)
	for _, v := range gptUsers {
		players[v.UUID] = v
	}

	return &UserManager{
		MaxPlayer:     global.MaxPlayer,
		TotalUsers:    0,
		Players:       players,
		SortedPlayers: make([]common.User, 0),
		Messages:      make([]models.Message, 0),
		PrevMessages:  make([][]models.Message, 0),
		Nicknames:     Nicknames,
	}
}

// AddPlayer adds player to players
func (um *UserManager) AddPlayerByUser(player common.User) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.Players[player.UUID] = player
}

// RemovePlayer removes player from players
func (um *UserManager) RemovePlayerByUserUUID(uuid string) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	delete(um.Players, uuid)
}

// GetPlayer gets player from players
func (um *UserManager) GetPlayer(player common.User) (common.User, bool) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	user, exists := um.Players[player.UUID]
	return user, exists
}

// GetPlayers gets players from players
func (um *UserManager) GetPlayers() map[string]common.User {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	return um.Players
}

// GetMaxPlayer gets max player
func (um *UserManager) GetMaxPlayer() int {
	return um.MaxPlayer
}

func (um *UserManager) GetSortedPlayers() []common.User {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	players := make([]common.User, 0)
	players = append(players, um.SortedPlayers...)
	return players
}

func GetMaxPlayer() int {
	return global.MaxPlayer
}

func (um *UserManager) GetAdminUser() common.User {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	return um.adminUser
}

// SetAdminUser
func (um *UserManager) SetAdminUser(adminUser common.User) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.adminUser = adminUser
}

// SetAdminUserByDefault
func (um *UserManager) SetAdminUserByDefault() {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.adminUser = GenerateAdminUser()
}

func (um *UserManager) GetGPTUsers() []common.User {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	return um.gptUsers
}

// SetGPTUsers
func (um *UserManager) SetGPTUsers(gptUsers []common.User) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.gptUsers = gptUsers
}

// SetGPTUsersByDefault
func (um *UserManager) SetGPTUsersByDefault() {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.gptUsers = GenerateGPTUsers(global.GPTNum)
}

// GenerateAdminUser: Generate admin user
func GenerateAdminUser() common.User {
	// um.mutex.Lock()
	// defer um.mutex.Unlock()

	adminUser := common.User{
		UUID:       "0",
		UserId:     0,
		UserName:   "server",
		NicknameId: 0,
		Role:       "admin",
		IsOnline:   false,
		PlayerType: "admin",
	}
	return adminUser
}

// GenerateGPTUsers: Generate gpt users
func GenerateGPTUsers(num int) []common.User {

	var gptUsers []common.User
	for i := 0; i < num; i++ {
		gptUsers = append(gptUsers, GenerateGPTUser(i+1000))
	}
	return gptUsers
}

// GenerateGPTUser: Generate gpt user
func GenerateGPTUser(id int) common.User {
	return common.User{
		UUID:       strconv.Itoa(id),
		UserId:     int64(id),
		UserName:   "",
		NicknameId: 999,
		Role:       "gpt",
		IsOnline:   false,
		PlayerType: "watcher",
	}
}

// SetGPTUser: Set gpt user
func (um *UserManager) SetGPTUser(index int, user common.User) {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	um.gptUsers[index] = user
}

// SetPlayersByUser: Set players by User data
func (um *UserManager) SetPlayersByUser(user common.User) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	for uuid, player := range um.Players {
		if player.UUID == user.UUID {
			um.Players[uuid] = user
			break
		}
	}
}

// SetPlayerByUser(): Set player by user
func (um *UserManager) SetPlayerByUser(user common.User) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.Players[user.UUID] = user
}

// SetAllUsersAsWatchers: Set all users as watchers
func (um *UserManager) SetAllUsersAsWatchers() {
	// um.mutex.Lock()
	// defer um.mutex.Unlock()
	for uuid, player := range um.Players {
		player.PlayerType = "watcher"
		um.Players[uuid] = player
	}
}

// ShuffleSortedUsersRandomly: Shuffle users randomly
func (um *UserManager) SetPlayersRandomlyShuffled(webSocketService *WebSocketService, gameState *GameState) {

	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rand := rand.New(src)
	// Get GPTUsers
	GPTUsers := um.GetGPTUsers()

	users := um.GetSortedPlayers()

	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
	gameState.SetGPTReadyNums()
	GPTReadyNums := gameState.GetGPTReadyNums()
	log.Println("GPTReadyNums:", GPTReadyNums)
	// Swap each GPTUser with the user at the corresponding position in GPTReadyNums.
	for index, gptUser := range GPTUsers {
		readyIndex := GPTReadyNums[index] - 1

		// Iterate through the users to find the one matching the current GPTUser.
		for i, user := range users {
			if user.UUID == gptUser.UUID {
				// Swap the found us다er with the one at readyIndex position.
				users[i], users[readyIndex] = users[readyIndex], users[i]
				break
			}
		}
	}

	// Set Random UserName and NicknameId and If user is GPTUser, set user to GPTUser
	um.ResetNicknameUsed()
	for _, user := range users {
		nicknameId, userName := um.GenerateRandomUsername()
		user.UserName = userName
		user.NicknameId = nicknameId
		um.SetPlayerByUser(user)
		um.SetSortedPlayerByUser(user)

		isGPTUser := false
		for j, gpt := range GPTUsers {
			if gpt.UUID == user.UUID {
				isGPTUser = true
				um.SetGPTUser(j, user)
				break
			}
		}
		if !isGPTUser {
			webSocketService.SetClientsByUserUUID(user)
		}
	}

}

func (um *UserManager) ExcludePlayersFromSelections(webSocketService *WebSocketService, gameState *GameState) (vote int, eliminatedPlayer common.User, remainingPlayerList []common.User) {

	userSelections := gameState.GetNowUserSelections()
	// GPTUsers := um.GetGPTUsers()
	playerList := um.GetSortedPlayers()

	// Make votes map for each player
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
		// for _, gpt := range GPTUsers {
		// 	selection := models.UserSelection{
		// 		User:      gpt,
		// 		Selection: eliminatedPlayer.UserName,
		// 		Reason:    "AI같다.",
		// 	}
		// 	userSelections = append(userSelections, selection)
		// }

	} else {
		remainingPlayerList = []common.User{}
	}

	// Set userSelections
	gameState.SetUserSelections(userSelections)
	return votes[eliminatedPlayer.UserName], eliminatedPlayer, remainingPlayerList
}

// GetNicknames: Get nicknames
func (um *UserManager) GetNicknames() []common.Nickname {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	return um.Nicknames
}

// ResetNicknameUsed: Reset nickname used
func (um *UserManager) ResetNicknameUsed() {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	for i := range um.Nicknames {
		um.Nicknames[i].IsUsed = false
	}
}

// random suffle nicknames and return not used nickname
func (um *UserManager) GenerateRandomUsername() (int, string) {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rnd := rand.New(src)

	nicknames := um.GetNicknames()

	perm := rnd.Perm(len(nicknames))
	for _, v := range perm {
		if !nicknames[v].IsUsed {
			nicknames[v].IsUsed = true
			um.SetNicknamesToUsed(v)
			return nicknames[v].Id, nicknames[v].Nickname
		}
	}
	return -1, ""
}

// SetNicknamesToUsed: Set nicknames to used
func (um *UserManager) SetNicknamesToUsed(index int) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.Nicknames[index].IsUsed = true
}

// GetMessages: Get messages
func (um *UserManager) GetMessages() []models.Message {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	return um.Messages
}

// AddMessage: Add messages
func (um *UserManager) AddMessage(message models.Message) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	if len(um.Messages) > 0 && um.Messages[len(um.Messages)-1] == message {
		return
	}
	um.Messages = append(um.Messages, message)
}

func (um *UserManager) GetTotalUsersByClients(webSocketService *WebSocketService) int {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	return len(webSocketService.GetClients())
}
func (um *UserManager) SetTotalUsers(totalUsers int) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.TotalUsers = totalUsers
}

// SetPrevMessages: Set prev messages
func (um *UserManager) AddPrevMessagesFromMessages() {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	um.PrevMessages = append(um.PrevMessages, um.Messages)
}

// ClearMessages: Clear messages
func (um *UserManager) ClearMessages() {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	um.Messages = make([]models.Message, 0)
}

// ResetPlayers
func (um *UserManager) ResetPlayers() {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	players := um.Players

	um.Players = make(map[string]common.User)
	for uuid, player := range players {
		if player.Role == "human" {
			player.PlayerType = "watcher"
			um.Players[uuid] = player
		}
	}
	um.SortedPlayers = make([]common.User, 0)

}

// GetPlayerByUUID
func (um *UserManager) GetPlayerByUUID(uuid string) (common.User, bool) {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	player, exists := um.Players[uuid]
	return player, exists
}

// AddSortedPlayerByUser
func (um *UserManager) AddSortedPlayerByUser(player common.User) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	// if User is in SortedPlayers, remove it
	for _, v := range um.SortedPlayers {
		if v.UUID == player.UUID {
			return
		}
	}
	um.SortedPlayers = append(um.SortedPlayers, player)
}

// SetSortedPlayers: Set sorted players
func (um *UserManager) SetSortedPlayers(players []common.User) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.SortedPlayers = players
}

// SetSortedPlayerByUser
func (um *UserManager) SetSortedPlayerByUser(player common.User) {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	for i, v := range um.SortedPlayers {
		if v.UUID == player.UUID {
			um.SortedPlayers[i] = player
			break
		}
	}
}

func (um *UserManager) GetPrevMessage() (models.Message, bool) {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	if len(um.Messages) > 0 {
		return um.Messages[len(um.Messages)-1], true
	}
	return models.Message{}, false
}
