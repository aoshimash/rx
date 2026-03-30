# Move Video Storage to Set Level — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move video attachment from LogEntry level to LogSet level, using `video_object_key` (pre-signed URL flow) consistently.

**Architecture:** Remove `video_object_key` from LogEntry, replace `video_url` with `video_object_key` on LogSet. Propagate changes through OpenAPI spec → generated code → domain → handler → store → frontend types.

**Tech Stack:** Go, OpenAPI (oapi-codegen), PostgreSQL, TypeScript

---

### Task 1: Database Migration

**Files:**
- Create: `api/migrations/000009_video_to_set_level.up.sql`
- Create: `api/migrations/000009_video_to_set_level.down.sql`

- [ ] **Step 1: Create UP migration**

```sql
-- api/migrations/000009_video_to_set_level.up.sql
ALTER TABLE log_sets DROP COLUMN video_url;
ALTER TABLE log_sets ADD COLUMN video_object_key VARCHAR(500);
ALTER TABLE log_entries DROP COLUMN video_object_key;
```

- [ ] **Step 2: Create DOWN migration**

```sql
-- api/migrations/000009_video_to_set_level.down.sql
ALTER TABLE log_entries ADD COLUMN video_object_key VARCHAR(500);
ALTER TABLE log_sets DROP COLUMN video_object_key;
ALTER TABLE log_sets ADD COLUMN video_url VARCHAR(2000);
```

- [ ] **Step 3: Commit**

```bash
cd api
git add migrations/000009_video_to_set_level.up.sql migrations/000009_video_to_set_level.down.sql
git commit -m "feat: add migration to move video from entry to set level"
```

---

### Task 2: OpenAPI Spec Update

**Files:**
- Modify: `api/openapi/openapi.yaml`

- [ ] **Step 1: Update LogEntry schema — remove video_object_key**

In `LogEntry` (around line 1400), remove these lines:
```yaml
        video_object_key:
          type: string
          maxLength: 500
          description: Object key for uploaded video in storage (from upload-url response)
```

- [ ] **Step 2: Update LogEntryCreate schema — remove video_object_key**

In `LogEntryCreate` (around line 1434), remove these lines:
```yaml
        video_object_key:
          type: string
          maxLength: 500
```

- [ ] **Step 3: Update LogSet schema — replace video_url with video_object_key**

In `LogSet` (around line 1465), replace:
```yaml
        video_url:
          type: string
          maxLength: 2000
```
with:
```yaml
        video_object_key:
          type: string
          maxLength: 500
          description: Object key for uploaded video in storage (from upload-url response)
```

- [ ] **Step 4: Update LogSetCreate schema — replace video_url with video_object_key**

In `LogSetCreate` (around line 1484), replace:
```yaml
        video_url:
          type: string
          maxLength: 2000
```
with:
```yaml
        video_object_key:
          type: string
          maxLength: 500
```

- [ ] **Step 5: Update VideoUploadURLResponse description**

In `VideoUploadURLResponse` (around line 1549), change:
```yaml
          description: Object key to be stored in LogEntry.video_object_key after successful upload
```
to:
```yaml
          description: Object key to be stored in LogSet.video_object_key after successful upload
```

- [ ] **Step 6: Update VideoDownloadURLRequest description**

In `VideoDownloadURLRequest` (around line 1563), change:
```yaml
          description: Object key from LogEntry.video_object_key
```
to:
```yaml
          description: Object key from LogSet.video_object_key
```

- [ ] **Step 7: Update upload-url endpoint description**

In `/videos/upload-url` (around line 511), change:
```yaml
        After successful upload, include the returned object_key in the LogEntry.
```
to:
```yaml
        After successful upload, include the returned object_key in the LogSet.
```

- [ ] **Step 8: Run code generation**

```bash
cd api && task generate
```

Expected: `pkg/openapi/server.gen.go` regenerated without errors.

