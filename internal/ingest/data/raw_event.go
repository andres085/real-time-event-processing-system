package data

import (
	"database/sql"
	"time"
)

type RawEvent struct {
	ID                int       `json:"id"`
	Timestamp         time.Time `json:"timestamp"`
	Source            string    `json:"source"`
	Method            string    `json:"method"`
	Endpoint          string    `json:"endpoint"`
	StatusCode        int32     `json:"status_code"`
	ResponseTimeMs    int64     `json:"response_time_ms"`
	RequestSizeBytes  int64     `json:"request_size_bytes"`
	ResponseSizeBytes int64     `json:"response_size_bytes"`
	UserAgent         string    `json:"user_agent"`
	IpAddress         string    `json:"ip_address"`
	Processed         bool      `json:"processed"`
	CreatedAt         time.Time `json:"created_at"`
}

type RawEventModel struct {
	DB *sql.DB
}

func (r RawEventModel) Insert(rawEvent *RawEvent) error {
	query := `
	INSERT INTO raw_events(timestamp, source, method, endpoint, status_code, response_time_ms, request_size_bytes, response_size_bytes, user_agent, ip_address, processed)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id, created_at`

	args := []any{
		rawEvent.Timestamp,
		rawEvent.Source,
		rawEvent.Method,
		rawEvent.Endpoint,
		rawEvent.StatusCode,
		rawEvent.ResponseTimeMs,
		rawEvent.RequestSizeBytes,
		rawEvent.ResponseSizeBytes,
		rawEvent.UserAgent,
		rawEvent.IpAddress,
		rawEvent.Processed,
	}

	return r.DB.QueryRow(query, args...).Scan(&rawEvent.ID, &rawEvent.CreatedAt)
}
