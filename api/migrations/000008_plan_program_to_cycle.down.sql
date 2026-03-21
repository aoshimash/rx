-- Restore program_id column on plans from cycle's program_id
ALTER TABLE plans ADD COLUMN IF NOT EXISTS program_id UUID REFERENCES programs(id);

-- Populate program_id from cycle
UPDATE plans
SET program_id = cycles.program_id
FROM cycles
WHERE plans.cycle_id = cycles.id;

CREATE INDEX IF NOT EXISTS idx_plans_program_id ON plans (program_id);
