ALTER TABLE logs
    DROP CONSTRAINT IF EXISTS uq_logs_program_id_session_name;
