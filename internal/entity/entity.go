package entity

const (
	RequestParamUserID = "user_id"
)

type UserRegisterRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type HealthMetricsRequest struct {
	UserId  int            `json:"user_id"`
	Metrics []HealthMetric `json:"metrics"`
}

type HealthMetric struct {
	MetricType  string `json:"metric_type"`
	MetricValue string `json:"metric_value"`
}
