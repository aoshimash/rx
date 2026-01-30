package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aoshimash/optel-workout/api/internal/domain"
	"github.com/aoshimash/optel-workout/api/internal/middleware"
	"github.com/aoshimash/optel-workout/api/internal/repository"
	"github.com/google/uuid"
)

// workoutRequest represents the request body for creating/updating a workout
type workoutRequest struct {
	Timestamp      string                `json:"timestamp"`
	SessionStart   *string               `json:"session_start,omitempty"`
	SessionEnd     *string               `json:"session_end,omitempty"`
	BodyWeightKg   *float64              `json:"body_weight_kg,omitempty"`
	FatigueLevel   *int                  `json:"fatigue_level,omitempty"`
	SleepHours     *float64              `json:"sleep_hours,omitempty"`
	ConditionNotes *string               `json:"condition_notes,omitempty"`
	ProgramNodeID  *string               `json:"program_node_id,omitempty"`
	ProgramContext []string              `json:"program_context,omitempty"`
	Notes          *string               `json:"notes,omitempty"`
	Entries        []workoutEntryRequest `json:"entries"`
}

// workoutEntryRequest represents a workout entry in the request body
type workoutEntryRequest struct {
	ExerciseID           string               `json:"exercise_id"`
	DisplayName          *string              `json:"display_name,omitempty"`
	EntryType            *string              `json:"entry_type,omitempty"`
	Sets                 int                  `json:"sets"`
	Reps                 int                  `json:"reps"`
	LoadKg               float64              `json:"load_kg"`
	RPE                  int                  `json:"rpe"`
	EntryStart           *string              `json:"entry_start,omitempty"`
	EntryEnd             *string              `json:"entry_end,omitempty"`
	PlannedRestSeconds   *int                 `json:"planned_rest_seconds,omitempty"`
	PerformedRestSeconds *int                 `json:"performed_rest_seconds,omitempty"`
	PerSetRestOverrides  []int                `json:"per_set_rest_overrides,omitempty"`
	ProgramNodeID        *string              `json:"program_node_id,omitempty"`
	PlanSnapshot         *planSnapshotRequest `json:"plan_snapshot,omitempty"`
	Notes                *string              `json:"notes,omitempty"`
	VideoObjectKey       *string              `json:"video_object_key,omitempty"`
}

// planSnapshotRequest represents a plan snapshot in the request body
type planSnapshotRequest struct {
	ProgramNodeID      *string  `json:"program_node_id,omitempty"`
	TargetSets         *int     `json:"target_sets,omitempty"`
	TargetReps         *int     `json:"target_reps,omitempty"`
	TargetRPE          *int     `json:"target_rpe,omitempty"`
	TargetLoadKg       *float64 `json:"target_load_kg,omitempty"`
	Percent1RM         *float64 `json:"percent_1rm,omitempty"`
	PlannedRestSeconds *int     `json:"planned_rest_seconds,omitempty"`
}

// WorkoutHandler handles Workout-related HTTP requests
type WorkoutHandler struct {
	repo         repository.WorkoutRepository
	exerciseRepo repository.ExerciseRepository
	programRepo  repository.ProgramRepository
}

// NewWorkoutHandler creates a new WorkoutHandler
func NewWorkoutHandler(repo repository.WorkoutRepository, exerciseRepo repository.ExerciseRepository, programRepo repository.ProgramRepository) *WorkoutHandler {
	return &WorkoutHandler{
		repo:         repo,
		exerciseRepo: exerciseRepo,
		programRepo:  programRepo,
	}
}

