package services

import (
	"liar-of-turing/common"
	"liar-of-turing/global"
	"liar-of-turing/models"
	"log"
	"math/rand"
	"sync"
	"time"
)

type GameStatus struct {
	IsStarted     bool
	IsOver        bool
	TurnsNum      int
	RoundNum      int
	MaxPlayer     int
	InfoIdx       int
	IsUsersVoting bool
}

type GameState struct {
	mutex        sync.Mutex
	Status       GameStatus
	Info         []models.Game
	GPTEntryNums []int
	GPTReadyNums []int
	questions    []string
}

func NewGameState() *GameState {
	gameTurnNum := global.GetGameTurnNum()
	gameRoundNum := global.GetGameRoundNum()
	gameMaxPlayer := global.GetMaxPlayer()
	questions := common.GetQuestions()
	gameStatus := GameStatus{
		IsStarted:     false,
		IsOver:        false,
		IsUsersVoting: false,
		TurnsNum:      gameTurnNum,
		RoundNum:      gameRoundNum,
		MaxPlayer:     gameMaxPlayer,
		InfoIdx:       0,
	}

	return &GameState{
		mutex:        sync.Mutex{},
		Status:       gameStatus,
		Info:         make([]models.Game, 0),
		GPTEntryNums: make([]int, 0),
		GPTReadyNums: make([]int, 0),
		questions:    questions,
	}
}

// SetGPTReadyNums: Get random GPTReadyNums(Timing of GPT ready) / [2, maxPlayer)
func (gs *GameState) SetGPTReadyNums() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rand := rand.New(src)

	GPTEntryNums := gs.GPTEntryNums
	// Create a slice with all possible numbers in the range [1, maxPlayer).
	possibleNums := make([]int, 0)
	for i := 1; i <= global.MaxPlayer; i++ {
		possibleNums = append(possibleNums, i)
	}

	// Shuffle the slice to randomize the order.
	rand.Shuffle(len(possibleNums), func(i, j int) {
		possibleNums[i], possibleNums[j] = possibleNums[j], possibleNums[i]
	})

	// Select the first gptNum numbers from the shuffled slice.
	GPTReadyNums := make([]int, global.GPTNum)
	// while items of GPTEntryNums not 0, set global.GPTNum

	for _, v := range possibleNums {
		for idx, w := range GPTEntryNums {
			if v > w && GPTReadyNums[idx] == 0 {
				GPTReadyNums[idx] = v
				break
			}
		}
	}
	for idx, GPTReadyNum := range GPTReadyNums {
		GPTEntryNum := GPTEntryNums[idx]
		if GPTEntryNum > GPTReadyNum {
			GPTReadyNums[idx] = common.Min(global.GPTNum, GPTEntryNum+1)
		}
	}

	gs.GPTReadyNums = GPTReadyNums
}

// GetGPTReadyNums: Get GPTReadyNums
func (gs *GameState) GetGPTReadyNums() []int {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.GPTReadyNums
}

// SetGPTEntryNums: selects 2 numbers from 1 to max_player-1
func (gs *GameState) SetGPTEntryNums() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rand := rand.New(src)

	// Create a slice with all possible numbers in the range [1, maxPlayer-1).
	possibleNums := make([]int, 0)
	for i := 0; i <= global.MaxPlayer; i++ {
		possibleNums = append(possibleNums, i)
	}

	// Shuffle the slice to randomize the order.
	rand.Shuffle(len(possibleNums), func(i, j int) {
		possibleNums[i], possibleNums[j] = possibleNums[j], possibleNums[i]
	})

	// Select the first gptNum numbers from the shuffled slice.
	entryNums := make([]int, global.GPTNum)
	copy(entryNums, possibleNums[:global.GPTNum])

	gs.GPTEntryNums = entryNums
}

// GetGPTEntryNums: Get GPTEntryNums
func (gs *GameState) GetGPTEntryNums() []int {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.GPTEntryNums
}

func (gs *GameState) GetAllGameInfo() []models.Game {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.Info
}

// SetQuestionsRandomly: Set questions randomly
func (gs *GameState) SetQuestionsRandomly() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rand := rand.New(src)

	questions := gs.questions
	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})
	gs.questions = questions
}

// GetQuestion
func (gs *GameState) GetQuestion() string {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.questions[0]
}

func (gs *GameState) SearchUserInUserSelections(user common.User) (models.UserSelection, bool) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	idx := gs.Status.InfoIdx
	for _, v := range gs.Info[idx].UserSelections {
		if v.User.UUID == user.UUID {
			return v, true
		}
	}
	return models.UserSelection{}, false
}

