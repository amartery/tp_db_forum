package http

import (
	"fmt"

	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// TODO: s.router.POST("/forum/create", DeliveryForum.ForumCreate)
// TODO: s.router.GET("/forum/{slug}/details", s.ForumDetails)
// TODO: s.router.POST("/forum/{slug}/create", s.ForumCreateBranch)
// TODO: s.router.GET("/forum/{slug}/users", s.CurrentForumUsers)
// TODO: s.router.GET("​/forum​/{slug}​/threads", s.ForumBranches)

func ForumCreate(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting ForumCreate")
	ans := fmt.Sprintf("ForumCreate!\n")
	utils.SendResponse(200, ans, ctx)
}

func ForumDetails(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting ForumDetails")
	slug, ok := ctx.UserValue("slug").(string)
	if !ok {
		utils.SendServerError("some err", ctx)
		return
	}
	ans := fmt.Sprintf("ForumDetails!\nslug: %s\n", slug)
	utils.SendResponse(200, ans, ctx)
}

func ForumCreateBranch(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting ForumCreateBranch")
	slug, ok := ctx.UserValue("slug").(string)
	if !ok {
		utils.SendServerError("some err", ctx)
		return
	}
	ans := fmt.Sprintf("ForumCreateBranch!\nslug: %s\n", slug)
	utils.SendResponse(200, ans, ctx)
}

func CurrentForumUsers(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting CurrentForumUsers")
	slug, ok := ctx.UserValue("slug").(string)
	if !ok {
		utils.SendServerError("some err", ctx)
		return
	}
	ans := fmt.Sprintf("CurrentForumUsers!\nslug: %s\n", slug)
	utils.SendResponse(200, ans, ctx)
}

func ForumBranches(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting ForumBranches")
	slug, ok := ctx.UserValue("slug").(string)
	if !ok {
		utils.SendServerError("some err", ctx)
		return
	}
	ans := fmt.Sprintf("ForumBranches!\nslug: %s\n", slug)
	utils.SendResponse(200, ans, ctx)
}
