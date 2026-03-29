package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

// planSessionEntryRequest represents a plan session entry in the request body
type planSessionEntryRequest struct {
	ExerciseName string                 `json:"exercise_name"`
	Order        int                    `json:"order"`
	Fields       map[string]interface{} `json:"fields,omitempty"`
	Notes        *string                `json:"notes,omitempty"`
}

// planSessionRequest represents a plan session in the request body
type planSessionRequest struct {
	SessionName     string                    `json:"session_name"`
	Order           int                       `json:"order"`
	Date            *string                   `json:"date,omitempty"`
	SourceProgramID *string                   `json:"source_program_id,omitempty"`
	SourceSessionID *string                   `json:"source_session_id,omitempty"`
	Entries         []planSessionEntryRequest `json:"entries,omitempty"`
}

// PlanHandler handles Plan-related HTTP requests
type PlanHandler struct {
	planRepo    repository.PlanRepository
	programRepo repository.ProgramRepository
}

// NewPlanHandler creates a new PlanHandler
func NewPlanHandler(planRepo repository.PlanRepository, programRepo repository.ProgramRepository) *PlanHandler {
	return &PlanHandler{
		planRepo:    planRepo,
		programRepo: programRepo,
	}
}

// parsePlanSession converts a plan session request to a domain PlanSession
func parsePlanSession(req planSessionRequest) (domain.PlanSession, error) {
	sess := domain.PlanSession{
		SessionName: req.SessionName,
		Order:       req.Order,
		Entries:     make([]domain.PlanSessionEntry, len(req.Entries)),
	}

	if req.Date != nil {
		var d domain.DateOnly
		if err := json.Unmarshal([]byte(`"`+*req.Date+`"`), &d); err != nil {
			return sess, &domain.ValidationError{
				Field:   "date",
				Message: "invalid date format (expected YYYY-MM-DD)",
			}
		}
		sess.Date = &d
	}

	if req.SourceProgramID != nil {
		pid, err := uuid.Parse(*req.SourceProgramID)
		if err != nil {
			return sess, &domain.ValidationError{
				Field:   "source_program_id",
				Message: "invalid UUID format: " + err.Error(),
			}
		}
		sess.SourceProgramID = &pid
	}

	if req.SourceSessionID != nil {
		sid, err := uuid.Parse(*req.SourceSessionID)
		if err != nil {
			return sess, &domain.ValidationError{
				Field:   "source_session_id",
				Message: "invalid UUID format: " + err.Error(),
			}
		}
		sess.SourceSessionID = &sid
	}

	for j, e := range req.Entries {
		sess.Entries[j] = domain.PlanSessionEntry{
			ExerciseName: e.ExerciseName,
			Order:        e.Order,
			Fields:       e.Fields,
			Notes:        e.Notes,
		}
	}

	return sess, nil
}

// parsePlanSessions converts a slice of plan session requests to domain PlanSessions
func parsePlanSessions(reqs []planSessionRequest) ([]domain.PlanSession, error) {
	result := make([]domain.PlanSession, len(reqs))
	for i, req := range reqs {
		sess, err := parsePlanSession(req)
		if err != nil {
			return nil, err
		}
		result[i] = sess
	}
	return result, nil
}

// GetPlan handles GET /plan
func (h *PlanHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	plan, err := h.planRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			middleware.WriteNotFoundError(w, "Plan not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve plan")
		return
	}

	writeJSON(w, http.StatusOK, plan)
}

