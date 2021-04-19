package usecase

import "github.com/amartery/tp_db_forum/internal/app/service"

type ServiceUsecase struct {
	repo service.Repository
}

func NewServiceUsecase(ServiceRepo service.Repository) *ServiceUsecase {
	return &ServiceUsecase{
		repo: ServiceRepo,
	}
}
