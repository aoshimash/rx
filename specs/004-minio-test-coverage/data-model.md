# Data Model: MinIO Integration Test Coverage

**Feature**: 004-minio-test-coverage  
**Date**: 2026-01-30

## Overview

This feature adds integration tests only - no domain model changes. This document describes the test infrastructure data structures.

## Test Infrastructure Entities

### TestMinIOConfig

Configuration for MinIO integration test environment.

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| Endpoint | string | MinIO server endpoint URL | `http://localhost:9000` |
| AccessKey | string | MinIO access key (from env) | `minioadmin` |
| SecretKey | string | MinIO secret key (from env) | `minioadmin` |
| Bucket | string | Test bucket name | `optel-test-videos` |
| Region | string | AWS region (required by SDK) | `us-east-1` |

**Environment Variables**:
- `MINIO_ENDPOINT` - Override endpoint (optional)
- `MINIO_ROOT_USER` - Access key
- `MINIO_ROOT_PASSWORD` - Secret key

### TestObject

Represents a test file uploaded during integration tests.

| Field | Type | Description |
|-------|------|-------------|
| ObjectKey | string | S3 object key (e.g., `videos/user123/uuid.mp4`) |
| ContentType | string | MIME type (e.g., `video/mp4`) |
| Content | []byte | File content for upload/verification |
| ContentLength | int64 | Size in bytes |

## State Transitions

### Test Bucket Lifecycle

```
[Not Exists] --CreateBucket--> [Exists/Empty] --UploadObject--> [Has Objects]
                                     ^                              |
                                     |                              |
                                     +------CleanupObjects----------+
```

### Integration Test Flow

```
[Test Start]
    |
    v
[Check MinIO Availability] --unavailable--> [t.Skip]
    |
    v (available)
[Ensure Bucket Exists]
    |
    v
[Run Test Cases]
    |
    v
[Cleanup Objects]
    |
    v
[Test End]
```

## Validation Rules

### Object Key Format (from existing provider)

Pattern: `videos/{user_id}/{uuid}.{ext}`
- `user_id`: alphanumeric, underscore, hyphen
- `uuid`: standard UUID format
- `ext`: file extension (mp4, mov, etc.)

### Content Type Validation (from existing provider)

- Must start with `video/`
- Examples: `video/mp4`, `video/quicktime`, `video/webm`

## No Domain Model Changes

This feature does not modify:
- `internal/domain/` entities
- `api/openapi/openapi.yaml` specification
- Database schemas
- API contracts

All changes are confined to test infrastructure.
