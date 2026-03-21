-- Migration: Redesign from 4-tier (Program → Cycle → Plan → Log)
--            to 3-tier (ProgramTemplate → Program → Log)

-- ── Program Templates (replaces old Programs) ────────────────────────────────

CREATE TABLE IF NOT EXISTS program_templates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    notes       TEXT,
    metadata    JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

CREATE INDEX idx_program_templates_archived_at ON program_templates (archived_at);

-- Program template entries (exercise prescriptions using RPE/percent_1rm)
CREATE TABLE IF NOT EXISTS program_template_entries (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_template_id UUID NOT NULL REFERENCES program_templates(id) ON DELETE CASCADE,
    "order"             INTEGER NOT NULL DEFAULT 0,
    exercise_name       VARCHAR(200) NOT NULL,
    sets                INTEGER,
    reps                INTEGER,
    rpe                 INTEGER,
    percent_1rm         DOUBLE PRECISION,
    notes               TEXT,
    metadata            JSONB
);

CREATE INDEX idx_program_template_entries_template_id ON program_template_entries (program_template_id);
CREATE INDEX idx_program_template_entries_order ON program_template_entries (program_template_id, "order");

-- ── Programs (new entity: concrete training program with sessions) ────────────

-- Drop old foreign key constraints and rename old programs table
-- (We drop old tables at the end after migrating data)

CREATE TABLE IF NOT EXISTS programs_new (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_template_id UUID REFERENCES program_templates(id) ON DELETE SET NULL,
    name                VARCHAR(200) NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed')),
    notes               TEXT,
    metadata            JSONB,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_programs_new_template_id ON programs_new (program_template_id);
CREATE INDEX idx_programs_new_status ON programs_new (status);

-- Program sessions (named training sessions within a program)
CREATE TABLE IF NOT EXISTS program_sessions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id   UUID NOT NULL REFERENCES programs_new(id) ON DELETE CASCADE,
    session_name VARCHAR(200) NOT NULL,
    "order"      INTEGER NOT NULL DEFAULT 0,
    date         DATE
);

CREATE INDEX idx_program_sessions_program_id ON program_sessions (program_id);
CREATE INDEX idx_program_sessions_order ON program_sessions (program_id, "order");

-- Program session entries (exercise prescriptions with absolute weights)
CREATE TABLE IF NOT EXISTS program_session_entries (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id    UUID NOT NULL REFERENCES program_sessions(id) ON DELETE CASCADE,
    "order"       INTEGER NOT NULL DEFAULT 0,
    exercise_name VARCHAR(200) NOT NULL,
    sets          INTEGER,
    reps          INTEGER,
    load_kg       DOUBLE PRECISION,
    rpe           INTEGER,
    notes         TEXT,
    metadata      JSONB
);

CREATE INDEX idx_program_session_entries_session_id ON program_session_entries (session_id);
CREATE INDEX idx_program_session_entries_order ON program_session_entries (session_id, "order");

-- ── Logs: replace plan_id with program_id + session_name ─────────────────────

ALTER TABLE logs
    ADD COLUMN IF NOT EXISTS program_id_new   UUID REFERENCES programs_new(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS session_name VARCHAR(200);

-- Drop the old plan_id foreign key and column
ALTER TABLE logs DROP COLUMN IF EXISTS plan_id;

-- Rename programs_new to programs (after dropping/renaming old tables)
-- First we need to handle the old programs table

-- Rename old tables out of the way
ALTER TABLE programs RENAME TO programs_old;
ALTER TABLE program_entries RENAME TO program_entries_old;
ALTER TABLE programs_new RENAME TO programs;

-- Rename the new program_id column
ALTER TABLE logs RENAME COLUMN program_id_new TO program_id;

-- Add index on logs.program_id
CREATE INDEX IF NOT EXISTS idx_logs_program_id ON logs (program_id);
CREATE INDEX IF NOT EXISTS idx_logs_session_name ON logs (session_name);

-- ── Drop old tables (cycles, plans, program_entries, programs_old) ─────────────

DROP TABLE IF EXISTS plan_entries CASCADE;
DROP TABLE IF EXISTS plans CASCADE;
DROP TABLE IF EXISTS cycles CASCADE;
DROP TABLE IF EXISTS program_entries_old CASCADE;
DROP TABLE IF EXISTS programs_old CASCADE;
