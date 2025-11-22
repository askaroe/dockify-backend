package location

import (
	"context"

	"github.com/askaroe/dockify-backend/internal/models"
	"github.com/askaroe/dockify-backend/pkg/psql"
)

type Location interface {
	Insert(ctx context.Context, req models.Location) error
	GetNearestUsers(ctx context.Context, latitude, longitude float64, radius int) ([]models.Location, error)
}

type location struct {
	db *psql.Client
}

func NewLocationRepository(db *psql.Client) Location {
	return &location{db: db}
}

func (l *location) Insert(ctx context.Context, req models.Location) error {
	query := `INSERT INTO locations (user_id, latitude, longitude) VALUES ($1, $2, $3)`
	_, err := l.db.Exec(ctx, query, req.UserId, req.Latitude, req.Longitude)
	return err
}

func (l *location) GetNearestUsers(ctx context.Context, latitude, longitude float64, radius int) ([]models.Location, error) {
	query := `SELECT DISTINCT ON (user_id) user_id, latitude, longitude, recorded_at
	FROM (
	  SELECT user_id, latitude, longitude, recorded_at,
	    2 * 6371000 * ASIN(SQRT(
	      POWER(SIN(RADIANS($1 - latitude) / 2), 2) +
	      COS(RADIANS($1)) * COS(RADIANS(latitude)) *
	      POWER(SIN(RADIANS($2 - longitude) / 2), 2)
	    )) AS distance
	  FROM locations
	) AS l
	WHERE distance <= $3
	ORDER BY user_id, distance ASC`
	rows, err := l.db.Query(ctx, query, latitude, longitude, radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var loc models.Location
		if err := rows.Scan(&loc.UserId, &loc.Latitude, &loc.Longitude, &loc.RecordedAt); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}

	return locations, nil
}
