package repository

import "github.com/jackc/pgx/v4/pgxpool"

type ThreadRepository struct {
	Con *pgxpool.Pool
}

func NewThreadRepository(con *pgxpool.Pool) *ThreadRepository {
	return &ThreadRepository{
		Con: con,
	}
}
