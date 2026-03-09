package health

import (
	"context"

	"github.com/askaroe/dockify-backend/internal/models"
	"github.com/askaroe/dockify-backend/pkg/psql"
)

type Health interface {
	GetMetricsByUserId(ctx context.Context, id int) ([]models.HealthMetrics, error)
	CreateHealthMetric(ctx context.Context, req models.HealthMetrics) (int, error)
	CreateHealthMetrics(ctx context.Context, req []models.HealthMetrics) error
}

type health struct {
	db *psql.Client
}

func NewHealthRepository(db *psql.Client) Health {
	return &health{db: db}
}

func (h *health) GetMetricsByUserId(ctx context.Context, id int) ([]models.HealthMetrics, error) {
	query := `SELECT id, user_id, metric_type, metric_value, recorded_at FROM health_metrics WHERE user_id = $1`
	rows, err := h.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.HealthMetrics
	for rows.Next() {
		var metric models.HealthMetrics
		err := rows.Scan(&metric.ID, &metric.UserId, &metric.MetricType, &metric.MetricValue, &metric.RecordedAt)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return metrics, nil
}

func (h *health) CreateHealthMetric(ctx context.Context, req models.HealthMetrics) (int, error) {
	query := `INSERT INTO health_metrics (user_id, metric_type, metric_value) VALUES ($1, $2, $3) RETURNING id`
	err := h.db.QueryRow(ctx, query, req.UserId, req.MetricType, req.MetricValue).Scan(&req.ID)
	if err != nil {
		return 0, err
	}
	return req.ID, nil
}

func (h *health) CreateHealthMetrics(ctx context.Context, req []models.HealthMetrics) error {
	query := `INSERT INTO health_metrics (user_id, metric_type, metric_value) VALUES ($1, $2, $3)`
	for _, metric := range req {
		_, err := h.db.Exec(ctx, query, metric.UserId, metric.MetricType, metric.MetricValue)
		if err != nil {
			return err
		}
	}
	return nil

}
