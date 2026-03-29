ALTER TABLE programs ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'created' CHECK (status IN ('created', 'ongoing', 'completed', 'cancelled'));
ALTER TABLE programs ADD COLUMN metadata JSONB;
CREATE INDEX idx_programs_status ON programs (status);
