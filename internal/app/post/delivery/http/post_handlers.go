package http

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/amartery/tp_db_forum/internal/app/forum"
	"github.com/amartery/tp_db_forum/internal/app/post"
	"github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/amartery/tp_db_forum/internal/app/thread"
	"github.com/amartery/tp_db_forum/internal/app/user"
	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type PostHandler struct {
	usecaseUser   user.Usecase
	usecasePost   post.Usecase
	usecaseForum  forum.Usecase
	usecaseThread thread.Usecase
	logger        *logrus.Logger
}

func NewPostHandler(p post.Usecase, u user.Usecase, f forum.Usecase, t thread.Usecase) *PostHandler {
	return &PostHandler{
		usecaseUser:   u,
		usecasePost:   p,
		usecaseForum:  f,
		usecaseThread: t,
		logger:        logrus.New(),
	}
}

// TODO: s.router.GET("​/post​/{id}​/details", s.PostDetailsGet)
// TODO: s.router.POST("​/post​/{id}​/details", s.PostDetailsUpdate)

func (p *PostHandler) PostDetailsGet(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting PostDetailsGet")
	id := ctx.UserValue("id").(string)

	post, err := p.usecasePost.GetPost(id)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find post with id: %v", id),
		}

		_ = json.NewEncoder(ctx).Encode(msg)
		return
	}

	postInfo := models.PostResponse{
		Post: post,
	}

	related := string(ctx.URI().QueryArgs().Peek("related"))
	if strings.Contains(related, "user") {
		author, err := p.usecaseUser.Get(post.Author)
		if err != nil {
			msg := utils.Message{
				Text: fmt.Sprintf("Can't find user with id #%v\n", post.Author),
			}

			ctx.SetStatusCode(fasthttp.StatusNotFound)
			err = json.NewEncoder(ctx).Encode(msg)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
			return
		}
		postInfo.Author = author
	}

	if strings.Contains(related, "thread") {
		thread, err := p.usecaseThread.GetThread(strconv.Itoa(post.Thread))
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			msg := utils.Message{
				Text: fmt.Sprintf("Can't find thread forum by slug: %v", post.Thread),
			}

			_ = json.NewEncoder(ctx).Encode(msg)
			return
		}
		postInfo.Thread = thread
	}

	if strings.Contains(related, "forum") {
		forum, err := p.usecaseForum.GetForumBySlug(post.Forum)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			msg := utils.Message{
				Text: fmt.Sprintf("Can't find forum by slug: %v", post.Forum),
			}

			_ = json.NewEncoder(ctx).Encode(msg)
			return
		}
		postInfo.Forum = forum
	}

	err = json.NewEncoder(ctx).Encode(postInfo)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}

func (p *PostHandler) PostDetailsUpdate(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting PostDetailsUpdate")
	id := ctx.UserValue("id").(string)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	post := &models.Post{}
	err = json.Unmarshal(ctx.PostBody(), &post)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if post.Message == "" {
		post, err := p.usecasePost.GetPost(id)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			msg := utils.Message{
				Text: fmt.Sprintf("Can't find post with id: %v", id),
			}

			_ = json.NewEncoder(ctx).Encode(msg)
			return
		}

		err = json.NewEncoder(ctx).Encode(post)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		return
	}
	post.ID = idInt

	post, err = p.usecasePost.UpdatePost(post)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find post with id: %v", id),
		}

		_ = json.NewEncoder(ctx).Encode(msg)
		return
	}

	err = json.NewEncoder(ctx).Encode(post)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}
