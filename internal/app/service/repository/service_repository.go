package repository

import (
	"context"

	"github.com/amartery/tp_db_forum/internal/app/service/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ServiceRepository struct {
	Con *pgxpool.Pool
}

func NewServiceRepository(con *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{
		Con: con,
	}
}

func (repository *ServiceRepository) ClearDB() error {
	query := `TRUNCATE Users_to_forums RESTART IDENTITY;
			  TRUNCATE Votes RESTART IDENTITY;
			  TRUNCATE Posts RESTART IDENTITY;
			  TRUNCATE Threads RESTART IDENTITY;
			  TRUNCATE Forum RESTART IDENTITY;
			  TRUNCATE Users RESTART IDENTITY;`

	_, err := repository.Con.Exec(context.Background(), query)
	if err != nil {
		return err
	}
	return nil
}

func (repository *ServiceRepository) GetStatusDB() (*models.Status, error) {
	queryUser := `SELECT COUNT(*) AS user_count FROM Users;`
	queryForum := `SELECT COUNT(*) AS forum_count FROM Forum;`
	queryThread := `SELECT COUNT(*) AS thread_count FROM Threads;`
	queryPost := `SELECT COUNT(*) AS post_count FROM Posts;`

	status := &models.Status{}

	err := repository.Con.QueryRow(context.Background(), queryUser).Scan(&status.User)
	if err != nil {
		return nil, err
	}
	err = repository.Con.QueryRow(context.Background(), queryForum).Scan(&status.Forum)
	if err != nil {
		return nil, err
	}
	err = repository.Con.QueryRow(context.Background(), queryThread).Scan(&status.Thread)
	if err != nil {
		return nil, err
	}
	err = repository.Con.QueryRow(context.Background(), queryPost).Scan(&status.Post)
	if err != nil {
		return nil, err
	}

	return status, nil
}
