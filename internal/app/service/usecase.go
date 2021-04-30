package service

import (
	"github.com/amartery/tp_db_forum/internal/app/service/models"
)

type Usecase interface {
	ClearDB() error
	GetStatusDB() (*models.Status, error)
}
