package post

import (
	"github.com/amartery/tp_db_forum/internal/app/post/models"
)

type Usecase interface {
	UpdatePost(post *models.Post) (*models.Post, error)
	GetPost(id string) (*models.Post, error)
}
