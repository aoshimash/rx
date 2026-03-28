-- Restore rpe columns
ALTER TABLE log_entries ADD COLUMN IF NOT EXISTS rpe INTEGER;
ALTER TABLE program_session_entries ADD COLUMN IF NOT EXISTS rpe INTEGER;

-- Restore program_template_id on programs
ALTER TABLE programs ADD COLUMN IF NOT EXISTS program_template_id UUID;

-- Restore program_templates table
CREATE TABLE IF NOT EXISTS program_templates (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    notes TEXT,
    metadata JSONB,
    weeks VARCHAR(50),
    days_per_week VARCHAR(50),
    created_by TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

-- Restore program_template_entries table
CREATE TABLE IF NOT EXISTS program_template_entries (
    id UUID PRIMARY KEY,
    program_template_id UUID NOT NULL REFERENCES program_templates(id) ON DELETE CASCADE,
    "order" INTEGER NOT NULL,
    exercise_name VARCHAR(200) NOT NULL,
    sets INTEGER,
    reps INTEGER,
    rpe INTEGER,
    percent_1rm DOUBLE PRECISION,
    notes TEXT,
    metadata JSONB
);
