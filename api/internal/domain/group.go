package domain

import "github.com/google/uuid"

// MaxGroupDepth is the maximum allowed nesting depth for program groups.
// Depth 0 = top-level group, depth 1 = one level of nesting (e.g., Block > Week).
const MaxGroupDepth = 2

// ProgramGroup represents a hierarchical grouping of sessions within a program.
type ProgramGroup struct {
	ID            uuid.UUID  `json:"id"`
	ProgramID     uuid.UUID  `json:"program_id"`
	ParentGroupID *uuid.UUID `json:"parent_group_id,omitempty"`
	Name          string     `json:"name"`
	Order         int        `json:"order"`
	Notes         *string    `json:"notes,omitempty"`
}
