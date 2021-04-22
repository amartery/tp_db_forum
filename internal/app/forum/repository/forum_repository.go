package repository

import (
	"context"
	"fmt"

	"github.com/amartery/tp_db_forum/internal/app/forum/models"
	usersModels "github.com/amartery/tp_db_forum/internal/app/user/models"
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
	query := `INSERT INTO Forum (title, user_nickname, slug, posts, threads) VALUES ($1, $2, $3, $4, $5)`
	_, err := repo.Con.Exec(
		context.Background(),
		query,
		forum.Tittle,
		forum.Nickname,
		forum.Slug,
		forum.Posts,
		forum.Threads,
	)
	return err
}

func (repo *ForumRepository) GetForumBySlug(slug string) (*models.Forum, error) {
	query := `SELECT title, user_nickname, slug, posts, threads FROM Forum WHERE slug = $1`
	forum := &models.Forum{}

	err := repo.Con.QueryRow(
		context.Background(),
		query,
		slug).Scan(
		&forum.Tittle,
		&forum.Nickname,
		&forum.Slug,
		&forum.Posts,
		&forum.Threads)

	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (repo *ForumRepository) GetUsersByForum(slug, since string, limit int, desc bool) ([]*usersModels.User, error) {
	query := fmt.Sprintf(`select u.nickname, u.fullname, u.about, u.email from users_to_forums
			left join users u on users_to_forums.nickname = u.nickname
			where users_to_forums.forum = '%s'`, slug)
	if desc && since != "" {
		query += fmt.Sprintf(` and u.nickname < '%s'`, since)
	} else if since != "" {
		query += fmt.Sprintf(` and u.nickname > '%s'`, since)
	}
	query += ` order by u.nickname `
	if desc {
		query += "desc"
	}
	query += fmt.Sprintf(` limit %d`, limit)
	rows, err := repo.Con.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*usersModels.User, 0)

	for rows.Next() {
		user := &usersModels.User{}
		err := rows.Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, user)
	}
	return users, nil
}
