-- Add a unique constraint on (program_id, session_name) to enforce one log per session.
-- NULL values are excluded from unique constraints in PostgreSQL, so logs without a session_name remain unrestricted.
ALTER TABLE logs
    ADD CONSTRAINT uq_logs_program_id_session_name UNIQUE (program_id, session_name);
