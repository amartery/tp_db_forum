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

func (u *UserUsecase) CheckIfUserExists(nickname string) (string, error) {
	return u.repo.CheckIfUserExists(nickname)
}

func (u *UserUsecase) Get(nickname string) (*models.User, error) {
	return u.repo.Get(nickname)
}

func (u *UserUsecase) Create(model *models.User) error {
	return u.repo.Create(model)
}

func (u *UserUsecase) GetUsersWithNicknameAndEmail(nickname, email string) (*[]models.User, error) {
	return u.repo.GetUsersWithNicknameAndEmail(nickname, email)
}

func (u *UserUsecase) Update(model *models.User) (*models.User, error) {
	return u.repo.Update(model)
}

func (u *UserUsecase) GetUserNicknameWithEmail(email string) (string, error) {
	return u.repo.GetUserNicknameWithEmail(email)
}
