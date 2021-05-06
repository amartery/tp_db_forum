package http

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/amartery/tp_db_forum/internal/app/forum"
	postModel "github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/app/thread"
	"github.com/amartery/tp_db_forum/internal/app/thread/models"
	"github.com/amartery/tp_db_forum/internal/app/user"
	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type ThreadHandler struct {
	usecaseThread thread.Usecase
	usecaseUser   user.Usecase
	logger        *logrus.Logger
}

func NewThreadHandler(t thread.Usecase, u user.Usecase) *ThreadHandler {
	return &ThreadHandler{
		usecaseThread: t,
		usecaseUser:   u,
		logger:        logrus.New(),
	}
}

// TODO: s.router.POST("​/thread​/{slug_or_id}​/create", s.CreatePostInBranch)
// TODO: s.router.GET("/thread​/{slug_or_id}​/details", s.BranchDetailsGet)
// TODO: s.router.POST("/thread​/{slug_or_id}​/details", s.BranchDetailsUpdate)
// TODO: s.router.GET("/thread/{slug_or_id}/posts", s.CurrentBranchPosts)
// TODO: s.router.POST("/thread​/{slug_or_id}​/vote", s.VoteForBranch)

func (handler *ThreadHandler) CreatePostInBranch(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting CreatePostInBranch")
	slugOrID := ctx.UserValue("slug_or_id").(string)

	thread, err := handler.usecaseThread.GetThreadIDAndForum(slugOrID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find post thread by id: %v", slugOrID),
		}

		_ = json.NewEncoder(ctx).Encode(msg)
		return
	}

	posts := make([]postModel.Post, 0)
	err = json.Unmarshal(ctx.PostBody(), &posts)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if len(posts) == 0 {
		ctx.SetStatusCode(fasthttp.StatusCreated)
		_ = json.NewEncoder(ctx).Encode(posts)
		return
	}

	for _, post := range posts {
		_, err = handler.usecaseUser.CheckIfUserExists(post.Author)
		if err != nil {
			msg := utils.Message{
				Text: fmt.Sprintf("Can't find post author by nickname: %v", post.Author),
			}

			ctx.SetStatusCode(fasthttp.StatusNotFound)
			err = json.NewEncoder(ctx).Encode(msg)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
			return
		}
	}

	err = handler.usecaseThread.CreatePosts(*thread, posts)
	if err != nil {
		if err == forum.ErrWrongParent {
			ctx.SetStatusCode(fasthttp.StatusConflict)
			msg := utils.Message{
				Text: "Parent post was created in another thread",
			}

			_ = json.NewEncoder(ctx).Encode(msg)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusConflict)
		msg := utils.Message{
			Text: "Parent post was created in another thread",
		}

		_ = json.NewEncoder(ctx).Encode(msg)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	_ = json.NewEncoder(ctx).Encode(posts)

}

func (handler *ThreadHandler) BranchDetailsGet(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting BranchDetailsGet")
	slugOrID := ctx.UserValue("slug_or_id").(string)

	threads, err := handler.usecaseThread.GetThread(slugOrID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find thread by slug: %v", slugOrID),
		}

		_ = json.NewEncoder(ctx).Encode(msg)
		return
	}

	err = json.NewEncoder(ctx).Encode(threads)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}

func (handler *ThreadHandler) BranchDetailsUpdate(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting BranchDetailsUpdate")
	slugOrID := ctx.UserValue("slug_or_id").(string)

	err := handler.usecaseThread.CheckThread(slugOrID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find thread by slug: %v", slugOrID),
		}

		_ = json.NewEncoder(ctx).Encode(msg)
		return
	}

	thread := &models.Thread{}
	err = json.Unmarshal(ctx.PostBody(), &thread)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if thread.Title == "" && thread.Message == "" {
		thread, err = handler.usecaseThread.GetThread(slugOrID)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	} else {
		thread, err = handler.usecaseThread.UpdateThread(slugOrID, thread)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
	}

	err = json.NewEncoder(ctx).Encode(thread)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}

func (handler *ThreadHandler) CurrentBranchPosts(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting CurrentBranchPosts")
	slugOrID := ctx.UserValue("slug_or_id").(string)

	err := handler.usecaseThread.CheckThread(slugOrID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find thread by slug: %v", slugOrID),
		}

		_ = json.NewEncoder(ctx).Encode(msg)
		return
	}

	limitParam := string(ctx.URI().QueryArgs().Peek("limit"))
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	sortParam := string(ctx.URI().QueryArgs().Peek("sort"))
	descParam := string(ctx.URI().QueryArgs().Peek("desc"))
	if descParam == "" {
		descParam = "false"
	}

	sinceParam := string(ctx.URI().QueryArgs().Peek("since"))

	posts, err := handler.usecaseThread.GetPosts(slugOrID, limit, sortParam, descParam, sinceParam)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(ctx).Encode(posts)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}

func (handler *ThreadHandler) VoteForBranch(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting VoteForBranch")
	slugOrID := ctx.UserValue("slug_or_id").(string)

	vote := &models.Vote{}
	err := json.Unmarshal(ctx.PostBody(), &vote)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	vote.Slug = slugOrID
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		id = 0
	}
	vote.ID = id

	thread, err := handler.usecaseThread.Vote(vote)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find thread by slug: %v", slugOrID),
		}

		_ = json.NewEncoder(ctx).Encode(msg)
		return
	}

	err = json.NewEncoder(ctx).Encode(thread)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}
