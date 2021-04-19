package usecase

import (
	"github.com/amartery/tp_db_forum/internal/app/user"
	"github.com/amartery/tp_db_forum/internal/app/user/models"
)

type UserUsecase struct {
	repo user.Repository
}

func NewUserUsecase(UserRepo user.Repository) *UserUsecase {
	return &UserUsecase{
		repo: UserRepo,
	}
}

func (u *UserUsecase) GetUserByNickname(nickname string) (*models.User, error) {
	user, err := u.repo.GetUserByNickname(nickname)
	return user, err
}
