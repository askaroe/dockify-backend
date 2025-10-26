package repository

import (
	psql "github.com/askaroe/dockify-backend/pkg/psql"
)

type Repository struct {
}

func NewRepository(client *psql.Client) *Repository {
	return &Repository{}
}
