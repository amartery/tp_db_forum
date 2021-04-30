package thread

import (
	postModel "github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/app/thread/models"
)

type Usecase interface {
	GetThread(slugOrID string) (*models.Thread, error)
	GetThreadIDAndForum(slugOrID string) (*models.Thread, error)
	CreatePosts(thread models.Thread, posts []postModel.Post) error

	UpdateThread(slugOrID string, thread *models.Thread) (*models.Thread, error)
	GetPosts(slugOrID string, limit int, sort string, order string, since string) ([]postModel.Post, error)
	Vote(vote *models.Vote) (*models.Thread, error)

	/////
	CreateThread(thread *models.Thread) error
	GetThreadsByForumSlug(slug, limit, since, desc string) (*[]models.Thread, error)

	//!!
	CheckThread(slugOrID string) error
	CheckForum(slug string) (string, error)
}
