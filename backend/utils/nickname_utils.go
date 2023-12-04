package utils

import (
	"encoding/json"
	"liarOfTuring/global"
	"liarOfTuring/models"
	"os"
)

func LoadNicknames() error {
	file, err := os.ReadFile("data/nicknames.json")
	if err != nil {
		return err
	}

	var nicknames []models.Nickname
	err = json.Unmarshal(file, &nicknames)
	if err != nil {
		return err
	}

	// 전역 변수에 데이터 저장
	// global.GlobalNicknames = nicknames
	global.SetGlobalNicknames(nicknames)

	// global.GlobalNicknames = []models.Nickname{{Nickname: "test", IsUsed: false}}
	return nil
}
