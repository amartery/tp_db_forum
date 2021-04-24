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

func (u *UserUsecase) CreateUser(user *models.User) error {
	err := u.repo.CreateUser(user)
	return err
}

func (u *UserUsecase) GetUserByEmailOrNickname(nickname, email string) ([]*models.User, error) {
	users, err := u.repo.GetUserByEmailOrNickname(nickname, email)
	return users, err
}

func (u *UserUsecase) UpdateUserInformation(user *models.User) error {
	err := u.repo.UpdateUserInformation(user)
	return err
}

func (u *UserUsecase) GetUserByEmail(email string) (*models.User, error) {
	user, err := u.repo.GetUserByEmail(email)
	return user, err
}
