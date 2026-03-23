# ProgramTemplate Editing Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `POST /program-templates/{id}/edit` endpoint that performs in-place update when no Programs reference the template, or creates a new versioned template + archives the old one when Programs exist.

**Architecture:** Schema-first (OpenAPI → Go code gen), then domain model, repository (memory + postgres), handler, route registration. Frontend changes are out of scope for this plan — API only.

**Tech Stack:** Go, PostgreSQL, chi router, pgx v5, OpenAPI 3.0, testify

**Spec:** `docs/superpowers/specs/2026-03-23-template-editing-design.md`

---

### Task 1: Database Migration — Add `source_template_id` Column

**Files:**
- Create: `api/migrations/000014_add_source_template_id.up.sql`
- Create: `api/migrations/000014_add_source_template_id.down.sql`

- [ ] **Step 1: Create up migration**

```sql
-- api/migrations/000014_add_source_template_id.up.sql
ALTER TABLE program_templates
  ADD COLUMN source_template_id UUID REFERENCES program_templates(id) ON DELETE SET NULL;
```

- [ ] **Step 2: Create down migration**

```sql
-- api/migrations/000014_add_source_template_id.down.sql
ALTER TABLE program_templates
  DROP COLUMN source_template_id;
```

- [ ] **Step 3: Commit**

```bash
git add api/migrations/000014_add_source_template_id.up.sql api/migrations/000014_add_source_template_id.down.sql
git commit -m "feat: add source_template_id column to program_templates"
```

---

### Task 2: Domain Model — Add `SourceTemplateID` Field

**Files:**
- Modify: `api/internal/domain/program_template.go:29-42` (ProgramTemplate struct)

- [ ] **Step 1: Add `SourceTemplateID` field to `ProgramTemplate` struct**

Add the field after `ArchivedAt`:

```go
SourceTemplateID *uuid.UUID `json:"source_template_id,omitempty"`
```

- [ ] **Step 2: Update the struct comment**

Change the comment on line 28 from:
```go
// ProgramTemplates are immutable after creation; use Archive/Unarchive to hide unused ones.
```
to:
```go
// ProgramTemplates can be edited via the Edit endpoint, which either updates in-place
// (no linked Programs) or creates a new version with source_template_id pointing to the original.
```

- [ ] **Step 3: Verify build**

Run: `cd api && go build ./...`
Expected: compiles without errors

- [ ] **Step 4: Commit**

```bash
git add api/internal/domain/program_template.go
git commit -m "feat: add SourceTemplateID field to ProgramTemplate domain model"
```

---

### Task 3: Repository Interface — Add `Update` Method

**Files:**
- Modify: `api/internal/repository/program_template.go`

- [ ] **Step 1: Add `Update` method to the interface**

Add after the `Create` method (line 14):

```go
// Update replaces the content of an existing ProgramTemplate (name, description, notes, metadata, weeks, days_per_week, entries).
// The template ID, created_at, created_by, and archived_at are preserved. updated_at is set to NOW().
// Entries are replaced entirely (old entries deleted, new entries inserted).
Update(ctx context.Context, tmpl *domain.ProgramTemplate) error
```

- [ ] **Step 2: Update the interface comment**

Change line 12 from:
```go
// ProgramTemplates are immutable after creation; use Archive/Unarchive instead of Update.
```
to:
```go
// ProgramTemplateRepository defines the interface for ProgramTemplate storage operations.
```

