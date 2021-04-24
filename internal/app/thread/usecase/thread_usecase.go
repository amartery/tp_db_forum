package usecase

import (
	"errors"
	"strconv"

	postModels "github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/app/thread"
	"github.com/amartery/tp_db_forum/internal/app/thread/models"
	"github.com/amartery/tp_db_forum/internal/app/user"
	"github.com/jackc/pgconn"
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

func (t *ThreadUsecase) FindThreadBySlug(slug string) (*models.Thread, error) {
	thread, err := t.repThread.FindThreadBySlug(slug)
	return thread, err
}

func (t *ThreadUsecase) CreateThread(thread *models.Thread) (*models.Thread, error) {
	thread, err := t.repThread.CreateThread(thread)
	return thread, err
}

func (t *ThreadUsecase) GetThreadsByForumSlug(slug, since, desc string, limit int) ([]*models.Thread, error) {
	threads, err := t.repThread.GetThreadsByForumSlug(slug, since, desc, limit)
	return threads, err
}

func (t *ThreadUsecase) CreatePost(posts []*postModels.Post, slugOrInt string) ([]*postModels.Post, error) {
	threadID, err := strconv.Atoi(slugOrInt)
	thread := &models.Thread{}
	if err != nil {
		thread, err = t.repThread.FindThreadBySlug(slugOrInt)
		if err != nil {
			return nil, err
		}
	} else {
		thread, err = t.repThread.FindThreadByID(threadID)
		if err != nil {
			return nil, err
		}
	}
	if len(posts) == 0 {
		return posts, nil
	}
	if len(posts) != 0 {
		_, err = t.repUser.GetUserByNickname(posts[0].Author)
		if err != nil {
			return nil, err
		}
	}
	for _, post := range posts {
		post.Thread = thread.ID
		if post.Parent != 0 {
			parentThreadID, err := t.repThread.CheckThreadID(post.Parent)
			if err != nil {
				return nil, errors.New("some_err")
			}
			if parentThreadID != post.Thread {
				return nil, errors.New("some_err")
			}
		}
		post.Forum = thread.Forum
	}
	posts, err = t.repThread.CreatePost(posts)
	return posts, err
}

func (t *ThreadUsecase) GetThreadBySLUGorID(slug_or_id string) (*models.Thread, error) {
	var thread *models.Thread
	threadID, err := strconv.Atoi(slug_or_id)
	if err != nil {
		thread, err = t.repThread.FindThreadBySlug(slug_or_id)
		if err != nil {
			return nil, err
		}
	} else {
		thread, err = t.repThread.FindThreadByID(threadID)
		if err != nil {
			return nil, err
		}
	}
	return thread, err
}

func (t *ThreadUsecase) UpdateTreads(slug_or_id string, th *models.Thread) (*models.Thread, error) {
	threadID, err := strconv.Atoi(slug_or_id)
	if err != nil {
		th.Slug = slug_or_id
		oldThread, errRep := t.repThread.FindThreadBySlug(slug_or_id)
		if errRep != nil {
			return nil, errRep
		}
		if th.Title == "" {
			th.Title = oldThread.Title
		}
		if th.Message == "" {
			th.Message = oldThread.Message
		}
		th, errRep = t.repThread.UpdateThreadBySlug(th)
		if errRep != nil {
			return nil, errRep
		}
	} else {
		th.ID = threadID
		oldThread, errRep := t.repThread.FindThreadByID(threadID)
		if errRep != nil {
			return nil, errRep
		}
		if th.Title == "" {
			th.Title = oldThread.Title
		}
		if th.Message == "" {
			th.Message = oldThread.Message
		}
		th, errRep = t.repThread.UpdateThreadByID(th)
		if errRep != nil {
			return nil, errRep
		}
	}
	return th, nil
}

func (t *ThreadUsecase) GetPosts(sort, since, slug_or_id string, limit int, desc bool) ([]*postModels.Post, error) {
	threadID, err := strconv.Atoi(slug_or_id)
	thread := &models.Thread{}
	if err != nil {
		thread, err = t.repThread.FindThreadBySlug(slug_or_id)
		if err != nil {
			return nil, err
		}
	} else {
		thread, err = t.repThread.FindThreadByID(threadID)
		if err != nil {
			return nil, err
		}
	}
	posts, err := t.repThread.GetPosts(limit, thread.ID, sort, since, desc)
	return posts, err
}

func (t *ThreadUsecase) CreateNewVote(vote *models.Vote, slug_or_id string) (*models.Thread, error) {
	threadID, err := strconv.Atoi(slug_or_id)
	var th *models.Thread
	if err != nil {
		th, err = t.repThread.FindThreadBySlug(slug_or_id)
		if err != nil {
			return nil, err
		}
	} else {
		th, err = t.repThread.FindThreadByID(threadID)
		if err != nil {
			return nil, err
		}
	}
	vote.ThreadID = th.ID
	err = t.repThread.CreateNewVote(vote)
	if err != nil {
		if err.(*pgconn.PgError).Code == "23503" {
			return nil, err
		}
		linesUpdated, err := t.repThread.UpdateVote(vote)
		if err != nil {
			return nil, err
		}
		if linesUpdated != 0 {
			th.Votes += 2 * vote.Voice
		}
		return th, err
	}
	th.Votes += vote.Voice
	return th, err
}
