package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// fieldGroupRepository implements FieldGroupRepository with PostgreSQL
type fieldGroupRepository struct {
	pool *pgxpool.Pool
}

// NewFieldGroupRepository creates a new PostgreSQL FieldGroup repository
func NewFieldGroupRepository(pool *pgxpool.Pool) repository.FieldGroupRepository {
	return &fieldGroupRepository{pool: pool}
}

func (r *fieldGroupRepository) List(ctx context.Context, userID string) ([]domain.FieldGroup, error) {
	query := `
		SELECT id, name, description, program_fields, log_fields, created_at, updated_at
		FROM field_groups
		WHERE user_id = $1
		ORDER BY name ASC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		slog.Error("Failed to list field groups", "userID", userID, "error", err)
		return nil, err
	}
	defer rows.Close()

	groups := make([]domain.FieldGroup, 0)
	for rows.Next() {
		fg, err := scanFieldGroup(rows)
		if err != nil {
			return nil, err
		}
		groups = append(groups, fg)
	}

	return groups, rows.Err()
}

func (r *fieldGroupRepository) GetByID(ctx context.Context, userID string, id uuid.UUID) (*domain.FieldGroup, error) {
	query := `
		SELECT id, name, description, program_fields, log_fields, created_at, updated_at
		FROM field_groups
		WHERE id = $1 AND user_id = $2
	`

	row := r.pool.QueryRow(ctx, query, id, userID)
	fg, err := scanFieldGroupRow(row)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to get field group", "id", id, "error", err)
		return nil, err
	}

	return &fg, nil
}

func (r *fieldGroupRepository) Create(ctx context.Context, userID string, fg *domain.FieldGroup) error {
	id := uuid.New()
	if fg.ID != uuid.Nil {
		id = fg.ID
	}

	programFieldsJSON, err := json.Marshal(fg.ProgramFields)
	if err != nil {
		return fmt.Errorf("marshal program_fields: %w", err)
	}
	logFieldsJSON, err := json.Marshal(fg.LogFields)
	if err != nil {
		return fmt.Errorf("marshal log_fields: %w", err)
	}

	query := `
		INSERT INTO field_groups (id, user_id, name, description, program_fields, log_fields, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err = r.pool.QueryRow(ctx, query,
		id,
		userID,
		fg.Name,
		fg.Description,
		programFieldsJSON,
		logFieldsJSON,
	).Scan(&fg.CreatedAt, &fg.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "field_groups_user_id_name_key") || strings.Contains(err.Error(), "duplicate key") {
			return &domain.DomainError{
				Code:    domain.ErrorCodeConflict,
				Message: "field group with this name already exists",
			}
		}
		slog.Error("Failed to create field group", "error", err)
		return err
	}

	fg.ID = id
	return nil
}

func (r *fieldGroupRepository) Update(ctx context.Context, userID string, fg *domain.FieldGroup) error {
	programFieldsJSON, err := json.Marshal(fg.ProgramFields)
	if err != nil {
		return fmt.Errorf("marshal program_fields: %w", err)
	}
	logFieldsJSON, err := json.Marshal(fg.LogFields)
	if err != nil {
		return fmt.Errorf("marshal log_fields: %w", err)
	}

	query := `
		UPDATE field_groups
		SET name = $2, description = $3, program_fields = $4, log_fields = $5, updated_at = NOW()
		WHERE id = $1 AND user_id = $6
		RETURNING updated_at
	`

	err = r.pool.QueryRow(ctx, query,
		fg.ID,
		fg.Name,
		fg.Description,
		programFieldsJSON,
		logFieldsJSON,
		userID,
	).Scan(&fg.UpdatedAt)
	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		if strings.Contains(err.Error(), "field_groups_user_id_name_key") || strings.Contains(err.Error(), "duplicate key") {
			return &domain.DomainError{
				Code:    domain.ErrorCodeConflict,
				Message: "field group with this name already exists",
			}
		}
		slog.Error("Failed to update field group", "id", fg.ID, "error", err)
		return err
	}

	return nil
}

func (r *fieldGroupRepository) Delete(ctx context.Context, userID string, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM field_groups WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		slog.Error("Failed to delete field group", "id", id, "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// scanFieldGroup scans a FieldGroup from a pgx.Rows row
func scanFieldGroup(rows pgx.Rows) (domain.FieldGroup, error) {
	var fg domain.FieldGroup
	var programFieldsRaw, logFieldsRaw []byte

	err := rows.Scan(
		&fg.ID,
		&fg.Name,
		&fg.Description,
		&programFieldsRaw,
		&logFieldsRaw,
		&fg.CreatedAt,
		&fg.UpdatedAt,
	)
	if err != nil {
		return fg, err
	}

	if err := json.Unmarshal(programFieldsRaw, &fg.ProgramFields); err != nil {
		return fg, fmt.Errorf("unmarshal program_fields for %s: %w", fg.ID, err)
	}
	if err := json.Unmarshal(logFieldsRaw, &fg.LogFields); err != nil {
		return fg, fmt.Errorf("unmarshal log_fields for %s: %w", fg.ID, err)
	}

	return fg, nil
}

// scanFieldGroupRow scans a FieldGroup from a pgx.Row
func scanFieldGroupRow(row pgx.Row) (domain.FieldGroup, error) {
	var fg domain.FieldGroup
	var programFieldsRaw, logFieldsRaw []byte

	err := row.Scan(
		&fg.ID,
		&fg.Name,
		&fg.Description,
		&programFieldsRaw,
		&logFieldsRaw,
		&fg.CreatedAt,
		&fg.UpdatedAt,
	)
	if err != nil {
		return fg, err
	}

	if len(programFieldsRaw) > 0 {
		if err := json.Unmarshal(programFieldsRaw, &fg.ProgramFields); err != nil {
			return fg, fmt.Errorf("unmarshal program_fields for %s: %w", fg.ID, err)
		}
	}
	if len(logFieldsRaw) > 0 {
		if err := json.Unmarshal(logFieldsRaw, &fg.LogFields); err != nil {
			return fg, fmt.Errorf("unmarshal log_fields for %s: %w", fg.ID, err)
		}
	}

	return fg, nil
}
