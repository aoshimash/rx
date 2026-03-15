package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

// telemetryStore implements TelemetryPointRepository with in-memory map storage
type telemetryStore struct {
	mu     sync.RWMutex
	points map[uuid.UUID]*domain.TelemetryPoint
}

// NewTelemetryPointRepository creates a new in-memory TelemetryPoint repository
func NewTelemetryPointRepository() repository.TelemetryPointRepository {
	return &telemetryStore{
		points: make(map[uuid.UUID]*domain.TelemetryPoint),
	}
}

func (s *telemetryStore) Create(ctx context.Context, point *domain.TelemetryPoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	point.ID = uuid.New()
	point.CreatedAt = now

	s.points[point.ID] = point
	return nil
}

func (s *telemetryStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.TelemetryPoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	point, exists := s.points[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	// Return a copy to prevent external modifications
	result := *point
	return &result, nil
}

func (s *telemetryStore) Update(ctx context.Context, point *domain.TelemetryPoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.points[point.ID]; !exists {
		return domain.ErrNotFound
	}

	s.points[point.ID] = point
	return nil
}

func (s *telemetryStore) Delete(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.points[id]; !exists {
		return domain.ErrNotFound
	}

	delete(s.points, id)
	return nil
}

func (s *telemetryStore) List(ctx context.Context, limit int, after string) ([]*domain.TelemetryPoint, string, bool, error) {
	return s.listWithFilter(ctx, "", nil, nil, limit, after)
}

func (s *telemetryStore) ListByMetricAndTimeRange(ctx context.Context, metricName string, timestampFrom, timestampTo *time.Time, limit int, after string) ([]*domain.TelemetryPoint, string, bool, error) {
	return s.listWithFilter(ctx, metricName, timestampFrom, timestampTo, limit, after)
}

func (s *telemetryStore) listWithFilter(ctx context.Context, metricName string, timestampFrom, timestampTo *time.Time, limit int, after string) ([]*domain.TelemetryPoint, string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Convert map to slice and filter
	points := make([]*domain.TelemetryPoint, 0, len(s.points))
	for _, p := range s.points {
		// Filter by metric name if provided
		if metricName != "" && p.MetricName != metricName {
			continue
		}
		// Filter by timestamp range if provided
		if timestampFrom != nil && p.Timestamp.Before(*timestampFrom) {
			continue
		}
		if timestampTo != nil && !p.Timestamp.Before(*timestampTo) {
			continue
		}
		points = append(points, p)
	}

	// Sort by timestamp (descending), then by ID for consistent ordering
	sort.Slice(points, func(i, j int) bool {
		if !points[i].Timestamp.Equal(points[j].Timestamp) {
			return points[i].Timestamp.After(points[j].Timestamp)
		}
		return points[i].ID.String() < points[j].ID.String()
	})

	// Decode cursor if provided
	var startIdx int
	if after != "" {
		cursorID, err := decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		// Find the index after the cursor
		for i, p := range points {
			if p.ID == cursorID {
				startIdx = i + 1
				break
			}
		}
	}

	// Apply limit
	endIdx := startIdx + limit
	if endIdx > len(points) {
		endIdx = len(points)
	}

	result := points[startIdx:endIdx]
	hasMore := endIdx < len(points)

	var nextCursor string
	if hasMore && len(result) > 0 {
		nextCursor = encodeCursor(result[len(result)-1].ID)
	}

	// Return copies to prevent external modifications
	copies := make([]*domain.TelemetryPoint, len(result))
	for i, p := range result {
		copy := *p
		copies[i] = &copy
	}

	return copies, nextCursor, hasMore, nil
}
