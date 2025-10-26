package services

import "github.com/askaroe/dockify-backend/internal/repository"

type Service struct {
}

func NewService(
	repo *repository.Repository) *Service {
	return &Service{}
}
