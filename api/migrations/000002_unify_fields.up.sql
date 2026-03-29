-- Add fields JSONB to program_session_entries
ALTER TABLE program_session_entries ADD COLUMN fields JSONB;

-- Migrate existing data: merge sets/reps/load_kg/metadata into fields
UPDATE program_session_entries
SET fields = (
    COALESCE(metadata, '{}'::jsonb)
    || CASE WHEN sets IS NOT NULL THEN jsonb_build_object('sets', sets) ELSE '{}'::jsonb END
    || CASE WHEN reps IS NOT NULL THEN jsonb_build_object('reps', reps) ELSE '{}'::jsonb END
    || CASE WHEN load_kg IS NOT NULL THEN jsonb_build_object('load_kg', load_kg) ELSE '{}'::jsonb END
);

-- Drop old columns
ALTER TABLE program_session_entries DROP COLUMN sets;
ALTER TABLE program_session_entries DROP COLUMN reps;
ALTER TABLE program_session_entries DROP COLUMN load_kg;
ALTER TABLE program_session_entries DROP COLUMN metadata;

-- Add fields JSONB to log_entries
ALTER TABLE log_entries ADD COLUMN fields JSONB;

-- Migrate existing data
UPDATE log_entries
SET fields = (
    COALESCE(metadata, '{}'::jsonb)
    || CASE WHEN sets IS NOT NULL THEN jsonb_build_object('sets', sets) ELSE '{}'::jsonb END
    || CASE WHEN reps IS NOT NULL THEN jsonb_build_object('reps', reps) ELSE '{}'::jsonb END
    || CASE WHEN load_kg IS NOT NULL THEN jsonb_build_object('load_kg', load_kg) ELSE '{}'::jsonb END
);

-- Drop old columns
ALTER TABLE log_entries DROP COLUMN sets;
ALTER TABLE log_entries DROP COLUMN reps;
ALTER TABLE log_entries DROP COLUMN load_kg;
ALTER TABLE log_entries DROP COLUMN metadata;

-- Add program_fields and log_fields to programs
ALTER TABLE programs ADD COLUMN program_fields JSONB;
ALTER TABLE programs ADD COLUMN log_fields JSONB;

-- Set defaults for existing programs
UPDATE programs
SET
    program_fields = '[{"name":"load_kg","type":"number"},{"name":"reps","type":"number"},{"name":"sets","type":"number"}]'::jsonb,
    log_fields = '[{"name":"load_kg","type":"number"},{"name":"reps","type":"number"},{"name":"sets","type":"number"}]'::jsonb
WHERE program_fields IS NULL;
