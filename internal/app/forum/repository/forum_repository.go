package repository

import (
	"context"

	"github.com/amartery/tp_db_forum/internal/app/forum/models"
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
