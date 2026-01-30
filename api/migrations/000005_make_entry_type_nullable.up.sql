-- Make entry_type nullable to support user-defined entry types
-- Previously: entry_type VARCHAR(50) NOT NULL
-- Now: entry_type VARCHAR(50) (nullable)

ALTER TABLE workout_entries ALTER COLUMN entry_type DROP NOT NULL;
