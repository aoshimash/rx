# Implementation Plan: PostgreSQL導入

**Branch**: `005-postgresql-setup` | **Date**: 2026-01-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-postgresql-setup/spec.md`

## Summary

OPTel Training APIのデータ永続化層としてPostgreSQL 17を導入する。既存のインメモリストレージからPostgreSQLへの移行を実現し、Testcontainersを使用した統合テスト環境を構築する。

## Technical Context

**Language/Version**: Go 1.25+  
**Primary Dependencies**: 
- pgx/v5 (PostgreSQLドライバー)
- golang-migrate/migrate (マイグレーション)
- testcontainers-go (テスト)
- cenkalti/backoff/v4 (リトライ)

**Storage**: PostgreSQL 17.x  
**Testing**: go test + Testcontainers  
**Target Platform**: Linux server (Docker container)  
**Project Type**: Monorepo - api/ component  
**Performance Goals**: 接続確立10秒以内、ヘルスチェック100%成功  
**Constraints**: コネクションプール最大10接続（開発環境）  
**Scale/Scope**: 開発環境向け、本番設定はスコープ外

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with OPTel Workout Constitution principles:

- **Dumb Backend**: ✅ データストレージの実装変更のみ、ビジネスロジックなし
- **Domain-Driven Schema-First**: ✅ ドメインモデルは変更なし、OpenAPI仕様も変更なし
- **Terminology**: ✅ 既存の用語を継続使用
- **Clean Architecture**: ✅ リポジトリパターン維持、PostgreSQL実装は新しいアダプター
- **Monorepo Structure**: ✅ api/コンポーネント内の変更、他コンポーネントに影響なし

## Project Structure

### Documentation (this feature)

```text
specs/005-postgresql-setup/
├── plan.md              # This file
├── research.md          # Phase 0 output - Technology decisions
├── data-model.md        # Phase 1 output - Database schema
├── quickstart.md        # Phase 1 output - Setup guide
├── contracts/           # Phase 1 output - Contract changes (none)
│   └── README.md
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
api/
├── cmd/server/
│   └── main.go                    # Update: DB initialization
├── internal/
│   ├── config/
│   │   └── config.go              # Update: DB config
│   ├── store/
│   │   ├── memory/                # Existing: unchanged
│   │   │   ├── exercise.go
│   │   │   ├── workout.go
│   │   │   ├── program.go
│   │   │   └── telemetry.go
│   │   └── postgres/              # New: PostgreSQL implementation
│   │       ├── db.go              # Connection pool management
│   │       ├── exercise.go
│   │       ├── exercise_test.go
│   │       ├── workout.go
│   │       ├── workout_test.go
│   │       ├── program.go
│   │       ├── program_test.go
│   │       ├── telemetry.go
│   │       └── telemetry_test.go
│   └── handler/
│       └── health.go              # Update: DB health check
├── migrations/                    # New: Migration files
│   ├── 000001_create_exercises.up.sql
│   ├── 000001_create_exercises.down.sql
│   ├── 000002_create_programs.up.sql
│   ├── 000002_create_programs.down.sql
│   ├── 000003_create_workouts.up.sql
│   ├── 000003_create_workouts.down.sql
│   ├── 000004_create_telemetry.up.sql
│   └── 000004_create_telemetry.down.sql
├── Makefile                       # Update: DB commands
└── go.mod                         # Update: New dependencies

docker-compose.yml                 # Update: PostgreSQL 17 (from 18)
```

**Structure Decision**: api/コンポーネント内に`internal/store/postgres/`を追加。既存の`memory/`実装と並行して使用可能。Clean Architectureのリポジトリパターンに従い、インターフェースは変更なし。

## Complexity Tracking

> No Constitution violations. Standard repository pattern implementation.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A | N/A | N/A |

## Phase Outputs

### Phase 0: Research Complete
- [research.md](./research.md) - Technology decisions documented

### Phase 1: Design Complete
- [data-model.md](./data-model.md) - Database schema design
- [contracts/README.md](./contracts/README.md) - No API changes needed
- [quickstart.md](./quickstart.md) - Setup and usage guide

### Phase 2: Tasks (Next Step)
- Run `/speckit.tasks` to generate task breakdown

## Dependencies Summary

| Package | Purpose | Version |
|---------|---------|---------|
| github.com/jackc/pgx/v5 | PostgreSQL driver | Latest |
| github.com/jackc/pgx/v5/pgxpool | Connection pooling | Latest |
| github.com/golang-migrate/migrate/v4 | Database migrations | Latest |
| github.com/cenkalti/backoff/v4 | Retry logic | Latest |
| github.com/testcontainers/testcontainers-go | Integration testing | Latest |
| github.com/testcontainers/testcontainers-go/modules/postgres | PostgreSQL module | Latest |

## Risk Assessment

| Risk | Mitigation |
|------|------------|
| Migration failure corrupts data | Transaction-based migrations with rollback support |
| Connection pool exhaustion | Configurable pool size, monitoring via health endpoint |
| Test flakiness with Testcontainers | Explicit wait strategies, container health checks |
| Breaking existing tests | Memory store remains available, gradual migration |

## Next Steps

1. Run `/speckit.tasks` to break down into implementable tasks
2. Implement tasks in priority order (P1 → P2 → P3)
3. Verify acceptance criteria from spec
