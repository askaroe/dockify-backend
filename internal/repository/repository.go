package repository

import (
	"github.com/askaroe/dockify-backend/internal/repository/health"
	"github.com/askaroe/dockify-backend/internal/repository/user"
	"github.com/askaroe/dockify-backend/pkg/psql"
)

type Repository struct {
	health.Health
	user.User
}

func NewRepository(client *psql.Client) *Repository {
	return &Repository{
		Health: health.NewHealthRepository(client),
		User:   user.NewUserRepository(client),
	}
}
