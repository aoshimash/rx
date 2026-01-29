# Data Model: PostgreSQL導入

**Feature**: 005-postgresql-setup  
**Date**: 2026-01-30

## Overview

このドキュメントでは、既存のドメインエンティティをPostgreSQLスキーマにマッピングする設計を定義する。

## Database Schema

### 1. exercises

Exerciseエンティティのテーブル定義。

```sql
CREATE TABLE IF NOT EXISTS exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    aliases TEXT[],           -- PostgreSQL配列型
    muscle_groups TEXT[],     -- PostgreSQL配列型
    load_increment DECIMAL(5,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexing
CREATE INDEX idx_exercises_name ON exercises(name);
CREATE INDEX idx_exercises_created_at ON exercises(created_at);
```

**Domain Mapping**:
| Domain Field | DB Column | Type | Notes |
|--------------|-----------|------|-------|
| ID | id | UUID | Primary key, auto-generated |
| Name | name | VARCHAR(255) | Required |
| Description | description | TEXT | Optional |
| Aliases | aliases | TEXT[] | PostgreSQL array |
| MuscleGroups | muscle_groups | TEXT[] | PostgreSQL array |
| LoadIncrement | load_increment | DECIMAL(5,2) | Nullable |
| CreatedAt | created_at | TIMESTAMPTZ | Auto-set |
| UpdatedAt | updated_at | TIMESTAMPTZ | Auto-updated |

### 2. programs

Programエンティティのテーブル定義。

```sql
CREATE TABLE IF NOT EXISTS programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_programs_name ON programs(name);
CREATE INDEX idx_programs_created_at ON programs(created_at);
```

### 3. program_nodes

ProgramNodeエンティティのテーブル定義（再帰的構造）。

```sql
CREATE TABLE IF NOT EXISTS program_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES program_nodes(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    node_type VARCHAR(50) NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    -- Prescription fields (for leaf nodes)
    exercise_id UUID REFERENCES exercises(id) ON DELETE SET NULL,
    target_sets INTEGER,
    target_reps INTEGER,
    target_rpe INTEGER CHECK (target_rpe IS NULL OR (target_rpe >= 1 AND target_rpe <= 10)),
    percent_1rm DECIMAL(5,2) CHECK (percent_1rm IS NULL OR (percent_1rm >= 0 AND percent_1rm <= 200)),
    planned_rest_seconds INTEGER,
    muscle_groups TEXT[],
    notes TEXT
);

CREATE INDEX idx_program_nodes_program_id ON program_nodes(program_id);
CREATE INDEX idx_program_nodes_parent_id ON program_nodes(parent_id);
CREATE INDEX idx_program_nodes_order ON program_nodes(program_id, parent_id, "order");
```

### 4. workouts

Workoutエンティティのテーブル定義。

```sql
CREATE TABLE IF NOT EXISTS workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL,
    session_start TIMESTAMPTZ,
    session_end TIMESTAMPTZ,
    body_weight_kg DECIMAL(5,2) CHECK (body_weight_kg IS NULL OR body_weight_kg > 0),
    fatigue_level INTEGER CHECK (fatigue_level IS NULL OR (fatigue_level >= 1 AND fatigue_level <= 10)),
    sleep_hours DECIMAL(4,2) CHECK (sleep_hours IS NULL OR (sleep_hours >= 0 AND sleep_hours <= 24)),
    condition_notes TEXT,
    program_node_id UUID REFERENCES program_nodes(id) ON DELETE SET NULL,
    program_context TEXT[],
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workouts_timestamp ON workouts(timestamp);
CREATE INDEX idx_workouts_created_at ON workouts(created_at);
CREATE INDEX idx_workouts_program_node_id ON workouts(program_node_id);
```

### 5. workout_entries

WorkoutEntryエンティティのテーブル定義。

