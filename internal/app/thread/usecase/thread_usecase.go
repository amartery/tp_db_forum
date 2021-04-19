package usecase

import "github.com/amartery/tp_db_forum/internal/app/thread"

type ThreadUsecase struct {
	repo thread.Repository
}

func NewThreadUsecase(ThreadRepo thread.Repository) *ThreadUsecase {
	return &ThreadUsecase{
		repo: ThreadRepo,
	}
}
