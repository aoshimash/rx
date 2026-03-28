# Domain Models

This package contains the core domain entities for Rx.

## Entities

### Program
A training program containing sessions with scheduled exercises.

**Required fields**: `id`, `name`, `status`, `sessions`, `created_at`, `updated_at`
**Optional fields**: `notes`, `metadata`

### ProgramSession
A named training day within a program.

**Required fields**: `id`, `program_id`, `session_name`, `order`
**Optional fields**: `date`, `entries`

### ProgramSessionEntry
An exercise prescription within a session.

**Required fields**: `id`, `session_id`, `order`, `exercise_name`
**Optional fields**: `sets`, `reps`, `load_kg`, `notes`, `metadata`

### Log
A completed training session recording what was actually performed.

**Required fields**: `id`, `performed_at`, `entries`, `created_at`, `updated_at`
**Optional fields**: `program_id`, `session_name`, `started_at`, `finished_at`, `notes`, `metadata`

### LogEntry
A single performed exercise entry within a log.

**Required fields**: `id`, `log_id`, `order`, `exercise_name`
**Optional fields**: `sets`, `reps`, `load_kg`, `notes`, `video_object_key`, `started_at`, `finished_at`, `metadata`

## Validation

All entities have corresponding validation functions:

- `ValidateProgram(p *Program) error`
- `ValidateProgramSession(s *ProgramSession) error`
- `ValidateProgramSessionEntry(e *ProgramSessionEntry) error`
- `ValidateProgramStatusTransition(from, to ProgramStatus) error`
- `ValidateLog(l *Log) error`
- `ValidateLogEntry(e *LogEntry) error`

### Common Validation Helpers

- `ValidateTimestamp(t time.Time) error` - Check not in future
- `ValidateTimeRange(field string, startedAt, finishedAt *time.Time) error` - Check start before end
- `ValidateRequiredString(field, value string) error` - Check non-empty
- `ValidateStringLength(field, value string, min, max int) error` - Check length bounds

## Error Types

- `ValidationError` - Field-level validation errors
- `DomainError` - Domain-level errors with error codes

See `errors.go` for error code constants.

## Testing

Run tests:
```bash
go test ./internal/domain/...
```
