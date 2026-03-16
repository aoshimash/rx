# Go Coding Standards

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

1. **Schema-First** — Write OpenAPI spec first, then generate code
2. **Explicit Errors** — No panics. Return errors explicitly.
3. **Table-Driven Tests** — Mandatory for all test cases
4. **Clear Comments** — Write clear, descriptive comments that explain what the code does

## Project Structure Details

### cmd/

Entry points for executables.

```go
// cmd/server/main.go
package main

import (
    "log/slog"
    "os"

    "github.com/aoshimash/rx/api/internal/config"
    "github.com/aoshimash/rx/api/internal/handler"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)

    cfg, err := config.Load()
    if err != nil {
        slog.Error("failed to load config", "error", err)
        os.Exit(1)
    }

    // ... setup and run server
}
```

### internal/domain/

Domain entities and business logic. See [DOMAIN_MODEL.md](DOMAIN_MODEL.md) for the full schema.

```go
// internal/domain/program.go
package domain

// Program represents a reusable, RPE-based training template.
// It contains no dates and no absolute weights.
type Program struct {
    ID      uuid.UUID
    Name    string
    Entries []ProgramEntry
    // ...
}
```

### internal/repository/

Repository interfaces (ports).

```go
// internal/repository/program.go
package repository

import (
    "context"

    "github.com/aoshimash/rx/api/internal/domain"
    "github.com/google/uuid"
)

type ProgramRepository interface {
    Create(ctx context.Context, p *domain.Program) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error)
    List(ctx context.Context) ([]*domain.Program, error)
    Update(ctx context.Context, p *domain.Program) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### internal/store/

Repository implementations (adapters).

```go
// internal/store/memory/program.go
package memory

import (
    "context"
    "sync"

    "github.com/aoshimash/rx/api/internal/domain"
    "github.com/google/uuid"
)

type ProgramStore struct {
    mu       sync.RWMutex
    programs map[uuid.UUID]*domain.Program
}

func NewProgramStore() *ProgramStore {
    return &ProgramStore{
        programs: make(map[uuid.UUID]*domain.Program),
    }
}
```

## Error Handling

### Domain Errors

```go
// internal/domain/errors.go
package domain

import "fmt"

type ErrorCode string

const (
    ErrNotFound     ErrorCode = "NOT_FOUND"
    ErrInvalidInput ErrorCode = "INVALID_INPUT"
    ErrConflict     ErrorCode = "CONFLICT"
    ErrInternal     ErrorCode = "INTERNAL"
)

type DomainError struct {
    Code    ErrorCode
    Message string
    Err     error
}

func (e *DomainError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
    return e.Err
}

// Helper constructors
func NewNotFoundError(msg string) *DomainError {
    return &DomainError{Code: ErrNotFound, Message: msg}
}

func NewInvalidInputError(msg string) *DomainError {
    return &DomainError{Code: ErrInvalidInput, Message: msg}
}
```

### HTTP Error Mapping

```go
// internal/handler/errors.go
package handler

import (
    "errors"
    "net/http"

    "github.com/aoshimash/rx/api/internal/domain"
)

func mapDomainErrorToHTTP(err error) int {
    var domainErr *domain.DomainError
    if errors.As(err, &domainErr) {
        switch domainErr.Code {
        case domain.ErrNotFound:
            return http.StatusNotFound
        case domain.ErrInvalidInput:
            return http.StatusBadRequest
        case domain.ErrConflict:
            return http.StatusConflict
        default:
            return http.StatusInternalServerError
        }
    }
    return http.StatusInternalServerError
}
```

## Testing

### Table-Driven Tests

```go
// internal/domain/plan_test.go
package domain_test

import (
    "testing"

    "github.com/aoshimash/rx/api/internal/domain"
)

func TestPlan_Validate(t *testing.T) {
    tests := []struct {
        name    string
        input   domain.Plan
        wantErr bool
        errCode domain.ErrorCode
    }{
        {
            name: "valid plan",
            input: domain.Plan{
                Name:    "Week 1",
                Entries: []domain.PlanEntry{},
            },
            wantErr: false,
        },
        {
            name: "missing name",
            input: domain.Plan{
                Name: "",
            },
            wantErr: true,
            errCode: domain.ErrInvalidInput,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.input.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
            if tt.wantErr && tt.errCode != "" {
                var domainErr *domain.DomainError
                if !errors.As(err, &domainErr) || domainErr.Code != tt.errCode {
                    t.Errorf("Validate() error code = %v, want %v", domainErr.Code, tt.errCode)
                }
            }
        })
    }
}
```

## Logging

### Structured Logging with slog

```go
package main

import (
    "log/slog"
    "os"
)

func main() {
    // JSON output for production
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)

    // Usage
    slog.Info("server started", "port", 8080)
    slog.Error("failed to process request",
        "plan_id", "plan-123",
        "error", err,
    )
}
```

## Linting

### golangci-lint Configuration

```yaml
# .golangci.yml
run:
  timeout: 5m
  go: "1.25"

linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - typecheck
    - gofmt
    - goimports
    - misspell
    - unconvert
    - unparam
    - nakedret
    - prealloc
    - exportloopref
    - gocritic
    - revive
    - gosec
    - bodyclose
    - noctx
    - errorlint
    - wrapcheck

linters-settings:
  govet:
    check-shadowing: true
  revive:
    rules:
      - name: blank-imports
      - name: context-as-argument
      - name: error-return
      - name: error-strings
      - name: exported
      - name: increment-decrement
      - name: var-declaration
  wrapcheck:
    ignoreSigs:
      - .Errorf(
      - errors.New(
      - errors.Join(

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - wrapcheck
        - gosec
```

## OpenAPI Code Generation

### oapi-codegen Usage

```bash
# Generate server code from OpenAPI spec
oapi-codegen -generate types,chi-server -package openapi \
  -o pkg/openapi/server.gen.go openapi/openapi.yaml
```

### Makefile Target

```makefile
.PHONY: generate
generate:
	oapi-codegen -generate types,chi-server -package openapi \
		-o pkg/openapi/server.gen.go openapi/openapi.yaml
```