- [ ] **Step 9: Commit**

```bash
cd api
git add openapi/openapi.yaml pkg/openapi/server.gen.go
git commit -m "feat: update OpenAPI spec to move video_object_key from LogEntry to LogSet"
```

---

### Task 3: Domain Model and Validation Changes

**Files:**
- Modify: `api/internal/domain/log_set.go`
- Modify: `api/internal/domain/log.go`
- Modify: `api/internal/domain/validation.go`

- [ ] **Step 1: Update LogSet struct**

In `api/internal/domain/log_set.go`, replace:
```go
	VideoURL  *string                `json:"video_url,omitempty"`
```
with:
```go
	VideoObjectKey *string                `json:"video_object_key,omitempty"`
```

- [ ] **Step 2: Update LogEntry struct**

In `api/internal/domain/log.go`, remove this line from the `LogEntry` struct:
```go
	VideoObjectKey *string                `json:"video_object_key,omitempty"`
```

- [ ] **Step 3: Update ValidateLogSet — change video_url to video_object_key**

In `api/internal/domain/validation.go` (around line 533), replace:
```go
	if s.VideoURL != nil {
		if err := ValidateStringLength("video_url", *s.VideoURL, 1, 2000); err != nil {
			return err
		}
	}
```
with:
```go
	if s.VideoObjectKey != nil {
		if err := ValidateStringLength("video_object_key", *s.VideoObjectKey, 1, 500); err != nil {
			return err
		}
	}
```

- [ ] **Step 4: Update ValidateLogEntry — remove video_object_key validation**

In `api/internal/domain/validation.go` (around line 611), remove these lines:
```go
	if e.VideoObjectKey != nil {
		if err := ValidateStringLength("video_object_key", *e.VideoObjectKey, 1, 500); err != nil {
			return err
		}
	}
```

- [ ] **Step 5: Verify compilation**

```bash
cd api && go build ./...
```

Expected: no compilation errors.

- [ ] **Step 6: Commit**

```bash
cd api
git add internal/domain/log_set.go internal/domain/log.go internal/domain/validation.go
git commit -m "feat: move VideoObjectKey from LogEntry to LogSet in domain model"
```

---

### Task 4: Update Domain Tests

**Files:**
- Modify: `api/internal/domain/log_set_test.go`
- Modify: `api/internal/domain/log_test.go`

- [ ] **Step 1: Update log_set_test.go — rename video_url tests to video_object_key**

In `api/internal/domain/log_set_test.go`, replace the two video test cases:

Replace:
```go
	t.Run("video_url too long", func(t *testing.T) {
		s := validLogSet()
		longURL := strings.Repeat("a", 2001)
		s.VideoURL = &longURL
		err := ValidateLogSet(s)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "video_url", ve.Field)
	})

	t.Run("empty video_url", func(t *testing.T) {
		s := validLogSet()
		empty := ""
		s.VideoURL = &empty
		err := ValidateLogSet(s)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "video_url", ve.Field)
	})
```

with:
```go
	t.Run("video_object_key too long", func(t *testing.T) {
		s := validLogSet()
		longKey := strings.Repeat("a", 501)
		s.VideoObjectKey = &longKey
		err := ValidateLogSet(s)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "video_object_key", ve.Field)
	})

	t.Run("empty video_object_key", func(t *testing.T) {
		s := validLogSet()
		empty := ""
		s.VideoObjectKey = &empty
		err := ValidateLogSet(s)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "video_object_key", ve.Field)
	})

	t.Run("valid video_object_key", func(t *testing.T) {
		s := validLogSet()
		key := "videos/user123/abc.mp4"
		s.VideoObjectKey = &key
		err := ValidateLogSet(s)
		require.NoError(t, err)
	})
```

- [ ] **Step 2: Update log_test.go — remove video_object_key tests from LogEntry**

In `api/internal/domain/log_test.go`, remove these two test cases:

