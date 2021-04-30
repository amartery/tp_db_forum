package repository

import (
	"context"
	"fmt"

	"github.com/amartery/tp_db_forum/internal/app/user"
	"github.com/amartery/tp_db_forum/internal/app/user/models"
	"github.com/jackc/pgx/v4"
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

func (u *UserRepository) CheckIfUserExists(nickname string) (string, error) {
	query := `SELECT nickname FROM users WHERE nickname = $1`

	err := u.Con.QueryRow(context.Background(), query, nickname).Scan(&nickname)
	if err != nil {
		return "", fmt.Errorf("user doesnt exist: %w", err)
	}
	return nickname, nil
}

func (u *UserRepository) Get(nickname string) (*models.User, error) {
	query := `SELECT id, nickname, fullname, email, about FROM users WHERE nickname = $1`
	model := &models.User{}
	err := u.Con.QueryRow(context.Background(), query, nickname).Scan(
		&model.ID,
		&model.Nickname,
		&model.Fullname,
		&model.Email,
		&model.About)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, user.ErrUserDoesntExists
		}
		return nil, fmt.Errorf("couldn't get user with nickname '%v'. Error: %w", nickname, err)
	}

	return model, nil
}

func (u *UserRepository) Create(model *models.User) error {
	_, err := u.Con.Exec(context.Background(),
		"INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)",
		model.Nickname, model.Fullname, model.Email, model.About,
	)

	if err != nil {
		return fmt.Errorf("couldn't insert user: %v. Error: %w", model, err)
	}

	return nil
}

func (u *UserRepository) GetUsersWithNicknameAndEmail(nickname, email string) (*[]models.User, error) {
	rows, err := u.Con.Query(context.Background(),
		`SELECT nickname, fullname, email, about FROM users
		WHERE nickname = $1 OR email = $2`,
		nickname, email,
	)
	if err != nil {
		return nil, fmt.Errorf(`couldn't get users with nickname '%v' and email '%v'. Error: %w`, nickname, email, err)
	}
	defer rows.Close()

	users := make([]models.User, 0, 2)
	user := &models.User{}
	for rows.Next() {
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
		if err != nil {
			return nil, fmt.Errorf(`couldn't get users with nickname '%v' and email '%v'. Error: %w`, nickname, email, err)
		}

		users = append(users, *user)
	}

	return &users, nil
}

func (u *UserRepository) Update(model *models.User) (*models.User, error) {
	userFromDB, err := u.Get(model.Nickname)
	if err != nil {
		fmt.Println("fff")
		return nil, err
	}

	if model.Fullname == nil {
		model.Fullname = userFromDB.Fullname
	}

	if model.Email == nil {
		model.Email = userFromDB.Email
	}

	if model.About == nil {
		model.About = userFromDB.About
	}

	_, err = u.Con.Exec(context.Background(),
		`UPDATE users SET fullname = $1, email = $2, about = $3
		WHERE id = $4`,
		model.Fullname, model.Email, model.About, userFromDB.ID,
	)

	if err != nil {
		return nil, user.ErrDataConflict
	}

	return model, nil
}

func (u *UserRepository) GetUserNicknameWithEmail(email string) (string, error) {
	var nickname string
	err := u.Con.QueryRow(context.Background(),
		"SELECT nickname FROM users WHERE email = $1",
		email,
	).Scan(&nickname)

	if err != nil {
		return "", fmt.Errorf(`couldn't get user nickname with email '%v'. Error: %w`, email, err)
	}

	return nickname, nil
}
