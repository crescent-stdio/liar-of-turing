package models

// GameState represents the state of the game.
type GameState struct {
	IsStarted bool
	IsOver    bool
}

func NewGameState() *GameState {
	return &GameState{
		IsStarted: false,
		IsOver:    false,
	}
}

func (gs *GameState) StartGame() {
	gs.IsStarted = true
	gs.IsOver = false
}

func (gs *GameState) EndGame() {
	gs.IsOver = true
}

func (gs *GameState) ResetGame() {
	gs.IsStarted = false
	gs.IsOver = false
}
