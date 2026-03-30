# Move Video Storage from LogEntry to LogSet Level

## Summary

Consolidate video attachment from two levels (LogEntry + LogSet) to LogSet only, using object storage (pre-signed URLs) consistently. LogEntry-level `video_object_key` is removed; LogSet gains `video_object_key` replacing the previous `video_url` direct-link field.

## Motivation

- Videos are recorded per-set, not per-exercise — the LogEntry-level field is unnecessary
- `LogSet.video_url` (direct URL) lacks access control and lifecycle management
- Unifying on `video_object_key` + pre-signed URL flow aligns with the existing `storage.Provider` infrastructure and the "Dumb Backend" philosophy

## Design

### Database Migration (000009)

**UP:**
```sql
ALTER TABLE log_sets DROP COLUMN video_url;
ALTER TABLE log_sets ADD COLUMN video_object_key VARCHAR(500);
ALTER TABLE log_entries DROP COLUMN video_object_key;
```

**DOWN:**
```sql
ALTER TABLE log_entries ADD COLUMN video_object_key VARCHAR(500);
ALTER TABLE log_sets DROP COLUMN video_object_key;
ALTER TABLE log_sets ADD COLUMN video_url VARCHAR(2000);
```

No data migration needed — no production data exists.

### OpenAPI Spec Changes

**LogSet schema:**
- Remove `video_url` (string, maxLength 2000)
- Add `video_object_key` (string, maxLength 500, description: "Object key for uploaded video in storage")

**LogSetCreate schema:**
- Remove `video_url` (string, maxLength 2000)
- Add `video_object_key` (string, maxLength 500)

**LogEntry schema:**
- Remove `video_object_key`

**LogEntryCreate schema:**
- Remove `video_object_key`

**VideoUploadURLResponse:**
- Update `object_key` description: "Object key to be stored in LogSet.video_object_key after successful upload"

**VideoDownloadURLRequest:**
- Update `object_key` description: "Object key from LogSet.video_object_key"

### Domain Model Changes

**`domain/log_set.go` — LogSet struct:**
- `VideoURL *string` → `VideoObjectKey *string` (json tag: `video_object_key`)

**`domain/log.go` — LogEntry struct:**
- Remove `VideoObjectKey *string`

### Handler Changes

**`handler/log.go`:**
- `logEntryRequest` struct: remove `VideoObjectKey`
- `logSetRequest` struct: `VideoURL` → `VideoObjectKey`
- Entry mapping: remove `VideoObjectKey` assignment
- Set mapping: `VideoURL` → `VideoObjectKey`

### Postgres Store Changes

**`store/postgres/log.go`:**
- Entry INSERT: remove `video_object_key` from column list and values (9 → 8 params)
- Entry SELECT: remove `video_object_key` from column list and Scan
- Set INSERT: `video_url` → `video_object_key`
- Set SELECT: `video_url` → `video_object_key`
- Applies to both single-entry and batch queries

### Frontend Type Changes

**`web/types/api.ts`:**
- `LogEntry`: remove `video_object_key`
- `LogEntryCreate`: remove `video_object_key`
- `LogSet` type (if defined): add `video_object_key?: string`
- `LogSetCreate` type (if defined): add `video_object_key?: string`

### No Changes Required

- `storage.Provider` interface — unchanged
- `storage/s3/provider.go` — unchanged
- `handler/video.go` — upload/download URL endpoints unchanged
- Video upload/download API endpoints — unchanged

## Affected Files

1. `api/migrations/000009_video_to_set_level.up.sql` (new)
2. `api/migrations/000009_video_to_set_level.down.sql` (new)
3. `api/openapi/openapi.yaml`
4. `api/internal/domain/log.go`
5. `api/internal/domain/log_set.go`
6. `api/internal/handler/log.go`
7. `api/internal/store/postgres/log.go`
8. `web/types/api.ts`
9. Generated code via `task generate`
