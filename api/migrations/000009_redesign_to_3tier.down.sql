-- Reverse migration: Restore 4-tier model from 3-tier
-- NOTE: This down migration cannot recover data that was dropped.
--       It only restores the schema structure.

-- Restore logs.plan_id
ALTER TABLE logs
    ADD COLUMN IF NOT EXISTS plan_id UUID,
    DROP COLUMN IF EXISTS program_id,
    DROP COLUMN IF EXISTS session_name;

-- Restore old tables
CREATE TABLE IF NOT EXISTS plans (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    notes       TEXT,
    metadata    JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS plan_entries (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id       UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    "order"       INTEGER NOT NULL DEFAULT 0,
    exercise_name VARCHAR(200) NOT NULL,
    sets          INTEGER,
    reps          INTEGER,
    load_kg       DOUBLE PRECISION,
    rpe           INTEGER,
    notes         TEXT,
    metadata      JSONB
);

CREATE TABLE IF NOT EXISTS cycles (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID NOT NULL,
    name       VARCHAR(200) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata   JSONB
);

CREATE TABLE IF NOT EXISTS program_entries (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id    UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    "order"       INTEGER NOT NULL DEFAULT 0,
    exercise_name VARCHAR(200) NOT NULL,
    sets          INTEGER,
    reps          INTEGER,
    rpe           INTEGER,
    percent_1rm   DOUBLE PRECISION,
    notes         TEXT,
    metadata      JSONB
);

-- Drop new tables
DROP TABLE IF EXISTS program_session_entries CASCADE;
DROP TABLE IF EXISTS program_sessions CASCADE;
DROP TABLE IF EXISTS program_template_entries CASCADE;
DROP TABLE IF EXISTS program_templates CASCADE;

-- Drop programs table (new version) and recreate old version
DROP TABLE IF EXISTS programs CASCADE;

CREATE TABLE IF NOT EXISTS programs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    notes       TEXT,
    metadata    JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);
