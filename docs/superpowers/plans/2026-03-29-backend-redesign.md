# Backend API Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restructure the backend to separate Program (template), Plan (execution queue), and Log (historical record) with ProgramGroup hierarchy, LogSet per-set tracking, and plan_snapshot on Log.

**Architecture:** Schema-first approach — update OpenAPI spec, run code generation, then implement handler/domain/repository/store changes. Each task produces a passing test suite. Database migrations are applied incrementally.

**Tech Stack:** Go, chi router, PostgreSQL, golang-migrate, OpenAPI 3.0.3 + oapi-codegen

**Spec:** `docs/superpowers/specs/2026-03-29-program-plan-log-redesign.md`

---

## File Structure

### New files
- `api/internal/domain/plan.go` — Plan, PlanSession, PlanSessionEntry entities (replace current DateOnly-only file)
- `api/internal/domain/group.go` — ProgramGroup entity
- `api/internal/domain/log_set.go` — LogSet entity
- `api/internal/repository/plan.go` — PlanRepository interface
- `api/internal/store/memory/plan.go` — In-memory PlanRepository
- `api/internal/store/postgres/plan.go` — PostgreSQL PlanRepository
- `api/internal/handler/plan.go` — Plan HTTP handlers
- `api/migrations/000003_add_program_groups.up.sql` / `.down.sql`
- `api/migrations/000004_remove_program_status.up.sql` / `.down.sql`
- `api/migrations/000005_create_plan_tables.up.sql` / `.down.sql`
- `api/migrations/000006_add_log_sets_and_plan_snapshot.up.sql` / `.down.sql`

### Modified files
- `api/openapi/openapi.yaml` — Add ProgramGroup, Plan, LogSet schemas; remove Program status
- `api/internal/domain/program.go` — Remove ProgramStatus, add Groups field
- `api/internal/domain/log.go` — Add PlanSnapshot to Log, add LogSet to LogEntry
- `api/internal/domain/validation.go` — Add validation for ProgramGroup, Plan, PlanSession, LogSet; remove status validation from Program
- `api/internal/repository/program.go` — Remove UpdateStatus, remove status filter from List
- `api/internal/repository/log.go` — Remove LoggedSession methods, add LogSet support
- `api/internal/store/postgres/program.go` — Add group CRUD, remove status operations
- `api/internal/store/memory/program.go` — Add group support, remove status operations
- `api/internal/store/postgres/log.go` — Add log_sets CRUD, plan_snapshot column
- `api/internal/store/memory/log.go` — Add log_sets, plan_snapshot support
- `api/internal/handler/program.go` — Remove status handler, add groups to create/update/get
- `api/internal/handler/log.go` — Add plan_snapshot, log_sets to create/update/get
- `api/cmd/server/main.go` — Add PlanRepository init, Plan routes

### Test files (new or modified)
- `api/internal/domain/group_test.go`
- `api/internal/domain/plan_test.go` (replace current)
- `api/internal/domain/log_set_test.go`
- `api/internal/domain/program_test.go` (update)
- `api/internal/domain/log_test.go` (update)
- `api/internal/store/memory/plan_test.go`
- `api/internal/store/memory/program_test.go` (update)
- `api/internal/store/memory/log_test.go` (update)
- `api/internal/handler/plan_test.go`

---

## Task 1: ProgramGroup Domain Entity & Validation

**Files:**
- Create: `api/internal/domain/group.go`
- Create: `api/internal/domain/group_test.go`
- Modify: `api/internal/domain/validation.go`
- Modify: `api/internal/domain/program.go`

- [ ] **Step 1: Write ProgramGroup entity**

```go
// api/internal/domain/group.go
package domain

import "github.com/google/uuid"

const MaxGroupDepth = 2

// ProgramGroup represents a hierarchical grouping of sessions within a program.
// Supports nesting up to MaxGroupDepth (e.g., Block > Week).
type ProgramGroup struct {
	ID            uuid.UUID  `json:"id"`
	ProgramID     uuid.UUID  `json:"program_id"`
	ParentGroupID *uuid.UUID `json:"parent_group_id,omitempty"`
	Name          string     `json:"name"`
	Order         int        `json:"order"`
	Notes         *string    `json:"notes,omitempty"`
}
```

- [ ] **Step 2: Write failing tests for ProgramGroup validation**

```go
// api/internal/domain/group_test.go
package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidateProgramGroup(t *testing.T) {
	validGroup := func() *ProgramGroup {
		return &ProgramGroup{
			ID:        uuid.New(),
			ProgramID: uuid.New(),
			Name:      "Week 1",
			Order:     0,
		}
	}

	t.Run("valid group", func(t *testing.T) {
		g := validGroup()
		if err := ValidateProgramGroup(g); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("nil group", func(t *testing.T) {
		if err := ValidateProgramGroup(nil); err == nil {
			t.Error("expected error for nil group")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		g := validGroup()
		g.Name = ""
		if err := ValidateProgramGroup(g); err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("name too long", func(t *testing.T) {
		g := validGroup()
		g.Name = string(make([]byte, 201))
		if err := ValidateProgramGroup(g); err == nil {
			t.Error("expected error for long name")
		}
	})

	t.Run("negative order", func(t *testing.T) {
		g := validGroup()
		g.Order = -1
		if err := ValidateProgramGroup(g); err == nil {
			t.Error("expected error for negative order")
		}
	})

	t.Run("notes too long", func(t *testing.T) {
		g := validGroup()
		long := string(make([]byte, 5001))
		g.Notes = &long
		if err := ValidateProgramGroup(g); err == nil {
			t.Error("expected error for long notes")
		}
	})
}

func TestValidateGroupDepth(t *testing.T) {
	programID := uuid.New()
	topGroup := ProgramGroup{ID: uuid.New(), ProgramID: programID, Name: "Block 1", Order: 0}
	nestedGroup := ProgramGroup{ID: uuid.New(), ProgramID: programID, Name: "Week 1", Order: 0, ParentGroupID: &topGroup.ID}

	t.Run("depth 0 is valid", func(t *testing.T) {
		groups := []ProgramGroup{topGroup}
		if err := ValidateGroupDepths(groups); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("depth 1 is valid", func(t *testing.T) {
		groups := []ProgramGroup{topGroup, nestedGroup}
		if err := ValidateGroupDepths(groups); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("depth 2 is rejected", func(t *testing.T) {
		deepGroup := ProgramGroup{ID: uuid.New(), ProgramID: programID, Name: "SubWeek", Order: 0, ParentGroupID: &nestedGroup.ID}
		groups := []ProgramGroup{topGroup, nestedGroup, deepGroup}
		if err := ValidateGroupDepths(groups); err == nil {
			t.Error("expected error for depth >= 2")
		}
	})
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd api && go test ./internal/domain/... -run "TestValidateProgramGroup|TestValidateGroupDepth" -v`
Expected: FAIL — `ValidateProgramGroup` and `ValidateGroupDepths` not defined

- [ ] **Step 4: Implement ValidateProgramGroup and ValidateGroupDepths**

Add to `api/internal/domain/validation.go`:

