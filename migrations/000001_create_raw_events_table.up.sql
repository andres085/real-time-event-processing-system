CREATE TABLE IF NOT EXISTS raw_events (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    timestamp TIMESTAMPTZ NOT NULL,
    source VARCHAR(100) NOT NULL,
    method VARCHAR(10) NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    status_code INTEGER NOT NULL,
    response_time_ms INTEGER NOT NULL,
    request_size_bytes INTEGER NOT NULL,
    response_size_bytes INTEGER NOT NULL,
    user_agent TEXT,
    ip_address INET,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
