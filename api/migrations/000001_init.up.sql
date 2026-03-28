-- Programs
CREATE TABLE IF NOT EXISTS programs (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(200) NOT NULL,
    status     VARCHAR(20) NOT NULL DEFAULT 'created' CHECK (status IN ('created', 'ongoing', 'completed', 'cancelled')),
    notes      TEXT,
    metadata   JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX programs_name_unique ON programs (name);
CREATE INDEX idx_programs_status ON programs (status);

-- Program sessions (named training sessions within a program)
CREATE TABLE IF NOT EXISTS program_sessions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id   UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    session_name VARCHAR(200) NOT NULL,
    "order"      INTEGER NOT NULL DEFAULT 0,
    date         DATE
);

CREATE INDEX idx_program_sessions_program_id ON program_sessions (program_id);
CREATE INDEX idx_program_sessions_order ON program_sessions (program_id, "order");

-- Program session entries (exercise prescriptions)
CREATE TABLE IF NOT EXISTS program_session_entries (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id    UUID NOT NULL REFERENCES program_sessions(id) ON DELETE CASCADE,
    "order"       INTEGER NOT NULL DEFAULT 0,
    exercise_name VARCHAR(200) NOT NULL,
    sets          INTEGER,
    reps          INTEGER,
    load_kg       DOUBLE PRECISION,
    notes         TEXT,
    metadata      JSONB
);

CREATE INDEX idx_program_session_entries_session_id ON program_session_entries (session_id);
CREATE INDEX idx_program_session_entries_order ON program_session_entries (session_id, "order");

-- Logs (records of actual training performed)
CREATE TABLE IF NOT EXISTS logs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id   UUID REFERENCES programs(id) ON DELETE SET NULL,
    session_name VARCHAR(200),
    performed_at TIMESTAMPTZ NOT NULL,
    started_at   TIMESTAMPTZ,
    finished_at  TIMESTAMPTZ,
    notes        TEXT,
    metadata     JSONB,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_logs_program_id_session_name UNIQUE (program_id, session_name)
);

CREATE INDEX idx_logs_performed_at ON logs (performed_at);
CREATE INDEX idx_logs_program_id ON logs (program_id);
CREATE INDEX idx_logs_session_name ON logs (session_name);

-- Log entries (performed exercises within a log)
CREATE TABLE IF NOT EXISTS log_entries (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    log_id           UUID NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    "order"          INTEGER NOT NULL DEFAULT 0,
    exercise_name    VARCHAR(200) NOT NULL,
    sets             INTEGER,
    reps             INTEGER,
    load_kg          DOUBLE PRECISION,
    notes            TEXT,
    video_object_key VARCHAR(500),
    started_at       TIMESTAMPTZ,
    finished_at      TIMESTAMPTZ,
    metadata         JSONB
);

CREATE INDEX idx_log_entries_log_id ON log_entries (log_id);
CREATE INDEX idx_log_entries_order ON log_entries (log_id, "order");
