# Implementation Plan: MinIO Integration Test Coverage

**Branch**: `004-minio-test-coverage` | **Date**: 2026-01-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-minio-test-coverage/spec.md`

## Summary

Introduce MinIO to the test environment to expand S3 storage provider integration test coverage. Connect the existing S3 provider implementation (`api/internal/storage/s3/`) with actual S3-compatible storage (MinIO) and verify pre-signed URL generation, object upload/download/delete operations. Separate from unit tests using build tags (`//go:build integration`) and automatically execute on all PRs in CI pipeline.

## Technical Context

**Language/Version**: Go 1.25+  
**Primary Dependencies**: aws-sdk-go-v2, chi, oapi-codegen  
**Storage**: MinIO (S3-compatible, defined in docker-compose.yml)  
**Testing**: standard testing package with build tags (`//go:build integration`)  
**Target Platform**: Linux server (API), GitHub Actions (CI)  
**Project Type**: Monorepo - API component (`api/`)  
**Performance Goals**: Integration tests complete in <10 seconds locally  
**Constraints**: Tests must skip gracefully when MinIO unavailable  
**Scale/Scope**: 3 storage provider methods to test, 1 CI workflow to update

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with OPTel Workout Constitution principles:

- **Dumb Backend**: ✅ Pass - This feature adds tests only, no business logic changes
- **Domain-Driven Schema-First**: ✅ Pass - No API contract changes, testing existing implementation
- **Terminology**: ✅ Pass - Uses standard testing terminology
- **Clean Architecture**: ✅ Pass - Tests verify existing architecture, no structural changes
- **Monorepo Structure**: ✅ Pass - Changes confined to `api/` component

All principles satisfied. No violations to justify.

## Project Structure

### Documentation (this feature)

```text
specs/004-minio-test-coverage/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (minimal - test infrastructure only)
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A - no API changes)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
api/
├── internal/
│   └── storage/
│       └── s3/
│           ├── provider.go           # Existing - no changes
│           ├── provider_test.go      # Existing unit tests
│           └── provider_integration_test.go  # NEW: Integration tests
├── internal/
│   └── handler/
│       └── video_integration_test.go # NEW: Handler integration tests
└── Makefile                          # Update: Add integration test target

.github/
└── workflows/
    └── api-ci.yml                    # Update: Add MinIO service and integration tests

docker-compose.yml                    # Existing - MinIO already defined
```

**Structure Decision**: Integration tests are placed alongside existing tests with `_integration_test.go` suffix and `//go:build integration` build tag. This follows Go conventions and keeps tests near the code they verify.

## Complexity Tracking

> No Constitution Check violations. This section is empty.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A | N/A | N/A |
