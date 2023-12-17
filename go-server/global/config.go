package global

var MaxPlayer = 4
var GameTurnNum = 1
var GameRoundNum = 1

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