```go
	t.Run("valid video_object_key", func(t *testing.T) {
		e := validEntry()
		key := "videos/test.mp4"
		e.VideoObjectKey = &key
		err := ValidateLogEntry(e)
		assert.NoError(t, err)
	})

	t.Run("empty video_object_key", func(t *testing.T) {
		e := validEntry()
		key := ""
		e.VideoObjectKey = &key
		err := ValidateLogEntry(e)
		assert.Error(t, err)
	})
```

- [ ] **Step 3: Run tests**

```bash
cd api && go test ./internal/domain/... -v
```

Expected: all tests pass.

- [ ] **Step 4: Commit**

```bash
cd api
git add internal/domain/log_set_test.go internal/domain/log_test.go
git commit -m "test: update domain tests for video_object_key on LogSet"
```

---

### Task 5: Handler Changes

**Files:**
- Modify: `api/internal/handler/log.go`

- [ ] **Step 1: Update logSetRequest struct**

In `api/internal/handler/log.go` (line 19), replace:
```go
	VideoURL  *string                `json:"video_url,omitempty"`
```
with:
```go
	VideoObjectKey *string                `json:"video_object_key,omitempty"`
```

- [ ] **Step 2: Update logEntryRequest struct — remove VideoObjectKey**

In `api/internal/handler/log.go` (line 29), remove:
```go
	VideoObjectKey *string                `json:"video_object_key,omitempty"`
```

- [ ] **Step 3: Update parseLogRequest — remove entry VideoObjectKey mapping**

In `api/internal/handler/log.go` (around line 114), change the entry construction from:
```go
		entry := domain.LogEntry{
			ExerciseName:   entryReq.ExerciseName,
			Order:          i,
			Fields:         entryReq.Fields,
			Notes:          entryReq.Notes,
			VideoObjectKey: entryReq.VideoObjectKey,
		}
```
to:
```go
		entry := domain.LogEntry{
			ExerciseName: entryReq.ExerciseName,
			Order:        i,
			Fields:       entryReq.Fields,
			Notes:        entryReq.Notes,
		}
```

- [ ] **Step 4: Update parseLogRequest — update set mapping**

In `api/internal/handler/log.go` (around line 126), change:
```go
				sets[j] = domain.LogSet{
					SetNumber: setReq.SetNumber,
					Fields:    setReq.Fields,
					VideoURL:  setReq.VideoURL,
					Notes:     setReq.Notes,
				}
```
to:
```go
				sets[j] = domain.LogSet{
					SetNumber:      setReq.SetNumber,
					Fields:         setReq.Fields,
					VideoObjectKey: setReq.VideoObjectKey,
					Notes:          setReq.Notes,
				}
```

- [ ] **Step 5: Verify compilation**

```bash
cd api && go build ./...
```

Expected: no compilation errors.

- [ ] **Step 6: Commit**

```bash
cd api
git add internal/handler/log.go
git commit -m "feat: update log handler for video_object_key on sets"
```

---

### Task 6: Postgres Store Changes

**Files:**
- Modify: `api/internal/store/postgres/log.go`

- [ ] **Step 1: Update insertEntries — remove video_object_key from entry INSERT**

