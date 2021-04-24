package http

import (
	"encoding/json"
	"fmt"

	"github.com/amartery/tp_db_forum/internal/app/user"
	"github.com/amartery/tp_db_forum/internal/app/user/models"
	"github.com/amartery/tp_db_forum/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	usecaseUser user.Usecase
	logger      *logrus.Logger
}

func NewUserHandler(userUsecase user.Usecase) *UserHandler {
	return &UserHandler{
		usecaseUser: userUsecase,
		logger:      logrus.New(),
	}
}

// TODO: s.router.POST("/user​/{nickname}​/create", s.CreateUser)
// TODO: s.router.GET("/user​/{nickname}​/profile, s.AboutUserGet)
// TODO: s.router.POST("/user/{nickname}/profile", s.AboutUserUpdate)

func (handler *UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting CreateUser")
	nickname, ok := ctx.UserValue("nickname").(string)
	if !ok {
		utils.SendServerError("some err", fasthttp.StatusInternalServerError, ctx)
		return
	}
	newUser := &models.User{}
	err := json.Unmarshal(ctx.PostBody(), newUser)
	if err != nil {
		fmt.Println(err)
	}
	newUser.Nickname = nickname
	err = handler.usecaseUser.CreateUser(newUser)
	if err != nil {
		alredyExictedUser, err := handler.usecaseUser.GetUserByEmailOrNickname(newUser.Nickname, newUser.Email)
		if err != nil {
			utils.SendServerError(err.Error(), fasthttp.StatusInternalServerError, ctx)
			return
		}
		utils.SendResponse(fasthttp.StatusConflict, alredyExictedUser, ctx)
		return
	}

	utils.SendResponse(fasthttp.StatusCreated, newUser, ctx)
}

func (handler *UserHandler) AboutUserGet(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting AboutUserGet")
	nickname, ok := ctx.UserValue("nickname").(string)
	if !ok {
		utils.SendServerError("some err", fasthttp.StatusInternalServerError, ctx)
		return
	}
	user, err := handler.usecaseUser.GetUserByNickname(nickname)
	if err != nil {
		msg := fmt.Sprintf("Can't find user with nickname %s", nickname)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}

	utils.SendResponse(fasthttp.StatusOK, user, ctx)
}

func (handler *UserHandler) AboutUserUpdate(ctx *fasthttp.RequestCtx) {
	logrus.Info("starting AboutUserUpdate")
	nickname, ok := ctx.UserValue("nickname").(string)
	if !ok {
		utils.SendServerError("some err", fasthttp.StatusInternalServerError, ctx)
		return
	}
	newUser := &models.User{}
	err := json.Unmarshal(ctx.PostBody(), newUser)
	if err != nil {
		fmt.Println(err)
	}
	newUser.Nickname = nickname

	us, err := handler.usecaseUser.GetUserByEmail(newUser.Email)
	if err == nil && us.Nickname != newUser.Nickname {
		msg := fmt.Sprintf("Can't find user with nickname %s", nickname)
		utils.SendServerError(msg, fasthttp.StatusConflict, ctx)
		return
	}
	oldUser, err := handler.usecaseUser.GetUserByNickname(nickname)
	if err != nil {
		msg := fmt.Sprintf("Can't find user with nickname %s", nickname)
		utils.SendServerError(msg, fasthttp.StatusNotFound, ctx)
		return
	}

	err = handler.usecaseUser.UpdateUserInformation(newUser)

	if newUser.FullName == "" {
		newUser.FullName = oldUser.FullName
	}
	if newUser.About == "" {
		newUser.About = oldUser.About
	}
	if newUser.Email == "" {
		newUser.Email = oldUser.Email
	}
	utils.SendResponse(fasthttp.StatusOK, newUser, ctx)
}
