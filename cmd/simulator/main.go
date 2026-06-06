package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andres085/real-time-event-processing-system/internal/simulator"
)

type config struct {
	port      int
	ingestURL string
	apiKey    string
	interval  time.Duration
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 6000, "Simulator server port")
	flag.StringVar(&cfg.ingestURL, "ingestURL", "http://localhost:4000/v1/ingest", "The ingest endpoint URL")
	flag.StringVar(&cfg.apiKey, "apiKey", "---", "The apiKey needed to send the request")
	flag.DurationVar(&cfg.interval, "interval", 5*time.Second, "Time between events")

	flag.Parse()

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ticker := time.NewTicker(cfg.interval)
	defer ticker.Stop()

	worker := simulator.New(cfg.ingestURL, cfg.apiKey, logger, ticker, httpClient)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info("Simulator started")
	logger.Info("With URL: ", "url", cfg.ingestURL)
	logger.Info("With Interval: ", "interval", cfg.interval)

	worker.Run(ctx)

	log.Printf("Simulator closed")
}
