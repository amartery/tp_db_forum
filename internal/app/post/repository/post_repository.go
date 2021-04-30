package repository

import (
	"context"
	"strconv"

	"github.com/amartery/tp_db_forum/internal/app/post/models"
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

func (repo *PostRepository) GetPost(id string) (*models.Post, error) {
	query := `SELECT author, created, forum, id, msg, thread, isEdited, parent FROM posts WHERE id = $1`
	post := &models.Post{}
	err := repo.Con.QueryRow(context.Background(), query, id).Scan(
		&post.Author,
		&post.Created,
		&post.Forum,
		&post.ID,
		&post.Message,
		&post.Thread,
		&post.IsEdited,
		&post.Parent)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (repo *PostRepository) UpdatePost(post *models.Post) (*models.Post, error) {
	postDB, err := repo.GetPost(strconv.Itoa(post.ID))
	if err != nil {
		return nil, err
	}

	if post.Message == postDB.Message {
		return postDB, nil
	}
	query := `UPDATE posts SET msg = $1, isEdited = true 
	          WHERE id = $2
	          RETURNING author, created, forum, id, msg, thread, isEdited, parent`
	err = repo.Con.QueryRow(context.Background(), query, post.Message, post.ID).Scan(
		&post.Author,
		&post.Created,
		&post.Forum,
		&post.ID,
		&post.Message,
		&post.Thread,
		&post.IsEdited,
		&post.Parent)

	if err != nil {
		return nil, err
	}
	return post, nil
}
