# Rx Go Standards - Detailed Reference

## Project Structure Details

### cmd/

Entry points for executables.

```go
// cmd/server/main.go
package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

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

Domain entities and business logic.

```go
// internal/domain/workout.go
package domain

import "time"

// Workout represents a completed unit of physical exertion.
// Contains telemetry data: intensity (RPE), volume (load), duration, and timestamp.
type Workout struct {
    ID          string
    Timestamp   time.Time
    Duration    int       // seconds
    Intensity   int       // RPE 1-10
    Volume      float64   // kg
    MuscleGroups []string
    Notes       string
    ProgramID   *string
    CreatedAt   time.Time
}
```

### internal/repository/

Repository interfaces (ports).

```go
// internal/repository/workout.go
package repository

import (
    "context"
    "time"
    "github.com/aoshimash/rx/api/internal/domain"
)

type WorkoutRepository interface {
    Create(ctx context.Context, w *domain.Workout) error
    GetByID(ctx context.Context, id string) (*domain.Workout, error)
    List(ctx context.Context, filter WorkoutFilter) ([]*domain.Workout, error)
    Delete(ctx context.Context, id string) error
}

type WorkoutFilter struct {
    From         *time.Time
    To           *time.Time
    MuscleGroups []string
    Limit        int
    Offset       int
}
```

### internal/store/

Repository implementations (adapters).

```go
// internal/store/memory/workout.go
package memory

import (
    "context"
    "sync"

    "github.com/aoshimash/rx/api/internal/domain"
    "github.com/aoshimash/rx/api/internal/repository"
)

type WorkoutStore struct {
    mu       sync.RWMutex
    workouts map[string]*domain.Workout
}

func NewWorkoutStore() *WorkoutStore {
    return &WorkoutStore{
        workouts: make(map[string]*domain.Workout),
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
// internal/domain/workout_test.go
package domain_test

import (
    "testing"
    "time"

    "github.com/aoshimash/rx/api/internal/domain"
)

func TestWorkout_Validate(t *testing.T) {
    tests := []struct {
        name    string
        input   domain.Workout
        wantErr bool
        errCode domain.ErrorCode
    }{
        {
            name: "valid workout",
            input: domain.Workout{
                Timestamp: time.Now().Add(-1 * time.Hour),
                Duration:  3600,
                Intensity: 7,
                Volume:    1000,
            },
            wantErr: false,
        },
        {
            name: "invalid intensity - too high",
            input: domain.Workout{
                Timestamp: time.Now().Add(-1 * time.Hour),
                Duration:  3600,
                Intensity: 11,
                Volume:    1000,
            },
            wantErr: true,
            errCode: domain.ErrInvalidInput,
        },
        {
            name: "invalid intensity - too low",
            input: domain.Workout{
                Timestamp: time.Now().Add(-1 * time.Hour),
                Duration:  3600,
                Intensity: 0,
                Volume:    1000,
            },
            wantErr: true,
            errCode: domain.ErrInvalidInput,
        },
        {
            name: "future timestamp",
            input: domain.Workout{
                Timestamp: time.Now().Add(1 * time.Hour),
                Duration:  3600,
                Intensity: 7,
                Volume:    1000,
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
    slog.Error("failed to process workout", 
        "workout_id", "wo-123",
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
