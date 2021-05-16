package utils

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func SendServerError(errorMessage string, code int, ctx *fasthttp.RequestCtx) {
	logrus.Error("Loging in SendServerError: " + errorMessage)
	ctx.SetStatusCode(code)
	if errorMessage != "" {
		msg := fmt.Sprintf(`{"messege": "%s"}`, errorMessage)
		ctx.SetBody([]byte(msg))
	}
}

func SendResponse(code int, data interface{}, ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(code)
	serializedData, err := json.Marshal(data)
	if err != nil {
		SendServerError(err.Error(), fasthttp.StatusInternalServerError, ctx)
		return
	}
	ctx.SetBody(serializedData)
}
