# AGENTS.md

This file provides guidance for AI coding agents working on this repository.

## Project Overview

**Rx** - A plan-driven training management system.

This repository (`rx`) manages Programs, Plans, and Logs with a plan-first philosophy.

## Key Principles

1. **"Dumb Backend"** - No business logic for "health." Strictly stores and retrieves data.
2. **Domain-Driven Schema-First** - Domain models define business logic, OpenAPI spec defines API contract. Code is generated from OpenAPI spec.

For details, see [docs/PHILOSOPHY.md](docs/PHILOSOPHY.md).

## Commands

### API (run from `api/`)

```bash
task generate   # OpenAPI spec → Go code (run after editing openapi/openapi.yaml)
task check      # generate + format + lint + test (run before committing)
task test       # unit tests with race detection
task lint       # golangci-lint
task format     # gofmt check
task run        # dev server (localhost:8080)
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
- **pre-push**: runs `task generate` and verifies no stale generated code when `api/openapi/openapi.yaml` is changed

**Never use `git commit --no-verify`.**

## Architecture

### API Layer Flow

```
openapi/openapi.yaml  →  task generate  →  pkg/openapi/server.gen.go (do not edit)
                                                    ↓
                                          internal/handler/      ← implements generated interface
                                                    ↓
                                          internal/domain/       ← validates business rules
                                                    ↓
                                          internal/repository/   ← interfaces (ports)
                                                    ↓
                                          internal/store/{memory,postgres}/  ← implementations
```

**Key constraint:** Handlers convert between OpenAPI types and domain types manually — there is no auto-mapping. Domain entities go through `Validate*()` before being passed to the repository.

### Error Flow

Domain errors (`internal/domain/errors.go`) → handler catches with type assertion → middleware helpers write HTTP response:

| Domain error | HTTP status |
|---|---|
| `ValidationError` | 400 |
| `DomainError{Code: NOT_FOUND}` | 404 |
| `DomainError{Code: CONFLICT}` | 409 |

### Auth

Middleware-based (`internal/middleware/auth.go`). `AuthProvider` interface with `StubProvider` for dev (accepts any Bearer token, uses value as userID). UserID is stored in request context; retrieve with `middleware.GetUserID(ctx)`.

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
