package thread

import "github.com/amartery/tp_db_forum/internal/app/thread/models"

type Usecase interface {
	FindThreadBySlug(slug string) (*models.Thread, error)
	CreateThread(thread *models.Thread) (*models.Thread, error)
	GetThreadsByForumSlug(slug, since, desc string, limit int) ([]*models.Thread, error)
}
