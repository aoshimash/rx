-- Restore original columns for program_session_entries
ALTER TABLE program_session_entries ADD COLUMN sets INTEGER;
ALTER TABLE program_session_entries ADD COLUMN reps INTEGER;
ALTER TABLE program_session_entries ADD COLUMN load_kg DOUBLE PRECISION;
ALTER TABLE program_session_entries ADD COLUMN metadata JSONB;

-- Restore data from fields
UPDATE program_session_entries
SET
    sets = (fields->>'sets')::INTEGER,
    reps = (fields->>'reps')::INTEGER,
    load_kg = (fields->>'load_kg')::DOUBLE PRECISION,
    metadata = fields - 'sets' - 'reps' - 'load_kg';

ALTER TABLE program_session_entries DROP COLUMN fields;

-- Restore original columns for log_entries
ALTER TABLE log_entries ADD COLUMN sets INTEGER;
ALTER TABLE log_entries ADD COLUMN reps INTEGER;
ALTER TABLE log_entries ADD COLUMN load_kg DOUBLE PRECISION;
ALTER TABLE log_entries ADD COLUMN metadata JSONB;

UPDATE log_entries
SET
    sets = (fields->>'sets')::INTEGER,
    reps = (fields->>'reps')::INTEGER,
    load_kg = (fields->>'load_kg')::DOUBLE PRECISION,
    metadata = fields - 'sets' - 'reps' - 'load_kg';

ALTER TABLE log_entries DROP COLUMN fields;

-- Remove program_fields and log_fields from programs
ALTER TABLE programs DROP COLUMN program_fields;
ALTER TABLE programs DROP COLUMN log_fields;
