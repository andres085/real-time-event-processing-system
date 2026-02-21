ALTER TABLE raw_events DROP CONSTRAINT IF EXISTS raw_events_timestamp_check;

ALTER TABLE raw_events DROP CONSTRAINT IF EXISTS raw_events_status_code_check;

ALTER TABLE raw_events DROP CONSTRAINT IF EXISTS raw_events_response_time_ms_check;

ALTER TABLE raw_events DROP CONSTRAINT IF EXISTS raw_events_request_size_bytes_check;

ALTER TABLE raw_events DROP CONSTRAINT IF EXISTS raw_events_response_size_bytes_check;


