-- Add archived_at column to programs for soft-archiving
ALTER TABLE programs ADD COLUMN archived_at TIMESTAMPTZ;

CREATE INDEX idx_programs_archived_at ON programs (archived_at);
