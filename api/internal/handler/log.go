package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

// logRequest represents the request body for creating/updating a log
type logRequest struct {
	ProgramID   *string           `json:"program_id,omitempty"`
	SessionName *string           `json:"session_name,omitempty"`
	PerformedAt string            `json:"performed_at"`
	StartedAt   *string           `json:"started_at,omitempty"`
	FinishedAt  *string           `json:"finished_at,omitempty"`
	Notes       *string           `json:"notes,omitempty"`
	Metadata    json.RawMessage   `json:"metadata,omitempty"`
	Entries     []logEntryRequest `json:"entries"`
}

// logEntryRequest represents a log entry in the request body
type logEntryRequest struct {
	ExerciseName   string                 `json:"exercise_name"`
	Fields         map[string]interface{} `json:"fields,omitempty"`
	Notes          *string                `json:"notes,omitempty"`
	VideoObjectKey *string                `json:"video_object_key,omitempty"`
	StartedAt      *string                `json:"started_at,omitempty"`
	FinishedAt     *string                `json:"finished_at,omitempty"`
}

// LogHandler handles Log-related HTTP requests
type LogHandler struct {
	repo        repository.LogRepository
	programRepo repository.ProgramRepository
}

// NewLogHandler creates a new LogHandler
func NewLogHandler(repo repository.LogRepository, programRepo repository.ProgramRepository) *LogHandler {
	return &LogHandler{
		repo:        repo,
		programRepo: programRepo,
	}
}

// parseLogRequest parses and validates a log request into a domain model
func (h *LogHandler) parseLogRequest(req *logRequest) (*domain.Log, error) {
	performedAt, err := time.Parse(time.RFC3339, req.PerformedAt)
	if err != nil {
		return nil, &domain.ValidationError{
			Field:   "performed_at",
			Message: "invalid timestamp format: " + err.Error(),
		}
	}

	log := &domain.Log{
		SessionName: req.SessionName,
		PerformedAt: performedAt,
		Notes:       req.Notes,
		Metadata:    req.Metadata,
	}

	if req.StartedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.StartedAt)
		if err != nil {
			return nil, &domain.ValidationError{
				Field:   "started_at",
				Message: "invalid timestamp format: " + err.Error(),
			}
		}
		log.StartedAt = &t
	}

	if req.FinishedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.FinishedAt)
		if err != nil {
			return nil, &domain.ValidationError{
				Field:   "finished_at",
				Message: "invalid timestamp format: " + err.Error(),
			}
		}
		log.FinishedAt = &t
	}

	if req.ProgramID != nil {
		programID, err := uuid.Parse(*req.ProgramID)
		if err != nil {
			return nil, &domain.ValidationError{
				Field:   "program_id",
				Message: "invalid UUID format: " + err.Error(),
			}
		}
		log.ProgramID = &programID
	}

	entries := make([]domain.LogEntry, len(req.Entries))
	for i, entryReq := range req.Entries {
		entry := domain.LogEntry{
			ExerciseName:   entryReq.ExerciseName,
			Order:          i,
			Fields:         entryReq.Fields,
			Notes:          entryReq.Notes,
			VideoObjectKey: entryReq.VideoObjectKey,
		}

		if entryReq.StartedAt != nil {
			t, err := time.Parse(time.RFC3339, *entryReq.StartedAt)
			if err != nil {
				return nil, &domain.ValidationError{
					Field:   "entries[].started_at",
					Message: "invalid timestamp format: " + err.Error(),
				}
			}
			entry.StartedAt = &t
		}

		if entryReq.FinishedAt != nil {
			t, err := time.Parse(time.RFC3339, *entryReq.FinishedAt)
			if err != nil {
				return nil, &domain.ValidationError{
					Field:   "entries[].finished_at",
					Message: "invalid timestamp format: " + err.Error(),
				}
			}
			entry.FinishedAt = &t
		}

		entries[i] = entry
	}
	log.Entries = entries

	return log, nil
}

