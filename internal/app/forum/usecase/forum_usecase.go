package usecase

import (
	"github.com/amartery/tp_db_forum/internal/app/forum"
	"github.com/amartery/tp_db_forum/internal/app/forum/models"
	userModel "github.com/amartery/tp_db_forum/internal/app/user/models"
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
	return usecase.repo.GetForumBySlug(slug)
}

func (usecase *ForumUsecase) GetUsersByForum(slug string, limit int, since string, desc string) (*[]userModel.User, error) {
	switch desc {
	case "true":
		desc = "DESC"
	case "false":
		desc = "ASC"
	}
	return usecase.repo.GetUsersByForum(slug, limit, since, desc)
}

func (usecase *ForumUsecase) CheckForum(slug string) (string, error) {
	return usecase.repo.CheckForum(slug)
}
