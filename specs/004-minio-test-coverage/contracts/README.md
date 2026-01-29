# Contracts: MinIO Integration Test Coverage

**Feature**: 004-minio-test-coverage  
**Date**: 2026-01-30

## No API Contract Changes

This feature adds integration tests only. No changes to:

- OpenAPI specification (`api/openapi/openapi.yaml`)
- API endpoints
- Request/response schemas
- Authentication mechanisms

## Tested Contracts

The integration tests verify the existing Storage Provider interface contract:

### storage.Provider Interface

```go
type Provider interface {
    // GenerateUploadURL creates a pre-signed URL for uploading a video file
    GenerateUploadURL(ctx context.Context, req UploadURLRequest) (*UploadURLResponse, error)

    // GenerateDownloadURL creates a pre-signed URL for downloading a video file
    GenerateDownloadURL(ctx context.Context, req DownloadURLRequest) (*DownloadURLResponse, error)

    // DeleteObject removes a file from storage
    DeleteObject(ctx context.Context, objectKey string) error

    // ValidateObjectKey checks if an object key has valid format
    ValidateObjectKey(objectKey string) bool
}
```

### Pre-signed URL Contract

Integration tests verify that pre-signed URLs:
1. Are valid and accessible
2. Support actual file upload/download operations
3. Respect content type constraints
4. Expire as configured

## CI Contract

GitHub Actions workflow must:
1. Start MinIO service container before tests
2. Wait for MinIO health check to pass
3. Run integration tests with `-tags=integration`
4. Report test results