```go
// ValidateProgramGroup validates a ProgramGroup entity.
func ValidateProgramGroup(g *ProgramGroup) error {
	if g == nil {
		return &ValidationError{Field: "program_group", Message: "program_group cannot be nil"}
	}
	if err := ValidateRequiredString("name", g.Name); err != nil {
		return err
	}
	if err := ValidateStringLength("name", g.Name, 1, 200); err != nil {
		return err
	}
	if g.Order < 0 {
		return &ValidationError{Field: "order", Message: "order must be greater than or equal to 0"}
	}
	if g.Notes != nil {
		if err := ValidateStringLength("notes", *g.Notes, 0, 5000); err != nil {
			return err
		}
	}
	return nil
}

// ValidateGroupDepths checks that no group exceeds MaxGroupDepth.
// depth=0 means top-level (no parent), depth=1 means one parent, depth>=2 is rejected.
func ValidateGroupDepths(groups []ProgramGroup) error {
	if len(groups) == 0 {
		return nil
	}
	parentMap := make(map[uuid.UUID]*uuid.UUID, len(groups))
	for i := range groups {
		parentMap[groups[i].ID] = groups[i].ParentGroupID
	}
	for _, g := range groups {
		depth := 0
		current := g.ParentGroupID
		for current != nil {
			depth++
			if depth >= MaxGroupDepth {
				return &ValidationError{
					Field:   "groups",
					Message: fmt.Sprintf("group '%s' exceeds maximum nesting depth of %d", g.Name, MaxGroupDepth),
				}
			}
			parent, exists := parentMap[*current]
			if !exists {
				break
			}
			current = parent
		}
	}
	return nil
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd api && go test ./internal/domain/... -run "TestValidateProgramGroup|TestValidateGroupDepth" -v`
Expected: PASS

- [ ] **Step 6: Add Groups field to Program entity**

In `api/internal/domain/program.go`, remove `ProgramStatus` type, constants, and `ValidateProgramStatusTransition`. Add `Groups` field to `Program`:

```go
// Program represents a reusable training template with hierarchical session organization.
// No lifecycle status — Programs are persistent templates.
type Program struct {
	ID            uuid.UUID        `json:"id"`
	Name          string           `json:"name"`
	Notes         *string          `json:"notes,omitempty"`
	ProgramFields []FieldDef       `json:"program_fields,omitempty"`
	LogFields     []FieldDef       `json:"log_fields,omitempty"`
	Groups        []ProgramGroup   `json:"groups,omitempty"`
	Sessions      []ProgramSession `json:"sessions,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}
```

Also add `GroupID` to `ProgramSession`:

```go
type ProgramSession struct {
	ID          uuid.UUID             `json:"id"`
	ProgramID   uuid.UUID             `json:"program_id"`
	GroupID     *uuid.UUID            `json:"group_id,omitempty"`
	SessionName string                `json:"session_name"`
	Order       int                   `json:"order"`
	Date        *DateOnly             `json:"date,omitempty"`
	Entries     []ProgramSessionEntry `json:"entries,omitempty"`
}
```

- [ ] **Step 7: Update ValidateProgram to remove status validation, add group validation**

Update `ValidateProgram` in `api/internal/domain/validation.go`:

```go
func ValidateProgram(p *Program) error {
	if p == nil {
		return &ValidationError{Field: "program", Message: "program cannot be nil"}
	}
	if err := ValidateRequiredString("name", p.Name); err != nil {
		return err
	}
	if err := ValidateStringLength("name", p.Name, 1, 200); err != nil {
		return err
	}
	if p.Notes != nil {
		if err := ValidateStringLength("notes", *p.Notes, 0, 5000); err != nil {
			return err
		}
	}

	// Validate groups
	for i := range p.Groups {
		if err := ValidateProgramGroup(&p.Groups[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("groups[%d]", i),
				Message: err.Error(),
			}
		}
	}
	if err := ValidateGroupDepths(p.Groups); err != nil {
		return err
	}

	if len(p.Sessions) == 0 {
		return &ValidationError{Field: "sessions", Message: "program must have at least one session"}
	}
	if len(p.Sessions) > 100 {
		return &ValidationError{Field: "sessions", Message: "program cannot have more than 100 sessions"}
	}

	seenSessions := make(map[string]struct{}, len(p.Sessions))
	for i := range p.Sessions {
		if err := ValidateProgramSession(&p.Sessions[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("sessions[%d]", i),
				Message: err.Error(),
			}
		}
		s := &p.Sessions[i]
		if _, exists := seenSessions[s.SessionName]; exists {
			return &ValidationError{
				Field:   fmt.Sprintf("sessions[%d]", i),
				Message: fmt.Sprintf("duplicate session_name: %s", s.SessionName),
			}
		}
		seenSessions[s.SessionName] = struct{}{}
	}
	return nil
}
```

- [ ] **Step 8: Fix existing program_test.go to remove status-related tests, add group tests**

Update `api/internal/domain/program_test.go` — remove tests for `ValidateProgramStatusTransition` and status field in `ValidateProgram`. Add test for valid program with groups.

- [ ] **Step 9: Run all domain tests**

Run: `cd api && go test ./internal/domain/... -v`
Expected: PASS

- [ ] **Step 10: Commit**

```bash
cd api && git add internal/domain/group.go internal/domain/group_test.go internal/domain/program.go internal/domain/program_test.go internal/domain/validation.go
git commit -m "feat: add ProgramGroup domain entity and remove Program status

Add ProgramGroup with hierarchical nesting (max_depth=2).
Remove ProgramStatus lifecycle from Program entity.
Add group_id to ProgramSession for group membership.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 2: Plan Domain Entity & Validation

**Files:**
- Modify: `api/internal/domain/plan.go` (replace DateOnly-only content)
- Create: `api/internal/domain/plan_test.go`
- Modify: `api/internal/domain/validation.go`

- [ ] **Step 1: Write Plan entities**

Replace `api/internal/domain/plan.go` — keep `DateOnly` type, add Plan entities:

```go
package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DateOnly is a date without time, formatted as "2006-01-02" in JSON.
type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

// PlanSessionEntry represents an exercise prescription within a plan session.
type PlanSessionEntry struct {
	ID           uuid.UUID              `json:"id"`
	SessionID    uuid.UUID              `json:"session_id"`
	Order        int                    `json:"order"`
	ExerciseName string                 `json:"exercise_name"`
	Fields       map[string]interface{} `json:"fields,omitempty"`
	Notes        *string                `json:"notes,omitempty"`
}

// PlanSession represents a single session in the user's execution plan.
type PlanSession struct {
	ID              uuid.UUID          `json:"id"`
	PlanID          uuid.UUID          `json:"plan_id"`
	SessionName     string             `json:"session_name"`
	Order           int                `json:"order"`
	Date            *DateOnly          `json:"date,omitempty"`
	SourceProgramID *uuid.UUID         `json:"source_program_id,omitempty"`
	SourceSessionID *uuid.UUID         `json:"source_session_id,omitempty"`
	Entries         []PlanSessionEntry `json:"entries,omitempty"`
}

// Plan represents the user's working execution queue of upcoming sessions.
// Exactly one Plan per user. Persistent but disposable — no lifecycle constraints.
type Plan struct {
	ID        uuid.UUID     `json:"id"`
	Name      *string       `json:"name,omitempty"`
	Notes     *string       `json:"notes,omitempty"`
	ProgramID *uuid.UUID    `json:"program_id,omitempty"`
	Sessions  []PlanSession `json:"sessions,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
