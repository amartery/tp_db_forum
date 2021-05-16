package thread

import (
	postModel "github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/app/thread/models"
)

type Repository interface {
	GetThreadBySlug(slug string) (*models.Thread, error)
	GetThreadByID(id int) (*models.Thread, error)
	GetThreadIDAndForum(slugOrID string) (*models.Thread, error)
	CreatePosts(thread models.Thread, posts []postModel.Post) error
	CheckThreadBySlug(slug string) (int, error)
	CheckThreadByID(id int) (int, error)
	UpdateThread(thread *models.Thread) (*models.Thread, error)
	GetPosts(slugOrID string, limit int, order string, since string) ([]postModel.Post, error)
	GetPostsTree(slugOrID string, limit int, order string, since string) ([]postModel.Post, error)
	GetPostsParentTree(slugOrID string, limit int, order string, since string) ([]postModel.Post, error)
	Vote(vote *models.Vote) (*models.Thread, error)
	CreateThread(thread *models.Thread) error
	GetThreadsByForumSlug(slug, limit, since, desc string) (*[]models.Thread, error)
	CheckForum(slug string) (string, error)
}
