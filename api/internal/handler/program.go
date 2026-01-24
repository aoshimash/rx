package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/aoshimash/optel-training/api/internal/middleware"
	"github.com/aoshimash/optel-training/api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ProgramHandler handles Program-related HTTP requests
type ProgramHandler struct {
	repo         repository.ProgramRepository
	exerciseRepo repository.ExerciseRepository
	workoutRepo  repository.WorkoutRepository
}

// NewProgramHandler creates a new ProgramHandler
func NewProgramHandler(repo repository.ProgramRepository, exerciseRepo repository.ExerciseRepository, workoutRepo repository.WorkoutRepository) *ProgramHandler {
	return &ProgramHandler{
		repo:         repo,
		exerciseRepo: exerciseRepo,
		workoutRepo:  workoutRepo,
	}
}

// CreateProgram handles POST /programs
func (h *ProgramHandler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request body (OpenAPI ProgramCreate type)
	var req struct {
		Name        string `json:"name"`
		Description *string `json:"description,omitempty"`
		RootNodes   []struct {
			Name              string  `json:"name"`
			NodeType          string  `json:"node_type"`
			Order             int     `json:"order"`
			Children          []json.RawMessage `json:"children,omitempty"` // Recursive structure
			ExerciseID        *string `json:"exercise_id,omitempty"`
			TargetSets        *int    `json:"target_sets,omitempty"`
			TargetReps        *int    `json:"target_reps,omitempty"`
			TargetRPE         *int    `json:"target_rpe,omitempty"`
			Percent1RM        *float64 `json:"percent_1rm,omitempty"`
			PlannedRestSeconds *int    `json:"planned_rest_seconds,omitempty"`
			MuscleGroups      []string `json:"muscle_groups,omitempty"`
			Notes             *string  `json:"notes,omitempty"`
		} `json:"root_nodes,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Convert to domain model
	program := &domain.Program{
		Name:        req.Name,
		Description: req.Description,
		RootNodes:   make([]domain.ProgramNode, len(req.RootNodes)),
	}

	// Convert root nodes recursively
	for i, nodeReq := range req.RootNodes {
		node, err := h.convertNode(ctx, nodeReq, nil)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid program node", map[string]interface{}{
				"field": "root_nodes",
				"index": i,
				"error": err.Error(),
			})
			return
		}
		program.RootNodes[i] = *node
	}

	// Validate
	if err := domain.ValidateProgram(program); err != nil {
		if ve, ok := err.(*domain.ValidationError); ok {
			middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
				"field":   ve.Field,
				"message": ve.Message,
			})
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Create
	if err := h.repo.Create(ctx, program); err != nil {
		middleware.WriteInternalError(w, "Failed to create program")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(program)
}

// convertNode converts a JSON node request to a domain ProgramNode recursively
func (h *ProgramHandler) convertNode(ctx context.Context, nodeReq struct {
	Name              string          `json:"name"`
	NodeType          string          `json:"node_type"`
	Order             int             `json:"order"`
	Children          []json.RawMessage `json:"children,omitempty"`
	ExerciseID        *string         `json:"exercise_id,omitempty"`
	TargetSets        *int            `json:"target_sets,omitempty"`
	TargetReps        *int            `json:"target_reps,omitempty"`
	TargetRPE         *int            `json:"target_rpe,omitempty"`
	Percent1RM        *float64        `json:"percent_1rm,omitempty"`
	PlannedRestSeconds *int           `json:"planned_rest_seconds,omitempty"`
	MuscleGroups      []string        `json:"muscle_groups,omitempty"`
	Notes             *string         `json:"notes,omitempty"`
}, parentID *uuid.UUID) (*domain.ProgramNode, error) {
	node := &domain.ProgramNode{
		Name:         nodeReq.Name,
		NodeType:     nodeReq.NodeType,
		Order:        nodeReq.Order,
		ParentID:     parentID,
		TargetSets:   nodeReq.TargetSets,
		TargetReps:   nodeReq.TargetReps,
		TargetRPE:    nodeReq.TargetRPE,
		Percent1RM:   nodeReq.Percent1RM,
		PlannedRestSeconds: nodeReq.PlannedRestSeconds,
		MuscleGroups: nodeReq.MuscleGroups,
		Notes:        nodeReq.Notes,
	}

	// Parse exercise_id if provided
	if nodeReq.ExerciseID != nil {
		exerciseID, err := uuid.Parse(*nodeReq.ExerciseID)
		if err != nil {
			return nil, err
		}
		// Validate Exercise exists
		if _, err := h.exerciseRepo.GetByID(ctx, exerciseID); err != nil {
			return nil, err
		}
		node.ExerciseID = &exerciseID
	}

	// Convert children recursively
	if len(nodeReq.Children) > 0 {
		node.Children = make([]domain.ProgramNode, len(nodeReq.Children))
		for i, childJSON := range nodeReq.Children {
			var childReq struct {
				Name              string          `json:"name"`
				NodeType          string          `json:"node_type"`
				Order             int             `json:"order"`
				Children          []json.RawMessage `json:"children,omitempty"`
				ExerciseID        *string         `json:"exercise_id,omitempty"`
				TargetSets        *int            `json:"target_sets,omitempty"`
				TargetReps        *int            `json:"target_reps,omitempty"`
				TargetRPE         *int            `json:"target_rpe,omitempty"`
				Percent1RM        *float64        `json:"percent_1rm,omitempty"`
				PlannedRestSeconds *int           `json:"planned_rest_seconds,omitempty"`
				MuscleGroups      []string        `json:"muscle_groups,omitempty"`
				Notes             *string         `json:"notes,omitempty"`
			}
			if err := json.Unmarshal(childJSON, &childReq); err != nil {
				return nil, err
			}
			child, err := h.convertNode(ctx, childReq, &node.ID)
			if err != nil {
				return nil, err
			}
			node.Children[i] = *child
		}
	}

	return node, nil
}

// GetProgram handles GET /programs/{id}
func (h *ProgramHandler) GetProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from path
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteValidationError(w, "Invalid program ID format", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	// Get from repository
	program, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(program)
}

// UpdateProgram handles PUT /programs/{id}
func (h *ProgramHandler) UpdateProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from path
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteValidationError(w, "Invalid program ID format", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	// Check if program exists
	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	// Decode request body (full replacement, same structure as Create)
	var req struct {
		Name        string `json:"name"`
		Description *string `json:"description,omitempty"`
		RootNodes   []struct {
			Name              string          `json:"name"`
			NodeType          string          `json:"node_type"`
			Order             int             `json:"order"`
			Children          []json.RawMessage `json:"children,omitempty"`
			ExerciseID        *string         `json:"exercise_id,omitempty"`
			TargetSets        *int            `json:"target_sets,omitempty"`
			TargetReps        *int            `json:"target_reps,omitempty"`
			TargetRPE         *int            `json:"target_rpe,omitempty"`
			Percent1RM        *float64        `json:"percent_1rm,omitempty"`
			PlannedRestSeconds *int           `json:"planned_rest_seconds,omitempty"`
			MuscleGroups      []string       `json:"muscle_groups,omitempty"`
			Notes             *string        `json:"notes,omitempty"`
		} `json:"root_nodes,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Update existing program (full replacement)
	existing.Name = req.Name
	existing.Description = req.Description
	existing.RootNodes = make([]domain.ProgramNode, len(req.RootNodes))

	// Convert root nodes recursively
	for i, nodeReq := range req.RootNodes {
		node, err := h.convertNode(ctx, nodeReq, nil)
		if err != nil {
			middleware.WriteValidationError(w, "Invalid program node", map[string]interface{}{
				"field": "root_nodes",
				"index": i,
				"error": err.Error(),
			})
			return
		}
		existing.RootNodes[i] = *node
	}

	// Validate
	if err := domain.ValidateProgram(existing); err != nil {
		if ve, ok := err.(*domain.ValidationError); ok {
			middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
				"field":   ve.Field,
				"message": ve.Message,
			})
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Update
	if err := h.repo.Update(ctx, existing); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to update program")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existing)
}

// DeleteProgram handles DELETE /programs/{id}
func (h *ProgramHandler) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from path
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteValidationError(w, "Invalid program ID format", map[string]interface{}{
			"id": idStr,
		})
		return
	}

	// Check if program exists
	program, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program")
		return
	}

	// Check for referential integrity: Program referenced by Workout (FR-026)
	// We need to check all ProgramNodes in the tree
	var programNodeIDs []uuid.UUID
	var collectNodeIDs func(nodes []domain.ProgramNode)
	collectNodeIDs = func(nodes []domain.ProgramNode) {
		for _, node := range nodes {
			programNodeIDs = append(programNodeIDs, node.ID)
			if len(node.Children) > 0 {
				collectNodeIDs(node.Children)
			}
		}
	}
	collectNodeIDs(program.RootNodes)

	// Check if any Workout references any ProgramNode in this Program
	for _, nodeID := range programNodeIDs {
		referencingWorkouts, err := h.workoutRepo.ListByProgramNodeID(ctx, nodeID)
		if err != nil {
			middleware.WriteInternalError(w, "Failed to check references")
			return
		}
		if len(referencingWorkouts) > 0 {
			middleware.WriteConflictError(w, "Cannot delete program - referenced by workouts", map[string]interface{}{
				"blocking_references": []map[string]interface{}{
					{
						"type":  "workout",
						"count": len(referencingWorkouts),
					},
				},
			})
			return
		}
	}

	// Delete (cascades to ProgramNode records per FR-026)
	if err := h.repo.Delete(ctx, id); err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to delete program")
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
