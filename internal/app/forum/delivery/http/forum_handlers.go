package http

import (
	"encoding/json"
	"fmt"

	"github.com/amartery/tp_db_forum/internal/app/forum"
	"github.com/amartery/tp_db_forum/internal/app/forum/models"
	"github.com/amartery/tp_db_forum/internal/app/user"
	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type ForumHandler struct {
	usecaseForum forum.Usecase
	usecaseUser  user.Usecase
	logger       *logrus.Logger
}

func NewForumHandler(forumUsecase forum.Usecase, userUsecase user.Usecase) *ForumHandler {
	return &ForumHandler{
		usecaseForum: forumUsecase,
		usecaseUser:  userUsecase,
		logger:       logrus.New(),
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
		utils.SendServerError("", fasthttp.StatusInternalServerError, ctx)
		return
	}
	ans := fmt.Sprintf("ForumDetails!\nslug: %s\n", slug)
	utils.SendResponse(200, ans, ctx)
}

func (f *ForumHandler) ForumCreateBranch(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting ForumCreateBranch")
	slug, ok := ctx.UserValue("slug").(string)
	if !ok {
		utils.SendServerError("", fasthttp.StatusInternalServerError, ctx)
		return
	}
	ans := fmt.Sprintf("ForumCreateBranch!\nslug: %s\n", slug)
	utils.SendResponse(200, ans, ctx)
}

func (f *ForumHandler) CurrentForumUsers(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting CurrentForumUsers")
	slug, ok := ctx.UserValue("slug").(string)
	if !ok {
		utils.SendServerError("", fasthttp.StatusInternalServerError, ctx)
		return
	}
	ans := fmt.Sprintf("CurrentForumUsers!\nslug: %s\n", slug)
	utils.SendResponse(200, ans, ctx)
}

func (f *ForumHandler) ForumBranches(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting ForumBranches")
	slug, ok := ctx.UserValue("slug").(string)
	if !ok {
		utils.SendServerError("", fasthttp.StatusInternalServerError, ctx)
		return
	}
	ans := fmt.Sprintf("ForumBranches!\nslug: %s\n", slug)
	utils.SendResponse(200, ans, ctx)
}
