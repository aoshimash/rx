# Tasks: MinIO Integration Test Coverage

**Input**: Design documents from `/specs/004-minio-test-coverage/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md

**Tests**: This feature is specifically about adding integration tests. Test tasks are the primary implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions (OPTel Monorepo)

- API component: `api/internal/`, `api/pkg/`, `api/cmd/`
- CI workflows: `.github/workflows/`
- Docker: `docker-compose.yml` (root)

---

## Phase 1: Setup (Test Infrastructure Foundation)

**Purpose**: Create shared test helper utilities for MinIO integration tests

- [ ] T001 Create test helper file with build tag in `api/internal/storage/s3/testhelper_integration_test.go`
- [ ] T002 Implement `skipIfMinIOUnavailable(t *testing.T)` helper function for graceful test skipping
- [ ] T003 Implement `setupTestBucket(ctx, client, bucket)` helper function for bucket creation
- [ ] T004 Implement `cleanupTestObjects(ctx, client, bucket, keys)` helper function for test cleanup
- [ ] T005 Implement `newTestMinIOClient(t *testing.T)` helper function to create configured S3 client

---

## Phase 2: Foundational (Makefile Integration)

**Purpose**: Enable running integration tests via Makefile

**⚠️ CRITICAL**: This enables local execution of integration tests

- [ ] T006 Add `test-integration` target to `api/Makefile` that runs `go test -v -tags=integration ./...`
- [ ] T007 Add `test-all` target to `api/Makefile` that runs both unit and integration tests
- [ ] T008 Update `api/Makefile` comments to document integration test usage

**Checkpoint**: Developers can now run `make test-integration` locally (requires MinIO via docker-compose)

---

## Phase 3: User Story 1 - S3 Storage Provider Integration Tests (Priority: P1) 🎯 MVP

**Goal**: Verify S3 provider correctly integrates with MinIO for pre-signed URL generation and object operations

**Independent Test**: Run `make test-integration` with MinIO running via `docker compose up -d minio`

### Implementation for User Story 1

- [ ] T009 [US1] Create integration test file with build tag in `api/internal/storage/s3/provider_integration_test.go`
- [ ] T010 [US1] Implement `TestIntegration_GenerateUploadURL_Success` - verify pre-signed upload URL generation
- [ ] T011 [US1] Implement `TestIntegration_GenerateUploadURL_ActualUpload` - verify actual file upload via pre-signed URL
- [ ] T012 [US1] Implement `TestIntegration_GenerateDownloadURL_Success` - verify pre-signed download URL generation
- [ ] T013 [US1] Implement `TestIntegration_GenerateDownloadURL_ActualDownload` - verify actual file download via pre-signed URL
- [ ] T014 [US1] Implement `TestIntegration_DeleteObject_Success` - verify object deletion works correctly
- [ ] T015 [US1] Implement `TestIntegration_DeleteObject_NonExistent` - verify deletion of non-existent object is idempotent
- [ ] T016 [US1] Add test cleanup in `t.Cleanup()` to remove uploaded objects after each test

**Checkpoint**: S3 provider integration tests pass with MinIO. Run `cd api && make test-integration`

---

## Phase 4: User Story 2 - Video Handler Integration Tests (Priority: P2)

**Goal**: Verify VideoHandler correctly integrates with S3 provider for end-to-end API operations

**Independent Test**: Run integration tests for handler package with MinIO running

### Implementation for User Story 2

- [ ] T017 [US2] Create integration test file with build tag in `api/internal/handler/video_integration_test.go`
- [ ] T018 [US2] Implement test helper `newTestVideoHandler(t *testing.T)` that creates handler with real MinIO provider
- [ ] T019 [US2] Implement `TestIntegration_GenerateVideoUploadURL_ValidRequest` - verify API returns valid pre-signed URL
- [ ] T020 [US2] Implement `TestIntegration_GenerateVideoUploadURL_InvalidContentType` - verify content type validation
- [ ] T021 [US2] Implement `TestIntegration_GenerateVideoDownloadURL_ValidRequest` - verify API returns valid pre-signed URL
- [ ] T022 [US2] Implement `TestIntegration_GenerateVideoDownloadURL_InvalidObjectKey` - verify object key validation
- [ ] T023 [US2] Add test cleanup to remove any uploaded objects after handler tests

**Checkpoint**: Video handler integration tests pass. Both US1 and US2 tests pass independently.

---

## Phase 5: User Story 3 - CI/CD Pipeline Integration Tests (Priority: P3)

**Goal**: Enable automatic execution of integration tests in GitHub Actions on all PRs

**Independent Test**: Create a PR and verify integration tests run in CI

### Implementation for User Story 3

- [ ] T024 [US3] Add MinIO service container configuration to `.github/workflows/api-ci.yml`
- [ ] T025 [US3] Add MinIO environment variables (MINIO_ROOT_USER, MINIO_ROOT_PASSWORD) to CI workflow
- [ ] T026 [US3] Add health check wait step for MinIO service in CI workflow
- [ ] T027 [US3] Add integration test step `make test-integration` to CI workflow after unit tests
- [ ] T028 [US3] Update CI workflow paths trigger to include new integration test files

**Checkpoint**: CI pipeline runs integration tests on PRs. All user stories complete.

---

## Phase 6: Polish & Documentation

**Purpose**: Final improvements and documentation

- [ ] T029 [P] Update `specs/004-minio-test-coverage/quickstart.md` with actual test execution commands
- [ ] T030 [P] Add integration test documentation to `api/README.md`
- [ ] T031 Verify all integration tests pass locally with `cd api && make test-integration`
- [ ] T032 Verify integration tests are skipped gracefully when MinIO is unavailable
- [ ] T033 Run `make check` to ensure all linting and existing tests still pass

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion
- **User Story 1 (Phase 3)**: Depends on Phase 1 and Phase 2
- **User Story 2 (Phase 4)**: Depends on Phase 1 and Phase 2 (can run parallel with US1)
- **User Story 3 (Phase 5)**: Depends on Phase 2 (can run parallel with US1/US2)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Phase 2 - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Phase 2 - Uses same test helpers, independently testable
- **User Story 3 (P3)**: Can start after Phase 2 - CI configuration is independent of test implementation

### Within Each User Story

- Test helper functions must exist before integration tests
- Each test should be independently runnable
- Cleanup should be registered in `t.Cleanup()` for reliability

### Parallel Opportunities

- T001-T005 (Setup): T003, T004 can run in parallel after T002
- T006-T008 (Foundational): All can run in parallel
- US1, US2, US3 can be worked on in parallel after Foundational phase
- T029, T030 (Polish): Can run in parallel

---

## Parallel Example: User Story 1

```bash
# After test helpers are ready (Phase 1), these tests can be written in parallel:
Task T010: "TestIntegration_GenerateUploadURL_Success"
Task T012: "TestIntegration_GenerateDownloadURL_Success"
Task T014: "TestIntegration_DeleteObject_Success"

