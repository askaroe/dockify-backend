package user

import (
	"context"
	"fmt"

	"github.com/askaroe/dockify-backend/internal/models"
	"github.com/askaroe/dockify-backend/pkg/psql"
)

type User interface {
	CreateUser(ctx context.Context, req models.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
}

type user struct {
	db *psql.Client
}

func NewUserRepository(db *psql.Client) User {
	return &user{db: db}
}

func (u *user) CreateUser(ctx context.Context, req models.User) (int, error) {
	query := `INSERT INTO users (username, first_name, last_name, email, password_hash) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := u.db.QueryRow(ctx, query, req.Username, req.FirstName, req.LastName, req.Email, req.PasswordHash).Scan(&req.ID)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return req.ID, nil
}

func (u *user) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	query := `SELECT id, username, first_name, last_name, email, password_hash, created_at FROM users WHERE email = $1`
	err := u.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

func (u *user) GetUserByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
	query := `SELECT id, username, first_name, last_name, email, password_hash, created_at FROM users WHERE id = $1`
	err := u.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}
