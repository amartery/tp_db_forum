package models

import (
	modelsForum "github.com/amartery/tp_db_forum/internal/app/forum/models"
	modelsThread "github.com/amartery/tp_db_forum/internal/app/thread/models"
	modelsUser "github.com/amartery/tp_db_forum/internal/app/user/models"
	"github.com/go-openapi/strfmt"
)

//easyjson:json
type Post struct {
	ID       int             `json:"id"`
	Author   string          `json:"author"`
	Message  string          `json:"message"`
	Parent   int             `json:"parent,omitempty"`
	Forum    string          `json:"forum"`
	Thread   int             `json:"thread"`
	Created  strfmt.DateTime `json:"created,omitempty"`
	IsEdited bool            `json:"isEdited"`
}

//easyjson:json
type PostResponse struct {
	Post   *Post                `json:"post"`
	Author *modelsUser.User     `json:"author,omitempty"`
	Thread *modelsThread.Thread `json:"thread,omitempty"`
	Forum  *modelsForum.Forum   `json:"forum,omitempty"`
}
