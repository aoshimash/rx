# Research: REST API for Domain Models

**Feature**: 002-rest-api
**Date**: 2026-01-24

## Research Topics

This document consolidates research findings for technical decisions in the REST API implementation.

---

## 1. Authentication Mechanism

**Question**: How should authentication be implemented for the API?

### Decision
Implement authentication as **application-level middleware** with a **pluggable provider pattern**. This allows the same docker image to work across different deployment environments (docker compose, Kubernetes, AWS) by switching authentication providers via configuration.

**Implementation Approach**:
- Define an `AuthProvider` interface for authentication logic
- Implement chi middleware that uses the configured provider
- Support multiple provider implementations:
  - **Stub Provider** (development): Simple header presence check
  - **JWT Provider** (production - home k8s): JWT token validation
  - **AWS Cognito Provider** (future - AWS): AWS Cognito integration
- Configuration via environment variables/ConfigMap:
  - `AUTH_PROVIDER=stub|jwt|cognito`
  - Provider-specific settings (JWT secret, Cognito region, etc.)

**Deployment Strategy**:
- **docker compose (local)**: Use stub provider (`AUTH_PROVIDER=stub`)
- **Kubernetes (home)**: Use JWT provider with ConfigMap/Secret
- **AWS (future)**: Use Cognito provider with IAM roles

### Rationale
- FR-028 requires authentication for all operations
- Application-level implementation works consistently across all deployment targets
- Pluggable providers enable environment-specific authentication without code changes
- Same docker image can be used in docker compose, Kubernetes, and AWS
- Avoids dependency on external components (reverse proxy, ingress controllers) that differ between environments

### Alternatives Considered
1. **Skip authentication entirely**: Rejected - violates FR-028
2. **Reverse proxy (Nginx/Traefik)**: Rejected - requires additional components in docker compose, different setup for k8s vs AWS
3. **Kubernetes Ingress authentication**: Rejected - doesn't work in docker compose, requires reconfiguration for AWS
4. **AWS API Gateway**: Rejected - doesn't work in docker compose or home k8s
5. **Full JWT implementation only**: Rejected - too complex for MVP, doesn't support easy local development

---

## 2. Update Semantics (PUT vs PATCH)

**Question**: Should updates use PUT (full replacement) or PATCH (partial update)?

### Decision
Use **PUT for full replacement** semantics:
- Client must send the complete resource
- Server replaces the entire record
- Simpler to implement and test

### Rationale
- Simpler implementation for MVP
- Avoids complexity of merge semantics
- Domain models are not deeply nested (except entries/nodes which are managed with parent)

### Alternatives Considered
1. **PATCH with JSON Merge Patch**: More flexible but adds complexity
2. **PATCH with JSON Patch (RFC 6902)**: Overkill for this use case
3. **Both PUT and PATCH**: Increases API surface without clear benefit at MVP stage

### Future Consideration
If partial updates become needed, PATCH with JSON Merge Patch (RFC 7396) can be added.

---

## 3. Cursor-Based Pagination

**Question**: How should cursor-based pagination be implemented?

### Decision
Use **opaque cursor** based on record ID:
- Query params: `limit` (default 100, max 100), `after` (cursor)
- Response includes: `data[]`, `next_cursor` (null if no more), `has_more`
- Cursor is base64-encoded last record ID

### Rationale
- Consistent results even with concurrent writes
- More efficient than offset for large datasets
- Simple to implement with UUID-based IDs

### Implementation Pattern
```go
type PaginatedResponse[T any] struct {
    Data       []T     `json:"data"`
    NextCursor *string `json:"next_cursor,omitempty"`
    HasMore    bool    `json:"has_more"`
}
```

### Alternatives Considered
1. **Offset-based**: Inconsistent with concurrent writes, less efficient at scale
2. **Keyset pagination with timestamp**: More complex, not needed for this use case

---

## 4. Concurrent Update Handling

**Question**: How should the system handle concurrent updates to the same record?

### Decision
Use **last-write-wins** for MVP:
- No optimistic locking
- Latest PUT overwrites previous state
- Simpler implementation

