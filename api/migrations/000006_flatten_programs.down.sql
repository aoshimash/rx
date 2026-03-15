-- Restore program_nodes table
CREATE TABLE IF NOT EXISTS program_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES program_nodes(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    node_type VARCHAR(50) NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    exercise_id UUID REFERENCES exercises(id) ON DELETE SET NULL,
    target_sets INTEGER,
    target_reps INTEGER,
    target_rpe INTEGER CHECK (target_rpe IS NULL OR (target_rpe >= 1 AND target_rpe <= 10)),
    percent_1rm DECIMAL(5,2) CHECK (percent_1rm IS NULL OR (percent_1rm >= 0 AND percent_1rm <= 200)),
    planned_rest_seconds INTEGER,
    muscle_groups TEXT[],
    notes TEXT
);

CREATE INDEX idx_program_nodes_program_id ON program_nodes(program_id);
CREATE INDEX idx_program_nodes_parent_id ON program_nodes(parent_id);
CREATE INDEX idx_program_nodes_order ON program_nodes(program_id, parent_id, "order");

-- Restore workouts FK to program_nodes
UPDATE workouts SET program_node_id = NULL;

ALTER TABLE workouts
    DROP CONSTRAINT IF EXISTS workouts_program_entry_id_fkey;
ALTER TABLE workouts
    ADD CONSTRAINT workouts_program_node_id_fkey
    FOREIGN KEY (program_node_id) REFERENCES program_nodes(id) ON DELETE SET NULL;

-- Restore workout_entries FK to program_nodes
UPDATE workout_entries SET program_node_id = NULL;

ALTER TABLE workout_entries
    DROP CONSTRAINT IF EXISTS workout_entries_program_entry_id_fkey;
ALTER TABLE workout_entries
    ADD CONSTRAINT workout_entries_program_node_id_fkey
    FOREIGN KEY (program_node_id) REFERENCES program_nodes(id) ON DELETE SET NULL;

-- Drop program_entries table
DROP TABLE IF EXISTS program_entries;
