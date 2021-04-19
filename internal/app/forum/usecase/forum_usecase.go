package usecase

import (
	"github.com/amartery/tp_db_forum/internal/app/forum"
	"github.com/amartery/tp_db_forum/internal/app/forum/models"
)

type ForumUsecase struct {
	repo forum.Repository
}

func NewForumUsecase(forumRepo forum.Repository) *ForumUsecase {
	return &ForumUsecase{
		repo: forumRepo,
	}
}

func (usecase *ForumUsecase) CreateForum(forum *models.Forum) error {
	err := usecase.repo.CreateForum(forum)
	return err
}

func (usecase *ForumUsecase) GetForumBySlug(slug string) (*models.Forum, error) {
	forum, err := usecase.repo.GetForumBySlug(slug)
	return forum, err
}
