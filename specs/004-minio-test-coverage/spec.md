# Feature Specification: MinIO Integration Test Coverage

**Feature Branch**: `004-minio-test-coverage`  
**Created**: 2026-01-30  
**Status**: Draft  
**Input**: User description: "Introduce MinIO to expand test coverage"

## Clarifications

### Session 2026-01-30

- Q: How should integration tests behave when MinIO is unavailable? → A: Skip tests with `t.Skip` and continue running other tests
- Q: How should integration tests be separated from unit tests? → A: Use build tags (`//go:build integration`)
- Q: When should integration tests run in CI environment? → A: Run on all PRs
- Q: How should the test bucket name be determined? → A: Use a fixed name (`optel-test-videos`)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - S3 Storage Provider Integration Tests (Priority: P1)

As a developer, I want to verify that the S3 storage provider correctly integrates with actual S3-compatible storage (MinIO). This ensures reliability of storage operations in production.

**Why this priority**: Storage operations are the foundation for video upload/download functionality. If these don't work correctly, users cannot use video features.

**Independent Test**: Start MinIO container and verify that object upload/download/delete operations work correctly through the S3 provider.

**Acceptance Scenarios**:

1. **Given** MinIO is running, **When** a pre-signed upload URL is generated and a video file is uploaded to that URL, **Then** the file is stored in MinIO
2. **Given** a file exists in MinIO, **When** a pre-signed download URL is generated and accessed, **Then** the file can be downloaded
3. **Given** a file exists in MinIO, **When** the object is deleted, **Then** the file is removed from MinIO

---

### User Story 2 - Video Handler Integration Tests (Priority: P2)

As a developer, I want to verify that VideoHandler correctly integrates with the S3 storage provider, ensuring end-to-end operation from API requests to storage operations.

**Why this priority**: Handler-level integration tests verify correct coupling between API layer and storage layer, enabling early detection of interface mismatches.

**Independent Test**: Using test MinIO and VideoHandler, verify that upload URL generation API and download URL generation API correctly integrate with the storage provider.

**Acceptance Scenarios**:

1. **Given** MinIO is running and VideoHandler is configured, **When** the upload URL generation API is called, **Then** a valid pre-signed URL is returned
2. **Given** a file exists in MinIO, **When** the download URL generation API is called, **Then** a valid pre-signed URL is returned

---

### User Story 3 - CI/CD Pipeline Integration Tests (Priority: P3)

As a developer, I want integration tests using MinIO to run automatically in CI/CD environment. This enables continuous verification that code changes don't break storage functionality.

**Why this priority**: Automated tests enable early detection of regressions and reduce manual testing burden.

**Independent Test**: Start MinIO service in GitHub Actions CI pipeline and verify that integration tests run automatically.

**Acceptance Scenarios**:

1. **Given** the CI pipeline is executed, **When** integration tests run, **Then** tests using MinIO succeed
2. **Given** any code changes exist, **When** a PR is created, **Then** integration tests are automatically executed on all PRs

---

### Edge Cases

- When MinIO is not running, tests are skipped with `t.Skip`
- Is error handling appropriate when network timeout occurs?
- Is error handling appropriate when bucket doesn't exist?
- Is an appropriate error returned when pre-signed URL expires?
- Do conflicts occur when multiple uploads happen simultaneously?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System can execute integration tests using MinIO as S3-compatible storage
- **FR-002**: System verifies that pre-signed upload URLs can be generated and files can actually be uploaded
- **FR-003**: System verifies that pre-signed download URLs can be generated and files can actually be downloaded
- **FR-004**: System verifies that object deletion operations work correctly
- **FR-005**: System automatically creates the fixed MinIO bucket (`optel-test-videos`) before test execution
- **FR-006**: System cleans up test data after test execution
- **FR-007**: Integration tests are skipped with `t.Skip` when MinIO is unavailable, and other tests continue running
- **FR-008**: CI/CD pipeline executes integration tests using MinIO service on all PRs

### Key Entities

- **Storage Provider**: Interface providing access to S3-compatible storage. Encapsulates operations like pre-signed URL generation and object deletion.
- **Test Bucket**: Storage bucket dedicated to integration tests (fixed name: `optel-test-videos`). Created/verified before test execution, objects cleaned up after tests.
- **Test Object**: Temporary file object used in tests. Used for verifying upload/download/delete operations.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All storage provider methods (GenerateUploadURL, GenerateDownloadURL, DeleteObject) are covered by integration tests
- **SC-002**: Integration tests complete within 10 seconds in local development environment
- **SC-003**: Integration tests run automatically in CI pipeline and can detect storage functionality regressions
- **SC-004**: Error messages on test failure contain sufficient information to identify the root cause
- **SC-005**: Test coverage achieves 80% or higher for storage-related code

## Assumptions

- MinIO service is already defined in docker-compose.yml (verified)
- Developers can use Docker in local environment
- Docker services can be used in CI/CD environment (GitHub Actions)
- Integration tests can be separated from unit tests (controlled by build tag `//go:build integration`)
- MinIO is compatible with AWS S3 API and works with existing S3 provider implementation
