package handlers

import (
	"net/http"

	"github.com/askaroe/dockify-backend/internal/entity"
	"github.com/askaroe/dockify-backend/internal/handlers/health"
	"github.com/askaroe/dockify-backend/internal/handlers/location"
	"github.com/askaroe/dockify-backend/internal/handlers/user"
	"github.com/askaroe/dockify-backend/internal/services"
	"github.com/askaroe/dockify-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
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
// @Success 200 {object} entity.RecommendationResponse
// @Router /api/v1/recommendation [get]
func GetRecommendation(c *gin.Context) {
	c.JSON(http.StatusOK, entity.RecommendationResponse{Recommendation: "Stay hydrated and take regular breaks during work!"})
}

// GetNearestHospitals godoc
// @Summary Get Nearest Hospitals
// @Description Returns a list of nearest hospitals to the provided location
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param request body entity.NearestHospitalsRequest true "Nearest hospitals request"
// @Success 200 {array} entity.Location
// @Failure 400 {object} entity.ErrorMessage "invalid request"
// @Router /api/v1/hospitals/nearest [post]
func GetNearestHospitals(c *gin.Context) {
	var req entity.NearestHospitalsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid request"})
		return
	}

	type loc entity.Location

	round6 := func(f float64) decimal.Decimal {
		return decimal.NewFromFloat(f).Round(6)
	}

	locations := []loc{
		{Longitude: round6(76.851248), Latitude: round6(43.222015)},
		{Longitude: round6(76.851258), Latitude: round6(43.222025)},
		{Longitude: round6(76.851268), Latitude: round6(43.222035)},
	}

	c.JSON(http.StatusOK, locations)
}
