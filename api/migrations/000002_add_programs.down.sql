-- Remove date from plan_entries
ALTER TABLE plan_entries DROP COLUMN IF EXISTS date;

-- Remove program_id from plans
DROP INDEX IF EXISTS idx_plans_program_id;
ALTER TABLE plans DROP COLUMN IF EXISTS program_id;

-- Drop program tables
DROP TABLE IF EXISTS program_entries;
DROP TABLE IF EXISTS programs;
