ALTER TABLE plan_sessions DROP COLUMN IF EXISTS field_group_id;
ALTER TABLE program_sessions DROP COLUMN IF EXISTS field_group_id;
DROP TABLE IF EXISTS field_groups;
