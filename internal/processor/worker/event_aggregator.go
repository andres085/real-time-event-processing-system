package data

import (
	"context"
	"log/slog"
	"os"
	"time"

	ingestdata "github.com/andres085/real-time-event-processing-system/internal/ingest/data"
	processdata "github.com/andres085/real-time-event-processing-system/internal/processor/data"
)

type EventAggregateWorker struct {
	Duration time.Duration
	Logger   *slog.Logger
	Model    ingestdata.RawEventModel
}

func NewEventAggregateWorker(
	duration time.Duration,
	logger *slog.Logger,
	model ingestdata.RawEventModel,
) EventAggregateWorker {
	return EventAggregateWorker{
		duration,
		logger,
		model,
	}
}

func (w EventAggregateWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.Duration)
	defer ticker.Stop()
	// for range ticker.C {
	// 	err := w.ProcessAggregations(ctx)
	// 	if err != nil {
	// 		w.Logger.Error(err.Error())
	// 		os.Exit(1)
	// 	}
	// }

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
	events, err := w.Model.GetRecordsByTimeLapse(w.Duration)

	if err != nil {
		w.Logger.Error("Fetch events :%w", err)
		return err
	}

	w.Logger.Info("fetched events", "count", len(events), "duration", w.Duration)

	aggregateData(events, w.Logger)

	return nil
}

func aggregateData(rawEvents []*ingestdata.RawEvent, logger *slog.Logger) {
	r := processdata.RequestVolumeAggregationData{}
	r.TimeBucket = time.Now()
	r.DurationMinutes = 1
	breakDownBySource := make(map[string]int32)
	breakDownByMethod := make(map[string]int32)

	for _, event := range rawEvents {
		r.TotalRequests++

		if _, ok := breakDownBySource[event.Source]; ok {
			breakDownBySource[event.Source]++
		} else {
			breakDownBySource[event.Source] = 0
		}

		if _, ok := breakDownByMethod[event.Method]; ok {
			breakDownByMethod[event.Method] = 0
		} else {
			breakDownByMethod[event.Method]++
		}
	}

	logger.Info("arregation result",
		"id", r.ID,
		"time_bucket", r.TimeBucket.Format(time.RFC3339),
		"duration_minutes", r.DurationMinutes,
		"total_requests", r.TotalRequests,
		"created_at", r.CreatedAt.Format(time.RFC3339),
	)
}