// parseWorkoutRequest parses and validates a workout request into a domain model
func (h *WorkoutHandler) parseWorkoutRequest(ctx context.Context, req *workoutRequest) (*domain.Workout, error) {
	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		return nil, &domain.ValidationError{
			Field:   "timestamp",
			Message: "invalid timestamp format: " + err.Error(),
		}
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
	}

	// Parse optional timestamps (return error if format is invalid)
	if req.SessionStart != nil {
		t, err := time.Parse(time.RFC3339, *req.SessionStart)
		if err != nil {
			return nil, &domain.ValidationError{
				Field:   "session_start",
				Message: "invalid timestamp format: " + err.Error(),
			}
		}
		workout.SessionStart = &t
	}
	if req.SessionEnd != nil {
		t, err := time.Parse(time.RFC3339, *req.SessionEnd)
		if err != nil {
			return nil, &domain.ValidationError{
				Field:   "session_end",
				Message: "invalid timestamp format: " + err.Error(),
			}
		}
		workout.SessionEnd = &t
	}

	// Parse program_node_id (return error if format is invalid)
	if req.ProgramNodeID != nil {
		id, err := uuid.Parse(*req.ProgramNodeID)
		if err != nil {
			return nil, &domain.ValidationError{
				Field:   "program_node_id",
				Message: "invalid UUID format: " + err.Error(),
			}
		}
		workout.ProgramNodeID = &id
	}

	// Convert entries
	entries, err := h.convertEntries(ctx, req.Entries)
	if err != nil {
		return nil, err
	}
	workout.Entries = entries

	return workout, nil
}

// convertEntries converts request entries to domain entries with validation
func (h *WorkoutHandler) convertEntries(ctx context.Context, entryReqs []workoutEntryRequest) ([]domain.WorkoutEntry, error) {
	entries := make([]domain.WorkoutEntry, len(entryReqs))

	for i, entryReq := range entryReqs {
		exerciseID, err := uuid.Parse(entryReq.ExerciseID)
		if err != nil {
			return nil, &domain.ValidationError{
				Field:   fmt.Sprintf("entries[%d].exercise_id", i),
				Message: "invalid exercise_id format",
			}
		}

		// Validate Exercise exists (FR-026)
		if _, err := h.exerciseRepo.GetByID(ctx, exerciseID); err != nil {
			if err == domain.ErrNotFound {
				return nil, &domain.ValidationError{
					Field:   fmt.Sprintf("entries[%d].exercise_id", i),
					Message: fmt.Sprintf("exercise not found: %s", entryReq.ExerciseID),
				}
			}
			return nil, err
		}

		entry := domain.WorkoutEntry{
			ExerciseID:           exerciseID,
			DisplayName:          entryReq.DisplayName,
			EntryType:            entryReq.EntryType, // nullable
			Sets:                 entryReq.Sets,
			Reps:                 entryReq.Reps,
			LoadKg:               entryReq.LoadKg,
			RPE:                  entryReq.RPE,
			PlannedRestSeconds:   entryReq.PlannedRestSeconds,
			PerformedRestSeconds: entryReq.PerformedRestSeconds,
			PerSetRestOverrides:  entryReq.PerSetRestOverrides,
			Notes:                entryReq.Notes,
			VideoObjectKey:       entryReq.VideoObjectKey,
			Order:                i,
		}

		// Parse optional timestamps (return error if format is invalid)
		if entryReq.EntryStart != nil {
			t, err := time.Parse(time.RFC3339, *entryReq.EntryStart)
			if err != nil {
				return nil, &domain.ValidationError{
					Field:   fmt.Sprintf("entries[%d].entry_start", i),
					Message: "invalid timestamp format: " + err.Error(),
				}
			}
			entry.EntryStart = &t
		}
		if entryReq.EntryEnd != nil {
			t, err := time.Parse(time.RFC3339, *entryReq.EntryEnd)
			if err != nil {
				return nil, &domain.ValidationError{
					Field:   fmt.Sprintf("entries[%d].entry_end", i),
					Message: "invalid timestamp format: " + err.Error(),
				}
			}
			entry.EntryEnd = &t
		}

		// Parse program_node_id (return error if format is invalid)
		if entryReq.ProgramNodeID != nil {
			id, err := uuid.Parse(*entryReq.ProgramNodeID)
			if err != nil {
				return nil, &domain.ValidationError{
					Field:   fmt.Sprintf("entries[%d].program_node_id", i),
					Message: "invalid UUID format: " + err.Error(),
				}
			}
			entry.ProgramNodeID = &id
		}

		// Convert plan_snapshot
		if entryReq.PlanSnapshot != nil {
			snapshot, err := convertPlanSnapshot(entryReq.PlanSnapshot, i)
			if err != nil {
				return nil, err
			}
			entry.PlanSnapshot = snapshot
		}

		entries[i] = entry
	}

	return entries, nil
}

