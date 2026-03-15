-- Plans (training plans, typically AI-generated)
CREATE TABLE IF NOT EXISTS plans (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    notes       TEXT,
    metadata    JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_plans_created_at ON plans (created_at);

-- Plan entries (exercise prescriptions within a plan)
CREATE TABLE IF NOT EXISTS plan_entries (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id       UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    "order"       INTEGER NOT NULL DEFAULT 0,
    exercise_name VARCHAR(200) NOT NULL,
    sets          INTEGER,
    reps          INTEGER,
    load_kg       DOUBLE PRECISION,
    rpe           INTEGER,
    notes         TEXT,
    metadata      JSONB
);

CREATE INDEX idx_plan_entries_plan_id ON plan_entries (plan_id);
CREATE INDEX idx_plan_entries_order ON plan_entries (plan_id, "order");

-- Logs (records of actual training performed)
CREATE TABLE IF NOT EXISTS logs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id      UUID REFERENCES plans(id) ON DELETE SET NULL,
    performed_at TIMESTAMPTZ NOT NULL,
    notes        TEXT,
    metadata     JSONB,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_logs_performed_at ON logs (performed_at);
CREATE INDEX idx_logs_plan_id ON logs (plan_id);

-- Log entries (performed exercises within a log)
CREATE TABLE IF NOT EXISTS log_entries (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    log_id           UUID NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    "order"          INTEGER NOT NULL DEFAULT 0,
    exercise_name    VARCHAR(200) NOT NULL,
    sets             INTEGER,
    reps             INTEGER,
    load_kg          DOUBLE PRECISION,
    rpe              INTEGER,
    notes            TEXT,
    video_object_key VARCHAR(500),
    metadata         JSONB
);

CREATE INDEX idx_log_entries_log_id ON log_entries (log_id);
CREATE INDEX idx_log_entries_order ON log_entries (log_id, "order");
