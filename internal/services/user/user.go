package user

import (
	"context"

	"github.com/askaroe/dockify-backend/internal/entity"
	"github.com/askaroe/dockify-backend/internal/models"
	"github.com/askaroe/dockify-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type User interface {
	Register(ctx context.Context, request entity.UserRegisterRequest) (int, error)
	Login(ctx context.Context, request entity.UserLoginRequest) (models.User, error)
}

type user struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) User {
	return &user{repo: repo}
}

func (u *user) Register(ctx context.Context, request entity.UserRegisterRequest) (int, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		return 0, err
	}
	userModel := models.User{
		Username:     request.Username,
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Email:        request.Email,
		PasswordHash: string(b),
	}

	return u.repo.User.CreateUser(ctx, userModel)
}

func (u *user) Login(ctx context.Context, request entity.UserLoginRequest) (models.User, error) {
	userModel, err := u.repo.User.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userModel.PasswordHash), []byte(request.Password))
	if err != nil {
		return models.User{}, err
	}

	return userModel, nil
}
