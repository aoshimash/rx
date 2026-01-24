---
name: optel-go-standards
description: Go coding standards for OPTel project. Covers project structure, error handling, testing, logging, and linting configuration. Use when writing Go code, creating new files in api/, reviewing PRs, setting up packages, or configuring golangci-lint.
---

# OPTel Go Standards

## Technology Stack

| Component | Choice |
|-----------|--------|
| Go Version | 1.25+ |
| HTTP Server | chi |
| OpenAPI | oapi-codegen (Schema-First) |
| Testing | standard testing |
| Logging | log/slog |
| Linter | golangci-lint (strict) |

## Project Structure

```
api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── domain/
│   ├── handler/
│   ├── repository/
│   └── store/
│       └── memory/
├── pkg/
│   └── openapi/          # Generated from OpenAPI spec
├── openapi/
│   └── openapi.yaml
├── go.mod
└── go.sum
```

## Key Principles

1. **Schema-First** - Write OpenAPI spec first, then generate code
2. **Explicit Errors** - No panics. Return errors explicitly.
3. **Table-Driven Tests** - Mandatory for all test cases
4. **Infrastructure Comments** - SRE perspective, not fitness terminology

## Error Handling

Use custom domain errors:

```go
type DomainError struct {
    Code    ErrorCode
    Message string
}

const (
    ErrNotFound     ErrorCode = "NOT_FOUND"
    ErrInvalidInput ErrorCode = "INVALID_INPUT"
    ErrConflict     ErrorCode = "CONFLICT"
)
```

## Additional Resources

For detailed guidelines, see [reference.md](reference.md).
