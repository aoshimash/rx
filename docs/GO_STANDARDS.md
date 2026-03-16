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
├── cmd/server/          # Entry point
├── internal/
│   ├── config/          # Configuration
│   ├── domain/          # Entities + business rules (validation, invariants)
│   ├── handler/         # HTTP handlers (OpenAPI → domain → response)
│   ├── repository/      # Repository interfaces (ports)
│   └── store/memory/    # Repository implementations (adapters)
├── pkg/openapi/         # Generated from OpenAPI spec (do not edit)
└── openapi/openapi.yaml # API specification (source of truth)
```

## Key Principles

1. **Schema-First** — Write OpenAPI spec first, then generate code. Never edit generated files.
2. **Explicit Errors** — No panics. Return errors explicitly using domain error types.
3. **Table-Driven Tests** — Mandatory for all test cases. No ad-hoc single test functions.
4. **Clear Comments** — Comments explain *what* and *why*, not just *how*.
