ALTER TABLE log_entries ADD COLUMN video_object_key VARCHAR(500);
ALTER TABLE log_sets DROP COLUMN video_object_key;
ALTER TABLE log_sets ADD COLUMN video_url VARCHAR(2000);
