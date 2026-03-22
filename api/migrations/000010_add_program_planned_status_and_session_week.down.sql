-- Revert planned status

-- Convert any 'planned' programs back to 'active' before re-adding constraint
UPDATE programs SET status = 'active' WHERE status = 'planned';

ALTER TABLE programs DROP CONSTRAINT IF EXISTS programs_status_check;
ALTER TABLE programs ADD CONSTRAINT programs_new_status_check CHECK (status IN ('active', 'completed'));
