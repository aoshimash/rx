DROP INDEX IF EXISTS idx_programs_status;
ALTER TABLE programs DROP COLUMN IF EXISTS status;
ALTER TABLE programs DROP COLUMN IF EXISTS metadata;