// GetNowUserSelections
func (gs *GameState) GetNowUserSelections() []models.UserSelection {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	if len(gs.Info) == 0 {
		return make([]models.UserSelection, 0)
	}
	return gs.Info[gs.Status.InfoIdx].UserSelections
}

// SetUserSelections: Set user selections

func (gs *GameState) SetUserSelections(userSelections []models.UserSelection) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	if len(gs.Info) == 0 {
		return
	}
	gs.Info[gs.Status.InfoIdx].UserSelections = userSelections
}

// AddUserSelection: Add user selection
func (gs *GameState) AddUserSelection(userSelection models.UserSelection) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	infoIdx := gs.Status.InfoIdx
	gs.Info[infoIdx].UserSelections = append(gs.Info[infoIdx].UserSelections, userSelection)
}

// SetGameInfo: Set game info
func (gs *GameState) SetGameInfo(gameInfo []models.Game) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Info = gameInfo
}

// GetGameTurnsLeft: Get game turns left
func (gs *GameState) GetStatus() GameStatus {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.Status
}

func (gs *GameState) SetIfGameStarted() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Status.IsStarted = true
	gs.Status.IsOver = false
}

func (gs *GameState) SetIfGameTotallyOver() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Status.IsStarted = false
	gs.Status.IsOver = true

	// TODO: It might be not necessary?
}

func (gs *GameState) SetIfRoundIsOver() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Status.IsStarted = false
	gs.Status.IsOver = true
	if len(gs.Info) == 0 {
		return
	}
	gs.Info[gs.Status.InfoIdx].TurnsLeft = 0
}

func (gs *GameState) SetIfResetRound(userManager *UserManager) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	// infoIndex is not changed
	infoIdx := gs.Status.InfoIdx
	gs.Status.IsStarted = false
	gs.Status.IsOver = false
	gs.Status.IsUsersVoting = false

	gs.Info[infoIdx].TurnsLeft = global.GameTurnNum * len(gs.Info[infoIdx].PlayerList)

	gs.Info[infoIdx].NowUserIndex = 0
	gs.Info[infoIdx].UserSelections = make([]models.UserSelection, 0)

	gs.Info[infoIdx].Messages = make([]models.Message, 0)
	gs.Info[infoIdx].PlayerList = userManager.GetSortedPlayers()
}

func (gs *GameState) SetIfGameTotallyReset(userManager *UserManager) {
	gs.Status.IsStarted = false
	gs.Status.IsOver = false
	gs.Status.InfoIdx = 0
	gs.Status.IsUsersVoting = false

	// Reset Game Info
	gs.Info = make([]models.Game, 0)

	// Reset GPTEntryNums & GPTReadyNums
	gs.SetGPTEntryNums()
	gs.SetGPTReadyNums()

	// Reset Players & Sort Players
	userManager.ResetPlayers()
}

func (gs *GameState) GetNextTurnPlayer() (common.User, bool) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	if len(gs.Info) == 0 {
		return common.User{}, false
	}
	return gs.Info[gs.Status.InfoIdx].PlayerList[gs.Info[gs.Status.InfoIdx].NowUserIndex], true
}

// InitializeRoundInfo: Initialize game info
func (gs *GameState) InitializeRoundInfo(userManager *UserManager, webSocketService *WebSocketService) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	maxPlayer := len(userManager.GetSortedPlayers())
	roundNum := gs.Status.RoundNum
	nowRound := len(gs.Info) + 1
	turnsLeft := gs.Status.TurnsNum * maxPlayer * roundNum
	gs.Status.InfoIdx = 0
	gs.Info = make([]models.Game, 0)

	game := models.Game{
		NowUserIndex:   0,
		MaxPlayer:      maxPlayer,
		PlayerList:     userManager.GetSortedPlayers(),
		Round:          nowRound,
		TurnsLeft:      turnsLeft,
		UserSelections: make([]models.UserSelection, 0),
		Messages:       make([]models.Message, 0),
	}

	gs.Info = append(gs.Info, game)
	gs.Status.InfoIdx = len(gs.Info) - 1

	gs.Status.IsStarted = true
	gs.Status.IsOver = false
	gs.Status.IsUsersVoting = false
}

// SetGameInfoPlayerList: Set game info player list
func (gs *GameState) SetGameInfoPlayerList(idx int, playerList []common.User) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Info[idx].PlayerList = playerList
}

