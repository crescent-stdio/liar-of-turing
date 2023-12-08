package global

import (
	"liarOfTuring/models"
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
