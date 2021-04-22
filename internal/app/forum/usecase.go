package forum

import (
	"github.com/amartery/tp_db_forum/internal/app/forum/models"
	usersModels "github.com/amartery/tp_db_forum/internal/app/user/models"
)

type Usecase interface {
	CreateForum(forum *models.Forum) error
	GetForumBySlug(slug string) (*models.Forum, error)
	GetUsersByForum(slug, since string, limit int, desc bool) ([]*usersModels.User, error)
}