// CreateLog handles POST /logs
func (h *LogHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req logRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	log, err := h.parseLogRequest(&req)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteInternalError(w, "Failed to parse log request")
		return
	}

	if err := domain.ValidateLog(log); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Validate program exists before creating log
	if log.ProgramID != nil {
		if _, err := h.programRepo.GetByID(ctx, *log.ProgramID); err != nil {
			if err == domain.ErrNotFound {
				middleware.WriteValidationError(w, "Program not found", map[string]interface{}{
					"program_id": log.ProgramID.String(),
				})
				return
			}
			middleware.WriteInternalError(w, "Failed to retrieve program")
			return
		}
	}

	// Enforce one log per (program_id, session_name) pair
	if log.ProgramID != nil && log.SessionName != nil {
		exists, err := h.repo.ExistsByProgramIDAndSessionName(ctx, *log.ProgramID, *log.SessionName)
		if err != nil {
			middleware.WriteInternalError(w, "Failed to check for duplicate log")
			return
		}
		if exists {
			middleware.WriteConflictError(w, "A log already exists for this session", map[string]interface{}{
				"program_id":   log.ProgramID.String(),
				"session_name": *log.SessionName,
			})
			return
		}
	}

	if err := h.repo.Create(ctx, log); err != nil {
		middleware.WriteInternalError(w, "Failed to create log")
		return
	}

	// Auto-transition program from "created" to "ongoing" after successful log creation
	if log.ProgramID != nil {
		program, err := h.programRepo.GetByID(ctx, *log.ProgramID)
		if err == nil && program.Status == domain.ProgramStatusCreated {
			_ = h.programRepo.UpdateStatus(ctx, *log.ProgramID, domain.ProgramStatusOngoing)
		}
	}

	writeJSON(w, http.StatusCreated, log)
}

// GetLog handles GET /logs/{id}
func (h *LogHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "log")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	log, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Log not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve log")
		return
	}

	writeJSON(w, http.StatusOK, log)
}

// UpdateLog handles PUT /logs/{id}
func (h *LogHandler) UpdateLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "log")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Log not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve log")
		return
	}

	var req logRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	updated, err := h.parseLogRequest(&req)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteInternalError(w, "Failed to parse log request")
		return
	}

	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt

	if err := domain.ValidateLog(updated); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Enforce one log per (program_id, session_name) pair, excluding the log being updated
	if updated.ProgramID != nil && updated.SessionName != nil {
		exists, err := h.repo.ExistsByProgramIDAndSessionNameExcluding(ctx, *updated.ProgramID, *updated.SessionName, updated.ID)
		if err != nil {
			middleware.WriteInternalError(w, "Failed to check for duplicate log")
			return
		}
		if exists {
			middleware.WriteConflictError(w, "A log already exists for this session", map[string]interface{}{
				"program_id":   updated.ProgramID.String(),
				"session_name": *updated.SessionName,
			})
			return
		}
	}

	if err := h.repo.Update(ctx, updated); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Log not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update log")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// DeleteLog handles DELETE /logs/{id}
func (h *LogHandler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "log")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if _, err := h.repo.GetByID(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Log not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve log")
		return
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Log not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete log")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListLogs handles GET /logs
func (h *LogHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr, 1, 100); err == nil {
			limit = parsedLimit
		}
	}
	after := r.URL.Query().Get("after")

	var programID *uuid.UUID
	if pidStr := r.URL.Query().Get("program_id"); pidStr != "" {
		pid, err := uuid.Parse(pidStr)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid program_id format", nil)
			return
		}
		programID = &pid
	}

	var performedAtFrom, performedAtTo *time.Time
	if fromStr := r.URL.Query().Get("performed_at_from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			performedAtFrom = &t
		}
	}
	if toStr := r.URL.Query().Get("performed_at_to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			performedAtTo = &t
		}
	}

	var logs []*domain.Log
	var nextCursor string
	var hasMore bool
	var err error

	if performedAtFrom != nil || performedAtTo != nil {
		logs, nextCursor, hasMore, err = h.repo.ListByPerformedAtRange(ctx, programID, performedAtFrom, performedAtTo, limit, after)
	} else {
		logs, nextCursor, hasMore, err = h.repo.List(ctx, programID, limit, after)
	}

	if err != nil {
		middleware.WriteInternalError(w, "Failed to list logs")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":        logs,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}
