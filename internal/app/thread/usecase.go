package thread

import (
	postModels "github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/app/thread/models"
)

type Usecase interface {
	FindThreadBySlug(slug string) (*models.Thread, error)
	CreateThread(thread *models.Thread) (*models.Thread, error)
	GetThreadsByForumSlug(slug, since, desc string, limit int) ([]*models.Thread, error)
	CreatePost(posts []*postModels.Post, slugOrInt string) ([]*postModels.Post, error)
	GetThreadBySLUGorID(slug_or_id string) (*models.Thread, error)
	UpdateTreads(slug_or_id string, th *models.Thread) (*models.Thread, error)
	GetPosts(sort, since, slug_or_id string, limit int, desc bool) ([]*postModels.Post, error)
	CreateNewVote(vote *models.Vote, slug_or_id string) (*models.Thread, error)
}
