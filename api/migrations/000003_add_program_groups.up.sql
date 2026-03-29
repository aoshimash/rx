CREATE TABLE IF NOT EXISTS program_groups (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id      UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    parent_group_id UUID REFERENCES program_groups(id) ON DELETE CASCADE,
    name            VARCHAR(200) NOT NULL,
    "order"         INTEGER NOT NULL DEFAULT 0,
    notes           TEXT
);

CREATE INDEX idx_program_groups_program_id ON program_groups (program_id);
CREATE INDEX idx_program_groups_parent ON program_groups (parent_group_id);
CREATE INDEX idx_program_groups_order ON program_groups (program_id, "order");

-- Add group_id to program_sessions
ALTER TABLE program_sessions ADD COLUMN group_id UUID REFERENCES program_groups(id) ON DELETE SET NULL;
CREATE INDEX idx_program_sessions_group_id ON program_sessions (group_id);
