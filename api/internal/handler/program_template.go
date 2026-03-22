package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
)

// programTemplateEntryRequest represents a program template entry in the request body
type programTemplateEntryRequest struct {
	ExerciseName string          `json:"exercise_name"`
	Order        int             `json:"order"`
	Sets         *int            `json:"sets,omitempty"`
	Reps         *int            `json:"reps,omitempty"`
	RPE          *int            `json:"rpe,omitempty"`
	Percent1RM   *float64        `json:"percent_1rm,omitempty"`
	Notes        *string         `json:"notes,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}

// ProgramTemplateHandler handles ProgramTemplate-related HTTP requests
type ProgramTemplateHandler struct {
	repo        repository.ProgramTemplateRepository
	programRepo repository.ProgramRepository
}

// NewProgramTemplateHandler creates a new ProgramTemplateHandler
func NewProgramTemplateHandler(repo repository.ProgramTemplateRepository, programRepo repository.ProgramRepository) *ProgramTemplateHandler {
	return &ProgramTemplateHandler{
		repo:        repo,
		programRepo: programRepo,
	}
}

// CreateProgramTemplate handles POST /program-templates
func (h *ProgramTemplateHandler) CreateProgramTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name        string                        `json:"name"`
		Description *string                       `json:"description,omitempty"`
		Notes       *string                       `json:"notes,omitempty"`
		Metadata    json.RawMessage               `json:"metadata,omitempty"`
		Weeks       *string                       `json:"weeks,omitempty"`
		DaysPerWeek *string                       `json:"days_per_week,omitempty"`
		Entries     []programTemplateEntryRequest `json:"entries,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	var createdBy *string
	if userID := middleware.GetUserID(ctx); userID != "" {
		createdBy = &userID
	}
	tmpl := &domain.ProgramTemplate{
		Name:        req.Name,
		Description: req.Description,
		Notes:       req.Notes,
		Metadata:    req.Metadata,
		Weeks:       req.Weeks,
		DaysPerWeek: req.DaysPerWeek,
		CreatedBy:   createdBy,
		Entries:     make([]domain.ProgramTemplateEntry, len(req.Entries)),
	}

	for i, e := range req.Entries {
		tmpl.Entries[i] = convertProgramTemplateEntry(e)
	}

	if err := domain.ValidateProgramTemplate(tmpl); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := h.repo.Create(ctx, tmpl); err != nil {
		middleware.WriteInternalError(w, "Failed to create program template")
		return
	}

	writeJSON(w, http.StatusCreated, tmpl)
}

func convertProgramTemplateEntry(req programTemplateEntryRequest) domain.ProgramTemplateEntry {
	return domain.ProgramTemplateEntry{
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

// GetProgramTemplate handles GET /program-templates/{id}
func (h *ProgramTemplateHandler) GetProgramTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program_template")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	tmpl, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program template not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program template")
		return
	}

	writeJSON(w, http.StatusOK, tmpl)
}

// ArchiveProgramTemplate handles POST /program-templates/{id}/archive
func (h *ProgramTemplateHandler) ArchiveProgramTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program_template")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if err := h.repo.Archive(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program template not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to archive program template")
		return
	}

	tmpl, err := h.repo.GetByID(ctx, id)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to retrieve program template")
		return
	}

	writeJSON(w, http.StatusOK, tmpl)
}

// UnarchiveProgramTemplate handles POST /program-templates/{id}/unarchive
func (h *ProgramTemplateHandler) UnarchiveProgramTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program_template")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if err := h.repo.Unarchive(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program template not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to unarchive program template")
		return
	}

	tmpl, err := h.repo.GetByID(ctx, id)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to retrieve program template")
		return
	}

	writeJSON(w, http.StatusOK, tmpl)
}

// DuplicateProgramTemplate handles POST /program-templates/{id}/duplicate
func (h *ProgramTemplateHandler) DuplicateProgramTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program_template")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	original, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program template not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program template")
		return
	}

	// Parse optional request body for custom name
	copyName := strings.TrimSpace(original.Name) + " (copy)"
	if r.ContentLength != 0 {
		var req struct {
			Name string `json:"name,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
			middleware.WriteValidationError(w, "Invalid request body", nil)
			return
		}
		if req.Name != "" {
			copyName = req.Name
		}
	}

	entries := make([]domain.ProgramTemplateEntry, len(original.Entries))
	for i, e := range original.Entries {
		entries[i] = domain.ProgramTemplateEntry{
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

	var duplicatedBy *string
	if userID := middleware.GetUserID(ctx); userID != "" {
		duplicatedBy = &userID
	}

	duplicate := &domain.ProgramTemplate{
		Name:        copyName,
		Description: original.Description,
		Notes:       original.Notes,
		Metadata:    copyMeta,
		Weeks:       original.Weeks,
		DaysPerWeek: original.DaysPerWeek,
		CreatedBy:   duplicatedBy,
		Entries:     entries,
	}

	if err := domain.ValidateProgramTemplate(duplicate); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := h.repo.Create(ctx, duplicate); err != nil {
		middleware.WriteInternalError(w, "Failed to duplicate program template")
		return
	}

	writeJSON(w, http.StatusCreated, duplicate)
}

// GenerateProgram handles POST /program-templates/{id}/generate
func (h *ProgramTemplateHandler) GenerateProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program_template")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	tmpl, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program template not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program template")
		return
	}

	var req struct {
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

	if req.TargetWeights == nil {
		req.TargetWeights = make(map[string]float64)
	}

	input := &domain.GenerateProgramInput{
		Name:           req.Name,
		TargetWeights:  req.TargetWeights,
		LoadIncrements: req.LoadIncrements,
	}

	program := domain.GenerateProgramFromTemplate(tmpl, input)

	if err := domain.ValidateProgram(program); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Generated program validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := h.programRepo.Create(ctx, program); err != nil {
		middleware.WriteInternalError(w, "Failed to create generated program")
		return
	}

	writeJSON(w, http.StatusCreated, program)
}

// DeleteProgramTemplate handles DELETE /program-templates/{id}
func (h *ProgramTemplateHandler) DeleteProgramTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program_template")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	if _, err := h.repo.GetByID(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program template not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program template")
		return
	}

	// Check referential integrity: ProgramTemplate referenced by Programs
	hasPrograms, err := h.programRepo.ExistsByProgramTemplateID(ctx, id)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to check references")
		return
	}
	if hasPrograms {
		middleware.WriteConflictError(w, "Cannot delete program template - referenced by programs", map[string]interface{}{
			"blocking_references": []map[string]interface{}{
				{
					"type": "program",
				},
			},
		})
		return
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program template not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete program template")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListProgramTemplates handles GET /program-templates
func (h *ProgramTemplateHandler) ListProgramTemplates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr, 1, 100); err == nil {
			limit = parsedLimit
		}
	}
	after := r.URL.Query().Get("after")
	includeArchived := r.URL.Query().Get("include_archived") == "true"

	templates, nextCursor, hasMore, err := h.repo.List(ctx, limit, after, includeArchived)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to list program templates")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":        templates,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}
