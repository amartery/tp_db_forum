package thread

import (
	postModels "github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/app/thread/models"
)

type Repository interface {
	FindThreadBySlug(slug string) (*models.Thread, error)
	FindThreadByID(threadID int) (*models.Thread, error)
	CreateThread(thread *models.Thread) (*models.Thread, error)
	GetThreadsByForumSlug(slug, since, desc string, limit int) ([]*models.Thread, error)
	CheckThreadID(parentID int) (int, error)
	CreatePost(posts []*postModels.Post) ([]*postModels.Post, error)
	UpdateThreadByID(thread *models.Thread) (*models.Thread, error)
	UpdateThreadBySlug(thread *models.Thread) (*models.Thread, error)
	GetPosts(limit, threadID int, sort, since string, desc bool) ([]*postModels.Post, error)
	CreateNewVote(vote *models.Vote) error
	UpdateVote(vote *models.Vote) (int, error)
}
