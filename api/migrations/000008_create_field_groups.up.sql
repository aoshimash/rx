-- Field Groups (user-level reusable field definitions)
CREATE TABLE IF NOT EXISTS field_groups (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        VARCHAR(200) NOT NULL,
    name           VARCHAR(200) NOT NULL,
    description    TEXT,
    program_fields JSONB NOT NULL DEFAULT '[]',
    log_fields     JSONB NOT NULL DEFAULT '[]',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_field_groups_user_id ON field_groups (user_id);

-- Add field_group_id to session tables
ALTER TABLE program_sessions
    ADD COLUMN field_group_id UUID REFERENCES field_groups(id) ON DELETE SET NULL;

ALTER TABLE plan_sessions
    ADD COLUMN field_group_id UUID REFERENCES field_groups(id) ON DELETE SET NULL;