- [ ] **Step 3: Verify build fails (memory + postgres don't implement `Update` yet)**

Run: `cd api && go build ./...`
Expected: compile errors in `memory` and `postgres` packages — "does not implement ProgramTemplateRepository (missing method Update)"

- [ ] **Step 4: Commit**

```bash
git add api/internal/repository/program_template.go
git commit -m "feat: add Update method to ProgramTemplateRepository interface"
```

---

### Task 4: Memory Store — Implement `Update`

**Files:**
- Modify: `api/internal/store/memory/program_template.go`
- Create: `api/internal/store/memory/program_template_test.go`

- [ ] **Step 1: Write test for Update (in-place, no linked programs)**

```go
// api/internal/store/memory/program_template_test.go
package memory

import (
	"context"
	"testing"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgramTemplateStore_Update(t *testing.T) {
	store := NewProgramTemplateRepository()
	ctx := context.Background()

	// Create a template
	tmpl := &domain.ProgramTemplate{
		Name: "Original",
		Entries: []domain.ProgramTemplateEntry{
			{Order: 0, ExerciseName: "Squat"},
		},
	}
	require.NoError(t, store.Create(ctx, tmpl))
	originalID := tmpl.ID
	originalCreatedAt := tmpl.CreatedAt

	// Update it
	tmpl.Name = "Updated"
	tmpl.Entries = []domain.ProgramTemplateEntry{
		{Order: 0, ExerciseName: "Bench Press"},
		{Order: 1, ExerciseName: "Deadlift"},
	}
	require.NoError(t, store.Update(ctx, tmpl))

	// Verify
	updated, err := store.GetByID(ctx, originalID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", updated.Name)
	assert.Equal(t, originalID, updated.ID)
	assert.Equal(t, originalCreatedAt, updated.CreatedAt)
	assert.True(t, updated.UpdatedAt.After(originalCreatedAt))
	assert.Len(t, updated.Entries, 2)
	assert.Equal(t, "Bench Press", updated.Entries[0].ExerciseName)
}

func TestProgramTemplateStore_Update_NotFound(t *testing.T) {
	store := NewProgramTemplateRepository()
	ctx := context.Background()

	tmpl := &domain.ProgramTemplate{
		Name: "Nonexistent",
	}
	// ID is zero-value UUID — not in store
	err := store.Update(ctx, tmpl)
	assert.Equal(t, domain.ErrNotFound, err)
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd api && go test ./internal/store/memory/... -run TestProgramTemplateStore_Update -v`
Expected: compile error — `Update` not defined

- [ ] **Step 3: Implement `Update` on memory store**

Add to `api/internal/store/memory/program_template.go`:

```go
func (s *programTemplateStore) Update(ctx context.Context, tmpl *domain.ProgramTemplate) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.templates[tmpl.ID]
	if !exists {
		return domain.ErrNotFound
	}

	now := time.Now()
	tmpl.CreatedAt = existing.CreatedAt
	tmpl.CreatedBy = existing.CreatedBy
	tmpl.ArchivedAt = existing.ArchivedAt
	tmpl.UpdatedAt = now

	for i := range tmpl.Entries {
		tmpl.Entries[i].ID = uuid.New()
		tmpl.Entries[i].ProgramTemplateID = tmpl.ID
	}

	s.templates[tmpl.ID] = tmpl
	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd api && go test ./internal/store/memory/... -run TestProgramTemplateStore_Update -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add api/internal/store/memory/program_template.go api/internal/store/memory/program_template_test.go
git commit -m "feat: implement Update on memory ProgramTemplate store"
```

---

### Task 5: Postgres Store — Implement `Update`

**Files:**
- Modify: `api/internal/store/postgres/program_template.go`

- [ ] **Step 1: Implement `Update` on postgres store**

Add to `api/internal/store/postgres/program_template.go`. Reuse existing `insertEntries` helper:

```go
func (r *programTemplateRepository) Update(ctx context.Context, tmpl *domain.ProgramTemplate) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Delete old entries
	_, err = tx.Exec(ctx, `DELETE FROM program_template_entries WHERE program_template_id = $1`, tmpl.ID)
	if err != nil {
		slog.Error("Failed to delete old program template entries", "id", tmpl.ID, "error", err)
		return err
	}

	// Update template fields
	query := `
		UPDATE program_templates
		SET name = $2, description = $3, notes = $4, metadata = $5,
		    weeks = $6, days_per_week = $7, source_template_id = $8, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	err = tx.QueryRow(ctx, query,
		tmpl.ID, tmpl.Name, tmpl.Description, tmpl.Notes, tmpl.Metadata,
		tmpl.Weeks, tmpl.DaysPerWeek, tmpl.SourceTemplateID,
	).Scan(&tmpl.UpdatedAt)
	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to update program template", "id", tmpl.ID, "error", err)
		return err
	}

	// Insert new entries
	if len(tmpl.Entries) > 0 {
		if err = r.insertEntries(ctx, tx, tmpl.ID, tmpl.Entries); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
```

- [ ] **Step 2: Update `GetByID` to scan `source_template_id`**

In the `GetByID` method, update the SQL query to include `source_template_id`:

```sql
SELECT id, name, description, notes, metadata, weeks, days_per_week, created_by, created_at, updated_at, archived_at, source_template_id
FROM program_templates
WHERE id = $1
```

Add `&tmpl.SourceTemplateID` to the `Scan()` call.

- [ ] **Step 3: Update `List` queries to scan `source_template_id`**

Both the `includeArchived` and default queries in `List` need the same change — add `source_template_id` to SELECT and to `Scan()`.

- [ ] **Step 4: Update `Create` to include `source_template_id`**

Update the `Create` INSERT query:

```sql
INSERT INTO program_templates (id, name, description, notes, metadata, weeks, days_per_week, created_by, source_template_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
RETURNING created_at, updated_at
```

Add `tmpl.SourceTemplateID` to the QueryRow args.

- [ ] **Step 5: Verify build**

Run: `cd api && go build ./...`
Expected: compiles without errors

- [ ] **Step 6: Commit**

```bash
git add api/internal/store/postgres/program_template.go
git commit -m "feat: implement Update on postgres ProgramTemplate store and add source_template_id to queries"
```

---

### Task 6: OpenAPI Spec — Add Edit Endpoint and `source_template_id`

**Files:**
- Modify: `api/openapi/openapi.yaml`

- [ ] **Step 1: Add `source_template_id` to `ProgramTemplate` schema**

In the `ProgramTemplate` schema (after `archived_at`, around line 675), add:

```yaml
        source_template_id:
          type: string
          format: uuid
          nullable: true
          description: ID of the template this was derived from (version lineage)
```

- [ ] **Step 2: Add the edit endpoint**

Add after the `/program-templates/{id}/generate` path (around line 183):

```yaml
  /program-templates/{id}/edit:
    parameters:
      - $ref: '#/components/parameters/ProgramTemplateId'
    post:
      summary: Edit program template
      description: |
        Edits a program template. If no Programs reference this template, it is updated in-place (200).
        If Programs reference it, a new template is created with source_template_id pointing to the
        original, and the original is archived (201).
      operationId: editProgramTemplate
      tags: [ProgramTemplate]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ProgramTemplateCreate'
      responses:
        '200':
          description: Template updated in-place (no linked programs)
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ProgramTemplate'
        '201':
          description: New template version created (linked programs exist, original archived)
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ProgramTemplate'
        '400':
          $ref: '#/components/responses/ValidationError'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '404':
          $ref: '#/components/responses/NotFound'
```

- [ ] **Step 3: Run code generation**

Run: `cd api && task generate`
Expected: `pkg/openapi/server.gen.go` regenerated with new `EditProgramTemplate` method in the server interface

- [ ] **Step 4: Commit**

```bash
git add api/openapi/openapi.yaml api/pkg/openapi/server.gen.go
git commit -m "feat: add edit endpoint and source_template_id to OpenAPI spec"
```

---

### Task 7: Handler — Implement `EditProgramTemplate`

**Files:**
- Modify: `api/internal/handler/program_template.go`
- Create: `api/internal/handler/program_template_test.go`

- [ ] **Step 1: Write tests for the edit handler**

```go
// api/internal/handler/program_template_test.go
package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/middleware"
	"github.com/aoshimash/rx/api/internal/store/memory"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupProgramTemplateTestRouter() (chi.Router, *ProgramTemplateHandler) {
	templateRepo := memory.NewProgramTemplateRepository()
	programRepo := memory.NewProgramRepository()

	h := NewProgramTemplateHandler(templateRepo, programRepo)

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(middleware.NewStubProvider()))
	r.Post("/program-templates", h.CreateProgramTemplate)
	r.Get("/program-templates/{id}", h.GetProgramTemplate)
	r.Post("/program-templates/{id}/edit", h.EditProgramTemplate)
	r.Post("/program-templates/{id}/archive", h.ArchiveProgramTemplate)
	r.Post("/program-templates/{id}/generate", h.GenerateProgram)
	return r, h
}

func createTestTemplate(t *testing.T, router chi.Router) map[string]interface{} {
	t.Helper()
	body := `{
		"name": "Test Template",
		"entries": [
			{"exercise_name": "Squat", "order": 0, "sets": 3, "reps": 5, "rpe": 8}
		]
	}`
	req := httptest.NewRequest(http.MethodPost, "/program-templates", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	return result
}

func TestEditProgramTemplate_InPlace(t *testing.T) {
	router, _ := setupProgramTemplateTestRouter()

	// Create template
	tmpl := createTestTemplate(t, router)
	id := tmpl["id"].(string)

	// Edit it (no linked programs → 200)
	editBody := `{
		"name": "Updated Template",
		"entries": [
			{"exercise_name": "Bench Press", "order": 0, "sets": 4, "reps": 6, "rpe": 7}
		]
	}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/edit", id), bytes.NewBufferString(editBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, id, result["id"])
	assert.Equal(t, "Updated Template", result["name"])
	entries := result["entries"].([]interface{})
	assert.Len(t, entries, 1)
	assert.Equal(t, "Bench Press", entries[0].(map[string]interface{})["exercise_name"])
}

func TestEditProgramTemplate_NewVersion(t *testing.T) {
	router, _ := setupProgramTemplateTestRouter()

	// Create template
	tmpl := createTestTemplate(t, router)
	id := tmpl["id"].(string)

	// Generate a program from it (creates a linked program)
	genBody := `{"target_weights": {"Squat": 100.0}}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/generate", id), bytes.NewBufferString(genBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Edit the template (linked programs → 201)
	editBody := `{
		"name": "Version 2",
		"entries": [
			{"exercise_name": "Deadlift", "order": 0, "sets": 1, "reps": 5, "rpe": 9}
		]
	}`
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/edit", id), bytes.NewBufferString(editBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-user")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.NotEqual(t, id, result["id"], "new version should have a new ID")
	assert.Equal(t, "Version 2", result["name"])
	assert.Equal(t, id, result["source_template_id"])

	// Verify old template is archived
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/program-templates/%s", id), nil)
	req.Header.Set("Authorization", "Bearer test-user")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var old map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &old))
	assert.NotNil(t, old["archived_at"], "old template should be archived")
}

