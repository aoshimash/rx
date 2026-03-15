package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aoshimash/optel-workout/api/internal/domain"
	"github.com/aoshimash/optel-workout/api/internal/middleware"
	"github.com/aoshimash/optel-workout/api/internal/repository"
	"github.com/google/uuid"
)

// programEntryRequest represents a program entry in the request body
type programEntryRequest struct {
	Name               string          `json:"name"`
	Order              int             `json:"order"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	ExerciseID         *string         `json:"exercise_id,omitempty"`
	TargetSets         *int            `json:"target_sets,omitempty"`
	TargetReps         *int            `json:"target_reps,omitempty"`
	TargetRPE          *int            `json:"target_rpe,omitempty"`
	Percent1RM         *float64        `json:"percent_1rm,omitempty"`
	PlannedRestSeconds *int            `json:"planned_rest_seconds,omitempty"`
	MuscleGroups       []string        `json:"muscle_groups,omitempty"`
	Notes              *string         `json:"notes,omitempty"`
}

// ProgramHandler handles Program-related HTTP requests
type ProgramHandler struct {
	repo         repository.ProgramRepository
	exerciseRepo repository.ExerciseRepository
	workoutRepo  repository.WorkoutRepository
}

// NewProgramHandler creates a new ProgramHandler
func NewProgramHandler(repo repository.ProgramRepository, exerciseRepo repository.ExerciseRepository, workoutRepo repository.WorkoutRepository) *ProgramHandler {
	return &ProgramHandler{
		repo:         repo,
		exerciseRepo: exerciseRepo,
		workoutRepo:  workoutRepo,
	}
}

// CreateProgram handles POST /programs
func (h *ProgramHandler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name        string               `json:"name"`
		Description *string              `json:"description,omitempty"`
		Entries     []programEntryRequest `json:"entries,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Convert to domain model
	program := &domain.Program{
		Name:        req.Name,
		Description: req.Description,
		Entries:     make([]domain.ProgramEntry, len(req.Entries)),
	}

	for i, entryReq := range req.Entries {
		entry, err := h.convertEntry(ctx, entryReq)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid program entry", map[string]interface{}{
				"field": "entries",
				"index": i,
				"error": err.Error(),
			})
			return
		}
		program.Entries[i] = *entry
	}

	// Validate
	if err := domain.ValidateProgram(program); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Create
	if err := h.repo.Create(ctx, program); err != nil {
		middleware.WriteInternalError(w, "Failed to create program")
		return
	}

	// Return response
	writeJSON(w, http.StatusCreated, program)
}

// convertEntry converts a request entry to a domain ProgramEntry
func (h *ProgramHandler) convertEntry(ctx context.Context, entryReq programEntryRequest) (*domain.ProgramEntry, error) {
	entry := &domain.ProgramEntry{
		Name:               entryReq.Name,
		Order:              entryReq.Order,
		Metadata:           entryReq.Metadata,
		TargetSets:         entryReq.TargetSets,
		TargetReps:         entryReq.TargetReps,
		TargetRPE:          entryReq.TargetRPE,
		Percent1RM:         entryReq.Percent1RM,
		PlannedRestSeconds: entryReq.PlannedRestSeconds,
		MuscleGroups:       entryReq.MuscleGroups,
		Notes:              entryReq.Notes,
	}

	// Parse exercise_id if provided
	if entryReq.ExerciseID != nil {
		exerciseID, err := uuid.Parse(*entryReq.ExerciseID)
		if err != nil {
			return nil, err
		}
		// Validate Exercise exists
		if _, err := h.exerciseRepo.GetByID(ctx, exerciseID); err != nil {
			return nil, err
		}
		entry.ExerciseID = &exerciseID
	}

	return entry, nil
}

// GetProgram handles GET /programs/{id}
func (h *ProgramHandler) GetProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract and parse ID from path
	id, err := parseUUIDParam(r, "id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	// Get from repository
	program, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	// Return response
	writeJSON(w, http.StatusOK, program)
}

// UpdateProgram handles PUT /programs/{id}
func (h *ProgramHandler) UpdateProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract and parse ID from path
	id, err := parseUUIDParam(r, "id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	// Check if program exists
	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	var req struct {
		Name        string               `json:"name"`
		Description *string              `json:"description,omitempty"`
		Entries     []programEntryRequest `json:"entries,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Update existing program (full replacement)
	existing.Name = req.Name
	existing.Description = req.Description
	existing.Entries = make([]domain.ProgramEntry, len(req.Entries))

	for i, entryReq := range req.Entries {
		entry, err := h.convertEntry(ctx, entryReq)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid program entry", map[string]interface{}{
				"field": "entries",
				"index": i,
				"error": err.Error(),
			})
			return
		}
		existing.Entries[i] = *entry
	}

	// Validate
	if err := domain.ValidateProgram(existing); err != nil {
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
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update program")
		return
	}

	// Return response
	writeJSON(w, http.StatusOK, existing)
}

// DeleteProgram handles DELETE /programs/{id}
func (h *ProgramHandler) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract and parse ID from path
	id, err := parseUUIDParam(r, "id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	// Check if program exists
	program, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	// Check for referential integrity: Program referenced by Workout (FR-026)
	for _, entry := range program.Entries {
		referencingWorkouts, err := h.workoutRepo.ListByProgramNodeID(ctx, entry.ID)
		if err != nil {
			middleware.WriteInternalError(w, "Failed to check references")
			return
		}
		if len(referencingWorkouts) > 0 {
			middleware.WriteConflictError(w, "Cannot delete program - referenced by workouts", map[string]interface{}{
				"blocking_references": []map[string]interface{}{
					{
						"type":  "workout",
						"count": len(referencingWorkouts),
					},
				},
			})
			return
		}
	}

	// Delete (cascades to ProgramEntry records)
	if err := h.repo.Delete(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete program")
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// ListPrograms handles GET /programs
func (h *ProgramHandler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse pagination parameters
	limit := 100 // default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr, 1, 100); err == nil {
			limit = parsedLimit
		}
	}
	after := r.URL.Query().Get("after")

	// List from repository (no filters per FR-030)
	programs, nextCursor, hasMore, err := h.repo.List(ctx, limit, after)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to list programs")
		return
	}

	// Return paginated response
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":        programs,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}
