package http

import (
	"encoding/json"
	"fmt"

	postModels "github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/app/thread"
	"github.com/amartery/tp_db_forum/internal/app/thread/models"
	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type ThreadHandler struct {
	usecaseThread thread.Usecase
	logger        *logrus.Logger
}

func NewThreadHandler(threadUsecase thread.Usecase) *ThreadHandler {
	return &ThreadHandler{
		usecaseThread: threadUsecase,
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
	slug_or_id, ok := ctx.UserValue("slug_or_id").(string)
	if !ok {
		utils.SendServerError("Can't get slug_or_id", fasthttp.StatusInternalServerError, ctx)
		return
	}

	posts := make([]*postModels.Post, 0)
	err := json.Unmarshal(ctx.PostBody(), &posts)
	if err != nil {
		utils.SendServerError(err.Error(), fasthttp.StatusInternalServerError, ctx)
		return
	}
	posts, err = handler.usecaseThread.CreatePost(posts, slug_or_id)
	if err != nil {
		var code int
		if err.Error() == "some_err" {
			code = fasthttp.StatusConflict
		} else {
			code = fasthttp.StatusNotFound
		}
		msg := "Can't find post author by nickname: "
		utils.SendServerError(msg, code, ctx)
		return
	}
	if posts == nil {
		return
	}
	utils.SendResponse(fasthttp.StatusCreated, posts, ctx)
}

func (handler *ThreadHandler) BranchDetailsGet(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting BranchDetailsGet")
	slug_or_id, ok := ctx.UserValue("slug_or_id").(string)
	if !ok {
		utils.SendServerError("Can't get slug_or_id", fasthttp.StatusInternalServerError, ctx)
		return
	}
	thread, err := handler.usecaseThread.GetThreadBySLUGorID(slug_or_id)
	if err != nil {
		msg := fmt.Sprintf("Can't find thread with slug %s", slug_or_id)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}
	utils.SendResponse(fasthttp.StatusOK, thread, ctx)
}

func (handler *ThreadHandler) BranchDetailsUpdate(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting BranchDetailsUpdate")
	slug_or_id, ok := ctx.UserValue("slug_or_id").(string)
	if !ok {
		utils.SendServerError("some err", fasthttp.StatusInternalServerError, ctx)
		return
	}
	updateReq := &models.UpdateRequest{}
	err := updateReq.UnmarshalJSON(ctx.Request.Body())
	if err != nil {
		utils.SendServerError(err.Error(), fasthttp.StatusInternalServerError, ctx)
		return
	}
	thread := &models.Thread{
		Title:   updateReq.Title,
		Message: updateReq.Message,
	}
	thread, err = handler.usecaseThread.UpdateTreads(slug_or_id, thread)
	if err != nil {
		msg := fmt.Sprintf("Can't find thread with slug %s", slug_or_id)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}
	utils.SendResponse(fasthttp.StatusOK, thread, ctx)
}

func (handler *ThreadHandler) CurrentBranchPosts(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting CurrentBranchPosts")
	slug_or_id, ok := ctx.UserValue("slug_or_id").(string)
	if !ok {
		utils.SendServerError("Can't get slug_or_id", fasthttp.StatusInternalServerError, ctx)
		return
	}
	since := string(ctx.QueryArgs().Peek("since"))
	limit, err := ctx.QueryArgs().GetUint("limit")
	if err != nil {
		utils.SendServerError("Can't get limit", fasthttp.StatusInternalServerError, ctx)
		return
	}
	sort := string(ctx.QueryArgs().Peek("sort"))
	desc := ctx.QueryArgs().GetBool("desc")

	posts, err := handler.usecaseThread.GetPosts(sort, since, slug_or_id, limit, desc)
	if err != nil {
		msg := fmt.Sprintf("Can't threads with forum slug %s", slug_or_id)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}
	utils.SendResponse(fasthttp.StatusOK, posts, ctx)
}

func (handler *ThreadHandler) VoteForBranch(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting VoteForBranch")
	slug_or_id, ok := ctx.UserValue("slug_or_id").(string)
	if !ok {
		utils.SendServerError("some err", fasthttp.StatusInternalServerError, ctx)
		return
	}
	vote := &models.Vote{}
	err := vote.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		utils.SendServerError(err.Error(), fasthttp.StatusInternalServerError, ctx)
		return
	}

	thread, err := handler.usecaseThread.CreateNewVote(vote, slug_or_id)
	if err != nil {
		msg := fmt.Sprintf("Can't find thread with slug %s", slug_or_id)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}

	utils.SendResponse(fasthttp.StatusOK, thread, ctx)
}
