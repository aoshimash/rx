package domain

import (
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
)

// ProgramTemplateEntry represents a single exercise prescription in a reusable training template.
// It uses RPE and percent_1rm for intensity instead of absolute weights.
type ProgramTemplateEntry struct {
	ID                uuid.UUID       `json:"id"`
	ProgramTemplateID uuid.UUID       `json:"program_template_id"`
	Order             int             `json:"order"`
	ExerciseName      string          `json:"exercise_name"`
	Sets              *int            `json:"sets,omitempty"`
	Reps              *int            `json:"reps,omitempty"`
	RPE               *int            `json:"rpe,omitempty"`
	Percent1RM        *float64        `json:"percent_1rm,omitempty"`
	Notes             *string         `json:"notes,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
}

// ProgramTemplate represents a reusable, RPE-based training template.
// It contains no dates and no absolute weights.
// ProgramTemplates are immutable after creation; use Archive/Unarchive to hide unused ones.
type ProgramTemplate struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	Notes       *string                `json:"notes,omitempty"`
	Metadata    json.RawMessage        `json:"metadata,omitempty"`
	Weeks       *string                `json:"weeks,omitempty"`
	DaysPerWeek *string                `json:"days_per_week,omitempty"`
	CreatedBy   *string                `json:"created_by,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	ArchivedAt  *time.Time             `json:"archived_at,omitempty"`
	Entries     []ProgramTemplateEntry `json:"entries,omitempty"`
}

// GenerateProgramInput holds the user-supplied inputs needed to generate a Program from a ProgramTemplate.
type GenerateProgramInput struct {
	// Name is the base name for the generated Program. If empty, the ProgramTemplate name is used.
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

// getSessionName extracts the session name from entry metadata ("session" key).
// Returns empty string if not present.
func getSessionName(e ProgramTemplateEntry) string {
	if e.Metadata == nil {
		return ""
	}
	var meta map[string]interface{}
	if err := json.Unmarshal(e.Metadata, &meta); err != nil {
		return ""
	}
	if s, ok := meta["session"].(string); ok {
		return s
	}
	return ""
}

type templateSessionGroup struct {
	sessionName string
	entries     []ProgramTemplateEntry
}

// groupTemplateBySession groups ProgramTemplateEntries by metadata.session.
// Entries are grouped in the order of first appearance.
func groupTemplateBySession(entries []ProgramTemplateEntry) []templateSessionGroup {
	var order []string
	groups := make(map[string][]ProgramTemplateEntry)

	for _, e := range entries {
		key := getSessionName(e)
		if _, exists := groups[key]; !exists {
			order = append(order, key)
		}
		groups[key] = append(groups[key], e)
	}

	result := make([]templateSessionGroup, len(order))
	for i, key := range order {
		result[i] = templateSessionGroup{
			sessionName: key,
			entries:     groups[key],
		}
	}
	return result
}

// buildGenerationMetadata creates the metadata JSON storing generation parameters.
func buildGenerationMetadata(input *GenerateProgramInput) json.RawMessage {
	generation := map[string]interface{}{
		"target_weights": input.TargetWeights,
	}
	if len(input.LoadIncrements) > 0 {
		generation["load_increments"] = input.LoadIncrements
	}
	meta := map[string]interface{}{
		"generation": generation,
	}
	data, _ := json.Marshal(meta)
	return data
}

// convertTemplateEntry converts a ProgramTemplateEntry into a ProgramSessionEntry, calculating load_kg.
func convertTemplateEntry(te ProgramTemplateEntry, input *GenerateProgramInput) ProgramSessionEntry {
	entry := ProgramSessionEntry{
		Order:        te.Order,
		ExerciseName: te.ExerciseName,
		Sets:         te.Sets,
		Reps:         te.Reps,
		RPE:          te.RPE,
		Notes:        te.Notes,
	}

	// Copy metadata if present
	if te.Metadata != nil {
		entry.Metadata = make(json.RawMessage, len(te.Metadata))
		copy(entry.Metadata, te.Metadata)
	}

	// Calculate load_kg based on target weights
	if targetWeight, ok := input.TargetWeights[te.ExerciseName]; ok {
		var loadKg float64
		if te.Percent1RM != nil {
			loadKg = *te.Percent1RM * targetWeight
		} else {
			loadKg = targetWeight
		}
		increment := input.LoadIncrements[te.ExerciseName]
		loadKg = RoundToIncrement(loadKg, increment)
		entry.LoadKg = &loadKg
	}

	return entry
}

// GenerateProgramFromTemplate converts a ProgramTemplate into a Program with embedded sessions.
// Entries are grouped by metadata.session — one ProgramSession per session group.
func GenerateProgramFromTemplate(template *ProgramTemplate, input *GenerateProgramInput) *Program {
	name := template.Name
	if input.Name != "" {
		name = input.Name
	}

	groups := groupTemplateBySession(template.Entries)
	generationMeta := buildGenerationMetadata(input)

	sessions := make([]ProgramSession, 0, len(groups))
	for i, g := range groups {
		sessionName := g.sessionName
		if sessionName == "" {
			sessionName = name
		}

		session := ProgramSession{
			SessionName: sessionName,
			Order:       i,
			Entries:     make([]ProgramSessionEntry, len(g.entries)),
		}

		for j, te := range g.entries {
			session.Entries[j] = convertTemplateEntry(te, input)
		}

		sessions = append(sessions, session)
	}

	templateID := template.ID
	return &Program{
		ProgramTemplateID: &templateID,
		Name:              name,
		Status:            ProgramStatusActive,
		Notes:             template.Notes,
		Metadata:          generationMeta,
		Sessions:          sessions,
	}
}
