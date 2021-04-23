package http

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amartery/tp_db_forum/internal/app/post"
	"github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type PostHandler struct {
	usecasePost post.Usecase
	logger      *logrus.Logger
}

func NewPostHandler(postUsecase post.Usecase) *PostHandler {
	return &PostHandler{
		usecasePost: postUsecase,
		logger:      logrus.New(),
	}
}

// TODO: s.router.GET("​/post​/{id}​/details", s.PostDetailsGet)  // ??? в сваггере написано что для ветки а не для поста
// TODO: s.router.POST("​/post​/{id}​/details", s.PostDetailsUpdate)

func (p *PostHandler) PostDetailsGet(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting PostDetailsGet")
	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		utils.SendServerError("Can't get id", fasthttp.StatusInternalServerError, ctx)
		return
	}

	relatedArr := string(ctx.QueryArgs().Peek("related"))
	relatedStr := strings.Split(relatedArr, ",")
	for len(relatedStr) < 3 {
		relatedStr = append(relatedStr, "")
	}

	postResponse, err := p.usecasePost.GetPost(id, relatedStr)
	if err != nil {
		msg := fmt.Sprintf("Can't find post with id %d", id)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}

	body, err := postResponse.MarshalJSON()
	if err != nil {
		utils.SendServerError(err.Error(), fasthttp.StatusInternalServerError, ctx)
	}

	utils.SendResponse(fasthttp.StatusOK, body, ctx)
}

func (p *PostHandler) PostDetailsUpdate(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting PostDetailsUpdate")
	id, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		utils.SendServerError("Can't get id", fasthttp.StatusInternalServerError, ctx)
		return
	}
	post := &models.Post{}
	err = post.UnmarshalJSON(ctx.Request.Body())
	if err != nil {
		utils.SendServerError(err.Error(), fasthttp.StatusInternalServerError, ctx)
		return
	}
	post.ID = id
	post, err = p.usecasePost.UpdatePost(post)
	if err != nil {
		msg := fmt.Sprintf("Can't find post with id %d", id)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}
	body, err := post.MarshalJSON()
	if err != nil {
		utils.SendServerError(err.Error(), fasthttp.StatusInternalServerError, ctx)
		return
	}
	utils.SendResponse(fasthttp.StatusOK, body, ctx)
}
