# Research: MinIO Integration Test Coverage

**Feature**: 004-minio-test-coverage  
**Date**: 2026-01-30

## Research Tasks

### 1. Go Integration Test Best Practices with Build Tags

**Decision**: Use `//go:build integration` build tag to separate integration tests from unit tests.

**Rationale**:
- Standard Go convention for test separation
- `go test ./...` runs only unit tests by default
- `go test -tags=integration ./...` explicitly runs integration tests
- Clear intent in code and CI configuration

**Alternatives Considered**:
- Environment variables (`INTEGRATION_TEST=true`) - Rejected: Less idiomatic for Go, requires runtime checks
- Separate test directories - Rejected: Moves tests away from code they verify, harder to maintain
- Test file naming convention only - Rejected: No built-in tooling support for selective execution

### 2. MinIO S3 Compatibility for Pre-signed URLs

**Decision**: Use existing S3 provider implementation with MinIO endpoint configuration.

**Rationale**:
- MinIO implements AWS S3 API including pre-signed URL operations
- Existing `api/internal/storage/s3/provider.go` already supports custom endpoints via `Config.Endpoint`
- No code changes needed in provider - only test configuration

**Alternatives Considered**:
- Mock S3 service - Rejected: Doesn't validate actual S3 API compatibility
- LocalStack - Rejected: Heavier than MinIO, more complex setup

### 3. GitHub Actions Service Container for MinIO

**Decision**: Use GitHub Actions service container with MinIO Docker image.

**Rationale**:
- Native GitHub Actions feature for running services alongside jobs
- MinIO official Docker image supports health checks
- Network accessible via `localhost:9000` in job steps
- docker-compose.yml configuration can be reused for consistency

**Implementation**:
```yaml
services:
  minio:
    image: minio/minio:latest
    ports:
      - 9000:9000
    env:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    options: >-
      --health-cmd "mc ready local"
      --health-interval 5s
      --health-timeout 5s
      --health-retries 5
```

**Alternatives Considered**:
- docker-compose in CI - Rejected: More complex workflow setup
- Dedicated MinIO instance - Rejected: Unnecessary infrastructure, cost

### 4. Test Bucket Management Strategy

**Decision**: Use fixed bucket name `optel-test-videos`, create if not exists, clean up objects after each test.

**Rationale**:
- Fixed name simplifies configuration (no dynamic bucket creation/deletion)
- Bucket creation is idempotent (safe to call multiple times)
- Object cleanup prevents test pollution while keeping bucket for next run
- Consistent with docker-compose.yml MinIO volume persistence

**Implementation Pattern**:
```go
func setupTestBucket(ctx context.Context, client *s3.Client, bucket string) error {
    _, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
        Bucket: aws.String(bucket),
    })
    if err != nil {
        var bne *types.BucketAlreadyExists
        var bno *types.BucketAlreadyOwnedByYou
        if errors.As(err, &bne) || errors.As(err, &bno) {
            return nil // Bucket already exists, OK
        }
        return err
    }
    return nil
}
```

**Alternatives Considered**:
- UUID bucket per test run - Rejected: Leaves orphaned buckets on failure, cleanup complexity
- Shared bucket without cleanup - Rejected: Tests may interfere with each other

### 5. MinIO Availability Detection for Test Skipping

**Decision**: Attempt MinIO connection in `TestMain` or test setup, skip tests with `t.Skip()` if unavailable.

**Rationale**:
- Graceful degradation when developer doesn't have Docker running
- Standard Go testing pattern
- Clear skip message helps developers understand why tests didn't run

**Implementation Pattern**:
```go
func skipIfMinIOUnavailable(t *testing.T) {
    t.Helper()
    endpoint := os.Getenv("MINIO_ENDPOINT")
    if endpoint == "" {
        endpoint = "http://localhost:9000"
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    // Try to connect to MinIO
    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion("us-east-1"),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            os.Getenv("MINIO_ROOT_USER"),
            os.Getenv("MINIO_ROOT_PASSWORD"),
            "",
        )),
    )
    if err != nil {
        t.Skipf("MinIO unavailable: %v", err)
    }
    
    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.BaseEndpoint = aws.String(endpoint)
        o.UsePathStyle = true
    })
    
    _, err = client.ListBuckets(ctx, &s3.ListBucketsInput{})
    if err != nil {
        t.Skipf("MinIO unavailable: %v", err)
    }
}
```

**Alternatives Considered**:
- Fail tests immediately - Rejected: Blocks development without Docker
- Global skip flag - Rejected: Less granular control

## Summary

All unknowns resolved. Key decisions:
1. Build tag `//go:build integration` for test separation
2. Existing S3 provider works with MinIO (no code changes)
3. GitHub Actions service container for CI
4. Fixed bucket `optel-test-videos` with object cleanup
5. Connection check with `t.Skip()` for graceful degradation
