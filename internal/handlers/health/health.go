package health

import (
	"net/http"
	"strconv"

	"github.com/askaroe/dockify-backend/internal/entity"
	"github.com/askaroe/dockify-backend/internal/services"
	"github.com/askaroe/dockify-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Health interface {
	CreateHealthMetrics(c *gin.Context)
	GetHealthMetrics(c *gin.Context)
}

type health struct {
	s      *services.Service
	logger *utils.Logger
}

func NewHealthHandler(s *services.Service, logger *utils.Logger) Health {
	return &health{s: s, logger: logger}
}

// CreateHealthMetrics
// @Summary Create health metrics
// @Description Create health metrics for a user
// @Tags Metrics
// @Accept json
// @Produce json
// @Param request body entity.HealthMetricsRequest true "Health metrics payload"
// @Success 201 {object} map[string]string "status message"
// @Failure 400 {object} entity.ErrorMessage "invalid request"
// @Failure 500 {object} entity.ErrorMessage "failed to create health metrics"
// @Router /api/v1/metrics [post]
func (h *health) CreateHealthMetrics(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.HealthMetricsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid request"})
		return
	}

	err := h.s.Health.CreateHealthMetric(ctx, req)
	if err != nil {
		h.logger.Errorf("CreateHealthMetrics error: %v", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to create health metrics"})
		return
	}

	err = h.s.Location.CreateLocation(ctx, req)
	if err != nil {
		h.logger.Errorf("CreateLocation error: %v", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to create location"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "health metrics created"})
}

// GetHealthMetrics
// @Summary Get health metrics
// @Description Retrieve health metrics for a given user by query parameter user_id
// @Tags Metrics
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {array} object "list of health metrics"
// @Failure 400 {object} entity.ErrorMessage "invalid user_id or missing parameter"
// @Failure 500 {object} entity.ErrorMessage "failed to get health metrics"
// @Router /api/v1/metrics [get]
func (h *health) GetHealthMetrics(c *gin.Context) {
	ctx := c.Request.Context()

	userIdParam := c.Query(entity.RequestParamUserID)
	if userIdParam == "" {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid request"})
		return
	}

	// Convert userIdParam to int
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		h.logger.Errorf("GetHealthMetrics error: %v", err)
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid request"})
	}

	metrics, err := h.s.Health.GetMetricsByUserId(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to get health metrics"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}
