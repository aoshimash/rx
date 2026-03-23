# ProgramTemplate Editing with Automatic Version Management

## Summary

Add the ability to edit ProgramTemplates with smart version management. When no Programs reference a template, edits are applied in-place. When Programs exist, a new version is created and the old template is automatically archived.

## Motivation

ProgramTemplates are currently immutable after creation. Users need to edit templates for:

- Fixing typos or mistakes in prescriptions
- Tweaking set/rep/RPE schemes for the next training cycle
- Customizing templates received from others (future sharing)

Without editing, users must delete and recreate templates, losing the link to any Programs generated from them.

## Design

### Data Model

Add one column to `program_templates`:

```sql
ALTER TABLE program_templates
  ADD COLUMN source_template_id UUID REFERENCES program_templates(id) ON DELETE SET NULL;
```

- Nullable — original templates and manually created ones have `NULL`
- `ON DELETE SET NULL` — if the source is deleted, the chain breaks gracefully

Domain struct change:

```go
type ProgramTemplate struct {
    // ...existing fields...
    SourceTemplateID *uuid.UUID `json:"source_template_id,omitempty"`
}
```

Repository interface addition:

```go
Update(ctx context.Context, tmpl *domain.ProgramTemplate) error
```

### API Endpoint

```
POST /program-templates/{id}/edit
```

Request body: same structure as `ProgramTemplateCreate` (name, description, notes, metadata, weeks, days_per_week, entries).

| Case | Status | Body |
|------|--------|------|
| No linked Programs → in-place update | 200 OK | Updated template |
| Linked Programs → new version + archive | 201 Created | New template (new ID) |
| Template not found | 404 | |
| Archived template | 400 | |
| Validation error | 400 | |

`POST` is used instead of `PUT` because the versioning case creates a new resource (new ID), which is not idempotent.

### OpenAPI Schema Changes

1. **New endpoint**: `POST /program-templates/{id}/edit` with `ProgramTemplateCreate` as request body
2. **Response schema**: Add `source_template_id` (nullable UUID) to the existing `ProgramTemplate` schema
3. **Request body**: Reuse `ProgramTemplateCreate` schema (no new request type needed)

### Backend Processing Flow

```
EditProgramTemplate(id, body):
  1. Get existing template (GetByID)
     → 404 if not found
     → 400 if archived

  2. Validate body (ValidateProgramTemplate)
     → 400 if invalid

  3. Check linked Programs (ExistsByProgramTemplateID)

  4a. No linked Programs:
      BEGIN TRANSACTION
      - DELETE FROM program_template_entries WHERE program_template_id = id
      - UPDATE program_templates SET name=..., updated_at=NOW() WHERE id = id
      - INSERT program_template_entries (new entries)
      COMMIT
      → return 200, updated template

  4b. Linked Programs exist:
      BEGIN TRANSACTION
      - Create new template (new UUID, source_template_id = old ID)
      - Create new entries with new UUIDs
      - Archive old template (set archived_at)
      COMMIT
      → return 201, new template
```

### Frontend UI Flow

**Edit button press:**

```
[Edit button]
    │
    ├── GET /programs?program_template_id={id}&limit=1
    │
    ├── data.length == 0
    │     → Navigate to editor
    │
    └── data.length > 0
          → Dialog: "This template is used by existing programs.
             Editing will create a new version."
            [Cancel] [Continue]
          → Continue → Navigate to editor
```

**Save:**

```
[Save]
    │
    ├── POST /program-templates/{id}/edit
    │
    ├── 200 OK → Stay on detail page (same ID)
    │
    └── 201 Created → Redirect to new template ID
```

**Editor:** Reuse existing template creation components (SessionAccordion etc.) with prefilled values from the existing template.

**Detail page:** When `source_template_id` is present, show a small "Derived from: {source template name}" link.

### Edit Behavior Summary

| Linked Programs | Old Template | New Template | User Experience |
|-----------------|-------------|-------------|-----------------|
| None | Updated in-place | N/A | Normal edit |
| Any (all statuses) | Archived | Created with source_template_id | New version with dialog |

### Linked Program Scope

All Programs with `program_template_id` matching the template are counted, regardless of status (created, ongoing, completed, cancelled).

## Out of Scope

- Export/Import for sharing
- Content hashing for integrity verification
- Version history list UI (only a "derived from" link)
- Batch editing
- Template merging

## References

- [GitHub Issue #141](https://github.com/aoshimash/rx/issues/141)
- `docs/DOMAIN_MODEL.md` — Three-tier lifecycle
- `docs/PHILOSOPHY.md` — "Dumb Backend" and API-First principles
