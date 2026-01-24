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

// WorkoutHandler handles Workout-related HTTP requests
type WorkoutHandler struct {
	repo            repository.WorkoutRepository
	exerciseRepo    repository.ExerciseRepository
	programRepo     repository.ProgramRepository
}

// NewWorkoutHandler creates a new WorkoutHandler
func NewWorkoutHandler(repo repository.WorkoutRepository, exerciseRepo repository.ExerciseRepository, programRepo repository.ProgramRepository) *WorkoutHandler {
	return &WorkoutHandler{
		repo:         repo,
		exerciseRepo: exerciseRepo,
		programRepo:  programRepo,
	}
}

// CreateWorkout handles POST /workouts
func (h *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request body (OpenAPI WorkoutCreate type)
	var req struct {
		Timestamp      string    `json:"timestamp"`
		SessionStart   *string   `json:"session_start,omitempty"`
		SessionEnd     *string   `json:"session_end,omitempty"`
		BodyWeightKg   *float64  `json:"body_weight_kg,omitempty"`
		FatigueLevel   *int      `json:"fatigue_level,omitempty"`
		SleepHours     *float64  `json:"sleep_hours,omitempty"`
		ConditionNotes *string   `json:"condition_notes,omitempty"`
		ProgramNodeID  *string   `json:"program_node_id,omitempty"`
		ProgramContext []string  `json:"program_context,omitempty"`
		Notes          *string   `json:"notes,omitempty"`
		Entries        []struct {
			ExerciseID           string    `json:"exercise_id"`
			DisplayName          *string   `json:"display_name,omitempty"`
			EntryType            string    `json:"entry_type"`
			Sets                 int       `json:"sets"`
			Reps                 int       `json:"reps"`
			LoadKg               float64   `json:"load_kg"`
			RPE                  int       `json:"rpe"`
			EntryStart           *string   `json:"entry_start,omitempty"`
			EntryEnd             *string   `json:"entry_end,omitempty"`
			PlannedRestSeconds   *int      `json:"planned_rest_seconds,omitempty"`
			PerformedRestSeconds *int      `json:"performed_rest_seconds,omitempty"`
			PerSetRestOverrides  []int     `json:"per_set_rest_overrides,omitempty"`
			ProgramNodeID        *string   `json:"program_node_id,omitempty"`
			PlanSnapshot         *struct {
				ProgramNodeID      *string  `json:"program_node_id,omitempty"`
				TargetSets         *int     `json:"target_sets,omitempty"`
				TargetReps         *int     `json:"target_reps,omitempty"`
				TargetRPE          *int     `json:"target_rpe,omitempty"`
				TargetLoadKg       *float64 `json:"target_load_kg,omitempty"`
				Percent1RM         *float64 `json:"percent_1rm,omitempty"`
				PlannedRestSeconds *int     `json:"planned_rest_seconds,omitempty"`
			} `json:"plan_snapshot,omitempty"`
			Notes *string `json:"notes,omitempty"`
		} `json:"entries"`
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
	workout := &domain.Workout{
		Timestamp:      timestamp,
		BodyWeightKg:   req.BodyWeightKg,
		FatigueLevel:   req.FatigueLevel,
		SleepHours:     req.SleepHours,
		ConditionNotes: req.ConditionNotes,
		ProgramContext: req.ProgramContext,
		Notes:          req.Notes,
		Entries:        make([]domain.WorkoutEntry, len(req.Entries)),
	}

	// Parse optional timestamps
	if req.SessionStart != nil {
		if t, err := time.Parse(time.RFC3339, *req.SessionStart); err == nil {
			workout.SessionStart = &t
		}
	}
	if req.SessionEnd != nil {
		if t, err := time.Parse(time.RFC3339, *req.SessionEnd); err == nil {
			workout.SessionEnd = &t
		}
	}

	// Parse program_node_id
	if req.ProgramNodeID != nil {
		if id, err := uuid.Parse(*req.ProgramNodeID); err == nil {
			workout.ProgramNodeID = &id
		}
	}

	// Convert entries
	for i, entryReq := range req.Entries {
		exerciseID, err := uuid.Parse(entryReq.ExerciseID)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid exercise_id format", map[string]interface{}{
				"field": "entries",
				"index": i,
			})
			return
		}

		// Validate Exercise exists (FR-026)
		if _, err := h.exerciseRepo.GetByID(ctx, exerciseID); err != nil {
			if err == domain.ErrNotFound {
				middleware.WriteValidationError(w, "Exercise not found", map[string]interface{}{
					"field":      "entries",
					"index":      i,
					"exercise_id": entryReq.ExerciseID,
				})
				return
			}
			middleware.WriteInternalError(w, "Failed to validate exercise")
			return
		}

		entry := domain.WorkoutEntry{
			ExerciseID:          exerciseID,
			DisplayName:         entryReq.DisplayName,
			EntryType:           entryReq.EntryType,
			Sets:                entryReq.Sets,
			Reps:                entryReq.Reps,
			LoadKg:              entryReq.LoadKg,
			RPE:                 entryReq.RPE,
			PlannedRestSeconds:  entryReq.PlannedRestSeconds,
			PerformedRestSeconds: entryReq.PerformedRestSeconds,
			PerSetRestOverrides: entryReq.PerSetRestOverrides,
			Notes:               entryReq.Notes,
			Order:               i,
		}

		// Parse optional timestamps
		if entryReq.EntryStart != nil {
			if t, err := time.Parse(time.RFC3339, *entryReq.EntryStart); err == nil {
				entry.EntryStart = &t
			}
		}
		if entryReq.EntryEnd != nil {
			if t, err := time.Parse(time.RFC3339, *entryReq.EntryEnd); err == nil {
				entry.EntryEnd = &t
			}
		}

		// Parse program_node_id
		if entryReq.ProgramNodeID != nil {
			if id, err := uuid.Parse(*entryReq.ProgramNodeID); err == nil {
				entry.ProgramNodeID = &id
			}
		}

		// Convert plan_snapshot
		if entryReq.PlanSnapshot != nil {
			snapshot := domain.PlanSnapshot{
				TargetSets:         entryReq.PlanSnapshot.TargetSets,
				TargetReps:         entryReq.PlanSnapshot.TargetReps,
				TargetRPE:          entryReq.PlanSnapshot.TargetRPE,
				TargetLoadKg:       entryReq.PlanSnapshot.TargetLoadKg,
				Percent1RM:         entryReq.PlanSnapshot.Percent1RM,
				PlannedRestSeconds: entryReq.PlanSnapshot.PlannedRestSeconds,
			}
			if entryReq.PlanSnapshot.ProgramNodeID != nil {
				if id, err := uuid.Parse(*entryReq.PlanSnapshot.ProgramNodeID); err == nil {
					snapshot.ProgramNodeID = &id
				}
			}
			entry.PlanSnapshot = &snapshot
		}

		workout.Entries[i] = entry
	}

	// Validate
	if err := domain.ValidateWorkout(workout); err != nil {
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
	if err := h.repo.Create(ctx, workout); err != nil {
		middleware.WriteInternalError(w, "Failed to create workout")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workout)
}

