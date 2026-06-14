package data

import (
	"database/sql"
	"encoding/json"
	ingestdata "github.com/andres085/real-time-event-processing-system/internal/ingest/data"
	"time"
)

type RequestVolumeAggregationData struct {
	ID                int
	TimeBucket        time.Time
	DurationMinutes   int
	TotalRequests     int
	BreakdownBySource json.RawMessage
	BreakdownByMethod json.RawMessage
	CreatedAt         time.Time
}

type RequestVolumeAggregationDataModel struct {
	DB *sql.DB
}

func (m RequestVolumeAggregationDataModel) Create(rawEvent ingestdata.RawEvent) {}
