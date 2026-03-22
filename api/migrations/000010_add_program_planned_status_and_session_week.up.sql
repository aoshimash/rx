-- Add 'planned' to program status

-- Update status check constraint to allow 'planned'
ALTER TABLE programs DROP CONSTRAINT IF EXISTS programs_new_status_check;
ALTER TABLE programs ADD CONSTRAINT programs_status_check CHECK (status IN ('active', 'completed', 'planned'));