// SetNextTurn: Set next turn
func (gs *GameState) SetNextTurnInfo() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Info[gs.Status.InfoIdx].NowUserIndex = (gs.Info[gs.Status.InfoIdx].NowUserIndex + 1) % gs.Info[gs.Status.InfoIdx].MaxPlayer
	gs.Info[gs.Status.InfoIdx].TurnsLeft--
}

// SetGameTurnNum
func (gs *GameState) SetGameTurnNum(turnNum int) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Status.TurnsNum = turnNum

	infoIdx := gs.Status.InfoIdx
	gs.Info[infoIdx].TurnsLeft = turnNum * len(gs.Info[infoIdx].PlayerList)
}

// SetGameRoundNum
func (gs *GameState) SetGameRoundNum(roundNum int) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Status.RoundNum = roundNum
}

// GetNowGameInfo
func (gs *GameState) GetNowGameInfo() models.Game {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	if len(gs.Info) == 0 {
		return models.Game{
			NowUserIndex:   0,
			MaxPlayer:      0,
			PlayerList:     make([]common.User, 0),
			Round:          gs.Status.RoundNum,
			TurnsLeft:      gs.Status.TurnsNum * gs.Status.RoundNum,
			UserSelections: make([]models.UserSelection, 0),
			Messages:       make([]models.Message, 0),
		}
	}
	return gs.Info[gs.Status.InfoIdx]
}

// CheckIsRoundOver: Check if round is over
func (gs *GameState) CheckIsRoundOver() bool {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	if len(gs.Info) == 0 {
		return false
	}
	log.Println("gs.Info[gs.Status.InfoIdx].TurnsLeft", gs.Info[gs.Status.InfoIdx].TurnsLeft)
	return gs.Info[gs.Status.InfoIdx].TurnsLeft == 0
}

// CheckIsGameStarted: Check if game is started
func (gs *GameState) CheckIsGameStarted() bool {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.Status.IsStarted
}

// CheckAllUserVote: Check if all user vote
func (gs *GameState) CheckAllUserVoted(userManager *UserManager) bool {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	if len(gs.Info) == 0 {
		return false
	}
	log.Println("len(gs.Info[gs.Status.InfoIdx].UserSelections)", len(gs.Info[gs.Status.InfoIdx].UserSelections))
	log.Println("len(gs.Info[gs.Status.InfoIdx].PlayerList)", len(gs.Info[gs.Status.InfoIdx].PlayerList))
	GPTNum := len(userManager.GetGPTUsers())
	return len(gs.Info[gs.Status.InfoIdx].UserSelections) >= len(gs.Info[gs.Status.InfoIdx].PlayerList)-GPTNum
}

// SetMaxPlayer
func (gs *GameState) SetMaxPlayer(maxPlayer int) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Status.MaxPlayer = maxPlayer
}

// CheckAllUserReady: Check if all user ready
func (gs *GameState) CheckAllUserReady(userManager *UserManager) bool {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	totalPlayerNum := len(userManager.GetSortedPlayers())
	if len(gs.Info) == 0 {
		return totalPlayerNum == gs.Status.MaxPlayer
	}
	return totalPlayerNum == gs.Info[gs.Status.InfoIdx].MaxPlayer
}

// CheckIsGameOver: Check if game is over
func (gs *GameState) CheckIsGameOver() bool {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.Status.IsOver
}

func (gs *GameState) GetHumanPlayerNum(userManager *UserManager) int {
	humanPlayerNum := 0
	for _, v := range userManager.GetSortedPlayers() {
		if v.Role == "human" {
			humanPlayerNum++
		}
	}
	return humanPlayerNum
}

// CheckAllHumanPlayerReady: Check if all human player ready
func (gs *GameState) CheckAllHumanPlayerReady(userManager *UserManager) bool {
	log.Println("CheckAllHumanPlayerReady")
	humanPlayerNum := gs.GetHumanPlayerNum(userManager)
	log.Println("humanPlayerNum", humanPlayerNum)
	GPTNum := len(userManager.GetGPTUsers())
	if len(gs.Info) == 0 {
		return humanPlayerNum == gs.Status.MaxPlayer-GPTNum
	}
	return humanPlayerNum == gs.Info[gs.Status.InfoIdx].MaxPlayer-GPTNum
}

// SetIsUsersVotingTrue: Set isUsersVoting true
func (gs *GameState) SetIsUsersVotingTrue() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Status.IsUsersVoting = true
}
