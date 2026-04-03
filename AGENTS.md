# AGENTS.md

This file provides guidance for AI coding agents working on this repository.

## Project Overview

**Rx** - A plan-driven training management system.

This repository (`rx`) manages Programs, Plans, and Logs with a plan-first philosophy.

## Key Principles

1. **"Dumb Backend"** - No business logic for "health." Strictly stores and retrieves data.
2. **Domain-Driven Schema-First** - Domain models define business logic, protobuf defines API contract. Code is generated from proto files via buf.

For details, see [docs/PHILOSOPHY.md](docs/PHILOSOPHY.md).

## Commands

### API (run from `api/`)

```bash
task generate   # proto → Go code (runs cd ../proto && buf generate)
task check      # generate + format + lint + test (run before committing)
task test       # unit tests with race detection
task lint       # golangci-lint
task format     # gofmt check
task run        # dev server (gRPC :9090 + HTTP :8080)
task migrate    # run DB migrations
```

Single test: `go test ./internal/domain/... -run TestRoundToIncrement -v`
Integration tests: `go test -tags=integration ./internal/store/postgres/... -v`

### Web (run from `web/`)

```bash
pnpm dev        # dev server (localhost:3000)
pnpm check      # biome lint + format check (run before committing)
pnpm check:fix  # fix biome issues
pnpm build      # production build
```

### Git Hooks (lefthook)

Hooks are managed by [lefthook](https://github.com/evilmartians/lefthook). Setup once after cloning:

```bash
aqua install    # installs lefthook and all other tools
lefthook install  # registers git hooks
```

If upgrading from the old `githooks/` setup: `git config --unset core.hooksPath` first.

- **pre-commit**: runs `task format`, `task lint`, `task test` for API changes; `pnpm check` for web changes
- **pre-push**: runs `buf lint`, `buf generate` (verifies no stale generated code), and `buf breaking` when `proto/**/*.proto` is changed

**Never use `git commit --no-verify`.**

## Architecture

### API Layer Flow

```
proto/rx/api/v1/*.proto  →  buf generate  →  pkg/gen/rx/api/v1/ (do not edit)
                                                      ↓
                                            internal/handler/      ← gRPC server implementations
                                                      ↓
                                            internal/domain/       ← validates business rules
                                                      ↓
                                            internal/repository/   ← interfaces (ports)
                                                      ↓
                                            internal/store/{memory,postgres}/  ← implementations
```

**Dual server:** gRPC server (port 9090) + gRPC-Gateway HTTP server (port 8080). REST paths are defined by `google.api.http` annotations in proto files.

**Key constraint:** Handlers convert between proto types and domain types manually via `internal/handler/convert.go` — there is no auto-mapping. Domain entities go through `Validate*()` before being passed to the repository.

### Error Flow

Domain errors (`internal/domain/errors.go`) → handler converts via `domainErrToStatus()` (`internal/handler/convert.go`) → gRPC status codes (mapped to HTTP by grpc-gateway):

| Domain error | gRPC code | HTTP status |
|---|---|---|
| `ValidationError` | `InvalidArgument` | 400 |
| `DomainError{Code: NOT_FOUND}` | `NotFound` | 404 |
| `DomainError{Code: CONFLICT}` | `AlreadyExists` | 409 |

### Auth

gRPC interceptor-based (`internal/middleware/grpc_auth.go`). `GRPCAuthProvider` interface with `GRPCStubProvider` for dev (accepts any Bearer token, uses value as userID). `UnaryAuthInterceptor` extracts Bearer token from gRPC metadata (HealthService is exempt). UserID is stored in context; retrieve with `middleware.GetUserID(ctx)` (`internal/middleware/context.go`).

### Pagination

All list endpoints use cursor-based pagination: `limit` (1–100) + `after` (base64-encoded cursor). Most endpoints use a UUID cursor; Plan endpoints use a `created_at|UUID` cursor to support `created_at` ordering. Repository `List()` returns `(entities, nextCursor string, hasMore bool, error)`.

### Storage

`STORAGE_TYPE=memory` (default) or `postgres`. Videos use pre-signed URLs — the server never handles file bytes, only issues upload/download URLs to S3/R2.

## Documentation

- [Philosophy](docs/PHILOSOPHY.md) - Core principles and constraints
- [Domain Model](docs/DOMAIN_MODEL.md) - Program / Plan / Log lifecycle
- [Go Standards](docs/GO_STANDARDS.md) - Go coding standards
- [Frontend Architecture](docs/FRONTEND_ARCHITECTURE.md) - Web/Mobile architecture and standards
- [Architecture](docs/ARCHITECTURE.md) - System architecture
- [Development Guide](docs/DEVELOPMENT.md) - Development setup

### Documentation Sync

A PostToolUse hook (`.claude/hooks/doc-sync-check.sh`) automatically reminds you when code changes may require documentation updates. The mapping:

| Code change | Check this document |
|---|---|
| `api/internal/domain/` (add/rename/remove entities or fields) | `docs/DOMAIN_MODEL.md` |
| `proto/rx/api/v1/*.proto` (schema changes) | `docs/DOMAIN_MODEL.md` |
| `web/app/**/page.tsx` (add/remove/rename routes) | `docs/WEB_UI_DESIGN.md` |
| `api/internal/handler/` (error handling changes) | AGENTS.md error flow table |
| `docs/PHILOSOPHY.md` terminology table | Verify against `api/internal/domain/` types |

Only structural changes (new entities, renamed fields, new routes) require doc updates — bug fixes and refactors within the same structure generally do not.
