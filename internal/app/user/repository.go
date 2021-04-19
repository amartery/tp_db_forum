package user

import "github.com/amartery/tp_db_forum/internal/app/user/models"

type Repository interface {
	GetUserByNickname(nickname string) (*models.User, error)
}
