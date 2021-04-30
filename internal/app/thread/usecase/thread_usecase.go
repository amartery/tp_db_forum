package usecase

import (
	"strconv"

	postModel "github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/app/thread"
	"github.com/amartery/tp_db_forum/internal/app/thread/models"
	"github.com/amartery/tp_db_forum/internal/app/user"
)

type ThreadUsecase struct {
	repThread thread.Repository
	repUser   user.Repository
}

func NewThreadUsecase(ThreadRepo thread.Repository, UserRepo user.Repository) *ThreadUsecase {
	return &ThreadUsecase{
		repThread: ThreadRepo,
		repUser:   UserRepo,
	}
}

func (t *ThreadUsecase) CreatePosts(thread models.Thread, posts []postModel.Post) error {
	return t.repThread.CreatePosts(thread, posts)
}

func (t *ThreadUsecase) GetPosts(slugOrID string, limit int, sort string, order string, since string) ([]postModel.Post, error) {
	switch order {
	case "true":
		order = "DESC"
	case "false":
		order = "ASC"
	}

	switch sort {
	case "flat":
		return t.repThread.GetPosts(slugOrID, limit, order, since)
	case "tree":
		return t.repThread.GetPostsTree(slugOrID, limit, order, since)
	case "parent_tree":
		return t.repThread.GetPostsParentTree(slugOrID, limit, order, since)
	default:
		return t.repThread.GetPosts(slugOrID, limit, order, since)
	}
}

func (t *ThreadUsecase) Vote(vote *models.Vote) (*models.Thread, error) {
	return t.repThread.Vote(vote)
}

func (t *ThreadUsecase) GetThread(slugOrID string) (*models.Thread, error) {
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		return t.repThread.GetThreadBySlug(slugOrID)
	}
	return t.repThread.GetThreadByID(id)
}

func (t *ThreadUsecase) GetThreadIDAndForum(slugOrID string) (*models.Thread, error) {
	return t.repThread.GetThreadIDAndForum(slugOrID)
}

func (t *ThreadUsecase) CheckThread(slugOrID string) error {
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		_, err = t.repThread.CheckThreadBySlug(slugOrID)
		return err
	}

	_, err = t.repThread.CheckThreadByID(id)
	return err
}

func (t *ThreadUsecase) UpdateThread(slugOrID string, thread *models.Thread) (*models.Thread, error) {
	thread.Slug = &slugOrID
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		id = 0
	}

	thread.ID = id
	return t.repThread.UpdateThread(thread)
}

func (t *ThreadUsecase) CreateThread(thread *models.Thread) error {
	return t.repThread.CreateThread(thread)
}

func (t *ThreadUsecase) GetThreadsByForumSlug(slug, limit, since, desc string) (*[]models.Thread, error) {
	return t.repThread.GetThreadsByForumSlug(slug, limit, since, desc)
}

func (t *ThreadUsecase) CheckForum(slug string) (string, error) {
	return t.repThread.CheckForum(slug)
}
