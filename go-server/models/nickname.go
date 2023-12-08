package models

type Nickname struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	IsUsed   bool   `json:"IsUsed"`
}
