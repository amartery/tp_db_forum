package http

import (
	"fmt"

	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// TODO: s.router.POST("​/thread​/{slug_or_id}​/create", s.CreatePostInBranch)
// TODO: s.router.GET("/thread​/{slug_or_id}​/details", s.BranchDetailsGet)
// TODO: s.router.POST("/thread​/{slug_or_id}​/details", s.BranchDetailsUpdate)
// TODO: s.router.GET("/thread/{slug_or_id}/posts", s.CurrentBranchPosts)
// TODO: s.router.POST("/thread​/{slug_or_id}​/vote", s.VoteForBranch)

func CreatePostInBranch(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting CreatePostInBranch")
	slug_or_id, ok := ctx.UserValue("slug_or_id").(string)
	if !ok {
		utils.SendServerError("some err", ctx)
		return
	}
	ans := fmt.Sprintf("CreatePostInBranch!\nslug_or_id: %s\n", slug_or_id)
	utils.SendResponse(200, ans, ctx)
}

func BranchDetailsGet(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting BranchDetailsGet")
	slug_or_id, ok := ctx.UserValue("slug_or_id").(string)
	if !ok {
		utils.SendServerError("some err", ctx)
		return
	}
	ans := fmt.Sprintf("BranchDetailsGet!\nslug_or_id: %s\n", slug_or_id)
	utils.SendResponse(200, ans, ctx)
}

func BranchDetailsUpdate(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting BranchDetailsUpdate")
	slug_or_id, ok := ctx.UserValue("slug_or_id").(string)
	if !ok {
		utils.SendServerError("some err", ctx)
		return
	}
	ans := fmt.Sprintf("BranchDetailsUpdate!\nslug_or_id: %s\n", slug_or_id)
	utils.SendResponse(200, ans, ctx)
}

func CurrentBranchPosts(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting CurrentBranchPosts")
	slug_or_id, ok := ctx.UserValue("slug_or_id").(string)
	if !ok {
		utils.SendServerError("some err", ctx)
		return
	}
	ans := fmt.Sprintf("CurrentBranchPosts!\nslug_or_id: %s\n", slug_or_id)
	utils.SendResponse(200, ans, ctx)
}

func VoteForBranch(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting VoteForBranch")
	slug_or_id, ok := ctx.UserValue("slug_or_id").(string)
	if !ok {
		utils.SendServerError("some err", ctx)
		return
	}
	ans := fmt.Sprintf("VoteForBranch!\nslug_or_id: %s\n", slug_or_id)
	utils.SendResponse(200, ans, ctx)
}
