// Command processor reads raw events from the database and applies
// the configured processing logic over configurable time windows.
package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	ingest "github.com/andres085/real-time-event-processing-system/internal/ingest/data"
	process "github.com/andres085/real-time-event-processing-system/internal/processor/data"
	processorWorker "github.com/andres085/real-time-event-processing-system/internal/processor/worker"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 3000, "Processor server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("REALTIMEPROCESSOR_DB_DSN"), "PosgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connections idle time")

	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	inputModel := ingest.NewModels(db)
	outputModel := process.NewModels(db)

	defer db.Close()
	logger.Info("Processor started")

	workerConfigs := []processorWorker.EventAggregateWorker{
		{
			Duration:    62 * time.Second,
			Logger:      logger,
			InputModel:  inputModel.RawEvents,
			OutputModel: outputModel.RequestVolumeAggregationData,
		},
		{
			Duration:    24 * time.Hour,
			Logger:      logger,
			InputModel:  inputModel.RawEvents,
			OutputModel: outputModel.RequestVolumeAggregationData,
		}, {
			Duration:    168 * time.Hour,
			Logger:      logger,
			InputModel:  inputModel.RawEvents,
			OutputModel: outputModel.RequestVolumeAggregationData,
		},
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	for _, worker := range workerConfigs {
		wg.Go(func() {
			worker.Start(ctx)
		})

	}

	<-ctx.Done()
	stop()
	logger.Info("Shutdown received, waiting for cleaning...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("All clean.")
	case <-shutdownCtx.Done():
		logger.Info("Timeout reached, forcing exit.")
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
