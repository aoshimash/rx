package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	repo     repository.ProgramRepository
	planRepo repository.PlanRepository
}

// NewProgramHandler creates a new ProgramHandler
func NewProgramHandler(repo repository.ProgramRepository, planRepo repository.PlanRepository) *ProgramHandler {
	return &ProgramHandler{
		repo:     repo,
		planRepo: planRepo,
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

// UpdateProgram handles PUT /programs/{id}
func (h *ProgramHandler) UpdateProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

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

	existing.Name = req.Name
	existing.Description = req.Description
	existing.Notes = req.Notes
	existing.Metadata = req.Metadata
	existing.Entries = make([]domain.ProgramEntry, len(req.Entries))

	for i, entryReq := range req.Entries {
		existing.Entries[i] = convertProgramEntry(entryReq)
	}

	if err := domain.ValidateProgram(existing); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := h.repo.Update(ctx, existing); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update program")
		return
	}

	writeJSON(w, http.StatusOK, existing)
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

	programs, nextCursor, hasMore, err := h.repo.List(ctx, limit, after)
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

	// Validate and save each plan
	savedPlans := make([]*domain.Plan, 0, len(plans))
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
