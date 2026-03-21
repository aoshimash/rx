-- Add optional started_at / finished_at timestamps to logs and log_entries
ALTER TABLE logs ADD COLUMN started_at TIMESTAMPTZ;
ALTER TABLE logs ADD COLUMN finished_at TIMESTAMPTZ;

ALTER TABLE log_entries ADD COLUMN started_at TIMESTAMPTZ;
ALTER TABLE log_entries ADD COLUMN finished_at TIMESTAMPTZ;
