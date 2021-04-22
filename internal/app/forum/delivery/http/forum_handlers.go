package http

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/amartery/tp_db_forum/internal/app/forum"
	"github.com/amartery/tp_db_forum/internal/app/forum/models"
	"github.com/amartery/tp_db_forum/internal/app/thread"
	threadModels "github.com/amartery/tp_db_forum/internal/app/thread/models"

	// usersModels "github.com/amartery/tp_db_forum/internal/app/user/models"
	"github.com/amartery/tp_db_forum/internal/app/user"
	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type ForumHandler struct {
	usecaseForum  forum.Usecase
	usecaseUser   user.Usecase
	usecaseThread thread.Usecase
	logger        *logrus.Logger
}

func NewForumHandler(forumUsecase forum.Usecase, userUsecase user.Usecase, threadUsecase thread.Usecase) *ForumHandler {
	return &ForumHandler{
		usecaseForum:  forumUsecase,
		usecaseUser:   userUsecase,
		usecaseThread: threadUsecase,
		logger:        logrus.New(),
	}
}

// TODO: s.router.POST("/forum/create", DeliveryForum.ForumCreate)
// TODO: s.router.GET("/forum/{slug}/details", s.ForumDetails)
// TODO: s.router.POST("/forum/{slug}/create", s.ForumCreateBranch)
// TODO: s.router.GET("/forum/{slug}/users", s.CurrentForumUsers)
// TODO: s.router.GET("​/forum​/{slug}​/threads", s.ForumBranches)

func (f *ForumHandler) ForumCreate(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting ForumCreate")
	newForum := &models.Forum{}
	err := json.Unmarshal(ctx.PostBody(), newForum)
	if err != nil {
		utils.SendServerError(err.Error(), fasthttp.StatusInternalServerError, ctx)
		return
	}
	_, err = f.usecaseUser.GetUserByNickname(newForum.Nickname)
	if err != nil {
		msg := fmt.Sprintf("Can't find user with nickname %s", newForum.Nickname)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}
	err = f.usecaseForum.CreateForum(newForum)
	if err != nil {
		alredyExicted, err := f.usecaseForum.GetForumBySlug(newForum.Slug)
		if err != nil {
			msg := fmt.Sprintf("Can't find user with nickname %s", newForum.Nickname)
			utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
			return
		}
		utils.SendResponse(fasthttp.StatusConflict, alredyExicted, ctx)
		return
	}
	utils.SendResponse(fasthttp.StatusCreated, newForum, ctx)
}

func (f *ForumHandler) ForumDetails(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting ForumDetails")
	slug, ok := ctx.UserValue("slug").(string)
	if !ok {
		utils.SendServerError("Can't get value from slug", fasthttp.StatusInternalServerError, ctx)
		return
	}
	forum, err := f.usecaseForum.GetForumBySlug(slug)
	if err != nil {
		msg := fmt.Sprintf("Can't find forum with slug %s", slug)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}
	utils.SendResponse(fasthttp.StatusOK, forum, ctx)
}

func (f *ForumHandler) ForumCreateBranch(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting ForumCreateBranch")
	slug, ok := ctx.UserValue("slug").(string)
	if !ok {
		utils.SendServerError("Can't get value from slug", fasthttp.StatusInternalServerError, ctx)
		return
	}
	thread := &threadModels.Thread{Forum: slug}
	err := thread.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		utils.SendServerError(err.Error(), fasthttp.StatusInternalServerError, ctx)
		return
	}
	alredyExicted, err := f.usecaseThread.FindThreadBySlug(thread.Slug)
	if err == nil && alredyExicted.Slug != "" {
		body, err := alredyExicted.MarshalJSON()
		if err != nil {
			utils.SendServerError(err.Error(), fasthttp.StatusInternalServerError, ctx)
			return
		}
		utils.SendResponse(fasthttp.StatusConflict, body, ctx) // ???? alredy json
		return
	}
	_, err = f.usecaseForum.GetForumBySlug(thread.Forum)
	if err != nil {
		msg := fmt.Sprintf("Can't find forum with slug %s", slug)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}
	newThread, err := f.usecaseThread.CreateThread(thread)
	if err != nil {
		utils.SendServerError("Can't create thread", fasthttp.StatusInternalServerError, ctx)
		return
	}
	body, err := newThread.MarshalJSON()
	if err != nil {
		utils.SendServerError(err.Error(), fasthttp.StatusInternalServerError, ctx)
		return
	}
	utils.SendResponse(fasthttp.StatusCreated, body, ctx)
}

func (f *ForumHandler) CurrentForumUsers(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting CurrentForumUsers")
	slug, ok := ctx.UserValue("slug").(string)
	if !ok {
		utils.SendServerError("Can't get value from slug", fasthttp.StatusInternalServerError, ctx)
		return
	}
	desc := ctx.QueryArgs().GetBool("desc")
	limit, err := ctx.QueryArgs().GetUint("limit")
	if err != nil {
		limit = 100
	}
	since := string(ctx.QueryArgs().Peek("since"))

	users, err := f.usecaseForum.GetUsersByForum(slug, since, limit, desc)
	if err != nil {
		msg := fmt.Sprintf("Can't find user with slug %s", slug)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}
	utils.SendResponse(fasthttp.StatusOK, users, ctx)
}

func (f *ForumHandler) ForumBranches(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting ForumBranches")
	slug, ok := ctx.UserValue("slug").(string)
	if !ok {
		utils.SendServerError("Can't get value from slug", fasthttp.StatusInternalServerError, ctx)
		return
	}
	limit, err := strconv.Atoi(string(ctx.FormValue("limit")))
	if err != nil {
		fmt.Println(err)
	}
	since := string(ctx.QueryArgs().Peek("since"))
	if err != nil {
		fmt.Println(err)
	}
	desc := string(ctx.FormValue("desc"))
	if err != nil {
		fmt.Println(err)
	}

	threads, err := f.usecaseThread.GetThreadsByForumSlug(slug, since, desc, limit)
	if err != nil {
		msg := fmt.Sprintf("Can't find thread with slug %s", slug)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}

	if len(threads) == 0 {
		_, err = f.usecaseForum.GetForumBySlug(slug)
		if err != nil {
			msg := fmt.Sprintf("Can't find thread with slug %s", slug)
			utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
			return
		}
	}
	utils.SendResponse(fasthttp.StatusOK, threads, ctx)
}
