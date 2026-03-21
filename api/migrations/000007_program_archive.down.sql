DROP INDEX IF EXISTS idx_programs_archived_at;
ALTER TABLE programs DROP COLUMN IF EXISTS archived_at;