```sql
CREATE TABLE IF NOT EXISTS workout_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    "order" INTEGER NOT NULL DEFAULT 0,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE RESTRICT,
    display_name VARCHAR(255),
    entry_type VARCHAR(50) NOT NULL,
    sets INTEGER NOT NULL CHECK (sets > 0),
    reps INTEGER NOT NULL CHECK (reps > 0),
    load_kg DECIMAL(6,2) NOT NULL CHECK (load_kg >= 0),
    rpe INTEGER NOT NULL CHECK (rpe >= 1 AND rpe <= 10),
    entry_start TIMESTAMPTZ,
    entry_end TIMESTAMPTZ,
    planned_rest_seconds INTEGER,
    performed_rest_seconds INTEGER,
    per_set_rest_overrides INTEGER[],
    program_node_id UUID REFERENCES program_nodes(id) ON DELETE SET NULL,
    -- Plan snapshot (embedded as JSONB for flexibility)
    plan_snapshot JSONB,
    notes TEXT,
    video_object_key VARCHAR(512)
);

CREATE INDEX idx_workout_entries_workout_id ON workout_entries(workout_id);
CREATE INDEX idx_workout_entries_exercise_id ON workout_entries(exercise_id);
CREATE INDEX idx_workout_entries_order ON workout_entries(workout_id, "order");
```

### 6. telemetry_points

TelemetryPointエンティティのテーブル定義。

```sql
CREATE TABLE IF NOT EXISTS telemetry_points (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL,
    metric_name VARCHAR(255) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(50) NOT NULL,
    workout_id UUID REFERENCES workouts(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_telemetry_points_timestamp ON telemetry_points(timestamp);
CREATE INDEX idx_telemetry_points_metric_name ON telemetry_points(metric_name);
CREATE INDEX idx_telemetry_points_workout_id ON telemetry_points(workout_id);
CREATE INDEX idx_telemetry_points_metric_time ON telemetry_points(metric_name, timestamp);
```

## Entity Relationships

```
┌─────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  exercises  │◄────│ workout_entries │────►│    workouts     │
└─────────────┘     └─────────────────┘     └─────────────────┘
       ▲                    │                        │
       │                    ▼                        ▼
       │            ┌─────────────────┐     ┌─────────────────┐
       └────────────│  program_nodes  │◄────│ telemetry_points│
                    └─────────────────┘     └─────────────────┘
                            │
                            ▼
                    ┌─────────────────┐
                    │    programs     │
                    └─────────────────┘
```

## Migration Strategy

### Migration Files

```
migrations/
├── 000001_create_exercises.up.sql
├── 000001_create_exercises.down.sql
├── 000002_create_programs.up.sql
├── 000002_create_programs.down.sql
├── 000003_create_workouts.up.sql
├── 000003_create_workouts.down.sql
├── 000004_create_telemetry.up.sql
└── 000004_create_telemetry.down.sql
```

### Migration Order

1. **exercises** - 独立エンティティ、外部キー依存なし
2. **programs** - 独立エンティティ、外部キー依存なし
3. **program_nodes** - programs, exercisesに依存
4. **workouts** - program_nodesに依存
5. **workout_entries** - workouts, exercises, program_nodesに依存
6. **telemetry_points** - workoutsに依存

### Rollback Strategy

各マイグレーションには対応するdownファイルを用意し、以下の順序でロールバック可能：

```bash
# 全てロールバック
make migrate-down

# 1ステップロールバック
make migrate-down-1
```

## Data Type Mapping

| Go Type | PostgreSQL Type | Notes |
|---------|-----------------|-------|
| uuid.UUID | UUID | gen_random_uuid()でデフォルト生成 |
| string | VARCHAR(n) / TEXT | 長さ制限に応じて選択 |
| *string | VARCHAR(n) / TEXT | NULLable |
| int | INTEGER | |
| *int | INTEGER | NULLable |
| float64 | DECIMAL(p,s) / DOUBLE PRECISION | 精度に応じて選択 |
| *float64 | DECIMAL(p,s) | NULLable |
| time.Time | TIMESTAMPTZ | タイムゾーン付き |
| *time.Time | TIMESTAMPTZ | NULLable |
| []string | TEXT[] | PostgreSQL配列型 |
| []int | INTEGER[] | PostgreSQL配列型 |
| struct (nested) | JSONB | 柔軟性のためJSONB使用 |

## Validation Constraints

データベースレベルで以下の制約を実装：

- **NOT NULL**: 必須フィールド
- **CHECK**: 数値範囲（RPE: 1-10、percent_1rm: 0-200等）
- **FOREIGN KEY**: 参照整合性（CASCADE / SET NULL / RESTRICT）
- **UNIQUE**: 一意性（必要に応じて追加）

ドメイン層のバリデーションと重複する部分もあるが、データベース制約はデータ整合性の最終防衛線として機能する。
