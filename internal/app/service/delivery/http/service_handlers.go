package http

import (
	"fmt"

	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// TODO: s.router.POST("​/service​/clear", s.ServiceClear)
// TODO: s.router.GET("​/service​/status", s.ServiceStatus)

func ServiceClear(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting ServiceClear")
	ans := fmt.Sprintf("ServiceClear!\n")
	utils.SendResponse(200, ans, ctx)
}

func ServiceStatus(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting ServiceStatus")
	ans := fmt.Sprintf("ServiceStatus!\n")
	utils.SendResponse(200, ans, ctx)
}
