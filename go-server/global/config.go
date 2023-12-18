package global

var MaxPlayer = 5
var GameRoundNum = 1
var GameTurnNum = 2

var GPTNum = 2

func GetGameTurnNum() int {
	return GameTurnNum
}

func SetGameTurnNum(gameTurnNum int) {
	GameRoundNum = gameTurnNum
}

func GetGameRoundNum() int {
	return GameRoundNum
}

func SetGameRoundNum(gameRoundNum int) {
	GameRoundNum = gameRoundNum
}

func GetMaxPlayer() int {
	return MaxPlayer
}

func SetMaxPlayer(maxPlayer int) {
	MaxPlayer = maxPlayer
}
