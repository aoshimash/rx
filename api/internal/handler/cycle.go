package handler

import (
	"net/http"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

// CycleHandler handles Cycle-related HTTP requests
type CycleHandler struct {
	repo     repository.CycleRepository
	planRepo repository.PlanRepository
}

// NewCycleHandler creates a new CycleHandler
func NewCycleHandler(repo repository.CycleRepository, planRepo repository.PlanRepository) *CycleHandler {
	return &CycleHandler{
		repo:     repo,
		planRepo: planRepo,
	}
}

// GetCycle handles GET /cycles/{id}
func (h *CycleHandler) GetCycle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "cycle")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	cycle, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Cycle not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve cycle")
		return
	}

	writeJSON(w, http.StatusOK, cycle)
}

// ListCycles handles GET /cycles
func (h *CycleHandler) ListCycles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr, 1, 100); err == nil {
			limit = parsedLimit
		}
	}
	after := r.URL.Query().Get("after")

	var cycles []*domain.Cycle
	var nextCursor string
	var hasMore bool
	var err error

	if programIDStr := r.URL.Query().Get("program_id"); programIDStr != "" {
		programID, parseErr := uuid.Parse(programIDStr)
		if parseErr != nil {
			middleware.WriteValidationError(w, "Invalid program_id format", nil)
			return
		}
		cycles, nextCursor, hasMore, err = h.repo.ListByProgramID(ctx, programID, limit, after)
	} else {
		cycles, nextCursor, hasMore, err = h.repo.List(ctx, limit, after)
	}

	if err != nil {
		middleware.WriteInternalError(w, "Failed to list cycles")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":        cycles,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}

// DeleteCycle handles DELETE /cycles/{id}
func (h *CycleHandler) DeleteCycle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "cycle")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if _, err := h.repo.GetByID(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Cycle not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve cycle")
		return
	}

	// Check for referential integrity: Cycle referenced by Plans
	planCount, err := h.planRepo.CountByCycleID(ctx, id)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to check references")
		return
	}
	if planCount > 0 {
		middleware.WriteConflictError(w, "Cannot delete cycle - referenced by plans", map[string]interface{}{
			"blocking_references": []map[string]interface{}{
				{
					"type":  "plan",
					"count": planCount,
				},
			},
		})
		return
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Cycle not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete cycle")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
