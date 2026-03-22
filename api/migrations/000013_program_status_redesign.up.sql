-- Drop old constraint first so UPDATE can set new values
ALTER TABLE programs DROP CONSTRAINT IF EXISTS programs_status_check;

-- Migrate existing status values
UPDATE programs SET status = 'ongoing' WHERE status = 'active';
UPDATE programs SET status = 'created' WHERE status = 'planned';

-- Add new constraint
ALTER TABLE programs ADD CONSTRAINT programs_status_check CHECK (status IN ('created', 'ongoing', 'completed', 'cancelled'));

-- Update default
ALTER TABLE programs ALTER COLUMN status SET DEFAULT 'created';
