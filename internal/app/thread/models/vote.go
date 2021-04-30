package models

//easyjson:json
type Vote struct {
	UserID   int    `json:"-"`
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
	ID       int    `json:"id"`
	Slug     string `json:"slug"`
}
