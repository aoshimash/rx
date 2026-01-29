# Contracts: PostgreSQL導入

**Feature**: 005-postgresql-setup  
**Date**: 2026-01-30

## API Contract Changes

**No API contract changes required.**

この機能はインフラストラクチャ層（データストレージ）の変更であり、APIコントラクト（OpenAPI仕様）に影響を与えない。

### Rationale

1. **Repository Pattern**: 既存のリポジトリインターフェースは変更なし
2. **Domain Models**: ドメインモデルは変更なし
3. **HTTP Handlers**: ハンドラーはリポジトリインターフェースを使用しており、実装の詳細に依存しない
4. **API Responses**: レスポンス形式は変更なし

### Health Endpoint Enhancement

`/health`エンドポイントのレスポンスに`database`フィールドが追加されるが、これは下位互換性のある追加であり、OpenAPI仕様の更新は任意。

```yaml
# Optional addition to openapi.yaml
components:
  schemas:
    HealthResponse:
      type: object
      properties:
        status:
          type: string
          enum: [healthy, unhealthy]
        database:
          type: string
          enum: [connected, disconnected]
```

## Internal Contracts

### Database Connection Interface

```go
// internal/store/postgres/db.go
type DB interface {
    // Pool returns the underlying connection pool
    Pool() *pgxpool.Pool
    
    // Ping checks database connectivity
    Ping(ctx context.Context) error
    
    // Close closes all connections
    Close()
}
```

### Repository Interfaces (Unchanged)

既存のリポジトリインターフェース（`internal/repository/`）は変更なし：

- `ExerciseRepository`
- `WorkoutRepository`
- `ProgramRepository`
- `TelemetryPointRepository`

PostgreSQL実装はこれらのインターフェースを満たす新しい実装として追加される。
