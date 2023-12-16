package utils

import (
	"encoding/json"
	"liar-of-turing/common"
	"os"
)

// LoadNicknames: load nicknames from json file
func LoadNicknames() error {
	file, err := os.ReadFile("data/nicknames.json")
	if err != nil {
		return err
	}

	var nicknames []common.Nickname
	err = json.Unmarshal(file, &nicknames)
	if err != nil {
		return err
	}

	common.SetNicknames(nicknames)

	// global.GlobalNicknames = []models.Nickname{{Nickname: "test", IsUsed: false}}
	return nil
}
