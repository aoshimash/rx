package memory

import (
	"context"
	"testing"
	"time"

	"github.com/aoshimash/optel-workout/api/internal/domain"
	"github.com/google/uuid"
)

func TestTelemetryPointRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		point   *domain.TelemetryPoint
		wantErr bool
	}{
		{
			name: "create valid telemetry point",
			point: &domain.TelemetryPoint{
				Timestamp:  time.Now(),
				MetricName: "heart_rate",
				Value:      72.0,
				Unit:       "bpm",
			},
			wantErr: false,
		},
		{
			name: "create telemetry point with workout_id",
			point: &domain.TelemetryPoint{
				Timestamp:  time.Now(),
				MetricName: "heart_rate",
				Value:      72.0,
				Unit:       "bpm",
				WorkoutID:  uuidPtr(uuid.New()),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTelemetryPointRepository()
			ctx := context.Background()

			err := repo.Create(ctx, tt.point)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.point.ID == uuid.Nil {
					t.Error("Create() did not generate ID")
				}
				if tt.point.CreatedAt.IsZero() {
					t.Error("Create() did not set CreatedAt")
				}
			}
		})
	}
}

func TestTelemetryPointRepository_GetByID(t *testing.T) {
	repo := NewTelemetryPointRepository().(*telemetryStore)
	ctx := context.Background()

	id := uuid.New()
	point := &domain.TelemetryPoint{
		ID:         id,
		Timestamp:  time.Now(),
		MetricName: "heart_rate",
		Value:      72.0,
		Unit:       "bpm",
		CreatedAt:  time.Now(),
	}
	repo.points[id] = point

	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if got.ID != id {
		t.Errorf("GetByID() got ID = %v, want %v", got.ID, id)
	}
	if got.MetricName != "heart_rate" {
		t.Errorf("GetByID() got MetricName = %v, want heart_rate", got.MetricName)
	}

	// Verify it's a copy
	originalValue := got.Value
	got.Value = 100.0
	recheck, _ := repo.GetByID(ctx, id)
	if recheck.Value != originalValue {
		t.Error("GetByID() did not return a copy")
	}
}

func TestTelemetryPointRepository_ListByMetricAndTimeRange(t *testing.T) {
	repo := NewTelemetryPointRepository().(*telemetryStore)
	ctx := context.Background()

	now := time.Now()
	point1 := &domain.TelemetryPoint{
		ID:         uuid.New(),
		Timestamp:  now.Add(-2 * time.Hour),
		MetricName: "heart_rate",
		Value:      70.0,
		Unit:       "bpm",
		CreatedAt:  time.Now(),
	}
	point2 := &domain.TelemetryPoint{
		ID:         uuid.New(),
		Timestamp:  now.Add(-1 * time.Hour),
		MetricName: "heart_rate",
		Value:      75.0,
		Unit:       "bpm",
		CreatedAt:  time.Now(),
	}
	point3 := &domain.TelemetryPoint{
		ID:         uuid.New(),
		Timestamp:  now,
		MetricName: "temperature",
		Value:      36.5,
		Unit:       "celsius",
		CreatedAt:  time.Now(),
	}
	repo.points[point1.ID] = point1
	repo.points[point2.ID] = point2
	repo.points[point3.ID] = point3

	from := now.Add(-90 * time.Minute)
	to := now.Add(-30 * time.Minute)

	points, _, _, err := repo.ListByMetricAndTimeRange(ctx, "heart_rate", &from, &to, 10, "")
	if err != nil {
		t.Fatalf("ListByMetricAndTimeRange() error = %v", err)
	}

	if len(points) != 1 {
		t.Errorf("ListByMetricAndTimeRange() got %d points, want 1", len(points))
	}
	if points[0].ID != point2.ID {
		t.Errorf("ListByMetricAndTimeRange() got point ID = %v, want %v", points[0].ID, point2.ID)
	}
}

// Helper function
func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}
