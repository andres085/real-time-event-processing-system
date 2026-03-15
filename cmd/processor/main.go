package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type application struct {
	config config
	logger *slog.Logger
}

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

	defer db.Close()
	log.Printf("Processor started")

	// processMinuteAggregation()
	rawEvents, err := getRecordsByTimeLapse(db, 62*time.Second)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	for _, event := range rawEvents {
		fmt.Printf("Event => %v", event)
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

func getRecordsByTimeLapse(db *sql.DB, duration time.Duration) ([]*RawEvent, error) {
	query := `
  		SELECT id, timestamp, source, method, endpoint, status_code, 
               response_time_ms, request_size_bytes, response_size_bytes,
               user_agent, ip_address, processed, created_at
        FROM raw_events 
        WHERE timestamp >= NOW() - $1::interval
        ORDER BY timestamp ASC
	`

	rows, err := db.Query(query, duration.String())
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
			&rawEvent.IpAddress,
			&rawEvent.Processed,
			&rawEvent.CreatedAt,
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

func processMinuteAggregation() {
	ticker := time.NewTicker(1 * time.Minute)

	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("Processing one minute records")
	}
}

// func processDailyAggregation() {

// }

// func processWeeklyAggregation() {

// }

// func processMonthlyAggregation() {

// }
