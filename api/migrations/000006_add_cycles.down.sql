DROP INDEX IF EXISTS idx_plans_cycle_id;
ALTER TABLE plans DROP COLUMN IF EXISTS cycle_id;
DROP INDEX IF EXISTS idx_cycles_program_id;
DROP TABLE IF EXISTS cycles;
