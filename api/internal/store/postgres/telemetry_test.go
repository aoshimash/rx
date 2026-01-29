package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/google/uuid"
)

func TestTelemetryPointRepository_Create(t *testing.T) {
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewTelemetryPointRepository(pool)

	point := &domain.TelemetryPoint{
		Timestamp:  time.Now(),
		MetricName: "heart_rate",
		Value:      72.0,
		Unit:       "bpm",
	}

	err = repo.Create(ctx, point)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if point.ID == uuid.Nil {
		t.Error("Create() did not set ID")
	}
	if point.CreatedAt.IsZero() {
		t.Error("Create() did not set CreatedAt")
	}
}

func TestTelemetryPointRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewTelemetryPointRepository(pool)

	point := &domain.TelemetryPoint{
		Timestamp:  time.Now(),
		MetricName: "heart_rate",
		Value:      72.0,
		Unit:       "bpm",
	}
	if err := repo.Create(ctx, point); err != nil {
		t.Fatalf("Failed to create telemetry point: %v", err)
	}

	got, err := repo.GetByID(ctx, point.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.ID != point.ID {
		t.Errorf("GetByID() ID = %v, want %v", got.ID, point.ID)
	}
	if got.MetricName != point.MetricName {
		t.Errorf("GetByID() MetricName = %v, want %v", got.MetricName, point.MetricName)
	}
}

func TestTelemetryPointRepository_Update(t *testing.T) {
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewTelemetryPointRepository(pool)

	point := &domain.TelemetryPoint{
		Timestamp:  time.Now(),
		MetricName: "heart_rate",
		Value:      72.0,
		Unit:       "bpm",
	}
	if err := repo.Create(ctx, point); err != nil {
		t.Fatalf("Failed to create telemetry point: %v", err)
	}

	point.Value = 75.0
	if err := repo.Update(ctx, point); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	got, err := repo.GetByID(ctx, point.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.Value != 75.0 {
		t.Errorf("GetByID() Value = %v, want 75.0", got.Value)
	}
}

func TestTelemetryPointRepository_Delete(t *testing.T) {
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewTelemetryPointRepository(pool)

	point := &domain.TelemetryPoint{
		Timestamp:  time.Now(),
		MetricName: "heart_rate",
		Value:      72.0,
		Unit:       "bpm",
	}
	if err := repo.Create(ctx, point); err != nil {
		t.Fatalf("Failed to create telemetry point: %v", err)
	}

	if err := repo.Delete(ctx, point.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.GetByID(ctx, point.ID)
	if err != domain.ErrNotFound {
		t.Errorf("GetByID() after Delete() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestTelemetryPointRepository_List(t *testing.T) {
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewTelemetryPointRepository(pool)

	// Create multiple telemetry points
	for i := 0; i < 5; i++ {
		point := &domain.TelemetryPoint{
			Timestamp:  time.Now().Add(time.Duration(i) * time.Minute),
			MetricName: "heart_rate",
			Value:      float64(70 + i),
			Unit:       "bpm",
		}
		if err := repo.Create(ctx, point); err != nil {
			t.Fatalf("Failed to create telemetry point: %v", err)
		}
	}

	points, _, hasMore, err := repo.List(ctx, 3, "")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(points) != 3 {
		t.Errorf("List() returned %d points, want 3", len(points))
	}
	if !hasMore {
		t.Error("List() hasMore = false, want true")
	}
}

func TestTelemetryPointRepository_ListByMetricAndTimeRange(t *testing.T) {
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewTelemetryPointRepository(pool)

	now := time.Now()
	// Create points with different metrics and timestamps
	for i := 0; i < 5; i++ {
		point := &domain.TelemetryPoint{
			Timestamp:  now.Add(time.Duration(i) * time.Hour),
			MetricName: "heart_rate",
			Value:      float64(70 + i),
			Unit:       "bpm",
		}
		if err := repo.Create(ctx, point); err != nil {
			t.Fatalf("Failed to create telemetry point: %v", err)
		}
	}

	// Create points with different metric
	for i := 0; i < 3; i++ {
		point := &domain.TelemetryPoint{
			Timestamp:  now.Add(time.Duration(i) * time.Hour),
			MetricName: "temperature",
			Value:      float64(36 + i),
			Unit:       "celsius",
		}
		if err := repo.Create(ctx, point); err != nil {
			t.Fatalf("Failed to create telemetry point: %v", err)
		}
	}

	// Filter by metric and time range
	from := now.Add(-1 * time.Hour)
	to := now.Add(3 * time.Hour)
	points, _, _, err := repo.ListByMetricAndTimeRange(ctx, "heart_rate", &from, &to, 10, "")
	if err != nil {
		t.Fatalf("ListByMetricAndTimeRange() error = %v", err)
	}

	if len(points) == 0 {
		t.Error("ListByMetricAndTimeRange() returned 0 points, want at least 1")
	}

	// Verify all points have correct metric
	for _, p := range points {
		if p.MetricName != "heart_rate" {
			t.Errorf("ListByMetricAndTimeRange() returned point with metric %v, want heart_rate", p.MetricName)
		}
	}
}
