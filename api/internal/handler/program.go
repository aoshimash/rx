package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

// programGroupRequest represents a program group in the request body
type programGroupRequest struct {
	Name          string  `json:"name"`
	Order         int     `json:"order"`
	ParentGroupID *string `json:"parent_group_id,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

// programSessionEntryRequest represents a program session entry in the request body
type programSessionEntryRequest struct {
	ExerciseName string                 `json:"exercise_name"`
	Order        int                    `json:"order"`
	Fields       map[string]interface{} `json:"fields,omitempty"`
	Notes        *string                `json:"notes,omitempty"`
}

// programSessionRequest represents a program session in the request body
type programSessionRequest struct {
	SessionName string                       `json:"session_name"`
	Order       int                          `json:"order"`
	GroupID     *string                      `json:"group_id,omitempty"`
	Date        *string                      `json:"date,omitempty"`
	Entries     []programSessionEntryRequest `json:"entries,omitempty"`
}

// ProgramHandler handles Program-related HTTP requests
type ProgramHandler struct {
	repo repository.ProgramRepository
}

// NewProgramHandler creates a new ProgramHandler
func NewProgramHandler(repo repository.ProgramRepository) *ProgramHandler {
	return &ProgramHandler{
		repo: repo,
	}
}

// parseGroups converts group request objects to domain ProgramGroup slice
func parseGroups(groups []programGroupRequest) ([]domain.ProgramGroup, error) {
	result := make([]domain.ProgramGroup, len(groups))
	for i, g := range groups {
		group := domain.ProgramGroup{
			ID:    uuid.New(),
			Name:  g.Name,
			Order: g.Order,
			Notes: g.Notes,
		}
		if g.ParentGroupID != nil {
			pid, err := uuid.Parse(*g.ParentGroupID)
			if err != nil {
				return nil, &domain.ValidationError{
					Field:   "groups[].parent_group_id",
					Message: "invalid UUID format: " + err.Error(),
				}
			}
			group.ParentGroupID = &pid
		}
		result[i] = group
	}
	return result, nil
}

// parseSessions converts session request objects to domain ProgramSession slice
func parseSessions(sessions []programSessionRequest) ([]domain.ProgramSession, error) {
	result := make([]domain.ProgramSession, len(sessions))
	for i, sessReq := range sessions {
		sess := domain.ProgramSession{
			SessionName: sessReq.SessionName,
			Order:       sessReq.Order,
			Entries:     make([]domain.ProgramSessionEntry, len(sessReq.Entries)),
		}

		if sessReq.GroupID != nil {
			gid, err := uuid.Parse(*sessReq.GroupID)
			if err != nil {
				return nil, &domain.ValidationError{
					Field:   "sessions[].group_id",
					Message: "invalid UUID format: " + err.Error(),
				}
			}
			sess.GroupID = &gid
		}

		if sessReq.Date != nil {
			var d domain.DateOnly
			if err := json.Unmarshal([]byte(`"`+*sessReq.Date+`"`), &d); err != nil {
				return nil, &domain.ValidationError{
					Field:   "sessions[].date",
					Message: "invalid date format (expected YYYY-MM-DD)",
				}
			}
			sess.Date = &d
		}

		for j, e := range sessReq.Entries {
			sess.Entries[j] = domain.ProgramSessionEntry{
				ExerciseName: e.ExerciseName,
				Order:        e.Order,
				Fields:       e.Fields,
				Notes:        e.Notes,
			}
		}

		result[i] = sess
	}
	return result, nil
}

// CreateProgram handles POST /programs
func (h *ProgramHandler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name          string                  `json:"name"`
		Notes         *string                 `json:"notes,omitempty"`
		ProgramFields []domain.FieldDef       `json:"program_fields,omitempty"`
		LogFields     []domain.FieldDef       `json:"log_fields,omitempty"`
		Groups        []programGroupRequest   `json:"groups,omitempty"`
		Sessions      []programSessionRequest `json:"sessions,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	exists, err := h.repo.ExistsByName(ctx, req.Name)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to check program name")
		return
	}
	if exists {
		middleware.WriteConflictError(w, "A program with this name already exists", map[string]interface{}{
			"field": "name",
		})
		return
	}

	groups, err := parseGroups(req.Groups)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Invalid groups", nil)
		return
	}

	sessions, err := parseSessions(req.Sessions)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Invalid sessions", nil)
		return
	}

	program := &domain.Program{
		Name:          req.Name,
		Notes:         req.Notes,
		ProgramFields: req.ProgramFields,
		LogFields:     req.LogFields,
		Groups:        groups,
		Sessions:      sessions,
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
		Name          string                  `json:"name"`
		Notes         *string                 `json:"notes,omitempty"`
		ProgramFields []domain.FieldDef       `json:"program_fields,omitempty"`
		LogFields     []domain.FieldDef       `json:"log_fields,omitempty"`
		Groups        []programGroupRequest   `json:"groups,omitempty"`
		Sessions      []programSessionRequest `json:"sessions,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Check name uniqueness if name changed
	if req.Name != existing.Name {
		nameExists, err := h.repo.ExistsByName(ctx, req.Name)
		if err != nil {
			middleware.WriteInternalError(w, "Failed to check program name")
			return
		}
		if nameExists {
			middleware.WriteConflictError(w, "A program with this name already exists", map[string]interface{}{
				"field": "name",
			})
			return
		}
	}

	groups, err := parseGroups(req.Groups)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Invalid groups", nil)
		return
	}

	sessions, err := parseSessions(req.Sessions)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Invalid sessions", nil)
		return
	}

	program := &domain.Program{
		ID:            existing.ID,
		Name:          req.Name,
		Notes:         req.Notes,
		ProgramFields: req.ProgramFields,
		LogFields:     req.LogFields,
		Groups:        groups,
		Sessions:      sessions,
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

	if err := h.repo.Update(ctx, program); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update program")
		return
	}

	writeJSON(w, http.StatusOK, program)
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
