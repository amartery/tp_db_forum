package models

import (
	modelsForum "github.com/amartery/tp_db_forum/internal/app/forum/models"
	modelsThread "github.com/amartery/tp_db_forum/internal/app/thread/models"
	modelsUser "github.com/amartery/tp_db_forum/internal/app/user/models"
)

type Post struct {
	ID       int    `json:"id"`
	Parent   int    `json:"parent"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited"`
	Forum    string `json:"forum"`
	Thread   int    `json:"thread"`
	Created  string `json:"created"`
}

//easyjson:json
type PostResponse struct {
	Post   *Post                `json:"post"`
	User   *modelsUser.User     `json:"author,omitempty"`
	Forum  *modelsForum.Forum   `json:"forum,omitempty"`
	Thread *modelsThread.Thread `json:"thread,omitempty"`
}
