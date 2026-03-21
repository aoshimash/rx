-- Create cycles table
CREATE TABLE IF NOT EXISTS cycles (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID NOT NULL REFERENCES programs(id),
    name       VARCHAR(200) NOT NULL,
    notes      TEXT,
    metadata   JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cycles_program_id ON cycles (program_id);

-- Add cycle_id to plans
ALTER TABLE plans ADD COLUMN cycle_id UUID REFERENCES cycles(id);
CREATE INDEX idx_plans_cycle_id ON plans (cycle_id);
