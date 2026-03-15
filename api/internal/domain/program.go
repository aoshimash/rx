package domain

import (
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
)

// ProgramEntry represents a single exercise prescription in a reusable training template.
// It uses RPE and percent_1rm for intensity instead of absolute weights.
type ProgramEntry struct {
	ID           uuid.UUID       `json:"id"`
	ProgramID    uuid.UUID       `json:"program_id"`
	Order        int             `json:"order"`
	ExerciseName string          `json:"exercise_name"`
	Sets         *int            `json:"sets,omitempty"`
	Reps         *int            `json:"reps,omitempty"`
	RPE          *int            `json:"rpe,omitempty"`
	Percent1RM   *float64        `json:"percent_1rm,omitempty"`
	Notes        *string         `json:"notes,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}

// Program represents a reusable, RPE-based training template.
// It contains no dates and no absolute weights.
type Program struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Notes       *string         `json:"notes,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Entries     []ProgramEntry  `json:"entries,omitempty"`
}

// ConvertProgramToPlanInput holds the user-supplied inputs needed to convert a Program into a Plan.
type ConvertProgramToPlanInput struct {
	// Name for the created Plan. If empty, the Program name is used.
	Name string `json:"name,omitempty"`
	// TargetWeights maps exercise_name to target weight in kg (e.g. 1RM or working weight).
	// For entries with percent_1rm: load_kg = percent_1rm * target_weight, rounded to increment.
	// For entries without percent_1rm: load_kg = target_weight (direct copy).
	TargetWeights map[string]float64 `json:"target_weights"`
	// LoadIncrements maps exercise_name to weight increment in kg for rounding.
	// If not specified for an exercise, weights are rounded to 0.1 kg precision.
	LoadIncrements map[string]float64 `json:"load_increments,omitempty"`
}

// RoundToIncrement rounds a weight to the nearest multiple of the given increment.
// If increment is <= 0, falls back to 0.1 kg precision (RoundLoad).
func RoundToIncrement(weight float64, increment float64) float64 {
	if increment <= 0 {
		return RoundLoad(weight)
	}
	return math.Round(weight/increment) * increment
}

// ConvertProgramToPlan converts a Program into a Plan using the provided inputs.
func ConvertProgramToPlan(program *Program, input *ConvertProgramToPlanInput) *Plan {
	planName := program.Name
	if input.Name != "" {
		planName = input.Name
	}

	programID := program.ID
	plan := &Plan{
		ProgramID:   &programID,
		Name:        planName,
		Description: program.Description,
		Notes:       program.Notes,
		Entries:     make([]PlanEntry, len(program.Entries)),
	}

	// Copy metadata if present
	if program.Metadata != nil {
		plan.Metadata = make(json.RawMessage, len(program.Metadata))
		copy(plan.Metadata, program.Metadata)
	}

	for i, pe := range program.Entries {
		entry := PlanEntry{
			Order:        pe.Order,
			ExerciseName: pe.ExerciseName,
			Sets:         pe.Sets,
			Reps:         pe.Reps,
			RPE:          pe.RPE,
			Notes:        pe.Notes,
		}

		// Copy metadata if present
		if pe.Metadata != nil {
			entry.Metadata = make(json.RawMessage, len(pe.Metadata))
			copy(entry.Metadata, pe.Metadata)
		}

		// Calculate load_kg based on target weights
		if targetWeight, ok := input.TargetWeights[pe.ExerciseName]; ok {
			var loadKg float64
			if pe.Percent1RM != nil {
				// RPE/percent_1rm-based: calculate from percentage
				loadKg = *pe.Percent1RM * targetWeight
			} else {
				// Direct weight: copy as-is
				loadKg = targetWeight
			}

			// Round to increment
			increment := input.LoadIncrements[pe.ExerciseName]
			loadKg = RoundToIncrement(loadKg, increment)
			entry.LoadKg = &loadKg
		}

		plan.Entries[i] = entry
	}

	return plan
}
