package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID           int        `json:"id"`
	Username     string     `json:"username"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	CreatedAt    *time.Time `json:"created_at"`
}

type HealthMetrics struct {
	ID          int        `json:"id"`
	UserId      int        `json:"user_id"`
	MetricType  string     `json:"metric_type"`
	MetricValue string     `json:"metric_value"`
	RecordedAt  *time.Time `json:"recorded_at"`
}

type Location struct {
	ID         int             `json:"id"`
	UserId     int             `json:"user_id"`
	Latitude   decimal.Decimal `json:"latitude"`
	Longitude  decimal.Decimal `json:"longitude"`
	RecordedAt *time.Time      `json:"recorded_at"`
}
