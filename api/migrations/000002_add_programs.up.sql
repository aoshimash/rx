-- Programs (reusable RPE-based training templates)
CREATE TABLE IF NOT EXISTS programs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    notes       TEXT,
    metadata    JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_programs_created_at ON programs (created_at);

-- Program entries (exercise prescriptions within a program template)
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

CREATE INDEX idx_program_entries_program_id ON program_entries (program_id);
CREATE INDEX idx_program_entries_order ON program_entries (program_id, "order");

-- Add program_id to plans (optional reference to source program)
ALTER TABLE plans ADD COLUMN program_id UUID REFERENCES programs(id) ON DELETE SET NULL;
CREATE INDEX idx_plans_program_id ON plans (program_id);

-- Add date to plan_entries (optional per-entry date)
ALTER TABLE plan_entries ADD COLUMN date DATE;
