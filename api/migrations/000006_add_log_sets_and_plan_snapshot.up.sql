-- Add plan_snapshot to logs
ALTER TABLE logs ADD COLUMN plan_snapshot JSONB;

-- Create log_sets table
CREATE TABLE IF NOT EXISTS log_sets (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id   UUID NOT NULL REFERENCES log_entries(id) ON DELETE CASCADE,
    set_number INTEGER NOT NULL CHECK (set_number >= 1),
    fields     JSONB NOT NULL,
    video_url  VARCHAR(2000),
    notes      TEXT
);

CREATE INDEX idx_log_sets_entry_id ON log_sets (entry_id);
CREATE INDEX idx_log_sets_order ON log_sets (entry_id, set_number);

-- Drop the unique constraint on (program_id, session_name) as it's too restrictive
-- for the new Plan-based model where logs are linked via plan_snapshot
ALTER TABLE logs DROP CONSTRAINT IF EXISTS uq_logs_program_id_session_name;