```

- [ ] **Step 2: Write failing tests for Plan validation**

```go
// api/internal/domain/plan_test.go
package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidatePlan(t *testing.T) {
	validPlan := func() *Plan {
		return &Plan{
			ID: uuid.New(),
			Sessions: []PlanSession{
				{
					ID:          uuid.New(),
					PlanID:      uuid.New(),
					SessionName: "Day 1",
					Order:       0,
					Entries: []PlanSessionEntry{
						{
							ID:           uuid.New(),
							SessionID:    uuid.New(),
							Order:        0,
							ExerciseName: "Back Squat",
						},
					},
				},
			},
		}
	}

	t.Run("valid plan", func(t *testing.T) {
		p := validPlan()
		if err := ValidatePlan(p); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("nil plan", func(t *testing.T) {
		if err := ValidatePlan(nil); err == nil {
			t.Error("expected error for nil plan")
		}
	})

	t.Run("empty sessions is valid", func(t *testing.T) {
		p := &Plan{ID: uuid.New()}
		if err := ValidatePlan(p); err != nil {
			t.Errorf("plan with no sessions should be valid, got %v", err)
		}
	})

	t.Run("name too long", func(t *testing.T) {
		p := validPlan()
		long := string(make([]byte, 201))
		p.Name = &long
		if err := ValidatePlan(p); err == nil {
			t.Error("expected error for long name")
		}
	})

	t.Run("session with empty name", func(t *testing.T) {
		p := validPlan()
		p.Sessions[0].SessionName = ""
		if err := ValidatePlan(p); err == nil {
			t.Error("expected error for empty session name")
		}
	})

	t.Run("entry with empty exercise name", func(t *testing.T) {
		p := validPlan()
		p.Sessions[0].Entries[0].ExerciseName = ""
		if err := ValidatePlan(p); err == nil {
			t.Error("expected error for empty exercise name")
		}
	})
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd api && go test ./internal/domain/... -run TestValidatePlan -v`
Expected: FAIL — `ValidatePlan` not defined

- [ ] **Step 4: Implement Plan validation**

Add to `api/internal/domain/validation.go`:

```go
// ValidatePlanSessionEntry validates a PlanSessionEntry entity.
func ValidatePlanSessionEntry(e *PlanSessionEntry) error {
	if e == nil {
		return &ValidationError{Field: "plan_session_entry", Message: "entry cannot be nil"}
	}
	if err := ValidateRequiredString("exercise_name", e.ExerciseName); err != nil {
		return err
	}
	if err := ValidateStringLength("exercise_name", e.ExerciseName, 1, 200); err != nil {
		return err
	}
	if e.Order < 0 {
		return &ValidationError{Field: "order", Message: "order must be greater than or equal to 0"}
	}
	if e.Notes != nil {
		if err := ValidateStringLength("notes", *e.Notes, 0, 2000); err != nil {
			return err
		}
	}
	return nil
}

// ValidatePlanSession validates a PlanSession entity.
func ValidatePlanSession(s *PlanSession) error {
	if s == nil {
		return &ValidationError{Field: "plan_session", Message: "session cannot be nil"}
	}
	if err := ValidateRequiredString("session_name", s.SessionName); err != nil {
		return err
	}
	if err := ValidateStringLength("session_name", s.SessionName, 1, 200); err != nil {
		return err
	}
	if s.Order < 0 {
		return &ValidationError{Field: "order", Message: "order must be greater than or equal to 0"}
	}
	for i := range s.Entries {
		if err := ValidatePlanSessionEntry(&s.Entries[i]); err != nil {
			return &ValidationError{Field: fmt.Sprintf("entries[%d]", i), Message: err.Error()}
		}
	}
	return nil
}

// ValidatePlan validates a Plan entity.
func ValidatePlan(p *Plan) error {
	if p == nil {
		return &ValidationError{Field: "plan", Message: "plan cannot be nil"}
	}
	if p.Name != nil {
		if err := ValidateStringLength("name", *p.Name, 0, 200); err != nil {
			return err
		}
	}
	if p.Notes != nil {
		if err := ValidateStringLength("notes", *p.Notes, 0, 5000); err != nil {
			return err
		}
	}
	if len(p.Sessions) > 200 {
		return &ValidationError{Field: "sessions", Message: "plan cannot have more than 200 sessions"}
	}
	for i := range p.Sessions {
		if err := ValidatePlanSession(&p.Sessions[i]); err != nil {
			return &ValidationError{Field: fmt.Sprintf("sessions[%d]", i), Message: err.Error()}
		}
	}
	return nil
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd api && go test ./internal/domain/... -run TestValidatePlan -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
cd api && git add internal/domain/plan.go internal/domain/plan_test.go internal/domain/validation.go
git commit -m "feat: add Plan domain entity and validation

Add Plan, PlanSession, PlanSessionEntry entities.
Plan is a working execution queue (one per user, no lifecycle).
Validation for name length, session limits, entry fields.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 3: LogSet Domain Entity & Log Changes

**Files:**
- Create: `api/internal/domain/log_set.go`
- Create: `api/internal/domain/log_set_test.go`
- Modify: `api/internal/domain/log.go`
- Modify: `api/internal/domain/validation.go`
- Modify: `api/internal/domain/log_test.go`

- [ ] **Step 1: Write LogSet entity**

```go
// api/internal/domain/log_set.go
package domain

import "github.com/google/uuid"

// LogSet represents a single performed set within a log entry.
type LogSet struct {
	ID        uuid.UUID              `json:"id"`
	EntryID   uuid.UUID              `json:"entry_id"`
	SetNumber int                    `json:"set_number"`
	Fields    map[string]interface{} `json:"fields"`
	VideoURL  *string                `json:"video_url,omitempty"`
	Notes     *string                `json:"notes,omitempty"`
}
```

- [ ] **Step 2: Add PlanSnapshot and Sets to Log/LogEntry**

In `api/internal/domain/log.go`, add `PlanSnapshot` field to `Log` and `Sets` field to `LogEntry`:

```go
type LogEntry struct {
	ID             uuid.UUID              `json:"id"`
	LogID          uuid.UUID              `json:"log_id"`
	Order          int                    `json:"order"`
	ExerciseName   string                 `json:"exercise_name"`
	Fields         map[string]interface{} `json:"fields,omitempty"`
	Notes          *string                `json:"notes,omitempty"`
	VideoObjectKey *string                `json:"video_object_key,omitempty"`
	StartedAt      *time.Time             `json:"started_at,omitempty"`
	FinishedAt     *time.Time             `json:"finished_at,omitempty"`
	Sets           []LogSet               `json:"sets,omitempty"`
}

type Log struct {
	ID           uuid.UUID       `json:"id"`
	ProgramID    *uuid.UUID      `json:"program_id,omitempty"`
	SessionName  *string         `json:"session_name,omitempty"`
	PerformedAt  time.Time       `json:"performed_at"`
	StartedAt    *time.Time      `json:"started_at,omitempty"`
	FinishedAt   *time.Time      `json:"finished_at,omitempty"`
	Notes        *string         `json:"notes,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	PlanSnapshot json.RawMessage `json:"plan_snapshot,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Entries      []LogEntry      `json:"entries"`
}
```

- [ ] **Step 3: Write failing tests for LogSet validation**

```go
// api/internal/domain/log_set_test.go
package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidateLogSet(t *testing.T) {
	validSet := func() *LogSet {
		return &LogSet{
			ID:        uuid.New(),
			EntryID:   uuid.New(),
			SetNumber: 1,
			Fields:    map[string]interface{}{"reps": float64(5), "weight_kg": float64(100)},
		}
	}

	t.Run("valid set", func(t *testing.T) {
		s := validSet()
		if err := ValidateLogSet(s); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("nil set", func(t *testing.T) {
		if err := ValidateLogSet(nil); err == nil {
			t.Error("expected error for nil set")
		}
	})

	t.Run("zero set number", func(t *testing.T) {
		s := validSet()
		s.SetNumber = 0
		if err := ValidateLogSet(s); err == nil {
			t.Error("expected error for zero set_number")
		}
	})

	t.Run("nil fields", func(t *testing.T) {
		s := validSet()
		s.Fields = nil
		if err := ValidateLogSet(s); err == nil {
			t.Error("expected error for nil fields")
		}
	})

	t.Run("notes too long", func(t *testing.T) {
		s := validSet()
		long := string(make([]byte, 2001))
		s.Notes = &long
		if err := ValidateLogSet(s); err == nil {
			t.Error("expected error for long notes")
		}
	})
}
```

- [ ] **Step 4: Run tests to verify they fail**

Run: `cd api && go test ./internal/domain/... -run TestValidateLogSet -v`
Expected: FAIL — `ValidateLogSet` not defined

- [ ] **Step 5: Implement LogSet validation**

Add to `api/internal/domain/validation.go`:

```go
// ValidateLogSet validates a LogSet entity.
func ValidateLogSet(s *LogSet) error {
	if s == nil {
		return &ValidationError{Field: "log_set", Message: "log_set cannot be nil"}
	}
	if s.SetNumber < 1 {
		return &ValidationError{Field: "set_number", Message: "set_number must be >= 1"}
	}
	if s.Fields == nil {
		return &ValidationError{Field: "fields", Message: "fields is required"}
	}
	if s.VideoURL != nil {
		if err := ValidateStringLength("video_url", *s.VideoURL, 1, 2000); err != nil {
			return err
		}
	}
	if s.Notes != nil {
		if err := ValidateStringLength("notes", *s.Notes, 0, 2000); err != nil {
			return err
		}
	}
	return nil
}
```

- [ ] **Step 6: Update ValidateLogEntry to validate nested sets**

In `api/internal/domain/validation.go`, add to the end of `ValidateLogEntry`:

```go
	// Validate nested sets
	for i := range e.Sets {
		if err := ValidateLogSet(&e.Sets[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("sets[%d]", i),
				Message: err.Error(),
			}
		}
	}
```

- [ ] **Step 7: Run all domain tests**

Run: `cd api && go test ./internal/domain/... -v`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
cd api && git add internal/domain/log_set.go internal/domain/log_set_test.go internal/domain/log.go internal/domain/log_test.go internal/domain/validation.go
git commit -m "feat: add LogSet entity and plan_snapshot to Log

Add LogSet for per-set tracking with video and notes.
Add PlanSnapshot field to Log for planned vs actual analysis.
Add Sets field to LogEntry for nested set data.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 4: Database Migrations

**Files:**
- Create: `api/migrations/000003_add_program_groups.up.sql`
- Create: `api/migrations/000003_add_program_groups.down.sql`
- Create: `api/migrations/000004_remove_program_status.up.sql`
- Create: `api/migrations/000004_remove_program_status.down.sql`
- Create: `api/migrations/000005_create_plan_tables.up.sql`
- Create: `api/migrations/000005_create_plan_tables.down.sql`
- Create: `api/migrations/000006_add_log_sets_and_plan_snapshot.up.sql`
- Create: `api/migrations/000006_add_log_sets_and_plan_snapshot.down.sql`

- [ ] **Step 1: Write ProgramGroup migration**

```sql
-- api/migrations/000003_add_program_groups.up.sql
CREATE TABLE IF NOT EXISTS program_groups (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id      UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    parent_group_id UUID REFERENCES program_groups(id) ON DELETE CASCADE,
    name            VARCHAR(200) NOT NULL,
    "order"         INTEGER NOT NULL DEFAULT 0,
    notes           TEXT
);

CREATE INDEX idx_program_groups_program_id ON program_groups (program_id);
CREATE INDEX idx_program_groups_parent ON program_groups (parent_group_id);
CREATE INDEX idx_program_groups_order ON program_groups (program_id, "order");

-- Add group_id to program_sessions
ALTER TABLE program_sessions ADD COLUMN group_id UUID REFERENCES program_groups(id) ON DELETE SET NULL;
CREATE INDEX idx_program_sessions_group_id ON program_sessions (group_id);
```

```sql
-- api/migrations/000003_add_program_groups.down.sql
ALTER TABLE program_sessions DROP COLUMN IF EXISTS group_id;
DROP TABLE IF EXISTS program_groups;
```

- [ ] **Step 2: Write Program status removal migration**

```sql
-- api/migrations/000004_remove_program_status.up.sql
DROP INDEX IF EXISTS idx_programs_status;
ALTER TABLE programs DROP COLUMN IF EXISTS status;
ALTER TABLE programs DROP COLUMN IF EXISTS metadata;
```

```sql
-- api/migrations/000004_remove_program_status.down.sql
ALTER TABLE programs ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'created' CHECK (status IN ('created', 'ongoing', 'completed', 'cancelled'));
ALTER TABLE programs ADD COLUMN metadata JSONB;
CREATE INDEX idx_programs_status ON programs (status);
```

- [ ] **Step 3: Write Plan tables migration**

```sql
-- api/migrations/000005_create_plan_tables.up.sql
CREATE TABLE IF NOT EXISTS plans (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    VARCHAR(200) NOT NULL,
    name       VARCHAR(200),
    notes      TEXT,
    program_id UUID REFERENCES programs(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- One plan per user
CREATE UNIQUE INDEX idx_plans_user_id ON plans (user_id);

CREATE TABLE IF NOT EXISTS plan_sessions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id           UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    session_name      VARCHAR(200) NOT NULL,
    "order"           INTEGER NOT NULL DEFAULT 0,
    date              DATE,
    source_program_id UUID REFERENCES programs(id) ON DELETE SET NULL,
    source_session_id UUID
);

CREATE INDEX idx_plan_sessions_plan_id ON plan_sessions (plan_id);
CREATE INDEX idx_plan_sessions_order ON plan_sessions (plan_id, "order");

CREATE TABLE IF NOT EXISTS plan_session_entries (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id    UUID NOT NULL REFERENCES plan_sessions(id) ON DELETE CASCADE,
    "order"       INTEGER NOT NULL DEFAULT 0,
    exercise_name VARCHAR(200) NOT NULL,
    fields        JSONB,
    notes         TEXT
);

CREATE INDEX idx_plan_session_entries_session_id ON plan_session_entries (session_id);
CREATE INDEX idx_plan_session_entries_order ON plan_session_entries (session_id, "order");
```

```sql
-- api/migrations/000005_create_plan_tables.down.sql
DROP TABLE IF EXISTS plan_session_entries;
DROP TABLE IF EXISTS plan_sessions;
DROP TABLE IF EXISTS plans;
```

- [ ] **Step 4: Write LogSet and plan_snapshot migration**

```sql
-- api/migrations/000006_add_log_sets_and_plan_snapshot.up.sql
-- Add plan_snapshot to logs
ALTER TABLE logs ADD COLUMN plan_snapshot JSONB;

-- Create log_sets table
CREATE TABLE IF NOT EXISTS log_sets (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id   UUID NOT NULL REFERENCES log_entries(id) ON DELETE CASCADE,
    set_number INTEGER NOT NULL CHECK (set_number >= 1),
    fields     JSONB NOT NULL,
    video_url  VARCHAR(2000),
    notes      TEXT
);

CREATE INDEX idx_log_sets_entry_id ON log_sets (entry_id);
CREATE INDEX idx_log_sets_order ON log_sets (entry_id, set_number);

-- Drop the unique constraint on (program_id, session_name) as it's too restrictive
-- for the new Plan-based model where logs are linked via plan_snapshot
ALTER TABLE logs DROP CONSTRAINT IF EXISTS uq_logs_program_id_session_name;
```

```sql
-- api/migrations/000006_add_log_sets_and_plan_snapshot.down.sql
ALTER TABLE logs ADD CONSTRAINT uq_logs_program_id_session_name UNIQUE (program_id, session_name);
DROP TABLE IF EXISTS log_sets;
ALTER TABLE logs DROP COLUMN IF EXISTS plan_snapshot;
```

- [ ] **Step 5: Verify migration files are correctly numbered**

Run: `ls -la api/migrations/`
Expected: 000001 through 000006 (up and down pairs)

- [ ] **Step 6: Commit**

```bash
git add api/migrations/
git commit -m "feat: add database migrations for redesign

Migration 3: program_groups table + group_id on sessions
Migration 4: remove status and metadata from programs
Migration 5: plan, plan_sessions, plan_session_entries tables
Migration 6: log_sets table + plan_snapshot on logs

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 5: Plan Repository Interface & Memory Store

**Files:**
- Create: `api/internal/repository/plan.go`
- Create: `api/internal/store/memory/plan.go`
- Create: `api/internal/store/memory/plan_test.go`

- [ ] **Step 1: Write PlanRepository interface**

```go
// api/internal/repository/plan.go
package repository

import (
	"context"

	"github.com/aoshimash/rx/api/internal/domain"
)

// PlanRepository defines the interface for Plan storage operations.
// Each user has at most one Plan.
type PlanRepository interface {
	// GetByUserID retrieves the user's Plan including all sessions and entries.
	// Returns domain.ErrNotFound if no Plan exists.
	GetByUserID(ctx context.Context, userID string) (*domain.Plan, error)

	// Create stores a new Plan. Returns conflict error if user already has a Plan.
	Create(ctx context.Context, userID string, plan *domain.Plan) error

	// Update replaces the Plan content (name, notes, sessions).
	Update(ctx context.Context, userID string, plan *domain.Plan) error

	// Delete removes the user's Plan and all sessions/entries.
	Delete(ctx context.Context, userID string) error

	// AddSessions appends sessions to the user's Plan.
	AddSessions(ctx context.Context, userID string, sessions []domain.PlanSession) error

	// UpdateSession updates a specific PlanSession.
	UpdateSession(ctx context.Context, userID string, session *domain.PlanSession) error

	// DeleteSession removes a specific PlanSession from the Plan.
	DeleteSession(ctx context.Context, userID string, sessionID string) error
}
```

- [ ] **Step 2: Write failing memory store test**

```go
// api/internal/store/memory/plan_test.go
package memory

import (
	"context"
	"testing"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

func TestPlanStore_CreateAndGet(t *testing.T) {
	store := NewPlanRepository()
	ctx := context.Background()
	userID := "user-1"

	plan := &domain.Plan{
		Sessions: []domain.PlanSession{
			{
				SessionName: "Day 1",
				Order:       0,
				Entries: []domain.PlanSessionEntry{
					{Order: 0, ExerciseName: "Squat"},
				},
			},
		},
	}

	if err := store.Create(ctx, userID, plan); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if plan.ID == uuid.Nil {
		t.Error("expected plan ID to be set")
	}

	got, err := store.GetByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("GetByUserID failed: %v", err)
	}
	if got.ID != plan.ID {
		t.Errorf("expected ID %s, got %s", plan.ID, got.ID)
	}
	if len(got.Sessions) != 1 {
		t.Errorf("expected 1 session, got %d", len(got.Sessions))
	}
	if len(got.Sessions[0].Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(got.Sessions[0].Entries))
	}
}

func TestPlanStore_CreateConflict(t *testing.T) {
	store := NewPlanRepository()
	ctx := context.Background()
	userID := "user-1"

	plan := &domain.Plan{Sessions: []domain.PlanSession{{SessionName: "Day 1", Order: 0}}}
	if err := store.Create(ctx, userID, plan); err != nil {
		t.Fatalf("first Create failed: %v", err)
	}

	plan2 := &domain.Plan{Sessions: []domain.PlanSession{{SessionName: "Day 2", Order: 0}}}
	err := store.Create(ctx, userID, plan2)
	if err == nil {
		t.Error("expected conflict error for second Create")
	}
}

func TestPlanStore_Delete(t *testing.T) {
	store := NewPlanRepository()
	ctx := context.Background()
	userID := "user-1"

	plan := &domain.Plan{Sessions: []domain.PlanSession{{SessionName: "Day 1", Order: 0}}}
	_ = store.Create(ctx, userID, plan)

	if err := store.Delete(ctx, userID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := store.GetByUserID(ctx, userID)
	if err == nil {
		t.Error("expected not found after delete")
	}
}

func TestPlanStore_GetNotFound(t *testing.T) {
	store := NewPlanRepository()
	ctx := context.Background()

	_, err := store.GetByUserID(ctx, "nonexistent")
	if err == nil {
		t.Error("expected not found error")
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd api && go test ./internal/store/memory/... -run TestPlanStore -v`
Expected: FAIL — `NewPlanRepository` not defined

- [ ] **Step 4: Implement in-memory Plan store**

```go
// api/internal/store/memory/plan.go
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

type planStore struct {
	mu    sync.RWMutex
	plans map[string]*domain.Plan // keyed by userID
}

func NewPlanRepository() *planStore {
	return &planStore{plans: make(map[string]*domain.Plan)}
}

func (s *planStore) GetByUserID(_ context.Context, userID string) (*domain.Plan, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plan, ok := s.plans[userID]
	if !ok {
		return nil, &domain.DomainError{Code: domain.ErrCodeNotFound, Message: "plan not found"}
	}
	return copyPlan(plan), nil
}

func (s *planStore) Create(_ context.Context, userID string, plan *domain.Plan) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.plans[userID]; exists {
		return &domain.DomainError{Code: domain.ErrCodeConflict, Message: "user already has a plan"}
	}

	now := time.Now()
	plan.ID = uuid.New()
	plan.CreatedAt = now
	plan.UpdatedAt = now

	for i := range plan.Sessions {
		plan.Sessions[i].ID = uuid.New()
		plan.Sessions[i].PlanID = plan.ID
		for j := range plan.Sessions[i].Entries {
			plan.Sessions[i].Entries[j].ID = uuid.New()
			plan.Sessions[i].Entries[j].SessionID = plan.Sessions[i].ID
		}
	}

	s.plans[userID] = copyPlan(plan)
	return nil
}

func (s *planStore) Update(_ context.Context, userID string, plan *domain.Plan) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.plans[userID]
	if !ok {
		return &domain.DomainError{Code: domain.ErrCodeNotFound, Message: "plan not found"}
	}

	plan.ID = existing.ID
	plan.CreatedAt = existing.CreatedAt
	plan.UpdatedAt = time.Now()

	for i := range plan.Sessions {
		if plan.Sessions[i].ID == uuid.Nil {
			plan.Sessions[i].ID = uuid.New()
		}
		plan.Sessions[i].PlanID = plan.ID
		for j := range plan.Sessions[i].Entries {
			if plan.Sessions[i].Entries[j].ID == uuid.Nil {
				plan.Sessions[i].Entries[j].ID = uuid.New()
			}
			plan.Sessions[i].Entries[j].SessionID = plan.Sessions[i].ID
		}
	}

	s.plans[userID] = copyPlan(plan)
	return nil
}

func (s *planStore) Delete(_ context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.plans[userID]; !ok {
		return &domain.DomainError{Code: domain.ErrCodeNotFound, Message: "plan not found"}
	}
	delete(s.plans, userID)
	return nil
}

func (s *planStore) AddSessions(_ context.Context, userID string, sessions []domain.PlanSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	plan, ok := s.plans[userID]
	if !ok {
		return &domain.DomainError{Code: domain.ErrCodeNotFound, Message: "plan not found"}
	}

	for i := range sessions {
		sessions[i].ID = uuid.New()
		sessions[i].PlanID = plan.ID
		for j := range sessions[i].Entries {
			sessions[i].Entries[j].ID = uuid.New()
			sessions[i].Entries[j].SessionID = sessions[i].ID
		}
	}

	plan.Sessions = append(plan.Sessions, sessions...)
	plan.UpdatedAt = time.Now()
	return nil
}

func (s *planStore) UpdateSession(_ context.Context, userID string, session *domain.PlanSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	plan, ok := s.plans[userID]
	if !ok {
		return &domain.DomainError{Code: domain.ErrCodeNotFound, Message: "plan not found"}
	}

	for i := range plan.Sessions {
		if plan.Sessions[i].ID == session.ID {
			session.PlanID = plan.ID
			for j := range session.Entries {
				if session.Entries[j].ID == uuid.Nil {
					session.Entries[j].ID = uuid.New()
				}
				session.Entries[j].SessionID = session.ID
			}
			plan.Sessions[i] = *session
			plan.UpdatedAt = time.Now()
			return nil
		}
	}
	return &domain.DomainError{Code: domain.ErrCodeNotFound, Message: "session not found in plan"}
}

func (s *planStore) DeleteSession(_ context.Context, userID string, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	plan, ok := s.plans[userID]
	if !ok {
		return &domain.DomainError{Code: domain.ErrCodeNotFound, Message: "plan not found"}
	}

	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return &domain.ValidationError{Field: "session_id", Message: "invalid UUID"}
	}

	for i := range plan.Sessions {
		if plan.Sessions[i].ID == sid {
			plan.Sessions = append(plan.Sessions[:i], plan.Sessions[i+1:]...)
			plan.UpdatedAt = time.Now()
			return nil
		}
	}
	return &domain.DomainError{Code: domain.ErrCodeNotFound, Message: "session not found in plan"}
}

func copyPlan(p *domain.Plan) *domain.Plan {
	cp := *p
	if p.Sessions != nil {
		cp.Sessions = make([]domain.PlanSession, len(p.Sessions))
		for i, s := range p.Sessions {
			cp.Sessions[i] = s
			if s.Entries != nil {
				cp.Sessions[i].Entries = make([]domain.PlanSessionEntry, len(s.Entries))
				copy(cp.Sessions[i].Entries, s.Entries)
			}
		}
	}
	return &cp
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd api && go test ./internal/store/memory/... -run TestPlanStore -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
cd api && git add internal/repository/plan.go internal/store/memory/plan.go internal/store/memory/plan_test.go
git commit -m "feat: add PlanRepository interface and in-memory implementation

PlanRepository supports single-plan-per-user constraint.
In-memory store with full CRUD and session-level operations.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 6: Update Program Repository & Memory Store

**Files:**
- Modify: `api/internal/repository/program.go`
- Modify: `api/internal/store/memory/program.go`
- Modify: `api/internal/store/memory/program_test.go`

- [ ] **Step 1: Update ProgramRepository interface**

Remove `UpdateStatus` and `status` parameter from `List`. Add group support:

```go
type ProgramRepository interface {
	Create(ctx context.Context, program *domain.Program) error
	Update(ctx context.Context, program *domain.Program) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit int, after string) ([]*domain.Program, string, bool, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
}
```

- [ ] **Step 2: Update memory program store**

Remove `UpdateStatus`, remove `status` from `List`, add group handling in `copyProgram`. Store `Groups` on the `Program` struct.

- [ ] **Step 3: Update memory program tests**

Remove status-related tests. Update `List` call sites to remove status parameter.

- [ ] **Step 4: Run memory store tests**

Run: `cd api && go test ./internal/store/memory/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd api && git add internal/repository/program.go internal/store/memory/program.go internal/store/memory/program_test.go
git commit -m "refactor: remove status from ProgramRepository, add group support

Remove UpdateStatus method and status filter from List.
Update memory store to handle Groups on Program entity.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 7: Update Log Repository & Memory Store for LogSet

**Files:**
- Modify: `api/internal/repository/log.go`
- Modify: `api/internal/store/memory/log.go`
- Modify: `api/internal/store/memory/log_test.go`

- [ ] **Step 1: Update LogRepository interface**

Remove `ListLoggedSessionsByProgramID`, `ExistsByProgramIDAndSessionName`, `ExistsByProgramIDAndSessionNameExcluding`. The unique constraint on `(program_id, session_name)` is dropped.

```go
type LogRepository interface {
	Create(ctx context.Context, log *domain.Log) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Log, error)
	Update(ctx context.Context, log *domain.Log) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, programID *uuid.UUID, limit int, after string) ([]*domain.Log, string, bool, error)
	ListByPerformedAtRange(ctx context.Context, programID *uuid.UUID, performedAtFrom, performedAtTo *time.Time, limit int, after string) ([]*domain.Log, string, bool, error)
}
```

- [ ] **Step 2: Update memory log store**

Remove `ListLoggedSessionsByProgramID`, `ExistsByProgramIDAndSessionName`, `ExistsByProgramIDAndSessionNameExcluding`. Add `PlanSnapshot` and `Sets` handling in `copyLog`/`copyLogEntry`. Store `LogSet` entries within `LogEntry`.

- [ ] **Step 3: Update memory log tests**

Remove tests for removed methods. Add test for creating/reading a Log with `PlanSnapshot` and `LogSet` entries.

- [ ] **Step 4: Run memory store tests**

Run: `cd api && go test ./internal/store/memory/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd api && git add internal/repository/log.go internal/store/memory/log.go internal/store/memory/log_test.go
git commit -m "refactor: simplify LogRepository, add LogSet and plan_snapshot support

Remove logged-sessions and unique session name methods.
Add PlanSnapshot (JSON) and LogSet (per-set data) to Log/LogEntry.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 8: Update OpenAPI Spec

**Files:**
- Modify: `api/openapi/openapi.yaml`

- [ ] **Step 1: Update Program schemas**

Remove `status` from `Program`, `ProgramCreate`, and `ProgramUpdate` schemas. Remove `metadata` from all Program schemas. Add `ProgramGroup` schema and `groups` field to `Program`. Add `group_id` to `ProgramSession` and `ProgramSessionCreate`.

- [ ] **Step 2: Add Plan schemas and endpoints**

Add `Plan`, `PlanCreate`, `PlanUpdate`, `PlanSession`, `PlanSessionCreate`, `PlanSessionEntry`, `PlanSessionEntryCreate` schemas. Add endpoints: `GET /plan`, `POST /plan`, `PUT /plan`, `DELETE /plan`, `POST /plan/sessions`, `PUT /plan/sessions/{id}`, `DELETE /plan/sessions/{id}`, `POST /plan/expand-program/{program_id}`.

- [ ] **Step 3: Update Log schemas**

Add `plan_snapshot` to `Log` and `LogCreate`/`LogUpdate`. Add `LogSet` and `LogSetCreate` schemas. Add `sets` to `LogEntry` and `LogEntryCreate`.

- [ ] **Step 4: Remove deprecated endpoints**

Remove `PATCH /programs/{id}/status` and `GET /programs/{id}/logged-sessions`. Remove `status` query parameter from `GET /programs`.

- [ ] **Step 5: Run code generation**

Run: `cd api && task generate`
Expected: Generated code in `pkg/openapi/server.gen.go` updated successfully

- [ ] **Step 6: Commit**

```bash
cd api && git add openapi/openapi.yaml pkg/
git commit -m "feat: update OpenAPI spec for redesign

Add ProgramGroup, Plan, LogSet schemas and endpoints.
Remove Program status lifecycle and metadata.
Remove /programs/{id}/status and /programs/{id}/logged-sessions.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 9: Update Program Handler

**Files:**
- Modify: `api/internal/handler/program.go`
- Modify: `api/internal/handler/program_test.go`

- [ ] **Step 1: Remove UpdateProgramStatus and ListLoggedSessions handlers**

Delete `UpdateProgramStatus` and `ListLoggedSessions` methods from `ProgramHandler`.

- [ ] **Step 2: Update CreateProgram handler**

Remove `status` and `metadata` from request parsing. Add `groups` to request/response. Set no default status.

- [ ] **Step 3: Update UpdateProgram handler**

Remove `status` and `metadata`. Add `groups`.

- [ ] **Step 4: Update GetProgram and ListPrograms handlers**

Remove `status` from response. Remove `status` filter from ListPrograms. Add `groups` to response.

- [ ] **Step 5: Update ProgramHandler constructor**

Remove `logRepo` dependency if only used for `ListLoggedSessions`.

- [ ] **Step 6: Run handler tests**

Run: `cd api && go test ./internal/handler/... -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
cd api && git add internal/handler/program.go internal/handler/program_test.go
git commit -m "refactor: update Program handler for redesign

Remove status and metadata handling.
Remove UpdateProgramStatus and ListLoggedSessions handlers.
Add ProgramGroup support to create/update/get.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 10: Plan Handler

**Files:**
- Create: `api/internal/handler/plan.go`
- Create: `api/internal/handler/plan_test.go`

- [ ] **Step 1: Write Plan handler struct and constructor**

```go
// api/internal/handler/plan.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

type PlanHandler struct {
	planRepo    repository.PlanRepository
	programRepo repository.ProgramRepository
}

func NewPlanHandler(planRepo repository.PlanRepository, programRepo repository.ProgramRepository) *PlanHandler {
	return &PlanHandler{planRepo: planRepo, programRepo: programRepo}
}
```

- [ ] **Step 2: Write GetPlan handler**

```go
func (h *PlanHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	plan, err := h.planRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		handleDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, plan)
}
```

- [ ] **Step 3: Write CreatePlan handler**

```go
func (h *PlanHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req struct {
		Name      *string              `json:"name"`
		Notes     *string              `json:"notes"`
		ProgramID *uuid.UUID           `json:"program_id"`
		Sessions  []planSessionRequest `json:"sessions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "invalid request body")
		return
	}

	plan := &domain.Plan{
		Name:      req.Name,
		Notes:     req.Notes,
		ProgramID: req.ProgramID,
	}
	for _, s := range req.Sessions {
		plan.Sessions = append(plan.Sessions, s.toDomain())
	}

	if err := domain.ValidatePlan(plan); err != nil {
		handleValidationError(w, err)
		return
	}

	if err := h.planRepo.Create(r.Context(), userID, plan); err != nil {
		handleDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, plan)
}
```

- [ ] **Step 4: Write UpdatePlan, DeletePlan handlers**

Implement `UpdatePlan` (PUT /plan) and `DeletePlan` (DELETE /plan) following the same pattern.

- [ ] **Step 5: Write ExpandProgram handler**

```go
func (h *PlanHandler) ExpandProgram(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	programID, err := parseUUIDParam(r, "program_id")
	if err != nil {
		middleware.WriteValidationError(w, "invalid program_id")
		return
	}

	program, err := h.programRepo.GetByID(r.Context(), programID)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	// Convert program sessions to plan sessions
	var sessions []domain.PlanSession
	for _, ps := range program.Sessions {
		planSession := domain.PlanSession{
			SessionName:     ps.SessionName,
			Order:           ps.Order,
			Date:            ps.Date,
			SourceProgramID: &program.ID,
			SourceSessionID: &ps.ID,
		}
		for _, e := range ps.Entries {
			planSession.Entries = append(planSession.Entries, domain.PlanSessionEntry{
				Order:        e.Order,
				ExerciseName: e.ExerciseName,
				Fields:       e.Fields,
				Notes:        e.Notes,
			})
		}
		sessions = append(sessions, planSession)
	}

	if err := h.planRepo.AddSessions(r.Context(), userID, sessions); err != nil {
		handleDomainError(w, err)
		return
	}

	// Return updated plan
	plan, err := h.planRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		handleDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, plan)
}
```

- [ ] **Step 6: Write plan session CRUD handlers**

Implement `AddPlanSessions` (POST /plan/sessions), `UpdatePlanSession` (PUT /plan/sessions/{id}), `DeletePlanSession` (DELETE /plan/sessions/{id}).

- [ ] **Step 7: Add planSessionRequest helper types**

```go
type planSessionEntryRequest struct {
	ExerciseName string                 `json:"exercise_name"`
	Order        int                    `json:"order"`
	Fields       map[string]interface{} `json:"fields,omitempty"`
	Notes        *string                `json:"notes,omitempty"`
}

type planSessionRequest struct {
	SessionName     string                    `json:"session_name"`
	Order           int                       `json:"order"`
	Date            *string                   `json:"date,omitempty"`
	SourceProgramID *uuid.UUID                `json:"source_program_id,omitempty"`
	SourceSessionID *uuid.UUID                `json:"source_session_id,omitempty"`
	Entries         []planSessionEntryRequest `json:"entries,omitempty"`
}

func (r planSessionRequest) toDomain() domain.PlanSession {
	s := domain.PlanSession{
		SessionName:     r.SessionName,
		Order:           r.Order,
		SourceProgramID: r.SourceProgramID,
		SourceSessionID: r.SourceSessionID,
	}
	if r.Date != nil {
		d := domain.DateOnly{}
		_ = d.UnmarshalJSON([]byte(`"` + *r.Date + `"`))
		s.Date = &d
	}
	for _, e := range r.Entries {
		s.Entries = append(s.Entries, domain.PlanSessionEntry{
			Order:        e.Order,
			ExerciseName: e.ExerciseName,
			Fields:       e.Fields,
			Notes:        e.Notes,
		})
	}
	return s
}
```

- [ ] **Step 8: Write basic handler tests**

Test `GetPlan` (404 when no plan), `CreatePlan` (success + 409 conflict), `DeletePlan`.

- [ ] **Step 9: Run handler tests**

Run: `cd api && go test ./internal/handler/... -v`
Expected: PASS

- [ ] **Step 10: Commit**

```bash
cd api && git add internal/handler/plan.go internal/handler/plan_test.go
git commit -m "feat: add Plan HTTP handlers

GET/POST/PUT/DELETE /plan for plan CRUD.
POST /plan/sessions, PUT/DELETE /plan/sessions/{id} for session ops.
POST /plan/expand-program/{program_id} to expand program into plan.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 11: Update Log Handler for LogSet & PlanSnapshot

**Files:**
- Modify: `api/internal/handler/log.go`

- [ ] **Step 1: Add LogSet and plan_snapshot to CreateLog request parsing**

Update the request struct and parsing in `CreateLog` to accept `plan_snapshot` (JSON) and `entries[].sets` (array of LogSet create requests).

- [ ] **Step 2: Add LogSet to Log response serialization**

Ensure `GetLog` and `ListLogs` include `sets` in each entry and `plan_snapshot` on the log.

- [ ] **Step 3: Update UpdateLog similarly**

- [ ] **Step 4: Remove ListLoggedSessions handler**

Delete `ListLoggedSessions` method.

- [ ] **Step 5: Run handler tests**

Run: `cd api && go test ./internal/handler/... -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
cd api && git add internal/handler/log.go
git commit -m "feat: add LogSet and plan_snapshot to Log handler

Support per-set data (fields, video_url, notes) in log entries.
Add plan_snapshot field for planned vs actual tracking.
Remove ListLoggedSessions handler.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 12: Update Router & Main

**Files:**
- Modify: `api/cmd/server/main.go`

- [ ] **Step 1: Add PlanRepository initialization**

In `main()`, add `planRepo` initialization for both memory and postgres backends:

```go
var planRepo repository.PlanRepository
// ... inside if/else for storage type:
// memory:
planRepo = memory.NewPlanRepository()
// postgres:
planRepo = postgresstore.NewPlanRepository(db.Pool())
```

- [ ] **Step 2: Initialize PlanHandler**

```go
planHandler := handler.NewPlanHandler(planRepo, programRepo)
```

- [ ] **Step 3: Update router — remove old routes, add Plan routes**

```go
// Remove these:
// r.Patch("/programs/{id}/status", programHandler.UpdateProgramStatus)
// r.Get("/programs/{id}/logged-sessions", programHandler.ListLoggedSessions)

// Add Plan routes:
r.Get("/plan", planHandler.GetPlan)
r.Post("/plan", planHandler.CreatePlan)
r.Put("/plan", planHandler.UpdatePlan)
r.Delete("/plan", planHandler.DeletePlan)
r.Post("/plan/sessions", planHandler.AddPlanSessions)
r.Put("/plan/sessions/{id}", planHandler.UpdatePlanSession)
r.Delete("/plan/sessions/{id}", planHandler.DeletePlanSession)
r.Post("/plan/expand-program/{program_id}", planHandler.ExpandProgram)
```

- [ ] **Step 4: Update ProgramHandler constructor**

If `logRepo` dependency was removed from `ProgramHandler`, update the constructor call.

- [ ] **Step 5: Verify compilation**

Run: `cd api && go build ./...`
Expected: Build succeeds

- [ ] **Step 6: Run all tests**

Run: `cd api && task check`
Expected: generate + format + lint + test all pass

- [ ] **Step 7: Commit**

```bash
cd api && git add cmd/server/main.go
git commit -m "feat: wire Plan routes and update router

Add Plan CRUD and session management routes.
Remove /programs/{id}/status and /programs/{id}/logged-sessions routes.
Initialize PlanRepository in both memory and postgres backends.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 13: PostgreSQL Store Implementations

**Files:**
- Modify: `api/internal/store/postgres/program.go`
- Create: `api/internal/store/postgres/plan.go`
- Modify: `api/internal/store/postgres/log.go`

- [ ] **Step 1: Update postgres program store**

Remove `UpdateStatus` method. Remove `status` parameter from `List`. Add `program_groups` CRUD within `Create`/`Update`/`GetByID`. Load groups in `GetByID` and join with sessions via `group_id`.

- [ ] **Step 2: Write postgres Plan store**

Implement `PlanRepository` interface with SQL queries for `plans`, `plan_sessions`, `plan_session_entries` tables. Use transactions for `Create`/`Update`. Enforce one-plan-per-user via the unique index on `user_id`.

- [ ] **Step 3: Update postgres log store**

Add `plan_snapshot` column handling in `Create`/`Update`/`GetByID`. Add `log_sets` CRUD within `Create`/`Update`/`GetByID` (insert/select/delete alongside log_entries). Remove `ListLoggedSessionsByProgramID`, `ExistsByProgramIDAndSessionName`, `ExistsByProgramIDAndSessionNameExcluding`.

- [ ] **Step 4: Run integration tests (if postgres is available)**

Run: `cd api && go test -tags=integration ./internal/store/postgres/... -v`
Expected: PASS (or skip if no postgres)

- [ ] **Step 5: Commit**

```bash
cd api && git add internal/store/postgres/
git commit -m "feat: update PostgreSQL stores for redesign

Update program store: groups support, remove status.
Add plan store: full CRUD with one-per-user constraint.
Update log store: plan_snapshot and log_sets support.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 14: Final Integration & Cleanup

**Files:**
- Modify: `api/internal/seed/seed.go`
- Modify: various test files

- [ ] **Step 1: Update seed data**

Update `api/internal/seed/seed.go` to create Programs without status, with groups. Create a sample Plan for dev user.

- [ ] **Step 2: Run full check**

Run: `cd api && task check`
Expected: generate + format + lint + test all pass

- [ ] **Step 3: Verify API starts and responds**

Run: `cd api && go build ./cmd/server && echo "Build OK"`
Expected: Build OK

- [ ] **Step 4: Commit**

```bash
cd api && git add .
git commit -m "chore: update seed data and final cleanup for redesign

Update seed to create Programs with groups and no status.
Add sample Plan for development testing.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Follow-up: Frontend Redesign (separate plan)

After this backend plan is complete, a separate frontend plan will cover:
1. Update TypeScript types (`web/types/api.ts`)
2. Add Plan API client and hooks
3. Replace Dashboard with Plan as top page
4. Update Program pages (remove status UI, add group editor)
5. Update Log pages (add LogSet UI, plan_snapshot display)
6. Update Sidebar navigation
