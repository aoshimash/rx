package domain

import (
	"github.com/google/uuid"
	"time"
)

// TelemetryPoint represents a time-series metric data point.
type TelemetryPoint struct {
	ID         uuid.UUID  `json:"id"`
	Timestamp  time.Time  `json:"timestamp"`
	MetricName string     `json:"metric_name"`
	Value      float64    `json:"value"`
	Unit       string     `json:"unit"`
	WorkoutID  *uuid.UUID `json:"workout_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}
