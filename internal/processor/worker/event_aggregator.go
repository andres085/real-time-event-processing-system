package data

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	ingestdata "github.com/andres085/real-time-event-processing-system/internal/ingest/data"
	processdata "github.com/andres085/real-time-event-processing-system/internal/processor/data"
)

type EventAggregateWorker struct {
	Duration    time.Duration
	Logger      *slog.Logger
	InputModel  ingestdata.RawEventModel
	OutputModel processdata.RequestVolumeAggregationDataModel
}

func NewEventAggregateWorker(
	duration time.Duration,
	logger *slog.Logger,
	inputModel ingestdata.RawEventModel,
	outputModel processdata.RequestVolumeAggregationDataModel,
) EventAggregateWorker {
	return EventAggregateWorker{
		duration,
		logger,
		inputModel,
		outputModel,
	}
}

func (w EventAggregateWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.Duration)
	defer ticker.Stop()

	w.Logger.Info("Worker for duration started with", "duration", w.Duration)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := w.ProcessAggregations(ctx)
			if err != nil {
				w.Logger.Error(err.Error())
				os.Exit(1)
			}
		}
	}
}

func (w EventAggregateWorker) ProcessAggregations(ctx context.Context) error {
	events, err := w.InputModel.GetRecordsByTimeLapse(w.Duration)

	if err != nil {
		w.Logger.Error("fetch events", "error", err.Error())
		return err
	}

	w.Logger.Info("fetched events", "count", len(events), "duration", w.Duration)

	a, err := aggregateData(events, w.Duration)
	if err != nil {
		w.Logger.Error("failed to aggregate data", "error", err.Error())
		return err
	}

	err = w.OutputModel.Insert(a)
	if err != nil {
		w.Logger.Error("failed to insert data", "error", err.Error())
		return err
	}

	return nil
}

func aggregateData(
	rawEvents []*ingestdata.RawEvent,
	duration time.Duration,
) (*processdata.RequestVolumeAggregationData, error) {
	r := &processdata.RequestVolumeAggregationData{}
	r.TimeBucket = time.Now()
	r.DurationMinutes = int(duration.Minutes())

	breakDownBySource := make(map[string]int32)
	breakDownByMethod := make(map[string]int32)

	for _, event := range rawEvents {
		r.TotalRequests++
		breakDownBySource[event.Source]++
		breakDownByMethod[event.Method]++
	}

	breakDownBySourceJSON, err := json.Marshal(breakDownBySource)
	if err != nil {
		return nil, err
	}
	breakDownByMethodJSON, err := json.Marshal(breakDownByMethod)
	if err != nil {
		return nil, err
	}
	r.BreakdownBySource = breakDownBySourceJSON
	r.BreakdownByMethod = breakDownByMethodJSON

	return r, nil
}
