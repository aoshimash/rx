# Research: PostgreSQL導入

**Feature**: 005-postgresql-setup  
**Date**: 2026-01-30

## 1. PostgreSQL Driver for Go

### Decision: pgx (jackc/pgx)

### Rationale
- **Native PostgreSQL protocol**: pgxは純粋なGoで書かれたPostgreSQLドライバーで、libpqに依存しない
- **Performance**: database/sql経由より直接使用の方がパフォーマンスが良い
- **Connection pooling**: pgxpoolパッケージで組み込みのコネクションプールを提供
- **Modern features**: PostgreSQL固有の機能（COPY、LISTEN/NOTIFY、カスタム型）をサポート
- **Active maintenance**: 最もアクティブにメンテナンスされているGoのPostgreSQLドライバー
- **stdlib compatibility**: pgx/v5/stdlibでdatabase/sql互換インターフェースも提供

### Alternatives Considered
| Library | Pros | Cons | Decision |
|---------|------|------|----------|
| lib/pq | database/sql標準、広く使用 | メンテナンス停止状態、パフォーマンス劣る | ❌ |
| pgx | 高性能、アクティブ開発、機能豊富 | 学習コスト | ✅ 採用 |
| go-pg | ORM機能込み | 過剰、Dumb Backend原則に反する | ❌ |

## 2. Migration Tool

### Decision: golang-migrate/migrate

### Rationale
- **CLI & Library**: CLIツールとGoライブラリの両方として使用可能
- **Wide adoption**: Goエコシステムで最も広く採用されているマイグレーションツール
- **Database support**: PostgreSQL、MySQL、SQLite等多数のDBをサポート
- **Simple format**: `{version}_{title}.up.sql` / `{version}_{title}.down.sql` の単純なファイル形式
- **Embedded support**: Go 1.16+のembed機能でバイナリにマイグレーションを埋め込み可能
- **Atomic migrations**: トランザクションベースのマイグレーションをサポート

### Alternatives Considered
| Tool | Pros | Cons | Decision |
|------|------|------|----------|
| golang-migrate | シンプル、CLI/ライブラリ両方、広く採用 | 機能最小限 | ✅ 採用 |
| goose | 人気、Go関数マイグレーション可能 | 設定複雑 | ❌ |
| dbmate | 言語非依存、シンプル | Go統合弱い | ❌ |
| Atlas | 宣言的、スキーマ比較 | 学習コスト高い、過剰 | ❌ |

## 3. Connection Pooling Strategy

### Decision: pgxpool (pgx組み込み)

### Rationale
- **Built-in**: pgxに組み込まれており、追加依存なし
- **Configurable**: 最大/最小接続数、接続寿命、ヘルスチェックが設定可能
- **Context-aware**: Go contextとの統合が優れている
- **Connection reuse**: 効率的な接続再利用とアイドル接続管理

### Configuration for Development
```go
config := pgxpool.Config{
    MaxConns:          10,              // 最大接続数（開発環境）
    MinConns:          2,               // 最小接続数
    MaxConnLifetime:   time.Hour,       // 接続の最大寿命
    MaxConnIdleTime:   30 * time.Minute, // アイドル接続のタイムアウト
    HealthCheckPeriod: 1 * time.Minute, // ヘルスチェック間隔
}
```

## 4. Retry Strategy with Exponential Backoff

### Decision: Custom implementation with cenkalti/backoff

### Rationale
- **Well-tested**: cenkalti/backoffは広く使用されている堅牢なライブラリ
- **Configurable**: 初期間隔、最大間隔、最大リトライ回数を設定可能
- **Jitter support**: 複数クライアントからの同時リトライを分散

### Configuration
```go
backoff := backoff.NewExponentialBackOff()
backoff.InitialInterval = 1 * time.Second
backoff.MaxInterval = 30 * time.Second
backoff.MaxElapsedTime = 2 * time.Minute
// MaxRetries: 5
```

### Alternatives Considered
| Approach | Pros | Cons | Decision |
|----------|------|------|----------|
| cenkalti/backoff | 堅牢、設定柔軟、Jitter対応 | 追加依存 | ✅ 採用 |
| Custom implementation | 依存なし | 車輪の再発明、テスト不足リスク | ❌ |
| No retry | シンプル | 一時的障害に弱い | ❌ |

