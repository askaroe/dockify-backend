package health

import (
	"context"
	"fmt"

	"github.com/askaroe/dockify-backend/internal/entity"
	"github.com/askaroe/dockify-backend/internal/models"
	"github.com/askaroe/dockify-backend/internal/repository"
)

type Health interface {
	GetMetricsByUserId(ctx context.Context, id int) (models.HealthMetrics, error)
	CreateHealthMetric(ctx context.Context, req entity.HealthMetricsRequest) error
}

type health struct {
	repo *repository.Repository
}

func NewHealthService(repo *repository.Repository) Health {
	return &health{repo: repo}
}

func (h *health) GetMetricsByUserId(ctx context.Context, id int) (models.HealthMetrics, error) {
	return h.repo.Health.GetMetricsByUserId(ctx, id)
}

func (h *health) CreateHealthMetric(ctx context.Context, req entity.HealthMetricsRequest) error {
	var metricsModel []models.HealthMetrics

	for _, metric := range req.Metrics {
		metricsModel = append(metricsModel, models.HealthMetrics{
			UserId:      req.UserId,
			MetricType:  metric.MetricType,
			MetricValue: metric.MetricValue,
		})
	}

	err := h.repo.Health.CreateHealthMetrics(ctx, metricsModel)
	if err != nil {
		return fmt.Errorf("create health metrics: %w", err)
	}

	return nil
}
