package location

import (
	"context"
	"fmt"

	"github.com/askaroe/dockify-backend/internal/entity"
	"github.com/askaroe/dockify-backend/internal/models"
	"github.com/askaroe/dockify-backend/internal/repository"
)

type Location interface {
	CreateLocation(ctx context.Context, req entity.HealthMetricsRequest) error
	GetNearestUsers(ctx context.Context, request entity.NearestUsersRequest) ([]entity.NearestUsersResponse, error)
}

type location struct {
	repo *repository.Repository
}

func NewLocationService(repo *repository.Repository) Location {
	return &location{repo: repo}
}

func (l *location) CreateLocation(ctx context.Context, req entity.HealthMetricsRequest) error {
	locationRecord := models.Location{
		Latitude:  req.Location.Latitude,
		Longitude: req.Location.Longitude,
		UserId:    req.UserId,
	}

	err := l.repo.Location.Insert(ctx, locationRecord)
	if err != nil {
		return fmt.Errorf("create location: %w", err)
	}

	return nil
}

func (l *location) GetNearestUsers(ctx context.Context, request entity.NearestUsersRequest) ([]entity.NearestUsersResponse, error) {
	locations, err := l.repo.Location.GetNearestUsers(ctx, request.Latitude, request.Longitude, request.Radius)
	if err != nil {
		return nil, fmt.Errorf("get nearest users: %w", err)
	}

	var response []entity.NearestUsersResponse
	for _, loc := range locations {
		if loc.UserId != request.UserId {
			userLocation := entity.Location{
				Latitude:  loc.Latitude,
				Longitude: loc.Longitude,
			}
			response = append(response, entity.NearestUsersResponse{
				UserID:   loc.UserId,
				Location: userLocation,
			})
		}
	}

	return response, nil
}
