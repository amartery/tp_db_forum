package repository

import (
	"context"

	"github.com/amartery/tp_db_forum/internal/app/user/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	Con *pgxpool.Pool
}

func NewUserRepository(con *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		Con: con,
	}
}

func (u *UserRepository) GetUserByNickname(nickname string) (*models.User, error) {
	query := `SELECT nickname, fullname, about, email FROM Users WHERE nickname = $1`
	user := &models.User{}

	err := u.Con.QueryRow(
		context.Background(),
		query,
		nickname).Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)

	if err != nil {
		return nil, err
	}
	return user, nil
}
