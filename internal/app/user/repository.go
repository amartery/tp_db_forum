package user

import (
	"fmt"

	"github.com/amartery/tp_db_forum/internal/app/user/models"
)

var (
	ErrUserDoesntExists = fmt.Errorf("user exists")
	ErrDataConflict     = fmt.Errorf("data conflict")
)

type Repository interface {
	Get(nickname string) (*models.User, error)
	CheckIfUserExists(nickname string) (string, error)
	Create(model *models.User) error
	GetUsersWithNicknameAndEmail(nickname, email string) (*[]models.User, error)
	Update(model *models.User) (*models.User, error)
	GetUserNicknameWithEmail(email string) (string, error)
}
