CREATE TABLE IF NOT EXISTS plans (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    VARCHAR(200) NOT NULL,
    name       VARCHAR(200),
    notes      TEXT,
    program_id UUID REFERENCES programs(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- One plan per user
CREATE UNIQUE INDEX idx_plans_user_id ON plans (user_id);

CREATE TABLE IF NOT EXISTS plan_sessions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id           UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    session_name      VARCHAR(200) NOT NULL,
    "order"           INTEGER NOT NULL DEFAULT 0,
    date              DATE,
    source_program_id UUID REFERENCES programs(id) ON DELETE SET NULL,
    source_session_id UUID
);

CREATE INDEX idx_plan_sessions_plan_id ON plan_sessions (plan_id);
CREATE INDEX idx_plan_sessions_order ON plan_sessions (plan_id, "order");

CREATE TABLE IF NOT EXISTS plan_session_entries (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id    UUID NOT NULL REFERENCES plan_sessions(id) ON DELETE CASCADE,
    "order"       INTEGER NOT NULL DEFAULT 0,
    exercise_name VARCHAR(200) NOT NULL,
    fields        JSONB,
    notes         TEXT
);

CREATE INDEX idx_plan_session_entries_session_id ON plan_session_entries (session_id);
CREATE INDEX idx_plan_session_entries_order ON plan_session_entries (session_id, "order");
