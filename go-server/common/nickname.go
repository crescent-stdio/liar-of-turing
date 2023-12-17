package common

type Nickname struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	IsUsed   bool   `json:"IsUsed"`
}

var Nicknames []Nickname

func GetNicknames() []Nickname {
	return Nicknames
}

func SetNicknames(nicknames []Nickname) {
	Nicknames = nicknames
}
