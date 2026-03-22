-- Drop new constraint first so UPDATE can revert values
ALTER TABLE programs DROP CONSTRAINT IF EXISTS programs_status_check;

-- Revert status values
UPDATE programs SET status = 'active' WHERE status IN ('ongoing', 'cancelled');
UPDATE programs SET status = 'planned' WHERE status = 'created';

-- Revert constraint
ALTER TABLE programs ADD CONSTRAINT programs_status_check CHECK (status IN ('active', 'completed', 'planned'));

-- Revert default
ALTER TABLE programs ALTER COLUMN status SET DEFAULT 'active';
