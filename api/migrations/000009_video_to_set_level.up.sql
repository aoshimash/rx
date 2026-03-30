ALTER TABLE log_sets DROP COLUMN video_url;
ALTER TABLE log_sets ADD COLUMN video_object_key VARCHAR(500);
ALTER TABLE log_entries DROP COLUMN video_object_key;
