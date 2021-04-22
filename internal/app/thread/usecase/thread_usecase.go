package usecase

import (
	"github.com/amartery/tp_db_forum/internal/app/thread"
	"github.com/amartery/tp_db_forum/internal/app/thread/models"
)

type ThreadUsecase struct {
	repo thread.Repository
}

func NewThreadUsecase(ThreadRepo thread.Repository) *ThreadUsecase {
	return &ThreadUsecase{
		repo: ThreadRepo,
	}
}

func (t *ThreadUsecase) FindThreadBySlug(slug string) (*models.Thread, error) {
	thread, err := t.repo.FindThreadBySlug(slug)
	return thread, err
}

func (t *ThreadUsecase) CreateThread(thread *models.Thread) (*models.Thread, error) {
	thread, err := t.repo.CreateThread(thread)
	return thread, err
}

func (t *ThreadUsecase) GetThreadsByForumSlug(slug, since, desc string, limit int) ([]*models.Thread, error) {
	threads, err := t.repo.GetThreadsByForumSlug(slug, since, desc, limit)
	return threads, err
}
