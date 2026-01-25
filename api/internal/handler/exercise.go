package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/aoshimash/optel-training/api/internal/middleware"
	"github.com/aoshimash/optel-training/api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ExerciseHandler handles Exercise-related HTTP requests
type ExerciseHandler struct {
	repo        repository.ExerciseRepository
	workoutRepo repository.WorkoutRepository
}

// NewExerciseHandler creates a new ExerciseHandler
func NewExerciseHandler(repo repository.ExerciseRepository, workoutRepo repository.WorkoutRepository) *ExerciseHandler {
	return &ExerciseHandler{
		repo:        repo,
		workoutRepo: workoutRepo,
	}
}

// CreateExercise handles POST /exercises
func (h *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.Info("CreateExercise request", "method", r.Method, "path", r.URL.Path)

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
		slog.Error("Failed to create exercise", "error", err)
		middleware.WriteInternalError(w, "Failed to create exercise")
		return
	}

	slog.Info("Exercise created", "id", exercise.ID)
	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(exercise); err != nil {
		slog.Error("Failed to encode exercise response", "error", err)
	}
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
	if err := json.NewEncoder(w).Encode(exercise); err != nil {
		slog.Error("Failed to encode exercise response", "error", err)
	}
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
	if err := json.NewEncoder(w).Encode(existing); err != nil {
		slog.Error("Failed to encode exercise response", "error", err)
	}
}

// DeleteExercise handles DELETE /exercises/{id}
func (h *ExerciseHandler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
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
	if _, err := h.repo.GetByID(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Exercise not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve exercise")
		return
	}

	// Check for referential integrity: Exercise referenced by WorkoutEntry (FR-026)
	referencingWorkouts, err := h.workoutRepo.ListByExerciseID(ctx, id)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to check references")
		return
	}
	if len(referencingWorkouts) > 0 {
		middleware.WriteConflictError(w, "Cannot delete exercise - referenced by workout entries", map[string]interface{}{
			"blocking_references": []map[string]interface{}{
				{
					"type":  "workout_entry",
					"count": len(referencingWorkouts),
				},
			},
		})
		return
	}

	// Delete
	if err := h.repo.Delete(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Exercise not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete exercise")
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// ListExercises handles GET /exercises
func (h *ExerciseHandler) ListExercises(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse pagination parameters
	limit := 100 // default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr, 1, 100); err == nil {
			limit = parsedLimit
		}
	}
	after := r.URL.Query().Get("after")

	// List from repository
	exercises, nextCursor, hasMore, err := h.repo.List(ctx, limit, after)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to list exercises")
		return
	}

	// Return paginated response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"data":        exercises,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	}); err != nil {
		slog.Error("Failed to encode exercise list response", "error", err)
	}
}
