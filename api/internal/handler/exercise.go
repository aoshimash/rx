package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/aoshimash/optel-training/api/internal/middleware"
	"github.com/aoshimash/optel-training/api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ExerciseHandler handles Exercise-related HTTP requests
type ExerciseHandler struct {
	repo repository.ExerciseRepository
}

// NewExerciseHandler creates a new ExerciseHandler
func NewExerciseHandler(repo repository.ExerciseRepository) *ExerciseHandler {
	return &ExerciseHandler{repo: repo}
}

// CreateExercise handles POST /exercises
func (h *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request body (OpenAPI ExerciseCreate type)
	var req struct {
		Name          string   `json:"name"`
		Description   *string  `json:"description,omitempty"`
		Aliases       []string `json:"aliases,omitempty"`
		MuscleGroups  []string `json:"muscle_groups,omitempty"`
		LoadIncrement *float64 `json:"load_increment,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Convert to domain model
	exercise := &domain.Exercise{
		Name:          req.Name,
		Description:   req.Description,
		Aliases:       req.Aliases,
		MuscleGroups:  req.MuscleGroups,
		LoadIncrement: req.LoadIncrement,
	}

	// Validate
	if err := domain.ValidateExercise(exercise); err != nil {
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
	if err := h.repo.Create(ctx, exercise); err != nil {
		middleware.WriteInternalError(w, "Failed to create exercise")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exercise)
}

// GetExercise handles GET /exercises/{id}
func (h *ExerciseHandler) GetExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from path
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteValidationError(w, "Invalid exercise ID format", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	// Get from repository
	exercise, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Exercise not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve exercise")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exercise)
}

// UpdateExercise handles PUT /exercises/{id}
func (h *ExerciseHandler) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from path
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteValidationError(w, "Invalid exercise ID format", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	// Check if exercise exists
	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Exercise not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve exercise")
		return
	}

	// Decode request body (full replacement)
	var req struct {
		Name          string   `json:"name"`
		Description   *string  `json:"description,omitempty"`
		Aliases       []string `json:"aliases,omitempty"`
		MuscleGroups  []string `json:"muscle_groups,omitempty"`
		LoadIncrement *float64 `json:"load_increment,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Update existing exercise (full replacement)
	existing.Name = req.Name
	existing.Description = req.Description
	existing.Aliases = req.Aliases
	existing.MuscleGroups = req.MuscleGroups
	existing.LoadIncrement = req.LoadIncrement

	// Validate
	if err := domain.ValidateExercise(existing); err != nil {
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
			middleware.WriteNotFoundError(w, "Exercise not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update exercise")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existing)
}
