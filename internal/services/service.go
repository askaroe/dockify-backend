package services

import (
	"github.com/askaroe/dockify-backend/internal/gateway"
	"github.com/askaroe/dockify-backend/internal/repository"
	"github.com/askaroe/dockify-backend/internal/services/document"
	"github.com/askaroe/dockify-backend/internal/services/health"
	"github.com/askaroe/dockify-backend/internal/services/location"
	"github.com/askaroe/dockify-backend/internal/services/user"
)

type Service struct {
	health.Health
	user.User
	location.Location
	Document document.Document
	Gateway  *gateway.Gateway
}

func NewService(repo *repository.Repository, gw *gateway.Gateway) *Service {
	return &Service{
		Health:   health.NewHealthService(repo),
		User:     user.NewUserService(repo),
		Location: location.NewLocationService(repo),
		Document: document.NewDocumentService(repo),
		Gateway:  gw,
	}
}