// CreatePlan handles POST /plan
func (h *PlanHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	var req struct {
		Name     *string              `json:"name,omitempty"`
		Notes    *string              `json:"notes,omitempty"`
		Sessions []planSessionRequest `json:"sessions,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	plan := &domain.Plan{
		Name:  req.Name,
		Notes: req.Notes,
	}

	if len(req.Sessions) > 0 {
		sessions, err := parsePlanSessions(req.Sessions)
		if err != nil {
			if handleValidationError(w, err) {
				return
			}
			middleware.WriteValidationError(w, "Invalid sessions", nil)
			return
		}
		plan.Sessions = sessions
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

	if err := h.planRepo.Create(ctx, userID, plan); err != nil {
		if de, ok := err.(*domain.DomainError); ok && de.Code == domain.ErrorCodeConflict {
			middleware.WriteConflictError(w, "A plan already exists for this user", nil)
			return
		}
		middleware.WriteInternalError(w, "Failed to create plan")
		return
	}

	// Re-read to get generated IDs
	created, err := h.planRepo.GetByUserID(ctx, userID)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to retrieve created plan")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// UpdatePlan handles PUT /plan
func (h *PlanHandler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	var req struct {
		Name     *string              `json:"name,omitempty"`
		Notes    *string              `json:"notes,omitempty"`
		Sessions []planSessionRequest `json:"sessions,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	existing, err := h.planRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			middleware.WriteNotFoundError(w, "Plan not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve plan")
		return
	}

	plan := &domain.Plan{
		Name:  req.Name,
		Notes: req.Notes,
	}

	if req.Sessions != nil {
		sessions, err := parsePlanSessions(req.Sessions)
		if err != nil {
			if handleValidationError(w, err) {
				return
			}
			middleware.WriteValidationError(w, "Invalid sessions", nil)
			return
		}
		plan.Sessions = sessions
	} else {
		plan.Sessions = existing.Sessions
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

	if err := h.planRepo.Update(ctx, userID, plan); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			middleware.WriteNotFoundError(w, "Plan not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update plan")
		return
	}

	updated, err := h.planRepo.GetByUserID(ctx, userID)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to retrieve updated plan")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// DeletePlan handles DELETE /plan
func (h *PlanHandler) DeletePlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	if err := h.planRepo.Delete(ctx, userID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			middleware.WriteNotFoundError(w, "Plan not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete plan")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddPlanSessions handles POST /plan/sessions
func (h *PlanHandler) AddPlanSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	var req struct {
		Sessions []planSessionRequest `json:"sessions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if len(req.Sessions) == 0 {
		middleware.WriteValidationError(w, "At least one session is required", nil)
		return
	}

	sessions, err := parsePlanSessions(req.Sessions)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Invalid sessions", nil)
		return
	}

	// Validate each session
	for i := range sessions {
		if err := domain.ValidatePlanSession(&sessions[i]); err != nil {
			if handleValidationError(w, err) {
				return
			}
			middleware.WriteValidationError(w, "Validation failed", nil)
			return
		}
	}

	if err := h.planRepo.AddSessions(ctx, userID, sessions); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			middleware.WriteNotFoundError(w, "Plan not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to add sessions")
		return
	}

	plan, err := h.planRepo.GetByUserID(ctx, userID)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to retrieve updated plan")
		return
	}

	writeJSON(w, http.StatusOK, plan)
}

// UpdatePlanSession handles PUT /plan/sessions/{session_id}
func (h *PlanHandler) UpdatePlanSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	sessionID, err := parseUUIDParam(r, "session_id", "session")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	var req planSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	sess, err := parsePlanSession(req)
	if err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Invalid session", nil)
		return
	}
	sess.ID = sessionID

	if err := domain.ValidatePlanSession(&sess); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", nil)
		return
	}

	if err := h.planRepo.UpdateSession(ctx, userID, &sess); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			middleware.WriteNotFoundError(w, "Plan or session not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update session")
		return
	}

	writeJSON(w, http.StatusOK, &sess)
}

// DeletePlanSession handles DELETE /plan/sessions/{session_id}
func (h *PlanHandler) DeletePlanSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	sessionID, err := parseUUIDParam(r, "session_id", "session")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if err := h.planRepo.DeleteSession(ctx, userID, sessionID.String()); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			middleware.WriteNotFoundError(w, "Plan or session not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ExpandProgram handles POST /plan/expand-program/{program_id}
func (h *PlanHandler) ExpandProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	programID, err := parseUUIDParam(r, "program_id", "program")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	// Load the program
	program, err := h.programRepo.GetByID(ctx, programID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	// Convert program sessions to plan sessions
	planSessions := make([]domain.PlanSession, len(program.Sessions))
	for i, ps := range program.Sessions {
		entries := make([]domain.PlanSessionEntry, len(ps.Entries))
		for j, e := range ps.Entries {
			entries[j] = domain.PlanSessionEntry{
				Order:        e.Order,
				ExerciseName: e.ExerciseName,
				Fields:       e.Fields,
				Notes:        e.Notes,
			}
		}

		srcProgramID := program.ID
		srcSessionID := ps.ID
		planSessions[i] = domain.PlanSession{
			SessionName:     ps.SessionName,
			Order:           ps.Order,
			Date:            ps.Date,
			SourceProgramID: &srcProgramID,
			SourceSessionID: &srcSessionID,
			Entries:         entries,
		}
	}

	// Check if a plan already exists
	_, err = h.planRepo.GetByUserID(ctx, userID)
	if errors.Is(err, domain.ErrNotFound) {
		// Create a new plan with the program's sessions
		plan := &domain.Plan{
			Sessions: planSessions,
		}
		if err := h.planRepo.Create(ctx, userID, plan); err != nil {
			middleware.WriteInternalError(w, "Failed to create plan")
			return
		}
	} else if err != nil {
		middleware.WriteInternalError(w, "Failed to check existing plan")
		return
	} else {
		// Append sessions to existing plan
		if err := h.planRepo.AddSessions(ctx, userID, planSessions); err != nil {
			middleware.WriteInternalError(w, "Failed to add sessions to plan")
			return
		}
	}

	plan, err := h.planRepo.GetByUserID(ctx, userID)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to retrieve plan")
		return
	}

	writeJSON(w, http.StatusOK, plan)
}