func TestEditProgramTemplate_ArchivedTemplate(t *testing.T) {
	router, _ := setupProgramTemplateTestRouter()

	// Create and archive template
	tmpl := createTestTemplate(t, router)
	id := tmpl["id"].(string)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/archive", id), nil)
	req.Header.Set("Authorization", "Bearer test-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Try to edit archived template → 400
	editBody := `{"name": "Should Fail"}`
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/edit", id), bytes.NewBufferString(editBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-user")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEditProgramTemplate_NotFound(t *testing.T) {
	router, _ := setupProgramTemplateTestRouter()

	editBody := `{"name": "Nope"}`
	req := httptest.NewRequest(http.MethodPost, "/program-templates/00000000-0000-0000-0000-000000000000/edit", bytes.NewBufferString(editBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestEditProgramTemplate_ValidationError(t *testing.T) {
	router, _ := setupProgramTemplateTestRouter()

	tmpl := createTestTemplate(t, router)
	id := tmpl["id"].(string)

	// Empty name → validation error
	editBody := `{"name": ""}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/program-templates/%s/edit", id), bytes.NewBufferString(editBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd api && go test ./internal/handler/... -run TestEditProgramTemplate -v`
Expected: compile error — `EditProgramTemplate` not defined

- [ ] **Step 3: Implement `EditProgramTemplate` handler**

Add to `api/internal/handler/program_template.go`:

```go
// EditProgramTemplate handles POST /program-templates/{id}/edit
func (h *ProgramTemplateHandler) EditProgramTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUIDParam(r, "id", "program_template")
	if err != nil {
		middleware.WriteValidationError(w, err.Error(), nil)
		return
	}

	existing, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			middleware.WriteNotFoundError(w, "Program template not found")
			return
		}
		middleware.WriteInternalError(w, "Failed to retrieve program template")
		return
	}

	if existing.ArchivedAt != nil {
		middleware.WriteValidationError(w, "Cannot edit archived program template", nil)
		return
	}

	var req struct {
		Name        string                       `json:"name"`
		Description *string                      `json:"description,omitempty"`
		Notes       *string                      `json:"notes,omitempty"`
		Metadata    json.RawMessage              `json:"metadata,omitempty"`
		Weeks       *string                      `json:"weeks,omitempty"`
		DaysPerWeek *string                      `json:"days_per_week,omitempty"`
		Entries     []programTemplateEntryRequest `json:"entries,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, "Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	entries := make([]domain.ProgramTemplateEntry, len(req.Entries))
	for i, e := range req.Entries {
		entries[i] = convertProgramTemplateEntry(e)
	}

	// Build the new template content for validation
	newContent := &domain.ProgramTemplate{
		Name:        req.Name,
		Description: req.Description,
		Notes:       req.Notes,
		Metadata:    req.Metadata,
		Weeks:       req.Weeks,
		DaysPerWeek: req.DaysPerWeek,
		Entries:     entries,
	}

	if err := domain.ValidateProgramTemplate(newContent); err != nil {
		if handleValidationError(w, err) {
			return
		}
		middleware.WriteValidationError(w, "Validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	hasPrograms, err := h.programRepo.ExistsByProgramTemplateID(ctx, id)
	if err != nil {
		middleware.WriteInternalError(w, "Failed to check references")
		return
	}

	if !hasPrograms {
		// In-place update
		newContent.ID = existing.ID
		newContent.SourceTemplateID = existing.SourceTemplateID
		if err := h.repo.Update(ctx, newContent); err != nil {
			middleware.WriteInternalError(w, "Failed to update program template")
			return
		}
		// Re-fetch to get full state
		updated, err := h.repo.GetByID(ctx, newContent.ID)
		if err != nil {
			middleware.WriteInternalError(w, "Failed to retrieve updated program template")
			return
		}
		writeJSON(w, http.StatusOK, updated)
		return
	}

	// New version: create new template + archive old
	var createdBy *string
	if userID := middleware.GetUserID(ctx); userID != "" {
		createdBy = &userID
	}

	newTmpl := &domain.ProgramTemplate{
		Name:             req.Name,
		Description:      req.Description,
		Notes:            req.Notes,
		Metadata:         req.Metadata,
		Weeks:            req.Weeks,
		DaysPerWeek:      req.DaysPerWeek,
		CreatedBy:        createdBy,
		SourceTemplateID: &id,
		Entries:          entries,
	}

	if err := h.repo.Create(ctx, newTmpl); err != nil {
		middleware.WriteInternalError(w, "Failed to create new program template version")
		return
	}

	if err := h.repo.Archive(ctx, id); err != nil {
		middleware.WriteInternalError(w, "Failed to archive old program template")
		return
	}

	writeJSON(w, http.StatusCreated, newTmpl)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd api && go test ./internal/handler/... -run TestEditProgramTemplate -v`
Expected: all 5 tests PASS

- [ ] **Step 5: Commit**

```bash
git add api/internal/handler/program_template.go api/internal/handler/program_template_test.go
git commit -m "feat: implement EditProgramTemplate handler with version management"
```

---

### Task 8: Route Registration

**Files:**
- Modify: `api/cmd/server/main.go:130-138`

- [ ] **Step 1: Add the edit route**

Add after line 138 (the `Delete` route):

```go
r.Post("/program-templates/{id}/edit", programTemplateHandler.EditProgramTemplate)
```

- [ ] **Step 2: Verify build**

Run: `cd api && go build ./...`
Expected: compiles without errors

- [ ] **Step 3: Commit**

```bash
git add api/cmd/server/main.go
git commit -m "feat: register /program-templates/{id}/edit route"
```

---

### Task 9: Run Full Check

- [ ] **Step 1: Run full check suite**

Run: `cd api && task check`
Expected: generate + format + lint + test all pass

- [ ] **Step 2: Fix any issues found**

If `task generate` produced changes to `server.gen.go` that conflict with the handler (e.g., the generated interface now includes `EditProgramTemplate`), ensure the handler satisfies the interface. The generated interface from OpenAPI may use different parameter types — adapt the handler signature if needed.

- [ ] **Step 3: Final commit if needed**

```bash
git add -A
git commit -m "fix: resolve any check failures"
```
