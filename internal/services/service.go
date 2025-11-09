package services

import (
	"github.com/askaroe/dockify-backend/internal/repository"
	"github.com/askaroe/dockify-backend/internal/services/health"
	"github.com/askaroe/dockify-backend/internal/services/user"
)

type Service struct {
	health.Health
	user.User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Health: health.NewHealthService(repo),
		User:   user.NewUserService(repo),
	}
}
