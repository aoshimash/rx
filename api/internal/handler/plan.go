package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

// planEntryRequest represents a plan entry in the request body
type planEntryRequest struct {
	ExerciseName string          `json:"exercise_name"`
	Order        int             `json:"order"`
	Sets         *int            `json:"sets,omitempty"`
	Reps         *int            `json:"reps,omitempty"`
	LoadKg       *float64        `json:"load_kg,omitempty"`
	RPE          *int            `json:"rpe,omitempty"`
	Notes        *string         `json:"notes,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}

// PlanHandler handles Plan-related HTTP requests
type PlanHandler struct {
	repo    repository.PlanRepository
	logRepo repository.LogRepository
}

// NewPlanHandler creates a new PlanHandler
func NewPlanHandler(repo repository.PlanRepository, logRepo repository.LogRepository) *PlanHandler {
	return &PlanHandler{
		repo:    repo,
		logRepo: logRepo,
	}
}

// CreatePlan handles POST /plans
func (h *PlanHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		ProgramID   *string            `json:"program_id,omitempty"`
		CycleID     *string            `json:"cycle_id,omitempty"`
		Name        string             `json:"name"`
		Date        *domain.DateOnly   `json:"date,omitempty"`
		SessionName *string            `json:"session_name,omitempty"`
		Description *string            `json:"description,omitempty"`
		Notes       *string            `json:"notes,omitempty"`
		Metadata    json.RawMessage    `json:"metadata,omitempty"`
		Entries     []planEntryRequest `json:"entries,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	plan := &domain.Plan{
		Name:        req.Name,
		Date:        req.Date,
		SessionName: req.SessionName,
		Description: req.Description,
		Notes:       req.Notes,
		Metadata:    req.Metadata,
		Entries:     make([]domain.PlanEntry, len(req.Entries)),
	}

	if req.ProgramID != nil {
		programID, err := uuid.Parse(*req.ProgramID)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid program_id format", nil)
			return
		}
		plan.ProgramID = &programID
	}

	if req.CycleID != nil {
		cycleID, err := uuid.Parse(*req.CycleID)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid cycle_id format", nil)
			return
		}
		plan.CycleID = &cycleID
	}

	for i, entryReq := range req.Entries {
		plan.Entries[i] = convertPlanEntry(entryReq)
	}

	if err := domain.ValidatePlan(plan); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := h.repo.Create(ctx, plan); err != nil {
		middleware.WriteInternalError(w, "Failed to create plan")
		return
	}

	writeJSON(w, http.StatusCreated, plan)
}

func convertPlanEntry(req planEntryRequest) domain.PlanEntry {
	return domain.PlanEntry{
		ExerciseName: req.ExerciseName,
		Order:        req.Order,
		Sets:         req.Sets,
		Reps:         req.Reps,
		LoadKg:       req.LoadKg,
		RPE:          req.RPE,
		Notes:        req.Notes,
		Metadata:     req.Metadata,
	}
}

// GetPlan handles GET /plans/{id}
func (h *PlanHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "plan")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	plan, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Plan not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve plan")
		return
	}

	writeJSON(w, http.StatusOK, plan)
}

// UpdatePlan handles PUT /plans/{id}
func (h *PlanHandler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "plan")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Plan not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve plan")
		return
	}

	var req struct {
		Name        string             `json:"name"`
		Date        *domain.DateOnly   `json:"date,omitempty"`
		SessionName *string            `json:"session_name,omitempty"`
		Description *string            `json:"description,omitempty"`
		Notes       *string            `json:"notes,omitempty"`
		Metadata    json.RawMessage    `json:"metadata,omitempty"`
		Entries     []planEntryRequest `json:"entries,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	existing.Name = req.Name
	existing.Date = req.Date
	existing.SessionName = req.SessionName
	existing.Description = req.Description
	existing.Notes = req.Notes
	existing.Metadata = req.Metadata
	existing.Entries = make([]domain.PlanEntry, len(req.Entries))

	for i, entryReq := range req.Entries {
		existing.Entries[i] = convertPlanEntry(entryReq)
	}

	if err := domain.ValidatePlan(existing); err != nil {
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
			middleware.WriteNotFoundError(w, "Plan not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update plan")
		return
	}

	writeJSON(w, http.StatusOK, existing)
}

// DeletePlan handles DELETE /plans/{id}
func (h *PlanHandler) DeletePlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "plan")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if _, err := h.repo.GetByID(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Plan not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve plan")
		return
	}

	// Check for referential integrity: Plan referenced by Log
	referencingLogs, err := h.logRepo.ListByPlanID(ctx, id)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to check references")
		return
	}
	if len(referencingLogs) > 0 {
		middleware.WriteConflictError(w, "Cannot delete plan - referenced by logs", map[string]interface{}{
			"blocking_references": []map[string]interface{}{
				{
					"type":  "log",
					"count": len(referencingLogs),
				},
			},
		})
		return
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Plan not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete plan")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListPlans handles GET /plans
func (h *PlanHandler) ListPlans(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr, 1, 100); err == nil {
			limit = parsedLimit
		}
	}
	after := r.URL.Query().Get("after")

	var plans []*domain.Plan
	var nextCursor string
	var hasMore bool
	var err error

	if programIDStr := r.URL.Query().Get("program_id"); programIDStr != "" {
		programID, parseErr := uuid.Parse(programIDStr)
		if parseErr != nil {
			middleware.WriteValidationError(w, "Invalid program_id format", nil)
			return
		}
		plans, nextCursor, hasMore, err = h.repo.ListByProgramID(ctx, programID, limit, after)
	} else {
		plans, nextCursor, hasMore, err = h.repo.List(ctx, limit, after)
	}

	if err != nil {
		middleware.WriteInternalError(w, "Failed to list plans")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":        plans,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}
