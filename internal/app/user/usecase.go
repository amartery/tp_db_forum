package user

import "github.com/amartery/tp_db_forum/internal/app/user/models"

type Usecase interface {
	GetUserByNickname(nickname string) (*models.User, error)
	CreateUser(user *models.User) error
	GetUserByEmailOrNickname(nickname, email string) ([]*models.User, error)
	UpdateUserInformation(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}
