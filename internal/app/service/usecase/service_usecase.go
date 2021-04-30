package usecase

import (
	"github.com/amartery/tp_db_forum/internal/app/service"
	"github.com/amartery/tp_db_forum/internal/app/service/models"
)

type ServiceUsecase struct {
	repo service.Repository
}

func NewServiceUsecase(ServiceRepo service.Repository) *ServiceUsecase {
	return &ServiceUsecase{
		repo: ServiceRepo,
	}
}

func (usecase *ServiceUsecase) ClearDB() error {
	err := usecase.repo.ClearDB()
	return err
}

func (usecase *ServiceUsecase) GetStatusDB() (*models.Status, error) {
	status, err := usecase.repo.GetStatusDB()
	return status, err
}
