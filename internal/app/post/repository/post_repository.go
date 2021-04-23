package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/amartery/tp_db_forum/internal/app/post/models"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostRepository struct {
	Con *pgxpool.Pool
}

func NewPostRepository(con *pgxpool.Pool) *PostRepository {
	return &PostRepository{
		Con: con,
	}
}

func (repo *PostRepository) GetPost(postID int) (*models.Post, error) {
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Posts
			 WHERE id = $1
			 ORDER BY created, id`

	post := &models.Post{}
	t := &time.Time{}
	var parent sql.NullInt64
	err := repo.Con.QueryRow(
		context.Background(),
		query,
		postID).Scan(
		&post.ID,
		&parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		t)
	post.Created = strfmt.DateTime(t.UTC()).String()
	if parent.Valid {
		post.Parent = int(parent.Int64)
	}

	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *PostRepository) UpdatePostByID(post *models.Post) (*models.Post, error) {
	query := `UPDATE Posts
			  SET message = (CASE WHEN LTRIM($1) = '' THEN message ELSE $1 END),
    		      is_edited = (CASE WHEN LTRIM($1) = '' THEN false ELSE true END)
			  WHERE id = $2 RETURNING id, parent, author, message, is_edited, forum, thread, created`

	t := &time.Time{}
	var parent sql.NullInt64
	err := repo.Con.QueryRow(
		context.Background(),
		query,
		post.Message,
		post.ID).Scan(
		&post.ID,
		&parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		t)
	if err != nil {
		return nil, err
	}
	post.Created = strfmt.DateTime(t.UTC()).String()
	if parent.Valid {
		post.Parent = int(parent.Int64)
	}

	return post, nil
}
