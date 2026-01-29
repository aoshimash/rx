package postgres

import (
	"context"
	"log/slog"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/aoshimash/optel-training/api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// programRepository implements ProgramRepository with PostgreSQL
type programRepository struct {
	pool *pgxpool.Pool
}

// NewProgramRepository creates a new PostgreSQL Program repository
func NewProgramRepository(pool *pgxpool.Pool) repository.ProgramRepository {
	return &programRepository{pool: pool}
}

func (r *programRepository) Create(ctx context.Context, program *domain.Program) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	id := uuid.New()
	if program.ID != uuid.Nil {
		id = program.ID
	}

	// Insert program
	query := `
		INSERT INTO programs (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err = tx.QueryRow(ctx, query, id, program.Name, program.Description).Scan(
		&program.CreatedAt, &program.UpdatedAt,
	)
	if err != nil {
		slog.Error("Failed to create program", "error", err)
		return err
	}

	program.ID = id

	// Insert program nodes recursively
	if len(program.RootNodes) > 0 {
		err = r.insertNodes(ctx, tx, id, nil, program.RootNodes)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *programRepository) insertNodes(ctx context.Context, tx pgx.Tx, programID uuid.UUID, parentID *uuid.UUID, nodes []domain.ProgramNode) error {
	for i := range nodes {
		nodeID := uuid.New()
		if nodes[i].ID != uuid.Nil {
			nodeID = nodes[i].ID
		}

		query := `
			INSERT INTO program_nodes (
				id, program_id, parent_id, name, node_type, "order",
				exercise_id, target_sets, target_reps, target_rpe,
				percent_1rm, planned_rest_seconds, muscle_groups, notes
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		`

		_, err := tx.Exec(ctx, query,
			nodeID,
			programID,
			parentID,
			nodes[i].Name,
			nodes[i].NodeType,
			nodes[i].Order,
			nodes[i].ExerciseID,
			nodes[i].TargetSets,
			nodes[i].TargetReps,
			nodes[i].TargetRPE,
			nodes[i].Percent1RM,
			nodes[i].PlannedRestSeconds,
			nodes[i].MuscleGroups,
			nodes[i].Notes,
		)
		if err != nil {
			slog.Error("Failed to insert program node", "error", err)
			return err
		}

		nodes[i].ID = nodeID
		nodes[i].ProgramID = programID
		nodes[i].ParentID = parentID

		// Recursively insert children
		if len(nodes[i].Children) > 0 {
			err = r.insertNodes(ctx, tx, programID, &nodeID, nodes[i].Children)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *programRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error) {
	// Get program
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM programs
		WHERE id = $1
	`

	var program domain.Program
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&program.ID,
		&program.Name,
		&program.Description,
		&program.CreatedAt,
		&program.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to get program by ID", "id", id, "error", err)
		return nil, err
	}

	// Load root nodes for the program
	rootNodes, err := r.getNodesForProgram(ctx, id)
	if err != nil {
		return nil, err
	}
	program.RootNodes = rootNodes

	return &program, nil
}

