package repository

import "github.com/jackc/pgx/v4/pgxpool"

type PostRepository struct {
	Con *pgxpool.Pool
}

func NewPostRepository(con *pgxpool.Pool) *PostRepository {
	return &PostRepository{
		Con: con,
	}
}
