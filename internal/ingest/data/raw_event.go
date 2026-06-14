// Package data provides the database models and query methods
// for the ingest service, including raw event persistence.
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
	IPAddress         string    `json:"ip_address"`
}

type RawEventModel struct {
	DB *sql.DB
}

func (m RawEventModel) Insert(rawEvent *RawEvent) error {
	query := `
	INSERT INTO raw_events(timestamp, source, method, endpoint, status_code, response_time_ms, request_size_bytes, response_size_bytes, user_agent, ip_address)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id, timestamp`

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
		rawEvent.IPAddress,
	}

	return m.DB.QueryRow(query, args...).Scan(&rawEvent.ID, &rawEvent.Timestamp)
}

func (m RawEventModel) GetRecordsByTimeLapse(
	duration time.Duration,
) ([]*RawEvent, error) {
	query := `
  		SELECT id, timestamp, source, method, endpoint, status_code, 
               response_time_ms, request_size_bytes, response_size_bytes,
               user_agent, ip_address
        FROM raw_events 
        WHERE timestamp >= NOW() - $1::interval
        ORDER BY timestamp ASC
	`

	rows, err := m.DB.Query(query, duration.String())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rawEvents := []*RawEvent{}

	for rows.Next() {
		var rawEvent RawEvent

		err := rows.Scan(
			&rawEvent.ID,
			&rawEvent.Timestamp,
			&rawEvent.Source,
			&rawEvent.Method,
			&rawEvent.Endpoint,
			&rawEvent.StatusCode,
			&rawEvent.ResponseTimeMs,
			&rawEvent.RequestSizeBytes,
			&rawEvent.ResponseSizeBytes,
			&rawEvent.UserAgent,
			&rawEvent.IPAddress,
		)
		if err != nil {
			return nil, err
		}

		rawEvents = append(rawEvents, &rawEvent)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rawEvents, nil
}
