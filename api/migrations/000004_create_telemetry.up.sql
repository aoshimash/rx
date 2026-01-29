CREATE TABLE IF NOT EXISTS telemetry_points (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL,
    metric_name VARCHAR(255) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(50) NOT NULL,
    workout_id UUID REFERENCES workouts(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_telemetry_points_timestamp ON telemetry_points(timestamp);
CREATE INDEX idx_telemetry_points_metric_name ON telemetry_points(metric_name);
CREATE INDEX idx_telemetry_points_workout_id ON telemetry_points(workout_id);
CREATE INDEX idx_telemetry_points_metric_time ON telemetry_points(metric_name, timestamp);
