package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/aoshimash/optel-training/api/internal/middleware"
	"github.com/aoshimash/optel-training/api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TelemetryHandler handles TelemetryPoint-related HTTP requests
type TelemetryHandler struct {
	repo         repository.TelemetryPointRepository
	workoutRepo  repository.WorkoutRepository
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
					"field":     "workout_id",
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
		if ve, ok := err.(*domain.ValidationError); ok {
			middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
				"field":   ve.Field,
				"message": ve.Message,
			})
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(point)
}

// GetTelemetryPoint handles GET /telemetry/{id}
func (h *TelemetryHandler) GetTelemetryPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from path
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteValidationError(w, "Invalid telemetry point ID format", map[string]interface{}{
			"id": idStr,
		})
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(point)
}

// UpdateTelemetryPoint handles PUT /telemetry/{id}
func (h *TelemetryHandler) UpdateTelemetryPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from path
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteValidationError(w, "Invalid telemetry point ID format", map[string]interface{}{
			"id": idStr,
		})
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
					"field":     "workout_id",
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
		if ve, ok := err.(*domain.ValidationError); ok {
			middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
				"field":   ve.Field,
				"message": ve.Message,
			})
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existing)
}
