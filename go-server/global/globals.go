package global

import (
	"liar-of-turing/models"
)

type Global struct {
	GlobalNicknames []models.Nickname
}

var global = &Global{}

func GetGlobalNicknames() []models.Nickname {
	return global.GlobalNicknames
}

func SetGlobalNicknames(nicknames []models.Nickname) {
	global.GlobalNicknames = nicknames
}
