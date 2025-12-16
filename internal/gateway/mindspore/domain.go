package mindspore

const (
	PredictSleepEndpoint     = "/predict/sleep"
	PredictLifestyleEndpoint = "/predict/lifestyle"
)

type PredictSleepRequest struct {
	SleepDurationHours float64 `json:"sleep_duration_hours"`
	TimeInBedHours     float64 `json:"time_in_bed_hours"`
	HeartRate          int     `json:"heart_rate"`
	SleepEfficiency    float64 `json:"sleep_efficiency"`
	MovementsPerHour   float64 `json:"movements_per_hour"`
	SnoreTime          int     `json:"snore_time"`
	DayOfWeek          int     `json:"day_of_week"`
	HourStarted        int     `json:"hour_started"`
	NoteCoffee         int     `json:"note_coffee"`
	NoteTea            int     `json:"note_tea"`
	NoteWorkout        int     `json:"note_workout"`
	NoteStress         int     `json:"note_stress"`
	NoteAteLate        int     `json:"note_ate_late"`
}

type PredictSleepResponse struct {
	SleepQualityScore float64 `json:"sleep_quality_score"`
	SleepStage        string  `json:"sleep_stage"`
	Interpretation    string  `json:"interpretation"`
}

type PredictLifestyleRequest struct {
	Age                  int     `json:"age"`
	WeightKg             int     `json:"weight_kg"`
	HeightM              float64 `json:"height_m"`
	Bmi                  float64 `json:"bmi"`
	FatPercentage        float64 `json:"fat_percentage"`
	MaxBpm               int     `json:"max_bpm"`
	AvgBpm               int     `json:"avg_bpm"`
	RestingBpm           int     `json:"resting_bpm"`
	SessionDurationHours float64 `json:"session_duration_hours"`
	CaloriesBurned       int     `json:"calories_burned"`
	WorkoutFrequency     int     `json:"workout_frequency"`
	DailyCalories        int     `json:"daily_calories"`
	WaterIntakeLiters    float64 `json:"water_intake_liters"`
}

type PredictLifestyleResponse struct {
	LifestyleCategory string  `json:"lifestyle_category"`
	NextDayCalories   float64 `json:"next_day_calories"`
	HealthRiskScore   float64 `json:"health_risk_score"`
	Interpretation    string  `json:"interpretation"`
}
