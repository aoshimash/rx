# Quickstart: MinIO Integration Test Coverage

**Feature**: 004-minio-test-coverage  
**Date**: 2026-01-30

## Prerequisites

- Docker and Docker Compose installed
- Go 1.25+ installed (or use DevContainer)
- Repository cloned locally

## Local Development

### 1. Start MinIO

```bash
# From repository root
docker compose up -d minio

# Verify MinIO is running
docker compose ps minio
# Should show "healthy" status

# MinIO Console available at http://localhost:9001
# Login: minioadmin / minioadmin
```

### 2. Run Integration Tests

```bash
cd api

# Run only integration tests
go test -v -tags=integration ./internal/storage/s3/...

# Run all tests (unit + integration)
go test -v -tags=integration ./...

# Or use Makefile target (after implementation)
make test-integration
```

### 3. Run Without MinIO

If MinIO is not running, integration tests will be skipped:

```bash
# Stop MinIO
docker compose down minio

# Run tests - integration tests will skip
go test -v -tags=integration ./internal/storage/s3/...
# Output: "=== SKIP: TestXxx: MinIO unavailable..."
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MINIO_ENDPOINT` | `http://localhost:9000` | MinIO server endpoint |
| `MINIO_ROOT_USER` | `minioadmin` | MinIO access key |
| `MINIO_ROOT_PASSWORD` | `minioadmin` | MinIO secret key |

## CI Pipeline

Integration tests run automatically on all PRs:

```yaml
# .github/workflows/api-ci.yml additions
services:
  minio:
    image: minio/minio:latest
    ports:
      - 9000:9000
    env:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin

steps:
  - name: Run integration tests
    run: go test -v -tags=integration ./...
```

## Test Coverage

After implementation, verify coverage:

```bash
cd api

# Run with coverage
go test -v -tags=integration -coverprofile=coverage-integration.out ./internal/storage/s3/...

# View coverage report
go tool cover -html=coverage-integration.out
```

## Troubleshooting

### MinIO Connection Refused

```bash
# Check if MinIO is running
docker compose ps

# Check MinIO logs
docker compose logs minio

# Restart MinIO
docker compose restart minio
```

### Tests Skipped Unexpectedly

```bash
# Verify environment variables
echo $MINIO_ENDPOINT
echo $MINIO_ROOT_USER

# Test MinIO connectivity directly
curl http://localhost:9000/minio/health/live
```

### Bucket Creation Errors

MinIO console at `http://localhost:9001` can be used to manually inspect/create buckets if needed.

## Files to Implement

1. `api/internal/storage/s3/provider_integration_test.go` - S3 provider integration tests
2. `api/internal/handler/video_integration_test.go` - Video handler integration tests
3. `api/Makefile` - Add `test-integration` target
4. `.github/workflows/api-ci.yml` - Add MinIO service and integration test step
