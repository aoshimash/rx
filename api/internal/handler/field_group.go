package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
)

// FieldGroupHandler handles FieldGroup-related HTTP requests
type FieldGroupHandler struct {
	repo repository.FieldGroupRepository
}

// NewFieldGroupHandler creates a new FieldGroupHandler
func NewFieldGroupHandler(repo repository.FieldGroupRepository) *FieldGroupHandler {
	return &FieldGroupHandler{repo: repo}
}

type fieldDefRequest struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Options     []string `json:"options,omitempty"`
	Description string   `json:"description,omitempty"`
}

func toFieldDefs(reqs []fieldDefRequest) []domain.FieldDef {
	defs := make([]domain.FieldDef, len(reqs))
	for i, r := range reqs {
		defs[i] = domain.FieldDef{
			Name:        r.Name,
			Type:        r.Type,
			Options:     r.Options,
			Description: r.Description,
		}
	}
	return defs
}

// ListFieldGroups handles GET /field-groups
func (h *FieldGroupHandler) ListFieldGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	groups, err := h.repo.List(ctx, userID)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to list field groups")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": groups,
	})
}

// CreateFieldGroup handles POST /field-groups
func (h *FieldGroupHandler) CreateFieldGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	var req struct {
		Name          string            `json:"name"`
		Description   *string           `json:"description,omitempty"`
		ProgramFields []fieldDefRequest `json:"program_fields"`
		LogFields     []fieldDefRequest `json:"log_fields"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	fg := &domain.FieldGroup{
		Name:          req.Name,
		Description:   req.Description,
		ProgramFields: toFieldDefs(req.ProgramFields),
		LogFields:     toFieldDefs(req.LogFields),
	}

	if err := domain.ValidateFieldGroup(fg); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", nil)
		return
	}

	if err := h.repo.Create(ctx, userID, fg); err != nil {
		if de, ok := err.(*domain.DomainError); ok && de.Code == domain.ErrorCodeConflict {
			middleware.WriteConflictError(w, "A field group with this name already exists", nil)
			return
		}
		middleware.WriteInternalError(w, "Failed to create field group")
		return
	}

	writeJSON(w, http.StatusCreated, fg)
}

// GetFieldGroup handles GET /field-groups/{id}
func (h *FieldGroupHandler) GetFieldGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "field group")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	fg, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			middleware.WriteNotFoundError(w, "Field group not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve field group")
		return
	}

	writeJSON(w, http.StatusOK, fg)
}

// UpdateFieldGroup handles PUT /field-groups/{id}
func (h *FieldGroupHandler) UpdateFieldGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "field group")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			middleware.WriteNotFoundError(w, "Field group not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve field group")
		return
	}

	var req struct {
		Name          string            `json:"name"`
		Description   *string           `json:"description,omitempty"`
		ProgramFields []fieldDefRequest `json:"program_fields"`
		LogFields     []fieldDefRequest `json:"log_fields"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.ProgramFields = toFieldDefs(req.ProgramFields)
	existing.LogFields = toFieldDefs(req.LogFields)

	if err := domain.ValidateFieldGroup(existing); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", nil)
		return
	}

	if err := h.repo.Update(ctx, existing); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			middleware.WriteNotFoundError(w, "Field group not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update field group")
		return
	}

	writeJSON(w, http.StatusOK, existing)
}

// DeleteFieldGroup handles DELETE /field-groups/{id}
func (h *FieldGroupHandler) DeleteFieldGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "field group")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			middleware.WriteNotFoundError(w, "Field group not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete field group")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