// getNodesForProgram loads and builds the tree structure of nodes for a program
func (r *programRepository) getNodesForProgram(ctx context.Context, programID uuid.UUID) ([]domain.ProgramNode, error) {
	nodesQuery := `
		SELECT id, program_id, parent_id, name, node_type, "order",
		       exercise_id, target_sets, target_reps, target_rpe,
		       percent_1rm, planned_rest_seconds, muscle_groups, notes
		FROM program_nodes
		WHERE program_id = $1
		ORDER BY "order" ASC
	`

	rows, err := r.pool.Query(ctx, nodesQuery, programID)
	if err != nil {
		slog.Error("Failed to get program nodes", "programID", programID, "error", err)
		return nil, err
	}
	defer rows.Close()

	// Build node map and track parent-child relationships
	nodeMap := make(map[uuid.UUID]*domain.ProgramNode)
	childrenMap := make(map[uuid.UUID][]uuid.UUID) // parent ID -> child IDs
	var rootNodeIDs []uuid.UUID

	for rows.Next() {
		var node domain.ProgramNode
		err := rows.Scan(
			&node.ID,
			&node.ProgramID,
			&node.ParentID,
			&node.Name,
			&node.NodeType,
			&node.Order,
			&node.ExerciseID,
			&node.TargetSets,
			&node.TargetReps,
			&node.TargetRPE,
			&node.Percent1RM,
			&node.PlannedRestSeconds,
			&node.MuscleGroups,
			&node.Notes,
		)
		if err != nil {
			return nil, err
		}

		nodeMap[node.ID] = &node

		if node.ParentID == nil {
			rootNodeIDs = append(rootNodeIDs, node.ID)
		} else {
			childrenMap[*node.ParentID] = append(childrenMap[*node.ParentID], node.ID)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Build tree structure recursively (bottom-up approach ensures all children are included)
	var buildNode func(id uuid.UUID) domain.ProgramNode
	buildNode = func(id uuid.UUID) domain.ProgramNode {
		node := *nodeMap[id]
		childIDs := childrenMap[id]
		if len(childIDs) > 0 {
			node.Children = make([]domain.ProgramNode, 0, len(childIDs))
			for _, childID := range childIDs {
				node.Children = append(node.Children, buildNode(childID))
			}
		}
		return node
	}

	// Build root nodes with full tree structure
	rootNodes := make([]domain.ProgramNode, 0, len(rootNodeIDs))
	for _, rootID := range rootNodeIDs {
		rootNodes = append(rootNodes, buildNode(rootID))
	}

	return rootNodes, nil
}

func (r *programRepository) Update(ctx context.Context, program *domain.Program) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Update program
	query := `
		UPDATE programs
		SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err = tx.QueryRow(ctx, query, program.ID, program.Name, program.Description).Scan(&program.UpdatedAt)
	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to update program", "id", program.ID, "error", err)
		return err
	}

	// Delete existing nodes (CASCADE will handle children)
	_, err = tx.Exec(ctx, `DELETE FROM program_nodes WHERE program_id = $1`, program.ID)
	if err != nil {
		slog.Error("Failed to delete program nodes", "error", err)
		return err
	}

	// Insert new nodes
	if len(program.RootNodes) > 0 {
		err = r.insertNodes(ctx, tx, program.ID, nil, program.RootNodes)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *programRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM programs WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete program", "id", id, "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *programRepository) List(ctx context.Context, limit int, after string) ([]*domain.Program, string, bool, error) {
	var startID uuid.UUID
	if after != "" {
		var err error
		startID, err = decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
	}

	query := `
		SELECT id, name, description, created_at, updated_at
		FROM programs
		WHERE ($1::uuid IS NULL OR id > $1)
		ORDER BY id ASC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, startID, limit+1)
	if err != nil {
		slog.Error("Failed to list programs", "error", err)
		return nil, "", false, err
	}
	defer rows.Close()

	programs := make([]*domain.Program, 0, limit)

	for rows.Next() {
		var program domain.Program
		err := rows.Scan(
			&program.ID,
			&program.Name,
			&program.Description,
			&program.CreatedAt,
			&program.UpdatedAt,
		)
		if err != nil {
			return nil, "", false, err
		}

		programs = append(programs, &program)
	}

	if err := rows.Err(); err != nil {
		return nil, "", false, err
	}

	hasMore := len(programs) > limit
	if hasMore {
		programs = programs[:limit]
	}

	// Load RootNodes for each program (consistent with memory store behavior)
	for _, program := range programs {
		rootNodes, err := r.getNodesForProgram(ctx, program.ID)
		if err != nil {
			return nil, "", false, err
		}
		program.RootNodes = rootNodes
	}

	var nextCursor string
	if hasMore && len(programs) > 0 {
		// Use the last item in the returned set, not the extra item
		nextCursor = encodeCursor(programs[len(programs)-1].ID)
	}

	return programs, nextCursor, hasMore, nil
}
