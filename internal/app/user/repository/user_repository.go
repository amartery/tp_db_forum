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

func (u *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO Users (nickname, fullname, about, email)
			  VALUES ($1, $2, $3, $4)`

	_, err := u.Con.Exec(
		context.Background(),
		query,
		user.Nickname,
		user.FullName,
		user.About,
		user.Email)

	return err
}

func (u *UserRepository) GetUserByEmailOrNickname(nickname, email string) ([]*models.User, error) {
	query := `SELECT nickname, fullname, about, email FROM Users
			  WHERE nickname = $1 OR email = $2`

	users := make([]*models.User, 0)

	rows, err := u.Con.Query(
		context.Background(),
		query,
		nickname,
		email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(
			&user.Nickname,
			&user.FullName,
			&user.About,
			&user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserRepository) UpdateUserInformation(user *models.User) error {
	query := `UPDATE Users SET fullname = (CASE WHEN LTRIM($1) = '' THEN fullname ELSE $1 END), 
	          about = (CASE WHEN $2 = '' THEN about ELSE $2 END), 
			  email = (CASE WHEN LTRIM($3) = '' THEN email ELSE LTRIM($3) END)
              WHERE nickname = $4`

	_, err := u.Con.Exec(
		context.Background(),
		query,
		user.FullName,
		user.About,
		user.Email,
		user.Nickname)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT nickname, fullname, about, email FROM Users WHERE email = $1`
	user := &models.User{}

	err := u.Con.QueryRow(
		context.Background(),
		query,
		email).Scan(
		&user.Nickname,
		&user.FullName,
		&user.About,
		&user.Email)

	if err != nil {
		return nil, err
	}
	return user, nil
}
