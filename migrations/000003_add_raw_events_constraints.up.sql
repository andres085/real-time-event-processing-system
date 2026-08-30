ALTER TABLE raw_events ADD CONSTRAINT raw_events_status_code_check CHECK (status_code >= 100 AND status_code <= 599);

ALTER TABLE raw_events ADD CONSTRAINT raw_events_response_time_ms_check CHECK (response_time_ms >= 0);

ALTER TABLE raw_events ADD CONSTRAINT raw_events_request_size_bytes_check CHECK (request_size_bytes >= 0);

ALTER TABLE raw_events ADD CONSTRAINT raw_events_response_size_bytes_check CHECK (response_size_bytes >= 0);

