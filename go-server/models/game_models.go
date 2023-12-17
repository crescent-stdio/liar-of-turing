package models

import "liar-of-turing/common"

type Game struct {
	NowUserIndex int `json:"now_user_index"`
	MaxPlayer    int `json:"max_player"`
	// OnlineUserList []common.User   `json:"online_user_list"`
	PlayerList     []common.User   `json:"player_list"`
	Round          int             `json:"round"`
	TurnsLeft      int             `json:"turns_left"`
	UserSelections []UserSelection `json:"user_round_selections"`
	Messages       []Message       `json:"messages"`
}

type UserSelection struct {
	User      common.User `json:"user"`
	Timestamp int64       `json:"timestamp"`
	Selection string      `json:"selection"`
	Reason    string      `json:"reason"`
}
