package repository

import "github.com/jackc/pgx/v4/pgxpool"

type ServiceRepository struct {
	Con *pgxpool.Pool
}

func NewServiceRepository(con *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{
		Con: con,
	}
}
