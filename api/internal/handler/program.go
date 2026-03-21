package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

// programEntryRequest represents a program entry in the request body
type programEntryRequest struct {
	ExerciseName string          `json:"exercise_name"`
	Order        int             `json:"order"`
	Sets         *int            `json:"sets,omitempty"`
	Reps         *int            `json:"reps,omitempty"`
	RPE          *int            `json:"rpe,omitempty"`
	Percent1RM   *float64        `json:"percent_1rm,omitempty"`
	Notes        *string         `json:"notes,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}

// ProgramHandler handles Program-related HTTP requests
type ProgramHandler struct {
	repo      repository.ProgramRepository
	planRepo  repository.PlanRepository
	cycleRepo repository.CycleRepository
}

// NewProgramHandler creates a new ProgramHandler
func NewProgramHandler(repo repository.ProgramRepository, planRepo repository.PlanRepository, cycleRepo repository.CycleRepository) *ProgramHandler {
	return &ProgramHandler{
		repo:      repo,
		planRepo:  planRepo,
		cycleRepo: cycleRepo,
	}
}

// CreateProgram handles POST /programs
func (h *ProgramHandler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name        string                `json:"name"`
		Description *string               `json:"description,omitempty"`
		Notes       *string               `json:"notes,omitempty"`
		Metadata    json.RawMessage       `json:"metadata,omitempty"`
		Entries     []programEntryRequest `json:"entries,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	program := &domain.Program{
		Name:        req.Name,
		Description: req.Description,
		Notes:       req.Notes,
		Metadata:    req.Metadata,
		Entries:     make([]domain.ProgramEntry, len(req.Entries)),
	}

	for i, entryReq := range req.Entries {
		program.Entries[i] = convertProgramEntry(entryReq)
	}

	if err := domain.ValidateProgram(program); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := h.repo.Create(ctx, program); err != nil {
		middleware.WriteInternalError(w, "Failed to create program")
		return
	}

	writeJSON(w, http.StatusCreated, program)
}

func convertProgramEntry(req programEntryRequest) domain.ProgramEntry {
	return domain.ProgramEntry{
		ExerciseName: req.ExerciseName,
		Order:        req.Order,
		Sets:         req.Sets,
		Reps:         req.Reps,
		RPE:          req.RPE,
		Percent1RM:   req.Percent1RM,
		Notes:        req.Notes,
		Metadata:     req.Metadata,
	}
}

// GetProgram handles GET /programs/{id}
func (h *ProgramHandler) GetProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	program, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	writeJSON(w, http.StatusOK, program)
}

// ArchiveProgram handles POST /programs/{id}/archive
func (h *ProgramHandler) ArchiveProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if err := h.repo.Archive(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to archive program")
		return
	}

	program, err := h.repo.GetByID(ctx, id)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	writeJSON(w, http.StatusOK, program)
}

// UnarchiveProgram handles POST /programs/{id}/unarchive
func (h *ProgramHandler) UnarchiveProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if err := h.repo.Unarchive(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to unarchive program")
		return
	}

	program, err := h.repo.GetByID(ctx, id)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	writeJSON(w, http.StatusOK, program)
}

// DuplicateProgram handles POST /programs/{id}/duplicate
func (h *ProgramHandler) DuplicateProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	original, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	copyName := strings.TrimSpace(original.Name) + " (copy)"

	entries := make([]domain.ProgramEntry, len(original.Entries))
	for i, e := range original.Entries {
		entries[i] = domain.ProgramEntry{
			Order:        e.Order,
			ExerciseName: e.ExerciseName,
			Sets:         e.Sets,
			Reps:         e.Reps,
			RPE:          e.RPE,
			Percent1RM:   e.Percent1RM,
			Notes:        e.Notes,
		}
		if e.Metadata != nil {
			entries[i].Metadata = make([]byte, len(e.Metadata))
			copy(entries[i].Metadata, e.Metadata)
		}
	}

	var copyMeta json.RawMessage
	if original.Metadata != nil {
		copyMeta = make(json.RawMessage, len(original.Metadata))
		copy(copyMeta, original.Metadata)
	}

	duplicate := &domain.Program{
		Name:        copyName,
		Description: original.Description,
		Notes:       original.Notes,
		Metadata:    copyMeta,
		Entries:     entries,
	}

	if err := h.repo.Create(ctx, duplicate); err != nil {
		middleware.WriteInternalError(w, "Failed to duplicate program")
		return
	}

	writeJSON(w, http.StatusCreated, duplicate)
}

// DeleteProgram handles DELETE /programs/{id}
func (h *ProgramHandler) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if _, err := h.repo.GetByID(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	// Check for referential integrity: Program referenced by Cycles
	hasCycles, err := h.cycleRepo.ExistsByProgramID(ctx, id)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to check references")
		return
	}
	if hasCycles {
		middleware.WriteConflictError(w, "Cannot delete program - referenced by cycles", map[string]interface{}{
			"blocking_references": []map[string]interface{}{
				{
					"type": "cycle",
				},
			},
		})
		return
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete program")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListPrograms handles GET /programs
func (h *ProgramHandler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr, 1, 100); err == nil {
			limit = parsedLimit
		}
	}
	after := r.URL.Query().Get("after")
	includeArchived := r.URL.Query().Get("include_archived") == "true"

	programs, nextCursor, hasMore, err := h.repo.List(ctx, limit, after, includeArchived)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to list programs")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":        programs,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}

// ConvertToPlans handles POST /plans/from-program
func (h *ProgramHandler) ConvertToPlans(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		ProgramID      string             `json:"program_id"`
		Name           string             `json:"name,omitempty"`
		TargetWeights  map[string]float64 `json:"target_weights"`
		LoadIncrements map[string]float64 `json:"load_increments,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if req.ProgramID == "" {
		middleware.WriteValidationError(w, "program_id is required", nil)
		return
	}

	programID, err := parseUUIDString(req.ProgramID, "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	// Fetch the program
	program, err := h.repo.GetByID(ctx, programID)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	// Convert to plans (one per session)
	input := &domain.ConvertProgramToPlansInput{
		Name:           req.Name,
		TargetWeights:  req.TargetWeights,
		LoadIncrements: req.LoadIncrements,
	}

	plans := domain.ConvertProgramToPlans(program, input)

	// Validate all plans before saving any
	for _, plan := range plans {
		if err := domain.ValidatePlan(plan); err != nil {
			if handleValidationError(w, err) {
				return
			}
			middleware.WriteValidationError(w, "Generated plan validation failed", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
	}

	// Save all validated plans
	savedPlans := make([]*domain.Plan, 0, len(plans))
	for _, plan := range plans {
		if err := h.planRepo.Create(ctx, plan); err != nil {
			middleware.WriteInternalError(w, "Failed to create plan")
			return
		}
		savedPlans = append(savedPlans, plan)
	}

	writeJSON(w, http.StatusCreated, savedPlans)
}

// parseUUIDString parses a UUID from a string value
func parseUUIDString(s string, resourceName string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s ID format: %s", resourceName, s)
	}
	return id, nil
}
