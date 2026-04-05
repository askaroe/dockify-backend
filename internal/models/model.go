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

type Document struct {
	ID          string     `json:"id"`
	UserId      int        `json:"user_id"`
	FileName    string     `json:"file_name"`
	FilePath    string     `json:"file_path"`
	FileSize    int64      `json:"file_size"`
	ContentType string     `json:"content_type"`
	Summary     string     `json:"summary"`
	UploadedAt  *time.Time `json:"uploaded_at"`
}

type ChatMessage struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	DocID     *string    `json:"doc_id"`
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	CreatedAt *time.Time `json:"created_at"`
}
