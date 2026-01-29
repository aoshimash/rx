DROP INDEX IF EXISTS idx_workout_entries_order;
DROP INDEX IF EXISTS idx_workout_entries_exercise_id;
DROP INDEX IF EXISTS idx_workout_entries_workout_id;
DROP TABLE IF EXISTS workout_entries;

DROP INDEX IF EXISTS idx_workouts_program_node_id;
DROP INDEX IF EXISTS idx_workouts_created_at;
DROP INDEX IF EXISTS idx_workouts_timestamp;
DROP TABLE IF EXISTS workouts;
