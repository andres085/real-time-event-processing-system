package data

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/andres085/real-time-event-processing-system/internal/assert"
	ingestdata "github.com/andres085/real-time-event-processing-system/internal/ingest/data"
)

func TestAggregateData(t *testing.T) {
	testData := []*ingestdata.RawEvent{
		{Source: "web_dashboard", Method: "GET"},
		{Source: "mobile_app", Method: "POST"},
		{Source: "web_dashboard", Method: "GET"},
		{Source: "mobile_app", Method: "POST"},
		{Source: "web_dashboard", Method: "GET"},
	}

	duration := 62 * time.Second

	t.Run("AggregateData", func(t *testing.T) {
		aggregatedData, err := aggregateData(testData, duration)

		var breakdown map[string]int32
		json.Unmarshal(aggregatedData.BreakdownByMethod, &breakdown)

		assert.Equal(t, breakdown["GET"], int32(3))
		assert.Equal(t, aggregatedData.DurationMinutes, 1)
		assert.NilError(t, err)
	})
}
