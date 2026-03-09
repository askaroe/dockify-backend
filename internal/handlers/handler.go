package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/askaroe/dockify-backend/internal/entity"
	"github.com/askaroe/dockify-backend/internal/gateway/mindspore"
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
	s      *services.Service
	logger *utils.Logger
}

func NewHandler(logger *utils.Logger, s *services.Service) *Handler {
	return &Handler{
		User:     user.NewUserHandler(s, logger),
		Health:   health.NewHealthHandler(s, logger),
		Location: location.NewLocationHandler(s, logger),
		s:        s,
		logger:   logger,
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
// @Description Returns a recommendation string based on user metrics
// @Tags Recommendation
// @Accept json
// @Produce json
// @Param user_id query integer true "User ID for fetching metrics"
// @Success 200 {object} entity.RecommendationResponse
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object} entity.ErrorMessage
// @Router /api/v1/recommendation [get]
func (h *Handler) GetRecommendation(c *gin.Context) {
	ctx := c.Request.Context()
	userIDParam := c.Query("user_id")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "user_id is required"})
		return
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid user_id"})
		return
	}

	metrics, err := h.s.Health.GetMetricsByUserId(ctx, userID)
	if err != nil {
		h.logger.Errorf("Failed to get metrics: %v", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to get metrics"})
		return
	}

	var sleepReq mindspore.PredictSleepRequest
	var lifeReq mindspore.PredictLifestyleRequest

	for _, m := range metrics {
		valFloat, _ := strconv.ParseFloat(m.MetricValue, 64)
		valInt := int(valFloat)

		switch m.MetricType {
		// Lifestyle
		case "AGE":
			lifeReq.Age = valInt
		case "WEIGHT":
			lifeReq.WeightKg = valInt
		case "HEIGHT_M":
			lifeReq.HeightM = valFloat
		case "BMI":
			lifeReq.Bmi = valFloat
		case "FAT_PERCENTAGE":
			lifeReq.FatPercentage = valFloat
		case "MAX_BPM":
			lifeReq.MaxBpm = valInt
		case "AVG_BPM":
			lifeReq.AvgBpm = valInt
		case "resting_bpm":
			lifeReq.RestingBpm = valInt
		case "SESSION_DURATION_HOURS":
			lifeReq.SessionDurationHours = valFloat
		case "CALORIES":
			lifeReq.CaloriesBurned = valInt
		case "WORKOUT_FREQUENCY":
			lifeReq.WorkoutFrequency = valInt
		case "DAILY_CALORIES":
			lifeReq.DailyCalories = valInt
		case "WATER_INTAKE_LITERS":
			lifeReq.WaterIntakeLiters = valFloat

		// Sleep
		case "SLEEP_DURATION":
			sleepReq.SleepDurationHours = valFloat
		case "TIME_IN_BED_HOURS":
			sleepReq.TimeInBedHours = valFloat
		case "HEART_RATE":
			sleepReq.HeartRate = valInt
		case "SLEEP_EFFICIENCY":
			sleepReq.SleepEfficiency = valFloat
		case "MOVEMENTS_PER_HOUR":
			sleepReq.MovementsPerHour = valFloat
		case "SNORE_TIME":
			sleepReq.SnoreTime = valInt
		case "DAY_OF_WEEK":
			sleepReq.DayOfWeek = valInt
		case "HOUR_STARTED":
			sleepReq.HourStarted = valInt
		case "NOTE_COFFEE":
			sleepReq.NoteCoffee = valInt
		case "NOTE_TEA":
			sleepReq.NoteTea = valInt
		case "NOTE_WORKOUT":
			sleepReq.NoteWorkout = valInt
		case "NOTE_STRESS":
			sleepReq.NoteStress = valInt
		case "NOTE_ATE_LATE":
			sleepReq.NoteAteLate = valInt
		}
	}

	sleepRes, err1 := h.s.Gateway.MindSpore.PredictSleep(ctx, sleepReq)
	lifeRes, err2 := h.s.Gateway.MindSpore.PredictLifestyle(ctx, lifeReq)

	if err1 != nil {
		h.logger.Errorf("MindSpore PredictSleep error: %v", err1)
	}
	if err2 != nil {
		h.logger.Errorf("MindSpore PredictLifestyle error: %v", err2)
	}

	var recommendation string
	if err1 == nil && err2 == nil {
		recommendation = fmt.Sprintf("%s %s", lifeRes.Interpretation, sleepRes.Interpretation)
	} else {
		// Fallback recommendation if the ML service is down
		recommendation = "Stay hydrated and take regular breaks during work!"
	}

	c.JSON(http.StatusOK, entity.RecommendationResponse{Recommendation: recommendation})
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