// GetWorkout handles GET /workouts/{id}
func (h *WorkoutHandler) GetWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from path
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteValidationError(w, "Invalid workout ID format", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	// Get from repository
	workout, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Workout not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve workout")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)
}

// UpdateWorkout handles PUT /workouts/{id}
func (h *WorkoutHandler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from path
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteValidationError(w, "Invalid workout ID format", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	// Check if workout exists
	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Workout not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve workout")
		return
	}

	// Decode request body (full replacement, same structure as Create)
	var req struct {
		Timestamp      string    `json:"timestamp"`
		SessionStart   *string   `json:"session_start,omitempty"`
		SessionEnd     *string   `json:"session_end,omitempty"`
		BodyWeightKg   *float64  `json:"body_weight_kg,omitempty"`
		FatigueLevel   *int      `json:"fatigue_level,omitempty"`
		SleepHours     *float64  `json:"sleep_hours,omitempty"`
		ConditionNotes *string   `json:"condition_notes,omitempty"`
		ProgramNodeID  *string   `json:"program_node_id,omitempty"`
		ProgramContext []string  `json:"program_context,omitempty"`
		Notes          *string   `json:"notes,omitempty"`
		Entries        []struct {
			ExerciseID           string    `json:"exercise_id"`
			DisplayName          *string   `json:"display_name,omitempty"`
			EntryType            string    `json:"entry_type"`
			Sets                 int       `json:"sets"`
			Reps                 int       `json:"reps"`
			LoadKg               float64   `json:"load_kg"`
			RPE                  int       `json:"rpe"`
			EntryStart           *string   `json:"entry_start,omitempty"`
			EntryEnd             *string   `json:"entry_end,omitempty"`
			PlannedRestSeconds   *int      `json:"planned_rest_seconds,omitempty"`
			PerformedRestSeconds *int      `json:"performed_rest_seconds,omitempty"`
			PerSetRestOverrides  []int     `json:"per_set_rest_overrides,omitempty"`
			ProgramNodeID        *string   `json:"program_node_id,omitempty"`
			PlanSnapshot         *struct {
				ProgramNodeID      *string  `json:"program_node_id,omitempty"`
				TargetSets         *int     `json:"target_sets,omitempty"`
				TargetReps         *int     `json:"target_reps,omitempty"`
				TargetRPE          *int     `json:"target_rpe,omitempty"`
				TargetLoadKg       *float64 `json:"target_load_kg,omitempty"`
				Percent1RM         *float64 `json:"percent_1rm,omitempty"`
				PlannedRestSeconds *int     `json:"planned_rest_seconds,omitempty"`
			} `json:"plan_snapshot,omitempty"`
			Notes *string `json:"notes,omitempty"`
		} `json:"entries"`
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

	// Update existing workout (full replacement)
	existing.Timestamp = timestamp
	existing.BodyWeightKg = req.BodyWeightKg
	existing.FatigueLevel = req.FatigueLevel
	existing.SleepHours = req.SleepHours
	existing.ConditionNotes = req.ConditionNotes
	existing.ProgramContext = req.ProgramContext
	existing.Notes = req.Notes
	existing.Entries = make([]domain.WorkoutEntry, len(req.Entries))

	// Parse optional timestamps
	if req.SessionStart != nil {
		if t, err := time.Parse(time.RFC3339, *req.SessionStart); err == nil {
			existing.SessionStart = &t
		}
	} else {
		existing.SessionStart = nil
	}
	if req.SessionEnd != nil {
		if t, err := time.Parse(time.RFC3339, *req.SessionEnd); err == nil {
			existing.SessionEnd = &t
		}
	} else {
		existing.SessionEnd = nil
	}

	// Parse program_node_id
	if req.ProgramNodeID != nil {
		if id, err := uuid.Parse(*req.ProgramNodeID); err == nil {
			existing.ProgramNodeID = &id
		}
	} else {
		existing.ProgramNodeID = nil
	}

	// Convert entries (same logic as Create)
	for i, entryReq := range req.Entries {
		exerciseID, err := uuid.Parse(entryReq.ExerciseID)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid exercise_id format", map[string]interface{}{
				"field": "entries",
				"index": i,
			})
			return
		}

		// Validate Exercise exists
		if _, err := h.exerciseRepo.GetByID(ctx, exerciseID); err != nil {
			if err == domain.ErrNotFound {
				middleware.WriteValidationError(w, "Exercise not found", map[string]interface{}{
					"field":      "entries",
					"index":      i,
					"exercise_id": entryReq.ExerciseID,
				})
				return
			}
			middleware.WriteInternalError(w, "Failed to validate exercise")
			return
		}

		entry := domain.WorkoutEntry{
			ExerciseID:          exerciseID,
			DisplayName:         entryReq.DisplayName,
			EntryType:           entryReq.EntryType,
			Sets:                entryReq.Sets,
			Reps:                entryReq.Reps,
			LoadKg:              entryReq.LoadKg,
			RPE:                 entryReq.RPE,
			PlannedRestSeconds:  entryReq.PlannedRestSeconds,
			PerformedRestSeconds: entryReq.PerformedRestSeconds,
			PerSetRestOverrides: entryReq.PerSetRestOverrides,
			Notes:               entryReq.Notes,
			Order:               i,
		}

		// Parse optional timestamps
		if entryReq.EntryStart != nil {
			if t, err := time.Parse(time.RFC3339, *entryReq.EntryStart); err == nil {
				entry.EntryStart = &t
			}
		}
		if entryReq.EntryEnd != nil {
			if t, err := time.Parse(time.RFC3339, *entryReq.EntryEnd); err == nil {
				entry.EntryEnd = &t
			}
		}

		// Parse program_node_id
		if entryReq.ProgramNodeID != nil {
			if id, err := uuid.Parse(*entryReq.ProgramNodeID); err == nil {
				entry.ProgramNodeID = &id
			}
		}

		// Convert plan_snapshot
		if entryReq.PlanSnapshot != nil {
			snapshot := domain.PlanSnapshot{
				TargetSets:         entryReq.PlanSnapshot.TargetSets,
				TargetReps:         entryReq.PlanSnapshot.TargetReps,
				TargetRPE:          entryReq.PlanSnapshot.TargetRPE,
				TargetLoadKg:       entryReq.PlanSnapshot.TargetLoadKg,
				Percent1RM:         entryReq.PlanSnapshot.Percent1RM,
				PlannedRestSeconds: entryReq.PlanSnapshot.PlannedRestSeconds,
			}
			if entryReq.PlanSnapshot.ProgramNodeID != nil {
				if id, err := uuid.Parse(*entryReq.PlanSnapshot.ProgramNodeID); err == nil {
					snapshot.ProgramNodeID = &id
				}
			}
			entry.PlanSnapshot = &snapshot
		}

		existing.Entries[i] = entry
	}

	// Validate
	if err := domain.ValidateWorkout(existing); err != nil {
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
			middleware.WriteNotFoundError(w, "Workout not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update workout")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existing)
}

// DeleteWorkout handles DELETE /workouts/{id}
func (h *WorkoutHandler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from path
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteValidationError(w, "Invalid workout ID format", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	// Check if workout exists
	if _, err := h.repo.GetByID(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Workout not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve workout")
		return
	}

	// Delete (cascades to WorkoutEntry records per FR-026)
	if err := h.repo.Delete(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Workout not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete workout")
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
