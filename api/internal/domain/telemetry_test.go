package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateTelemetryPoint(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)
	workoutID := uuid.New()

	tests := []struct {
		name    string
		point   *TelemetryPoint
		wantErr bool
		errCode string
	}{
		{
			name:    "nil telemetry point",
			point:   nil,
			wantErr: true,
		},
		{
			name: "valid telemetry point with minimal fields",
			point: &TelemetryPoint{
				ID:         uuid.New(),
				Timestamp:  past,
				MetricName: "daily_volume_kg",
				Value:      5000.0,
				Unit:       "kg",
				CreatedAt:  now,
			},
			wantErr: false,
		},
		{
			name: "valid telemetry point with workout link",
			point: &TelemetryPoint{
				ID:         uuid.New(),
				Timestamp:  past,
				MetricName: "session_duration_minutes",
				Value:      60.5,
				Unit:       "minutes",
				WorkoutID:  &workoutID,
				CreatedAt:  now,
			},
			wantErr: false,
		},
		{
			name: "empty metric_name",
			point: &TelemetryPoint{
				ID:         uuid.New(),
				Timestamp:  past,
				MetricName: "",
				Value:      100.0,
				Unit:       "kg",
				CreatedAt:  now,
			},
			wantErr: true,
			errCode: ErrCodeMissingRequiredField,
		},
		{
			name: "metric_name too long",
			point: &TelemetryPoint{
				ID:         uuid.New(),
				Timestamp:  past,
				MetricName: stringWithLength(101),
				Value:      100.0,
				Unit:       "kg",
				CreatedAt:  now,
			},
			wantErr: true,
		},
		{
			name: "empty unit",
			point: &TelemetryPoint{
				ID:         uuid.New(),
				Timestamp:  past,
				MetricName: "metric",
				Value:      100.0,
				Unit:       "",
				CreatedAt:  now,
			},
			wantErr: true,
			errCode: ErrCodeMissingRequiredField,
		},
		{
			name: "unit too long",
			point: &TelemetryPoint{
				ID:         uuid.New(),
				Timestamp:  past,
				MetricName: "metric",
				Value:      100.0,
				Unit:       stringWithLength(51),
				CreatedAt:  now,
			},
			wantErr: true,
		},
		{
			name: "timestamp in future",
			point: &TelemetryPoint{
				ID:         uuid.New(),
				Timestamp:  future,
				MetricName: "metric",
				Value:      100.0,
				Unit:       "kg",
				CreatedAt:  now,
			},
			wantErr: true,
			errCode: ErrCodeInvalidTimestamp,
		},
		{
			name: "valid with different metric types",
			point: &TelemetryPoint{
				ID:         uuid.New(),
				Timestamp:  past,
				MetricName: "body_weight_kg",
				Value:      74.5,
				Unit:       "kg",
				CreatedAt:  now,
			},
			wantErr: false,
		},
		{
			name: "valid with negative value (allowed for metrics)",
			point: &TelemetryPoint{
				ID:         uuid.New(),
				Timestamp:  past,
				MetricName: "temperature_change",
				Value:      -2.5,
				Unit:       "celsius",
				CreatedAt:  now,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTelemetryPoint(tt.point)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTelemetryPoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errCode != "" {
				if domainErr, ok := err.(*DomainError); ok {
					if domainErr.Code != tt.errCode {
						t.Errorf("ValidateTelemetryPoint() error code = %v, want %v", domainErr.Code, tt.errCode)
					}
				}
			}
		})
	}
}
