package forum

import "github.com/amartery/tp_db_forum/internal/app/forum/models"

type Usecase interface {
	CreateForum(forum *models.Forum) error
}
