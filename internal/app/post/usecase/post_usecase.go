package usecase

import (
	"github.com/amartery/tp_db_forum/internal/app/forum"
	"github.com/amartery/tp_db_forum/internal/app/post"
	"github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/app/thread"
	"github.com/amartery/tp_db_forum/internal/app/user"
	"github.com/amartery/tp_db_forum/internal/pkg/utils"
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

func (pUse *PostUsecase) UpdatePost(post *models.Post) (*models.Post, error) {
	oldPost, err := pUse.repoPost.GetPost(post.ID)
	if err != nil {
		return nil, err
	}
	if oldPost.Message == post.Message {
		return oldPost, nil
	}
	post, err = pUse.repoPost.UpdatePostByID(post)
	return post, err
}

func (pUse *PostUsecase) GetPost(postID int, relatedStrs []string) (*models.PostResponse, error) {
	post, err := pUse.repoPost.GetPost(postID)
	if err != nil {
		return nil, err
	}
	postResponse := &models.PostResponse{Post: post}
	if utils.Find(relatedStrs, "user") {
		user, err := pUse.repoUser.GetUserByNickname(post.Author)
		if err != nil {
			return nil, err
		}
		postResponse.User = user
	}
	if utils.Find(relatedStrs, "forum") {
		forum, err := pUse.repoForum.GetForumBySlug(post.Forum)
		if err != nil {
			return nil, err
		}
		postResponse.Forum = forum
	}
	if utils.Find(relatedStrs, "thread") {
		thread, err := pUse.repoThread.FindThreadByID(post.Thread)
		if err != nil {
			return nil, err
		}
		postResponse.Thread = thread
	}
	return postResponse, err
}
