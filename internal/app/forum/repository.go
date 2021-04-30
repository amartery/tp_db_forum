package forum

import (
	"fmt"

	"github.com/amartery/tp_db_forum/internal/app/forum/models"
	userModel "github.com/amartery/tp_db_forum/internal/app/user/models"
)

var (
	ErrForumDoesntExists = fmt.Errorf("forum not exists")
	ErrDataConflict      = fmt.Errorf("data conflict")
	ErrWrongParent       = fmt.Errorf("wrong parent")
)

type Repository interface {
	CreateForum(forum *models.Forum) error
	GetForumBySlug(slug string) (*models.Forum, error)
	GetUsersByForum(slug string, limit int, since string, desc string) (*[]userModel.User, error)
	CheckForum(slug string) (string, error)
}
