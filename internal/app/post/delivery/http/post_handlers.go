package http

import (
	"fmt"

	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// TODO: s.router.GET("​/post​/{id}​/details", s.PostDetailsGet)  // ??? в сваггере написано что для ветки а не для поста
// TODO: s.router.POST("​/post​/{id}​/details", s.PostDetailsUpdate)

func PostDetailsGet(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting PostDetailsGet")
	id, ok := ctx.UserValue("id").(string)
	if !ok {
		utils.SendServerError("some err", ctx)
		return
	}
	ans := fmt.Sprintf("PostDetailsGet!\nid: %s\n", id)
	utils.SendResponse(200, ans, ctx)
}

func PostDetailsUpdate(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting PostDetailsUpdate")
	id, ok := ctx.UserValue("id").(string)
	if !ok {
		utils.SendServerError("some err", ctx)
		return
	}
	ans := fmt.Sprintf("PostDetailsUpdate!\nid: %s\n", id)
	utils.SendResponse(200, ans, ctx)
}
