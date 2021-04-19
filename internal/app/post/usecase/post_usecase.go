package usecase

import "github.com/amartery/tp_db_forum/internal/app/post"

type PostUsecase struct {
	repo post.Repository
}

func NewPostUsecase(PostRepo post.Repository) *PostUsecase {
	return &PostUsecase{
		repo: PostRepo,
	}
}
