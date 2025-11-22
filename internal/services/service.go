package services

import (
	"github.com/askaroe/dockify-backend/internal/repository"
	"github.com/askaroe/dockify-backend/internal/services/health"
	"github.com/askaroe/dockify-backend/internal/services/location"
	"github.com/askaroe/dockify-backend/internal/services/user"
)

type Service struct {
	health.Health
	user.User
	location.Location
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Health:   health.NewHealthService(repo),
		User:     user.NewUserService(repo),
		Location: location.NewLocationService(repo),
	}
}
