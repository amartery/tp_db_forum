package http

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/amartery/tp_db_forum/internal/app/forum"
	"github.com/amartery/tp_db_forum/internal/app/user"
	"github.com/amartery/tp_db_forum/internal/app/user/models"
	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	usecaseUser  user.Usecase
	usecaseForum forum.Usecase
	logger       *logrus.Logger
}

func NewUserHandler(u user.Usecase, f forum.Usecase) *UserHandler {
	return &UserHandler{
		usecaseUser:  u,
		usecaseForum: f,
		logger:       logrus.New(),
	}
}

// TODO: s.router.POST("/user​/{nickname}​/create", s.CreateUser)
// TODO: s.router.GET("/user​/{nickname}​/profile, s.AboutUserGet)
// TODO: s.router.POST("/user/{nickname}/profile", s.AboutUserUpdate)

func (handler *UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting CreateUser")
	nickname := ctx.UserValue("nickname").(string)

	user := &models.User{}
	err := json.Unmarshal(ctx.PostBody(), user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	user.Nickname = nickname

	err = handler.usecaseUser.Create(user)
	if err != nil {
		users, err := handler.usecaseUser.GetUsersWithNicknameAndEmail(nickname, *user.Email)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusConflict)
		err = json.NewEncoder(ctx).Encode(users)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	_ = json.NewEncoder(ctx).Encode(user)
}

func (handler *UserHandler) AboutUserGet(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting AboutUserGet")
	nickname := ctx.UserValue("nickname").(string)

	profile, err := handler.usecaseUser.Get(nickname)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		msg := utils.Message{
			Text: fmt.Sprintf("Can't find user with id #%v\n", nickname),
		}

		err = json.NewEncoder(ctx).Encode(msg)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		return
	}

	err = json.NewEncoder(ctx).Encode(profile)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) AboutUserUpdate(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting AboutUserUpdate")
	nickname := ctx.UserValue("nickname").(string)

	profile := &models.User{}
	err := json.Unmarshal(ctx.PostBody(), profile)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	profile.Nickname = nickname

	fullProfile, err := handler.usecaseUser.Update(profile)
	if err != nil {
		var msg utils.Message
		if errors.Is(err, user.ErrUserDoesntExists) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			msg = utils.Message{
				Text: fmt.Sprintf("Can't find user with id #%v\n", nickname),
			}
		} else if errors.Is(err, user.ErrDataConflict) {
			emailOwnerNickname, err := handler.usecaseUser.GetUserNicknameWithEmail(*profile.Email)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}

			ctx.SetStatusCode(fasthttp.StatusConflict)
			msg = utils.Message{
				Text: fmt.Sprintf("This email is already registered by user: %v", emailOwnerNickname),
			}
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(ctx).Encode(msg)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		return
	}

	err = json.NewEncoder(ctx).Encode(fullProfile)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}
