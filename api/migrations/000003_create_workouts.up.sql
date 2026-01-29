CREATE TABLE IF NOT EXISTS workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL,
    session_start TIMESTAMPTZ,
    session_end TIMESTAMPTZ,
    body_weight_kg DECIMAL(5,2) CHECK (body_weight_kg IS NULL OR body_weight_kg > 0),
    fatigue_level INTEGER CHECK (fatigue_level IS NULL OR (fatigue_level >= 1 AND fatigue_level <= 10)),
    sleep_hours DECIMAL(4,2) CHECK (sleep_hours IS NULL OR (sleep_hours >= 0 AND sleep_hours <= 24)),
    condition_notes TEXT,
    program_node_id UUID REFERENCES program_nodes(id) ON DELETE SET NULL,
    program_context TEXT[],
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workouts_timestamp ON workouts(timestamp);
CREATE INDEX idx_workouts_created_at ON workouts(created_at);
CREATE INDEX idx_workouts_program_node_id ON workouts(program_node_id);

CREATE TABLE IF NOT EXISTS workout_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    "order" INTEGER NOT NULL DEFAULT 0,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE RESTRICT,
    display_name VARCHAR(255),
    entry_type VARCHAR(50) NOT NULL,
    sets INTEGER NOT NULL CHECK (sets > 0),
    reps INTEGER NOT NULL CHECK (reps > 0),
    load_kg DECIMAL(6,2) NOT NULL CHECK (load_kg >= 0),
    rpe INTEGER NOT NULL CHECK (rpe >= 1 AND rpe <= 10),
    entry_start TIMESTAMPTZ,
    entry_end TIMESTAMPTZ,
    planned_rest_seconds INTEGER,
    performed_rest_seconds INTEGER,
    per_set_rest_overrides INTEGER[],
    program_node_id UUID REFERENCES program_nodes(id) ON DELETE SET NULL,
    plan_snapshot JSONB,
    notes TEXT,
    video_object_key VARCHAR(512)
);

CREATE INDEX idx_workout_entries_workout_id ON workout_entries(workout_id);
CREATE INDEX idx_workout_entries_exercise_id ON workout_entries(exercise_id);
CREATE INDEX idx_workout_entries_order ON workout_entries(workout_id, "order");
