-- Create flat program_entries table (replaces recursive program_nodes)
CREATE TABLE IF NOT EXISTS program_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    metadata JSONB,
    exercise_id UUID REFERENCES exercises(id) ON DELETE SET NULL,
    target_sets INTEGER,
    target_reps INTEGER,
    target_rpe INTEGER CHECK (target_rpe IS NULL OR (target_rpe >= 1 AND target_rpe <= 10)),
    percent_1rm DECIMAL(5,2) CHECK (percent_1rm IS NULL OR (percent_1rm >= 0 AND percent_1rm <= 200)),
    planned_rest_seconds INTEGER,
    muscle_groups TEXT[],
    notes TEXT
);

CREATE INDEX idx_program_entries_program_id ON program_entries(program_id);
CREATE INDEX idx_program_entries_order ON program_entries(program_id, "order");
CREATE INDEX idx_program_entries_metadata ON program_entries USING GIN(metadata);

-- Migrate exercise-level nodes from program_nodes to program_entries.
-- Uses a recursive CTE to propagate week/day ancestor names into metadata.
WITH RECURSIVE node_ancestors AS (
    -- Base: root nodes (no parent)
    SELECT
        id,
        program_id,
        parent_id,
        name,
        node_type,
        "order",
        exercise_id,
        target_sets,
        target_reps,
        target_rpe,
        percent_1rm,
        planned_rest_seconds,
        muscle_groups,
        notes,
        NULL::text AS week_name,
        NULL::text AS day_name
    FROM program_nodes
    WHERE parent_id IS NULL

    UNION ALL

    SELECT
        pn.id,
        pn.program_id,
        pn.parent_id,
        pn.name,
        pn.node_type,
        pn."order",
        pn.exercise_id,
        pn.target_sets,
        pn.target_reps,
        pn.target_rpe,
        pn.percent_1rm,
        pn.planned_rest_seconds,
        pn.muscle_groups,
        pn.notes,
        CASE WHEN na.node_type = 'week' THEN na.name ELSE na.week_name END AS week_name,
        CASE WHEN na.node_type = 'day'  THEN na.name ELSE na.day_name  END AS day_name
    FROM program_nodes pn
    JOIN node_ancestors na ON na.id = pn.parent_id
)
INSERT INTO program_entries (
    id, program_id, name, "order", metadata,
    exercise_id, target_sets, target_reps, target_rpe,
    percent_1rm, planned_rest_seconds, muscle_groups, notes
)
SELECT
    id,
    program_id,
    name,
    "order",
    CASE
        WHEN week_name IS NOT NULL OR day_name IS NOT NULL
            THEN jsonb_strip_nulls(jsonb_build_object('week', week_name, 'day', day_name))
        ELSE NULL
    END,
    exercise_id,
    target_sets,
    target_reps,
    target_rpe,
    percent_1rm,
    planned_rest_seconds,
    muscle_groups,
    notes
FROM node_ancestors
WHERE node_type = 'exercise'
   OR (node_type IS NULL AND NOT EXISTS (
       SELECT 1 FROM program_nodes child WHERE child.parent_id = node_ancestors.id
   ));

-- Update workouts: point program_node_id FK to program_entries.
-- Day-level nodes are not migrated, so set those references to NULL first.
UPDATE workouts
SET program_node_id = NULL
WHERE program_node_id IS NOT NULL
  AND program_node_id NOT IN (SELECT id FROM program_entries);

ALTER TABLE workouts
    DROP CONSTRAINT IF EXISTS workouts_program_node_id_fkey;
ALTER TABLE workouts
    ADD CONSTRAINT workouts_program_entry_id_fkey
    FOREIGN KEY (program_node_id) REFERENCES program_entries(id) ON DELETE SET NULL;

-- Update workout_entries: exercise-level nodes were migrated with same IDs, so FK remains valid.
ALTER TABLE workout_entries
    DROP CONSTRAINT IF EXISTS workout_entries_program_node_id_fkey;
ALTER TABLE workout_entries
    ADD CONSTRAINT workout_entries_program_entry_id_fkey
    FOREIGN KEY (program_node_id) REFERENCES program_entries(id) ON DELETE SET NULL;

-- Drop old program_nodes table (cascade removes its indexes)
DROP TABLE IF EXISTS program_nodes;
