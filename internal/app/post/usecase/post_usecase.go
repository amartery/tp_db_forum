package usecase

import (
	"github.com/amartery/tp_db_forum/internal/app/forum"
	"github.com/amartery/tp_db_forum/internal/app/post"
	"github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/app/thread"
	"github.com/amartery/tp_db_forum/internal/app/user"
)

type PostUsecase struct {
	repoPost   post.Repository
	repoUser   user.Repository
	repoForum  forum.Repository
	repoThread thread.Repository
}

func NewPostUsecase(p post.Repository, u user.Repository, f forum.Repository, t thread.Repository) *PostUsecase {
	return &PostUsecase{
		repoPost:   p,
		repoUser:   u,
		repoForum:  f,
		repoThread: t,
	}
}

func (pUse *PostUsecase) GetPost(id string) (*models.Post, error) {
	return pUse.repoPost.GetPost(id)
}

func (pUse *PostUsecase) UpdatePost(post *models.Post) (*models.Post, error) {
	return pUse.repoPost.UpdatePost(post)
}
