package entity

import "github.com/shopspring/decimal"

const (
	RequestParamUserID = "user_id"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

type UserRegisterRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type CreatedUserResponse struct {
	UserID int `json:"user_id"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type HealthMetricsRequest struct {
	UserId   int            `json:"user_id"`
	Location Location       `json:"location"`
	Metrics  []HealthMetric `json:"metrics"`
}

type HealthMetric struct {
	MetricType  string `json:"metric_type"`
	MetricValue string `json:"metric_value"`
}

type Location struct {
	Longitude decimal.Decimal `json:"longitude" example:"37.617396"`
	Latitude  decimal.Decimal `json:"latitude" example:"55.755825"`
}

type NearestUsersRequest struct {
	UserId    int     `json:"user_id"`
	Longitude float64 `json:"longitude" example:"37.617396"`
	Latitude  float64 `json:"latitude" example:"55.755825"`
	Radius    int     `json:"radius" example:"5000"` // in meters
}

type NearestHospitalsRequest struct {
	Longitude float64 `json:"longitude" example:"37.617396"`
	Latitude  float64 `json:"latitude" example:"55.755825"`
	Radius    int     `json:"radius" example:"5000"` // in meters
}

type NearestUsersResponse struct {
	UserID   int      `json:"user_id"`
	Location Location `json:"location"`
}

type RecommendationResponse struct {
	Recommendation string `json:"recommendation"`
}
