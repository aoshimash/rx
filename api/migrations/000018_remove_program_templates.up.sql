-- Remove ProgramTemplate feature entirely
DROP TABLE IF EXISTS program_template_entries CASCADE;
DROP TABLE IF EXISTS program_templates CASCADE;

-- Remove program_template_id from programs
ALTER TABLE programs DROP COLUMN IF EXISTS program_template_id;

-- Remove rpe from program_session_entries
ALTER TABLE program_session_entries DROP COLUMN IF EXISTS rpe;

-- Remove rpe from log_entries
ALTER TABLE log_entries DROP COLUMN IF EXISTS rpe;
