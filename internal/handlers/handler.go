package handlers

import (
	"net/http"

	"github.com/askaroe/dockify-backend/internal/handlers/health"
	"github.com/askaroe/dockify-backend/internal/handlers/location"
	"github.com/askaroe/dockify-backend/internal/handlers/user"
	"github.com/askaroe/dockify-backend/internal/services"
	"github.com/askaroe/dockify-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	user.User
	health.Health
	location.Location
}

func NewHandler(logger *utils.Logger, s *services.Service) *Handler {
	return &Handler{
		User:     user.NewUserHandler(s, logger),
		Health:   health.NewHealthHandler(s, logger),
		Location: location.NewLocationHandler(s, logger),
	}
}

// HealthCheck godoc
// @Summary Health Check (Live)
// @Description Returns the live status of the service
// @Tags Health
// @Produce json
// @Success 200 {string} string "health"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "health")
}

// GetRecommendation godoc
// @Summary Get Recommendation
// @Description Returns a recommendation string
// @Tags Recommendation
// @Produce json
// @Success 200 {string} string "recommendation"
// @Router /api/v1/recommendation [get]
func GetRecommendation(c *gin.Context) {
	c.JSON(http.StatusOK, "recommendation")
}
