package models

type Forum struct {
	Tittle   string `json:"title"`
	Nickname string `json:"user"`
	Slug     string `json:"slug"`
	Posts    int    `json:"posts"`
	Threads  int    `json:"threads"`
}