### Rationale
- Single-user MVP scenario
- Adding optimistic locking (ETag/If-Match) can be done later
- Reduces initial complexity

### Future Consideration
When multi-user support is needed:
- Add `version` field to entities
- Return ETag header with version
- Require `If-Match` header for updates
- Return 409 Conflict on version mismatch

---

## 5. Nested Resource Management

**Question**: How should WorkoutEntry and ProgramNode be managed?

### Decision
**Nested resources managed with parent**:

For WorkoutEntry:
- Created/updated as part of Workout (entries array in request body)
- No independent CRUD endpoints
- When Workout is deleted, entries are cascade deleted

For ProgramNode:
- Created/updated as part of Program (root_nodes array in request body)
- No independent CRUD endpoints
- When Program is deleted, nodes are cascade deleted

### Rationale
- Maintains data integrity (entry always belongs to workout)
- Simpler API surface
- Aligns with domain model (entries are embedded in Workout)

### API Pattern
```
# Workout with entries
POST   /api/v1/workouts           # Creates workout with entries
GET    /api/v1/workouts/{id}      # Returns workout with all entries
PUT    /api/v1/workouts/{id}      # Replaces workout including entries
DELETE /api/v1/workouts/{id}      # Deletes workout and all entries

# Program with nodes
POST   /api/v1/programs           # Creates program with nodes
GET    /api/v1/programs/{id}      # Returns program with full node tree
PUT    /api/v1/programs/{id}      # Replaces program including nodes
DELETE /api/v1/programs/{id}      # Deletes program and all nodes
```

---

## 6. Error Response Format

**Question**: What error response format should be used?

### Decision
Use consistent error response structure (already defined in OpenAPI):

```json
{
  "code": "NOT_FOUND",
  "message": "Workout not found",
  "details": {
    "id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### Error Codes and HTTP Status Mapping

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Missing or invalid authentication |
| `VALIDATION_ERROR` | 400 | Invalid input data |
| `NOT_FOUND` | 404 | Resource does not exist |
| `CONFLICT` | 409 | Referential integrity violation |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

### Rationale
- Consistent with existing Error schema in OpenAPI
- Machine-readable code for programmatic handling
- Human-readable message for debugging

---

## 7. Referential Integrity on Delete

**Question**: How should deletion handle references?

### Decision
**Reject deletion if referenced** (except parent-child):

| Entity | Action | Behavior |
|--------|--------|----------|
| Exercise | DELETE | 409 if any WorkoutEntry references it |
| Program | DELETE | 409 if any Workout references its node; cascade delete ProgramNodes |
| Workout | DELETE | Cascade delete WorkoutEntries |
| TelemetryPoint | DELETE | No dependencies, always allowed |

### Implementation
Before deleting Exercise:
1. Check if any WorkoutEntry has `exercise_id` matching
2. If found, return 409 with blocking references

### Rationale
- Prevents orphaned references
- Explicit cascade for parent-child (Workout→Entry, Program→Node)
- Aligns with FR-026

---

## 8. Payload Size Limits

**Question**: How should large payloads be handled?

### Decision
**Soft limits in validation**:
- Workout entries: max 500 entries per workout
- Program nodes: max 1000 nodes per program (total tree size)
- Request body: 10MB max (server configuration)

### Rationale
- Prevents abuse and memory issues
- Reasonable limits for practical use
- Can be adjusted based on real usage patterns

---

## Summary of Decisions

| Topic | Decision |
|-------|----------|
| Authentication | Middleware hook with stub, 401 for missing auth |
| Update Semantics | PUT (full replacement) |
| Pagination | Cursor-based with opaque base64 cursor |
| Concurrent Updates | Last-write-wins (no locking for MVP) |
| Nested Resources | Managed with parent, no independent endpoints |
| Error Format | Code + message + details object |
| Referential Integrity | Reject delete if referenced (409), cascade for parent-child |
| Payload Limits | 500 entries/workout, 1000 nodes/program |

All clarifications resolved. Proceeding to Phase 1: Design & Contracts.
