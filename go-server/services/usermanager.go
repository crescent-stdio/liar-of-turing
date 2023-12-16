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
	gptUser       []common.User
}

// NewUserManager is a constructor for UserManager
func NewUserManager() *UserManager {
	Nicknames := common.GetNicknames()
	players := make(map[string]common.User)

	// add admin user into players
	admin := GenerateAdminUser()
	players[admin.UUID] = admin

	// add gpt users into players
	gptUsers := GenerateGPTUsers(global.GPTNum)
	for _, v := range gptUsers {
		players[v.UUID] = v
	}

	return &UserManager{
		MaxPlayer:     global.MaxPlayer,
		TotalUsers:    0,
		Players:       make(map[string]common.User),
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

func GetAdminUser() common.User {
	return GenerateAdminUser()
}

func (um *UserManager) GetGPTUsers() []common.User {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	gptUsers := make([]common.User, 0)
	gptUsers = append(gptUsers, um.gptUser...)
	return gptUsers
}

// SetGPTUsers
func (um *UserManager) SetGPTUsers(gptUsers []common.User) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.gptUser = gptUsers
}

// GenerateAdminUser: Generate admin user
func GenerateAdminUser() common.User {
	return common.User{
		UUID:       "0",
		UserId:     0,
		UserName:   "server",
		NicknameId: 999,
		Role:       "admin",
		IsOnline:   false,
		PlayerType: "admin",
	}
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
		Role:       "player",
		IsOnline:   false,
		PlayerType: "watcher",
	}
}

// SetGPTUser: Set gpt user
func (um *UserManager) SetGPTUser(index int, user common.User) {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	um.gptUser[index] = user
}

// SetPlayersByUUID: Set players by User data
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

// SetAllUsersAsWatchers: Set all users as watchers
func (um *UserManager) SetAllUsersAsWatchers() {
	um.mutex.TryLock()
	for uuid, player := range um.Players {
		player.PlayerType = "watcher"
		um.Players[uuid] = player
	}
	um.mutex.Unlock()
}

// ShuffleSortedUsersRandomly: Shuffle users randomly
func (um *UserManager) SetRandomlyShuffledPlayers(webSocketService *WebSocketService, gameState *GameState) {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rand := rand.New(src)
	// Get GPTUsers
	GPTUsers := um.GetGPTUsers()

	users := um.GetSortedPlayers()

	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
	gameState.SetGPTReadyNums()
	GPTReadyNums := gameState.GetGPTReadyNums()
	// Swap each GPTUser with the user at the corresponding position in GPTReadyNums.
	for index, gptUser := range GPTUsers {
		readyIndex := GPTReadyNums[index]

		// Iterate through the users to find the one matching the current GPTUser.
		for i, user := range users {
			if user.UUID == gptUser.UUID {
				// Swap the found user with the one at readyIndex position.
				users[i], users[readyIndex] = users[readyIndex], users[i]
				break
			}
		}
	}

	// Set Random UserName and NicknameId and If user is GPTUser, set user to GPTUser
	for _, user := range users {
		nicknameId, userName := um.GenerateRandomUsername()
		user.UserName = userName
		user.NicknameId = nicknameId
		webSocketService.SetClientsByUserUUID(user)
		um.AddPlayerByUser(user)

		for j, gpt := range GPTUsers {
			if gpt.UUID == user.UUID {
				um.SetGPTUser(j, user)
				break
			}
		}

	}
	um.SetSortedPlayers(users)
}

func (um *UserManager) ExcludePlayersFromSelections(webSocketService *WebSocketService, gameState *GameState) (vote int, eliminatedPlayer common.User, remainingPlayerList []common.User) {

	userSelections := gameState.GetUserSelections()
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
	um.mutex.Unlock()
	defer um.mutex.Unlock()
	return um.Nicknames
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
			return nicknames[v].Id, nicknames[v].Nickname
		}
	}
	return -1, ""
}

// GetMessages: Get messages
func (um *UserManager) GetMessages() []models.Message {
	um.mutex.Unlock()
	defer um.mutex.Unlock()
	return um.Messages
}

// AddMessage: Add messages
func (um *UserManager) AddMessage(message models.Message) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.Messages = append(um.Messages, message)
}

func (um *UserManager) GetTotalUsersByClients(webSocketService *WebSocketService) int {
	um.mutex.Unlock()
	defer um.mutex.Unlock()
	return len(webSocketService.GetClients())
}
func (um *UserManager) SetTotalUsers(totalUsers int) {
	um.mutex.Unlock()
	defer um.mutex.Unlock()
	um.TotalUsers = totalUsers
}

// GetTotalUsers: Get total users
func (um *UserManager) GetAdminUser() common.User {
	um.mutex.Unlock()
	defer um.mutex.Unlock()
	return um.adminUser
}

// SetPrevMessages: Set prev messages
func (um *UserManager) AddPrevMessagesFromMessages() {
	um.mutex.Unlock()
	defer um.mutex.Unlock()
	um.PrevMessages = append(um.PrevMessages, um.Messages)
}

// ClearMessages: Clear messages
func (um *UserManager) ClearMessages() {
	um.mutex.Unlock()
	defer um.mutex.Unlock()
	um.Messages = make([]models.Message, 0)
}

// ResetPlayers
func (um *UserManager) ResetPlayers() {
	um.mutex.Unlock()
	defer um.mutex.Unlock()
	um.Players = make(map[string]common.User)
	um.SortedPlayers = make([]common.User, 0)
}

// GetPlayerByUUID
func (um *UserManager) GetPlayerByUUID(uuid string) (common.User, bool) {
	um.mutex.Unlock()
	defer um.mutex.Unlock()
	player, exists := um.Players[uuid]
	return player, exists
}

// AddSortedPlayer
func (um *UserManager) AddSortedPlayer(player common.User) {
	um.mutex.Unlock()
	defer um.mutex.Unlock()
	um.SortedPlayers = append(um.SortedPlayers, player)
}

func (um *UserManager) SetSortedPlayers(players []common.User) {
	um.mutex.Unlock()
	defer um.mutex.Unlock()
	um.SortedPlayers = players
}
