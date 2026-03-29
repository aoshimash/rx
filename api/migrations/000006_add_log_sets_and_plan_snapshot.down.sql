ALTER TABLE logs ADD CONSTRAINT uq_logs_program_id_session_name UNIQUE (program_id, session_name);
DROP TABLE IF EXISTS log_sets;
ALTER TABLE logs DROP COLUMN IF EXISTS plan_snapshot;
