package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func SendServerError(errorMessage string, ctx *fasthttp.RequestCtx) {
	logrus.Error(errors.New(errorMessage))
	ctx.SetStatusCode(http.StatusInternalServerError)
}

func SendResponse(code int, data interface{}, ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(code)
	serializedData, err := json.Marshal(data)
	if err != nil {
		SendServerError(err.Error(), ctx)
		return
	}
	ctx.SetBody(serializedData)
}

func SendResponseOK(data interface{}, ctx *fasthttp.RequestCtx) {
	SendResponse(200, data, ctx)
}
