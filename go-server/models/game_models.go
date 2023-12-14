package models

type Game struct {
	NowUserIndex   int             `json:"now_user_index"`
	MaxPlayer      int             `json:"max_player"`
	OnlineUserList []User          `json:"online_user_list"`
	PlayerList     []User          `json:"player_list"`
	TurnsLeft      int             `json:"turns_left"`
	UserSelections []UserSelection `json:"user_round_selections"`
	Messages       []Message       `json:"messages"`
}

type UserSelection struct {
	User      User   `json:"user"`
	Selection string `json:"selection"`
	Reason    string `json:"reason"`
}
