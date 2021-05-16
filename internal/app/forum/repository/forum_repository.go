package repository

import (
	"context"
	"fmt"

	"github.com/amartery/tp_db_forum/internal/app/forum/models"
	userModel "github.com/amartery/tp_db_forum/internal/app/user/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ForumRepository struct {
	Con *pgxpool.Pool
}

func NewForumRepository(con *pgxpool.Pool) *ForumRepository {
	return &ForumRepository{
		Con: con,
	}
}

func (repo *ForumRepository) CreateForum(forum *models.Forum) error {
	query := `INSERT INTO forums (slug, title, user_nickname) VALUES($1, $2, $3)`
	_, err := repo.Con.Exec(
		context.Background(),
		query,
		forum.Slug,
		forum.Tittle,
		forum.User,
	)
	return err
}

func (repo *ForumRepository) GetForumBySlug(slug string) (*models.Forum, error) {
	query := `SELECT slug, title, user_nickname, thread_count, post_count FROM forums WHERE slug = $1`
	forum := &models.Forum{}

	err := repo.Con.QueryRow(context.Background(), query, slug).Scan(
		&forum.Slug,
		&forum.Tittle,
		&forum.User,
		&forum.Threads,
		&forum.Posts)

	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (repo *ForumRepository) GetUsersByForum(slug string, limit int, since string, desc string) (*[]userModel.User, error) {

	var compare string
	if desc == "DESC" {
		compare = "<"
	} else {
		compare = ">"
	}

	var query string
	if since != "" {
		if limit != 0 {
			query = fmt.Sprintf(`SELECT u.about, u.email, u.fullname, u.nickname FROM users AS u
				JOIN forum_user AS fu ON u.nickname = fu.nickname
				WHERE fu.forum_slug = '%s' AND fu.nickname %v '%s'
				ORDER BY u.nickname %v
				LIMIT %v`, slug, compare, since, desc, limit)
		} else {
			query = fmt.Sprintf(`SELECT u.about, u.email, u.fullname, u.nickname FROM users AS u
				JOIN forum_user AS fu ON u.nickname = fu.nickname
				WHERE fu.forum_slug = '%s' AND fu.nickname %v '%s'
				ORDER BY u.nickname %v`, slug, compare, since, desc)
		}
	} else {
		if limit != 0 {
			query = fmt.Sprintf(`SELECT u.about, u.email, u.fullname, u.nickname FROM users AS u
				JOIN forum_user AS fu ON u.nickname = fu.nickname
				WHERE fu.forum_slug = '%s'
				ORDER BY u.nickname %v
				LIMIT %v`, slug, desc, limit)
		} else {
			query = fmt.Sprintf(`SELECT u.about, u.email, u.fullname, u.nickname FROM users AS u
				JOIN forum_user AS fu ON u.nickname = fu.nickname
				WHERE fu.forum_slug = '%s'
				ORDER BY u.nickname %v`, slug, desc)
		}
	}

	rows, err := repo.Con.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err != nil {
		fmt.Println("err1:", err)
		return nil, err
	}

	users := make([]userModel.User, 0, limit)
	user := userModel.User{}
	for rows.Next() {
		err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			fmt.Println("err2:", err)
			return nil, err
		}

		users = append(users, user)
	}
	return &users, nil
}

func (repo *ForumRepository) CheckForum(slug string) (string, error) {
	query := `SELECT slug FROM forums WHERE slug = $1`
	err := repo.Con.QueryRow(context.Background(), query, slug).Scan(&slug)

	if err != nil {
		return "", fmt.Errorf("couldn't get forum with slug '%v'. Error: %w", slug, err)
	}

	return slug, nil
}