# Then add actual operation tests:
Task T011: "TestIntegration_GenerateUploadURL_ActualUpload"
Task T013: "TestIntegration_GenerateDownloadURL_ActualDownload"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (test helpers)
2. Complete Phase 2: Foundational (Makefile targets)
3. Complete Phase 3: User Story 1 (S3 provider integration tests)
4. **STOP and VALIDATE**: Run `make test-integration` with MinIO
5. This MVP already provides value: developers can verify S3 provider works

### Incremental Delivery

1. Setup + Foundational → Test infrastructure ready
2. Add User Story 1 → S3 provider tests → Validate locally
3. Add User Story 2 → Handler tests → Validate locally
4. Add User Story 3 → CI integration → Validate via PR
5. Each story adds confidence without breaking previous functionality

### Suggested Commit Points

1. After Phase 1: "test(storage): add integration test helpers"
2. After Phase 2: "build(makefile): add integration test targets"
3. After US1: "test(storage): add S3 provider integration tests with MinIO"
4. After US2: "test(handler): add video handler integration tests"
5. After US3: "ci: add MinIO service for integration tests"
6. After Polish: "docs: update integration test documentation"

---

## Notes

- All integration test files MUST have `//go:build integration` as first line
- Tests use `t.Skip()` when MinIO is unavailable - never fail due to missing infrastructure
- Fixed bucket name `optel-test-videos` is used across all tests
- Test cleanup uses `t.Cleanup()` to ensure objects are removed even on test failure
- MinIO credentials default to `minioadmin/minioadmin` matching docker-compose.yml
