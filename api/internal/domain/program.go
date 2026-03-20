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

// ConvertProgramToPlansInput holds the user-supplied inputs needed to convert a Program into Plans.
type ConvertProgramToPlansInput struct {
	// Name base name for created Plans. If empty, the Program name is used.
	// For multi-session programs, session names are appended: "{Name} - {SessionName}".
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

// getSessionName extracts the "session" field from a ProgramEntry's metadata.
// Returns empty string if metadata is nil or has no "session" key.
func getSessionName(metadata json.RawMessage) string {
	if metadata == nil {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal(metadata, &m); err != nil {
		return ""
	}
	if s, ok := m["session"].(string); ok {
		return s
	}
	return ""
}

// sessionGroup holds ProgramEntries belonging to the same session.
type sessionGroup struct {
	name    string
	entries []ProgramEntry
}

// groupBySession groups ProgramEntries by their metadata.session value.
// Entries are grouped in the order of first appearance.
func groupBySession(entries []ProgramEntry) []sessionGroup {
	var order []string
	groups := make(map[string][]ProgramEntry)

	for _, e := range entries {
		name := getSessionName(e.Metadata)
		if _, exists := groups[name]; !exists {
			order = append(order, name)
		}
		groups[name] = append(groups[name], e)
	}

	result := make([]sessionGroup, len(order))
	for i, name := range order {
		result[i] = sessionGroup{name: name, entries: groups[name]}
	}
	return result
}

// buildConversionMetadata creates the metadata.conversion JSON storing conversion parameters.
func buildConversionMetadata(input *ConvertProgramToPlansInput) json.RawMessage {
	conversion := map[string]interface{}{
		"target_weights": input.TargetWeights,
	}
	if len(input.LoadIncrements) > 0 {
		conversion["load_increments"] = input.LoadIncrements
	}
	meta := map[string]interface{}{
		"conversion": conversion,
	}
	data, _ := json.Marshal(meta)
	return data
}

// convertEntry converts a ProgramEntry into a PlanEntry, calculating load_kg.
func convertEntry(pe ProgramEntry, input *ConvertProgramToPlansInput) PlanEntry {
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
			loadKg = *pe.Percent1RM * targetWeight
		} else {
			loadKg = targetWeight
		}
		increment := input.LoadIncrements[pe.ExerciseName]
		loadKg = RoundToIncrement(loadKg, increment)
		entry.LoadKg = &loadKg
	}

	return entry
}

// ConvertProgramToPlans converts a Program into one Plan per session.
// Entries are grouped by metadata.session. If no session metadata exists,
// all entries become a single Plan.
func ConvertProgramToPlans(program *Program, input *ConvertProgramToPlansInput) []*Plan {
	baseName := program.Name
	if input.Name != "" {
		baseName = input.Name
	}

	programID := program.ID
	groups := groupBySession(program.Entries)
	multiSession := len(groups) > 1

	conversionMeta := buildConversionMetadata(input)

	plans := make([]*Plan, 0, len(groups))
	for _, g := range groups {
		planName := baseName
		var sessionName *string
		if g.name != "" {
			sessionName = &g.name
			if multiSession {
				planName = baseName + " - " + g.name
			}
		}

		plan := &Plan{
			ProgramID:   &programID,
			Name:        planName,
			SessionName: sessionName,
			Description: program.Description,
			Notes:       program.Notes,
			Metadata:    conversionMeta,
			Entries:     make([]PlanEntry, len(g.entries)),
		}

		for i, pe := range g.entries {
			plan.Entries[i] = convertEntry(pe, input)
		}

		plans = append(plans, plan)
	}

	return plans
}