// convertPlanSnapshot converts a plan snapshot request to domain model
func convertPlanSnapshot(req *planSnapshotRequest, entryIndex int) (*domain.PlanSnapshot, error) {
	snapshot := domain.PlanSnapshot{
		TargetSets:         req.TargetSets,
		TargetReps:         req.TargetReps,
		TargetRPE:          req.TargetRPE,
		TargetLoadKg:       req.TargetLoadKg,
		Percent1RM:         req.Percent1RM,
		PlannedRestSeconds: req.PlannedRestSeconds,
	}

	if req.ProgramNodeID != nil {
		id, err := uuid.Parse(*req.ProgramNodeID)
		if err != nil {
			return nil, &domain.ValidationError{
				Field:   fmt.Sprintf("entries[%d].plan_snapshot.program_node_id", entryIndex),
				Message: "invalid UUID format: " + err.Error(),
			}
		}
		snapshot.ProgramNodeID = &id
	}

	return &snapshot, nil
}

// CreateWorkout handles POST /workouts
func (h *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request body
	var req workoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Parse and validate request
	workout, err := h.parseWorkoutRequest(ctx, &req)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteInternalError(w, "Failed to parse workout request")
		return
	}

	// Validate
	if err := domain.ValidateWorkout(workout); err != nil {
		if handleValidationError(w, err) {
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
	writeJSON(w, http.StatusCreated, workout)
}

// GetWorkout handles GET /workouts/{id}
func (h *WorkoutHandler) GetWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract and parse ID from path
	id, err := parseUUIDParam(r, "id", "workout")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
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
	writeJSON(w, http.StatusOK, workout)
}

// UpdateWorkout handles PUT /workouts/{id}
func (h *WorkoutHandler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract and parse ID from path
	id, err := parseUUIDParam(r, "id", "workout")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
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

	// Decode request body
	var req workoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Parse and validate request
	updated, err := h.parseWorkoutRequest(ctx, &req)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteInternalError(w, "Failed to parse workout request")
		return
	}

	// Preserve ID and timestamps from existing workout
	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt

	// Validate
	if err := domain.ValidateWorkout(updated); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Update
	if err := h.repo.Update(ctx, updated); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Workout not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update workout")
		return
	}

	// Return response
	writeJSON(w, http.StatusOK, updated)
}

// DeleteWorkout handles DELETE /workouts/{id}
func (h *WorkoutHandler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract and parse ID from path
	id, err := parseUUIDParam(r, "id", "workout")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
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

// ListWorkouts handles GET /workouts
func (h *WorkoutHandler) ListWorkouts(w http.ResponseWriter, r *http.Request) {
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
	var workouts []*domain.Workout
	var nextCursor string
	var hasMore bool
	var err error

	if timestampFrom != nil || timestampTo != nil {
		workouts, nextCursor, hasMore, err = h.repo.ListByTimestampRange(ctx, timestampFrom, timestampTo, limit, after)
	} else {
		workouts, nextCursor, hasMore, err = h.repo.List(ctx, limit, after)
	}

	if err != nil {
		middleware.WriteInternalError(w, "Failed to list workouts")
		return
	}

	// Return paginated response
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":        workouts,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}
