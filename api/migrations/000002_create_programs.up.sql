CREATE TABLE IF NOT EXISTS programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_programs_name ON programs(name);
CREATE INDEX idx_programs_created_at ON programs(created_at);

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
