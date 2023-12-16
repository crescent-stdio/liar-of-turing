package services

import (
	"liar-of-turing/common"
	"liar-of-turing/global"
	"liar-of-turing/models"
	"math/rand"
	"sync"
	"time"
)

type GameStatus struct {
	IsStarted bool
	IsOver    bool
	TurnsNum  int
	RoundNum  int
	MaxPlayer int
	InfoIdx   int
}

type GameState struct {
	mutex        sync.Mutex
	Status       GameStatus
	Info         []models.Game
	GPTEntryNums []int
	GPTReadyNums []int
}

func NewGameState() *GameState {
	gameTurnNum := global.GetGameTurnNum()
	gameRoundNum := global.GetGameRoundNum()
	gameMaxPlayer := global.GetMaxPlayer()
	gameStatus := GameStatus{
		IsStarted: false,
		IsOver:    false,
		TurnsNum:  gameTurnNum,
		RoundNum:  gameRoundNum,
		MaxPlayer: gameMaxPlayer,
		InfoIdx:   -1,
	}

	return &GameState{
		mutex:        sync.Mutex{},
		Status:       gameStatus,
		Info:         make([]models.Game, 0),
		GPTEntryNums: make([]int, 0),
		GPTReadyNums: make([]int, 0),
	}
}

// SetGPTReadyNums: Get random GPTReadyNums(Timing of GPT ready) / [2, maxPlayer)
func (gs *GameState) SetGPTReadyNums() {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rand := rand.New(src)

	// Create a slice with all possible numbers in the range [2, maxPlayer).
	possibleNums := make([]int, global.MaxPlayer-2)
	for i := 2; i < global.MaxPlayer; i++ {
		possibleNums[i] = i
	}

	// Shuffle the slice to randomize the order.
	rand.Shuffle(len(possibleNums), func(i, j int) {
		possibleNums[i], possibleNums[j] = possibleNums[j], possibleNums[i]
	})

	// Select the first gptNum numbers from the shuffled slice.
	GPTReadyNums := make([]int, global.GPTNum)
	copy(GPTReadyNums, possibleNums[:global.GPTNum])

	gs.mutex.Lock()
	defer gs.mutex.Unlock()
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
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rand := rand.New(src)

	// Create a slice with all possible numbers in the range [1, maxPlayer-1).
	possibleNums := make([]int, global.MaxPlayer-1)
	for i := 1; i < global.MaxPlayer-1; i++ {
		possibleNums[i] = i
	}

	// Shuffle the slice to randomize the order.
	rand.Shuffle(len(possibleNums), func(i, j int) {
		possibleNums[i], possibleNums[j] = possibleNums[j], possibleNums[i]
	})

	// Select the first gptNum numbers from the shuffled slice.
	entryNums := make([]int, global.GPTNum)
	copy(entryNums, possibleNums[:global.GPTNum])

	gs.mutex.Lock()
	defer gs.mutex.Unlock()
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

func (gs *GameState) SearchUserInUserSelections(idx int, user common.User) (models.UserSelection, bool) {
	for _, v := range gs.Info[idx].UserSelections {
		if v.User.UUID == user.UUID {
			return v, true
		}
	}
	return models.UserSelection{}, false
}

// GetUserSelections: Get user selections
func (gs *GameState) GetUserSelections() []models.UserSelection {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.Info[gs.Status.InfoIdx].UserSelections
}

// GetNowUserSelections
func (gs *GameState) GetNowUserSelections() []models.UserSelection {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.Info[gs.Status.InfoIdx].UserSelections
}

// SetUserSelections: Set user selections

func (gs *GameState) SetUserSelections(userSelections []models.UserSelection) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
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
}

func (gs *GameState) SetIfResetRound(userManager *UserManager) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	// infoIndex is not changed
	infoIdx := gs.Status.InfoIdx
	gs.Status.IsStarted = false
	gs.Status.IsOver = false
	gs.Info[infoIdx].TurnsLeft = global.GameTurnNum * len(gs.Info[infoIdx].PlayerList)

	gs.Info[infoIdx].NowUserIndex = 0
	gs.Info[infoIdx].UserSelections = make([]models.UserSelection, 0)

	gs.Info[infoIdx].Messages = make([]models.Message, 0)
	gs.Info[infoIdx].PlayerList = userManager.GetSortedPlayers()
}

func (gs *GameState) SetIfGameTotallyReset(userManager *UserManager) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.Status.IsStarted = false
	gs.Status.IsOver = false
	gs.Status.InfoIdx = 0
	gs.Info = make([]models.Game, 0)

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
	return gs.Info[gs.Status.InfoIdx]
}

// CheckIsRoundOver: Check if round is over
func (gs *GameState) CheckIsRoundOver() bool {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.Info[gs.Status.InfoIdx].TurnsLeft == 0
}

// CheckAllUserVote: Check if all user vote
func (gs *GameState) CheckAllUserVote() bool {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return len(gs.Info[gs.Status.InfoIdx].UserSelections) == len(gs.Info[gs.Status.InfoIdx].PlayerList)
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
	if len(gs.Info) == 0 {
		return len(userManager.GetPlayers()) == gs.Status.MaxPlayer
	}
	return len(userManager.GetPlayers()) == gs.Info[gs.Status.InfoIdx].MaxPlayer
}
