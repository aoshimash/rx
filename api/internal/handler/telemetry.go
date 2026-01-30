package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aoshimash/optel-workout/api/internal/domain"
	"github.com/aoshimash/optel-workout/api/internal/middleware"
	"github.com/aoshimash/optel-workout/api/internal/repository"
	"github.com/google/uuid"
)

// TelemetryHandler handles TelemetryPoint-related HTTP requests
type TelemetryHandler struct {
	repo        repository.TelemetryPointRepository
	workoutRepo repository.WorkoutRepository
}

// NewTelemetryHandler creates a new TelemetryHandler
func NewTelemetryHandler(repo repository.TelemetryPointRepository, workoutRepo repository.WorkoutRepository) *TelemetryHandler {
	return &TelemetryHandler{
		repo:        repo,
		workoutRepo: workoutRepo,
	}
}

// CreateTelemetryPoint handles POST /telemetry
func (h *TelemetryHandler) CreateTelemetryPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request body (OpenAPI TelemetryPointCreate type)
	var req struct {
		Timestamp  string  `json:"timestamp"`
		MetricName string  `json:"metric_name"`
		Value      float64 `json:"value"`
		Unit       string  `json:"unit"`
		WorkoutID  *string `json:"workout_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		middleware.WriteValidationError(w, "Invalid timestamp format", map[string]interface{}{
			"field": "timestamp",
			"error": err.Error(),
		})
		return
	}

	// Convert to domain model
	point := &domain.TelemetryPoint{
		Timestamp:  timestamp,
		MetricName: req.MetricName,
		Value:      req.Value,
		Unit:       req.Unit,
	}

	// Parse workout_id if provided and validate it exists
	if req.WorkoutID != nil {
		workoutID, err := uuid.Parse(*req.WorkoutID)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid workout_id format", map[string]interface{}{
				"field": "workout_id",
			})
			return
		}
		// Validate Workout exists
		if _, err := h.workoutRepo.GetByID(ctx, workoutID); err != nil {
			if err == domain.ErrNotFound {
				middleware.WriteValidationError(w, "Workout not found", map[string]interface{}{
					"field":      "workout_id",
					"workout_id": *req.WorkoutID,
				})
				return
			}
			middleware.WriteInternalError(w, "Failed to validate workout")
			return
		}
		point.WorkoutID = &workoutID
	}

	// Validate
	if err := domain.ValidateTelemetryPoint(point); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Create
	if err := h.repo.Create(ctx, point); err != nil {
		middleware.WriteInternalError(w, "Failed to create telemetry point")
		return
	}

	// Return response
	writeJSON(w, http.StatusCreated, point)
}

// GetTelemetryPoint handles GET /telemetry/{id}
func (h *TelemetryHandler) GetTelemetryPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract and parse ID from path
	id, err := parseUUIDParam(r, "id", "telemetry point")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	// Get from repository
	point, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Telemetry point not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve telemetry point")
		return
	}

	// Return response
	writeJSON(w, http.StatusOK, point)
}

// UpdateTelemetryPoint handles PUT /telemetry/{id}
func (h *TelemetryHandler) UpdateTelemetryPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract and parse ID from path
	id, err := parseUUIDParam(r, "id", "telemetry point")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	// Check if telemetry point exists
	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Telemetry point not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve telemetry point")
		return
	}

	// Decode request body (full replacement)
	var req struct {
		Timestamp  string  `json:"timestamp"`
		MetricName string  `json:"metric_name"`
		Value      float64 `json:"value"`
		Unit       string  `json:"unit"`
		WorkoutID  *string `json:"workout_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		middleware.WriteValidationError(w, "Invalid timestamp format", map[string]interface{}{
			"field": "timestamp",
			"error": err.Error(),
		})
		return
	}

	// Update existing telemetry point (full replacement)
	existing.Timestamp = timestamp
	existing.MetricName = req.MetricName
	existing.Value = req.Value
	existing.Unit = req.Unit

	// Parse workout_id if provided and validate it exists
	if req.WorkoutID != nil {
		workoutID, err := uuid.Parse(*req.WorkoutID)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid workout_id format", map[string]interface{}{
				"field": "workout_id",
			})
			return
		}
		// Validate Workout exists
		if _, err := h.workoutRepo.GetByID(ctx, workoutID); err != nil {
			if err == domain.ErrNotFound {
				middleware.WriteValidationError(w, "Workout not found", map[string]interface{}{
					"field":      "workout_id",
					"workout_id": *req.WorkoutID,
				})
				return
			}
			middleware.WriteInternalError(w, "Failed to validate workout")
			return
		}
		existing.WorkoutID = &workoutID
	} else {
		existing.WorkoutID = nil
	}

	// Validate
	if err := domain.ValidateTelemetryPoint(existing); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Update
	if err := h.repo.Update(ctx, existing); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Telemetry point not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update telemetry point")
		return
	}

	// Return response
	writeJSON(w, http.StatusOK, existing)
}

// DeleteTelemetryPoint handles DELETE /telemetry/{id}
func (h *TelemetryHandler) DeleteTelemetryPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract and parse ID from path
	id, err := parseUUIDParam(r, "id", "telemetry point")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	// Check if telemetry point exists
	if _, err := h.repo.GetByID(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Telemetry point not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve telemetry point")
		return
	}

	// Delete (no dependencies per FR-026)
	if err := h.repo.Delete(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Telemetry point not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete telemetry point")
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// ListTelemetryPoints handles GET /telemetry
func (h *TelemetryHandler) ListTelemetryPoints(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse pagination parameters
	limit := 100 // default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr, 1, 100); err == nil {
			limit = parsedLimit
		}
	}
	after := r.URL.Query().Get("after")

	// Parse filter parameters (FR-030)
	metricName := r.URL.Query().Get("metric_name")
	var timestampFrom, timestampTo *time.Time
	if fromStr := r.URL.Query().Get("timestamp_from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			timestampFrom = &t
		}
	}
	if toStr := r.URL.Query().Get("timestamp_to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			timestampTo = &t
		}
	}

	// List from repository (with filters if provided)
	var points []*domain.TelemetryPoint
	var nextCursor string
	var hasMore bool
	var err error

	if metricName != "" || timestampFrom != nil || timestampTo != nil {
		points, nextCursor, hasMore, err = h.repo.ListByMetricAndTimeRange(ctx, metricName, timestampFrom, timestampTo, limit, after)
	} else {
		points, nextCursor, hasMore, err = h.repo.List(ctx, limit, after)
	}

	if err != nil {
		middleware.WriteInternalError(w, "Failed to list telemetry points")
		return
	}

	// Return paginated response
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":        points,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}
