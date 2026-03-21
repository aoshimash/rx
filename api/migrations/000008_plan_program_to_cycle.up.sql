-- Migrate Plans from program_id to cycle_id.
-- For each unique program_id on plans that does not yet have a cycle,
-- create a legacy Cycle and set cycle_id on the affected plans.

-- Create legacy Cycles for existing plans that have program_id but no cycle_id
INSERT INTO cycles (id, program_id, name, created_at)
SELECT
    gen_random_uuid(),
    p.program_id,
    prog.name,
    NOW()
FROM (
    SELECT DISTINCT program_id
    FROM plans
    WHERE program_id IS NOT NULL
      AND cycle_id IS NULL
) p
JOIN programs prog ON prog.id = p.program_id
ON CONFLICT DO NOTHING;

-- Set cycle_id on plans that have program_id but no cycle_id
-- Match by program_id using the most recently created cycle for that program
UPDATE plans
SET cycle_id = (
    SELECT id
    FROM cycles
    WHERE cycles.program_id = plans.program_id
    ORDER BY created_at DESC
    LIMIT 1
)
WHERE program_id IS NOT NULL
  AND cycle_id IS NULL;

-- Drop the program_id column and its index from plans
DROP INDEX IF EXISTS idx_plans_program_id;
ALTER TABLE plans DROP COLUMN IF EXISTS program_id;
