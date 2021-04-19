package http

import (
	"fmt"

	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// TODO: s.router.POST("/user​/{nickname}​/create", s.CreateUser)
// TODO: s.router.GET("/user​/{nickname}​/profile, s.AboutUserGet)
// TODO: s.router.POST("/user/{nickname}/profile", s.AboutUserUpdate)

func CreateUser(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting CreateUser")
	nickname, ok := ctx.UserValue("nickname").(string)
	if !ok {
		utils.SendServerError("some err", fasthttp.StatusInternalServerError, ctx)
		return
	}
	ans := fmt.Sprintf("CreateUser!\nnickname: %s\n", nickname)
	utils.SendResponse(200, ans, ctx)
}

func AboutUserGet(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting AboutUserGet")
	nickname, ok := ctx.UserValue("nickname").(string)
	if !ok {
		utils.SendServerError("some err", fasthttp.StatusInternalServerError, ctx)
		return
	}
	ans := fmt.Sprintf("AboutUserGet!\nnickname: %s\n", nickname)
	utils.SendResponse(200, ans, ctx)
}

func AboutUserUpdate(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting AboutUserUpdate")
	nickname, ok := ctx.UserValue("nickname").(string)
	if !ok {
		utils.SendServerError("some err", fasthttp.StatusInternalServerError, ctx)
		return
	}
	ans := fmt.Sprintf("AboutUserUpdate!\nnickname: %s\n", nickname)
	utils.SendResponse(200, ans, ctx)
}
