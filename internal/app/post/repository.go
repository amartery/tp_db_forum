package post

import (
	"github.com/amartery/tp_db_forum/internal/app/post/models"
)

type Repository interface {
	UpdatePostByID(post *models.Post) (*models.Post, error)
	GetPost(postID int) (*models.Post, error)
}
