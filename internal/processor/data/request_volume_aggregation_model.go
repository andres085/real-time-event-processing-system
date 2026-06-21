package data

import (
	"context"
	"database/sql"
	"encoding/json"
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

func (m RequestVolumeAggregationDataModel) Insert(
	requestVolumeAggregationData *RequestVolumeAggregationData,
) error {
	query := `
	INSERT INTO request_volume_aggregations(time_bucket, duration_minutes, total_requests, breakdown_by_source, breakdown_by_method)
	VALUES ($1, $2, $3, $4, $5)`

	args := []any{
		requestVolumeAggregationData.TimeBucket,
		requestVolumeAggregationData.DurationMinutes,
		requestVolumeAggregationData.TotalRequests,
		requestVolumeAggregationData.BreakdownBySource,
		requestVolumeAggregationData.BreakdownByMethod,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)

	return err
}
