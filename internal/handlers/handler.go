package handlers

import (
	"github.com/askaroe/dockify-backend/internal/services"
	"github.com/askaroe/dockify-backend/pkg/utils"
	"net/http"

	"github.com/askaroe/dockify-backend/config"
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler(logger *utils.Logger, s *services.Service, cfg *config.Config) *Handler {
	return &Handler{}
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
