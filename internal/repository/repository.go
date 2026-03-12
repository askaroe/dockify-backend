package repository

import (
	"github.com/askaroe/dockify-backend/internal/repository/document"
	"github.com/askaroe/dockify-backend/internal/repository/health"
	"github.com/askaroe/dockify-backend/internal/repository/location"
	"github.com/askaroe/dockify-backend/internal/repository/user"
	"github.com/askaroe/dockify-backend/pkg/psql"
)

type Repository struct {
	health.Health
	user.User
	location.Location
	Document document.Document
}

func NewRepository(client *psql.Client) *Repository {
	return &Repository{
		Health:   health.NewHealthRepository(client),
		User:     user.NewUserRepository(client),
		Location: location.NewLocationRepository(client),
		Document: document.NewDocumentRepository(client),
	}
}
