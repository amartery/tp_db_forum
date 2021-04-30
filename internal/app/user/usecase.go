package user

import "github.com/amartery/tp_db_forum/internal/app/user/models"

type Usecase interface {
	Create(model *models.User) error
	Get(nickname string) (*models.User, error)
	CheckIfUserExists(nickname string) (string, error)
	GetUsersWithNicknameAndEmail(nickname, email string) (*[]models.User, error)
	Update(model *models.User) (*models.User, error)
	GetUserNicknameWithEmail(email string) (string, error)
}
