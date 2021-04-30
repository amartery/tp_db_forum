package http

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/amartery/tp_db_forum/internal/app/forum"
	forumModel "github.com/amartery/tp_db_forum/internal/app/forum/models"
	"github.com/amartery/tp_db_forum/internal/app/thread"
	threadModel "github.com/amartery/tp_db_forum/internal/app/thread/models"

	// usersModels "github.com/amartery/tp_db_forum/internal/app/user/models"
	"github.com/amartery/tp_db_forum/internal/app/user"
	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type ForumHandler struct {
	usecaseForum  forum.Usecase
	usecaseUser   user.Usecase
	usecaseThread thread.Usecase
	logger        *logrus.Logger
}

func NewForumHandler(forumUsecase forum.Usecase, userUsecase user.Usecase, threadUsecase thread.Usecase) *ForumHandler {
	return &ForumHandler{
		usecaseForum:  forumUsecase,
		usecaseUser:   userUsecase,
		usecaseThread: threadUsecase,
		logger:        logrus.New(),
	}
}

// TODO: s.router.POST("/forum/create", DeliveryForum.ForumCreate)
// TODO: s.router.GET("/forum/{slug}/details", s.ForumDetails)
// TODO: s.router.POST("/forum/{slug}/create", s.ForumCreateBranch)
// TODO: s.router.GET("/forum/{slug}/users", s.CurrentForumUsers)
// TODO: s.router.GET("​/forum​/{slug}​/threads", s.ForumBranches)

func (f *ForumHandler) ForumCreate(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting ForumCreate")
	forum := &forumModel.Forum{}
	err := json.Unmarshal(ctx.PostBody(), &forum)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	nickname, err := f.usecaseUser.CheckIfUserExists(forum.User)
	if err != nil {
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find user with id #%v\n", forum.User),
		}

		ctx.SetStatusCode(fasthttp.StatusNotFound)
		err = json.NewEncoder(ctx).Encode(msg)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		return
	}
	forum.User = nickname

	err = f.usecaseForum.CreateForum(forum)
	if err != nil {

		existingForum, err := f.usecaseForum.GetForumBySlug(forum.Slug)
		if err != nil {

			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusConflict)
		err = json.NewEncoder(ctx).Encode(existingForum)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	err = json.NewEncoder(ctx).Encode(forum)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}

func (f *ForumHandler) ForumDetails(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting ForumDetails")
	slug := ctx.UserValue("slug").(string)
	forum, err := f.usecaseForum.GetForumBySlug(slug)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find forum with slug: %v", slug),
		}

		_ = json.NewEncoder(ctx).Encode(msg)
		return
	}

	err = json.NewEncoder(ctx).Encode(forum)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}

func (f *ForumHandler) ForumCreateBranch(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting ForumCreateBranch")
	slug := ctx.UserValue("slug").(string)

	thread := &threadModel.Thread{}
	err := json.Unmarshal(ctx.PostBody(), thread)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	thread.Forum = slug

	nickname, err := f.usecaseUser.CheckIfUserExists(thread.Author)
	if err != nil {
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find user with id #%v\n", thread.Author),
		}

		ctx.SetStatusCode(fasthttp.StatusNotFound)
		err = json.NewEncoder(ctx).Encode(msg)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		return
	}
	thread.Author = nickname

	err = f.usecaseThread.CreateThread(thread)
	if err != nil {
		if err == forum.ErrForumDoesntExists {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			msg := utils.Message{
				Text: fmt.Sprintf("Can't find thread forum by slug: %v", thread.Forum),
			}

			_ = json.NewEncoder(ctx).Encode(msg)
			return
		}

		existedThread, err := f.usecaseThread.GetThread(*thread.Slug)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusConflict)
		_ = json.NewEncoder(ctx).Encode(existedThread)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	err = json.NewEncoder(ctx).Encode(thread)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}

func (f *ForumHandler) CurrentForumUsers(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting CurrentForumUsers")
	slug := ctx.UserValue("slug").(string)

	_, err := f.usecaseForum.CheckForum(slug)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find forum by slug: %v", slug),
		}

		_ = json.NewEncoder(ctx).Encode(msg)
		return
	}
	fmt.Println("ggg")
	limitParam := string(ctx.URI().QueryArgs().Peek("limit"))
	limit, err := strconv.Atoi(limitParam)
	if err != nil && limitParam != "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	descParam := string(ctx.URI().QueryArgs().Peek("desc"))
	if descParam == "" {
		descParam = "false"
	}

	sinceParam := string(ctx.URI().QueryArgs().Peek("since"))

	users, err := f.usecaseForum.GetUsersByForum(slug, limit, sinceParam, descParam)
	if err != nil {
		fmt.Println("err>:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(ctx).Encode(users)
	if err != nil {
		fmt.Println("err>>:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}

func (f *ForumHandler) ForumBranches(ctx *fasthttp.RequestCtx) {
	f.logger.Info("starting ForumBranches")
	slug := ctx.UserValue("slug").(string)

	_, err := f.usecaseForum.CheckForum(slug)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find forum by slug: %v", slug),
		}

		_ = json.NewEncoder(ctx).Encode(msg)
		return
	}

	limit := string(ctx.URI().QueryArgs().Peek("limit"))
	desc := string(ctx.URI().QueryArgs().Peek("desc"))
	since := string(ctx.URI().QueryArgs().Peek("since"))

	threads, err := f.usecaseThread.GetThreadsByForumSlug(slug, limit, since, desc)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(ctx).Encode(threads)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}