In `api/internal/store/postgres/log.go` (around line 86), replace the entry INSERT query and exec:
```go
		query := `
			INSERT INTO log_entries (
				id, log_id, "order", exercise_name,
				fields,
				notes, video_object_key, started_at, finished_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`

		_, err = tx.Exec(ctx, query,
			entryID,
			logID,
			entries[i].Order,
			entries[i].ExerciseName,
			fieldsJSON,
			entries[i].Notes,
			entries[i].VideoObjectKey,
			entries[i].StartedAt,
			entries[i].FinishedAt,
		)
```
with:
```go
		query := `
			INSERT INTO log_entries (
				id, log_id, "order", exercise_name,
				fields,
				notes, started_at, finished_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

		_, err = tx.Exec(ctx, query,
			entryID,
			logID,
			entries[i].Order,
			entries[i].ExerciseName,
			fieldsJSON,
			entries[i].Notes,
			entries[i].StartedAt,
			entries[i].FinishedAt,
		)
```

- [ ] **Step 2: Update insertEntries — rename video_url to video_object_key in set INSERT**

In the same function (around line 126), replace:
```go
			setQuery := `
				INSERT INTO log_sets (id, entry_id, set_number, fields, video_url, notes)
				VALUES ($1, $2, $3, $4, $5, $6)
			`
			_, err = tx.Exec(ctx, setQuery,
				setID,
				entryID,
				entries[i].Sets[k].SetNumber,
				setFieldsJSON,
				entries[i].Sets[k].VideoURL,
				entries[i].Sets[k].Notes,
			)
```
with:
```go
			setQuery := `
				INSERT INTO log_sets (id, entry_id, set_number, fields, video_object_key, notes)
				VALUES ($1, $2, $3, $4, $5, $6)
			`
			_, err = tx.Exec(ctx, setQuery,
				setID,
				entryID,
				entries[i].Sets[k].SetNumber,
				setFieldsJSON,
				entries[i].Sets[k].VideoObjectKey,
				entries[i].Sets[k].Notes,
			)
```

- [ ] **Step 3: Update getEntriesForLog — remove video_object_key from entry SELECT**

In `getEntriesForLog` (around line 199), replace the query:
```go
	query := `
		SELECT id, log_id, "order", exercise_name,
		       fields,
		       notes, video_object_key, started_at, finished_at
		FROM log_entries
		WHERE log_id = $1
		ORDER BY "order" ASC
	`
```
with:
```go
	query := `
		SELECT id, log_id, "order", exercise_name,
		       fields,
		       notes, started_at, finished_at
		FROM log_entries
		WHERE log_id = $1
		ORDER BY "order" ASC
	`
```

And replace the Scan (around line 218):
```go
		err := rows.Scan(
			&entry.ID,
			&entry.LogID,
			&entry.Order,
			&entry.ExerciseName,
			&fieldsRaw,
			&entry.Notes,
			&entry.VideoObjectKey,
			&entry.StartedAt,
			&entry.FinishedAt,
		)
```
with:
```go
		err := rows.Scan(
			&entry.ID,
			&entry.LogID,
			&entry.Order,
			&entry.ExerciseName,
			&fieldsRaw,
			&entry.Notes,
			&entry.StartedAt,
			&entry.FinishedAt,
		)
```

- [ ] **Step 4: Update getSetsForEntry — rename video_url to video_object_key**

In `getSetsForEntry` (around line 256), replace:
```go
	query := `
		SELECT id, entry_id, set_number, fields, video_url, notes
		FROM log_sets
		WHERE entry_id = $1
		ORDER BY set_number ASC
	`
```
with:
```go
	query := `
		SELECT id, entry_id, set_number, fields, video_object_key, notes
		FROM log_sets
		WHERE entry_id = $1
		ORDER BY set_number ASC
	`
```

And replace the Scan (around line 279):
```go
			&s.VideoURL,
```
with:
```go
			&s.VideoObjectKey,
```

- [ ] **Step 5: Update getEntriesForLogsBatch — remove video_object_key from entry SELECT**

In `getEntriesForLogsBatch` (around line 460), replace:
```go
	query := `
		SELECT id, log_id, "order", exercise_name,
		       fields,
		       notes, video_object_key, started_at, finished_at
		FROM log_entries
		WHERE log_id = ANY($1::uuid[])
		ORDER BY log_id, "order" ASC
	`
```
with:
```go
	query := `
		SELECT id, log_id, "order", exercise_name,
		       fields,
		       notes, started_at, finished_at
		FROM log_entries
		WHERE log_id = ANY($1::uuid[])
		ORDER BY log_id, "order" ASC
	`
```

And replace the Scan (around line 480):
```go
		if err := rows.Scan(
			&entry.ID, &entry.LogID, &entry.Order, &entry.ExerciseName,
			&fieldsRaw, &entry.Notes, &entry.VideoObjectKey, &entry.StartedAt, &entry.FinishedAt,
		); err != nil {
```
with:
```go
		if err := rows.Scan(
			&entry.ID, &entry.LogID, &entry.Order, &entry.ExerciseName,
			&fieldsRaw, &entry.Notes, &entry.StartedAt, &entry.FinishedAt,
		); err != nil {
```

- [ ] **Step 6: Update getSetsForEntriesBatch — rename video_url to video_object_key**

In `getSetsForEntriesBatch` (around line 518), replace:
```go
	query := `
		SELECT id, entry_id, set_number, fields, video_url, notes
		FROM log_sets
		WHERE entry_id = ANY($1::uuid[])
		ORDER BY entry_id, set_number ASC
	`
```
with:
```go
	query := `
		SELECT id, entry_id, set_number, fields, video_object_key, notes
		FROM log_sets
		WHERE entry_id = ANY($1::uuid[])
		ORDER BY entry_id, set_number ASC
	`
```

And replace the Scan (around line 536):
```go
		if err := rows.Scan(&s.ID, &s.EntryID, &s.SetNumber, &fieldsRaw, &s.VideoURL, &s.Notes); err != nil {
```
with:
```go
		if err := rows.Scan(&s.ID, &s.EntryID, &s.SetNumber, &fieldsRaw, &s.VideoObjectKey, &s.Notes); err != nil {
```

- [ ] **Step 7: Verify compilation**

```bash
cd api && go build ./...
```

Expected: no compilation errors.

- [ ] **Step 8: Commit**

```bash
cd api
git add internal/store/postgres/log.go
git commit -m "feat: update postgres store for video_object_key on sets"
```

---

### Task 7: Frontend Type Changes

**Files:**
- Modify: `web/types/api.ts`

- [ ] **Step 1: Remove video_object_key from LogEntry and LogEntryCreate**

In `web/types/api.ts`, remove `video_object_key?: string;` from both `LogEntry` (line 89) and `LogEntryCreate` (line 98).

- [ ] **Step 2: Add LogSet and LogSetCreate types**

After the `LogEntryCreate` interface (after line 101), add:

```typescript
// ============================================================================
// LogSet
// ============================================================================

export interface LogSet {
  id: string;
  entry_id: string;
  set_number: number;
  fields: Record<string, unknown>;
  video_object_key?: string;
  notes?: string;
}

export interface LogSetCreate {
  set_number: number;
  fields: Record<string, unknown>;
  video_object_key?: string;
  notes?: string;
}
```

Also add `sets?: LogSet[];` to the `LogEntry` interface and `sets?: LogSetCreate[];` to `LogEntryCreate`.

- [ ] **Step 3: Run web check**

```bash
cd web && pnpm check
```

Expected: no errors.

- [ ] **Step 4: Commit**

```bash
cd web
git add types/api.ts
git commit -m "feat: update frontend types for video_object_key on LogSet"
```

---

### Task 8: Run Full Check

- [ ] **Step 1: Run API checks**

```bash
cd api && task check
```

Expected: generate + format + lint + test all pass.

- [ ] **Step 2: Run web checks**

```bash
cd web && pnpm check
```

Expected: lint + format pass.

- [ ] **Step 3: Verify no stale references**

Search the entire codebase for leftover `video_url` or `VideoURL` references (excluding migration files and this plan):

```bash
cd /home/aoshima/dev/github/aoshimash/rx
grep -rn "video_url\|VideoURL" --include="*.go" --include="*.ts" --include="*.yaml" api/ web/ | grep -v migrations/ | grep -v plans/
```

Expected: no matches.
