CREATE TABLE request_volume_aggregations (
    id BIGSERIAL PRIMARY KEY,
    time_bucket TIMESTAMPTZ NOT NULL,
    duration_minutes INTEGER NOT NULL,
    total_requests INTEGER NOT NULL,
    breakdown_by_source JSONB,
    breakdown_by_method JSONB, 
    created_at TIMESTAMPTZ DEFAULT NOW(),
    client_id UUID NOT NULL REFERENCES clients(id)
);
