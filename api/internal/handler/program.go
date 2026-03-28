package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

// programSessionEntryRequest represents a program session entry in the request body
type programSessionEntryRequest struct {
	ExerciseName string          `json:"exercise_name"`
	Order        int             `json:"order"`
	Sets         *int            `json:"sets,omitempty"`
	Reps         *int            `json:"reps,omitempty"`
	LoadKg       *float64        `json:"load_kg,omitempty"`
	RPE          *int            `json:"rpe,omitempty"`
	Notes        *string         `json:"notes,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}

// programSessionRequest represents a program session in the request body
type programSessionRequest struct {
	SessionName string                       `json:"session_name"`
	Order       int                          `json:"order"`
	Date        *string                      `json:"date,omitempty"`
	Entries     []programSessionEntryRequest `json:"entries,omitempty"`
}

// ProgramHandler handles Program-related HTTP requests
type ProgramHandler struct {
	repo    repository.ProgramRepository
	logRepo repository.LogRepository
}

// NewProgramHandler creates a new ProgramHandler
func NewProgramHandler(repo repository.ProgramRepository, logRepo repository.LogRepository) *ProgramHandler {
	return &ProgramHandler{
		repo:    repo,
		logRepo: logRepo,
	}
}

// CreateProgram handles POST /programs
func (h *ProgramHandler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		ProgramTemplateID *string                 `json:"program_template_id,omitempty"`
		Name              string                  `json:"name"`
		Notes             *string                 `json:"notes,omitempty"`
		Metadata          json.RawMessage         `json:"metadata,omitempty"`
		Sessions          []programSessionRequest `json:"sessions,omitempty"`
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

	program := &domain.Program{
		Name:     req.Name,
		Status:   domain.ProgramStatusCreated,
		Notes:    req.Notes,
		Metadata: req.Metadata,
		Sessions: make([]domain.ProgramSession, len(req.Sessions)),
	}

	if req.ProgramTemplateID != nil {
		tid, err := uuid.Parse(*req.ProgramTemplateID)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid program_template_id format", nil)
			return
		}
		program.ProgramTemplateID = &tid
	}

	for i, sessReq := range req.Sessions {
		sess := domain.ProgramSession{
			SessionName: sessReq.SessionName,
			Order:       sessReq.Order,
			Entries:     make([]domain.ProgramSessionEntry, len(sessReq.Entries)),
		}

		if sessReq.Date != nil {
			var d domain.DateOnly
			if err := json.Unmarshal([]byte(`"`+*sessReq.Date+`"`), &d); err != nil {
				middleware.WriteValidationError(w, "Invalid date format in session (expected YYYY-MM-DD)", nil)
				return
			}
			sess.Date = &d
		}

		for j, e := range sessReq.Entries {
			sess.Entries[j] = domain.ProgramSessionEntry{
				ExerciseName: e.ExerciseName,
				Order:        e.Order,
				Sets:         e.Sets,
				Reps:         e.Reps,
				LoadKg:       e.LoadKg,
				RPE:          e.RPE,
				Notes:        e.Notes,
				Metadata:     e.Metadata,
			}
		}

		program.Sessions[i] = sess
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

	var programTemplateID *uuid.UUID
	if tmplIDStr := r.URL.Query().Get("program_template_id"); tmplIDStr != "" {
		tid, err := uuid.Parse(tmplIDStr)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid program_template_id format", nil)
			return
		}
		programTemplateID = &tid
	}

	status := r.URL.Query().Get("status")

	programs, nextCursor, hasMore, err := h.repo.List(ctx, limit, after, programTemplateID, status)
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

// UpdateProgramStatus handles PATCH /programs/{id}/status
func (h *ProgramHandler) UpdateProgramStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", nil)
		return
	}

	newStatus := domain.ProgramStatus(req.Status)

	program, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	if err := domain.ValidateProgramStatusTransition(program.Status, newStatus); err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if err := h.repo.UpdateStatus(ctx, id, newStatus); err != nil {
		middleware.WriteInternalError(w, "Failed to update program status")
		return
	}

	updated, err := h.repo.GetByID(ctx, id)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to retrieve updated program")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// ListLoggedSessions handles GET /programs/{id}/logged-sessions
func (h *ProgramHandler) ListLoggedSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	// Verify the program exists
	if _, err := h.repo.GetByID(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	sessions, err := h.logRepo.ListDistinctLoggedSessionsByProgramID(ctx, id)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to list logged sessions")
		return
	}

	if sessions == nil {
		sessions = []string{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"sessions": sessions,
	})
}
