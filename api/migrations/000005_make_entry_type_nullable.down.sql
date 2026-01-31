-- Revert entry_type to NOT NULL
-- Note: This will fail if there are NULL values in entry_type column

ALTER TABLE workout_entries ALTER COLUMN entry_type SET NOT NULL;
