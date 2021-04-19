package forum

import "github.com/amartery/tp_db_forum/internal/app/forum/models"

type Repository interface {
	CreateForum(forum *models.Forum) error
	GetForumBySlug(slug string) (*models.Forum, error)
}