## 5. Testcontainers for Go

### Decision: testcontainers/testcontainers-go

### Rationale
- **Official**: Testcontainers公式のGo実装
- **PostgreSQL module**: PostgreSQL専用モジュールで設定が簡単
- **Automatic cleanup**: テスト終了時にコンテナを自動クリーンアップ
- **Parallel test support**: 各テストで独立したコンテナを起動可能
- **CI/CD compatible**: GitHub Actions等のCI環境でDocker-in-Docker対応

### Usage Pattern
```go
func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
    ctx := context.Background()
    container, err := postgres.Run(ctx, "postgres:17")
    require.NoError(t, err)
    
    connStr, err := container.ConnectionString(ctx)
    require.NoError(t, err)
    
    pool, err := pgxpool.New(ctx, connStr)
    require.NoError(t, err)
    
    cleanup := func() {
        pool.Close()
        container.Terminate(ctx)
    }
    
    return pool, cleanup
}
```

## 6. Environment Variables Configuration

### Decision: Standard Go environment variables with validation

### Rationale
- **12-factor app**: 環境変数で設定を管理する原則に従う
- **No external dependencies**: 標準ライブラリのos.Getenvを使用
- **Validation at startup**: アプリケーション起動時にバリデーション

### Environment Variables
| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | Yes (prod) | - | PostgreSQL接続URL |
| `DB_HOST` | No | localhost | データベースホスト |
| `DB_PORT` | No | 5432 | データベースポート |
| `DB_USER` | No | optel | データベースユーザー |
| `DB_PASSWORD` | No | optel | データベースパスワード |
| `DB_NAME` | No | optel_training | データベース名 |
| `DB_SSLMODE` | No | disable | SSL モード (dev: disable) |
| `DB_MAX_CONNS` | No | 10 | 最大接続数 |
| `DB_MIN_CONNS` | No | 2 | 最小接続数 |

## 7. Health Check Implementation

### Decision: Dedicated /health endpoint with DB check

### Rationale
- **Existing pattern**: 現在のAPIに`/health`エンドポイントが存在する可能性
- **Database connectivity**: DB接続状態をヘルスチェックに含める
- **Kubernetes ready**: Liveness/Readiness probeに対応

### Implementation
```go
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    
    if err := h.pool.Ping(ctx); err != nil {
        // DB connection failed
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "unhealthy",
            "database": "disconnected",
        })
        return
    }
    
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
        "database": "connected",
    })
}
```

## 8. Project Structure for PostgreSQL

### Decision: Add postgres store alongside memory store

### Rationale
- **Repository pattern**: 既存のリポジトリインターフェースを維持
- **Incremental migration**: memory storeとpostgres storeを共存可能
- **Clean architecture**: インフラ層の実装として分離

### Directory Structure
```
api/internal/store/
├── memory/          # 既存のインメモリ実装
│   ├── exercise.go
│   ├── workout.go
│   ├── program.go
│   └── telemetry.go
└── postgres/        # 新規PostgreSQL実装
    ├── db.go        # DB接続・プール管理
    ├── exercise.go
    ├── workout.go
    ├── program.go
    └── telemetry.go

api/migrations/      # マイグレーションファイル
├── 000001_create_exercises.up.sql
├── 000001_create_exercises.down.sql
├── 000002_create_programs.up.sql
├── 000002_create_programs.down.sql
├── 000003_create_workouts.up.sql
├── 000003_create_workouts.down.sql
├── 000004_create_telemetry.up.sql
└── 000004_create_telemetry.down.sql
```

## Summary

| Area | Decision | Library/Tool |
|------|----------|--------------|
| PostgreSQL Driver | pgx v5 | github.com/jackc/pgx/v5 |
| Connection Pool | pgxpool | github.com/jackc/pgx/v5/pgxpool |
| Migration | golang-migrate | github.com/golang-migrate/migrate/v4 |
| Retry | Exponential Backoff | github.com/cenkalti/backoff/v4 |
| Testing | Testcontainers | github.com/testcontainers/testcontainers-go |
| Configuration | Environment Variables | Standard library |
