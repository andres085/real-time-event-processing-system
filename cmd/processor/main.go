// Command processor reads raw events from the database and applies
// the configured processing logic over configurable time windows.
package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	ingest "github.com/andres085/real-time-event-processing-system/internal/ingest/data"
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

	rawEventModel := ingest.NewModels(db)
	// eventAggregatorModel := processdata.NewModels(db)

	defer db.Close()
	log.Printf("Processor started")

	minuteDuration := 62 * time.Second
	// dailyDuration := 24 * time.Hour
	// weeklyDuration := 168 * time.Hour

	eventAggregateWorker := processorWorker.NewEventAggregateWorker(
		minuteDuration,
		logger,
		rawEventModel.RawEvents,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Go(func() {
		eventAggregateWorker.Start(ctx)
	})

	wg.Wait()
	// go StartWorker(logger, dailyDuration, db, getRecordsByTimeLapse)
	// go StartWorker(logger, weeklyDuration, db, getRecordsByTimeLapse)
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
