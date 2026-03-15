package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

const version = "1.0.0"

type config struct {
	port      int
	ingestURL string
	apiKey    string
	interval  time.Duration
}

type application struct {
	config config
	logger *slog.Logger
}

type Event struct {
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
}

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 6000, "Simulator server port")
	flag.StringVar(&cfg.ingestURL, "ingestURL", "http://localhost:4000/v1/ingest", "The ingest endpoint URL")
	flag.StringVar(&cfg.apiKey, "apiKey", "---", "The apiKey needed to send the request")
	flag.DurationVar(&cfg.interval, "interval", 5*time.Second, "Time between events")

	flag.Parse()

	log.Printf("Simulator started")
	log.Printf("URL: %s", cfg.ingestURL)
	log.Printf("Interval: %v", cfg.interval)
	// log.Printf("API Key: %s...", cfg.apiKey[:15])

	ticker := time.NewTicker(cfg.interval)
	defer ticker.Stop()

	sendEvent(cfg, generateEvent())

	for range ticker.C {
		event := generateEvent()
		sendEvent(cfg, event)
	}
}

func generateEvent() Event {
	endpoints := []string{
		"/api/users/profile",
		"/api/products/list",
		"/api/orders/create",
		"/api/search",
		"/api/analytics/track",
	}

	methods := []string{"GET", "GET", "GET", "POST", "PUT", "DELETE"}
	statusCodes := []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 201, 400, 404, 500}
	sources := []string{"mobile_app", "web_dashboard", "batch_service"}

	userAgents := []string{
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
		"Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	}

	ips := []string{
		"192.168.1.100",
		"192.168.1.101",
		"192.168.1.102",
		"10.0.0.50",
		"10.0.0.51",
		"172.16.0.100",
		"172.16.0.101",
	}

	method := methods[rand.Intn(len(methods))]

	var requestSize int64
	if method == "POST" || method == "PUT" {
		requestSize = int64(rand.Intn(5000) + 500)
	} else {
		requestSize = int64(rand.Intn(500) + 100)
	}

	var responseSize int64
	endpoint := endpoints[rand.Intn(len(endpoints))]

	switch  endpoint{
	case "/api/products/list":
		responseSize = int64(rand.Intn(50000) + 10000)
	case "/api/search":
		responseSize = int64(rand.Intn(30000) + 5000)
	default:
		responseSize = int64(rand.Intn(5000) + 500)
	}

	return Event{
		Timestamp:         time.Now().UTC(),
		Source:            sources[rand.Intn(len(sources))],
		Method:            method,
		Endpoint:          endpoint,
		StatusCode:        int32(statusCodes[rand.Intn(len(statusCodes))]),
		ResponseTimeMs:    int64(rand.Intn(450) + 50),
		RequestSizeBytes:  requestSize,
		ResponseSizeBytes: responseSize,
		UserAgent:         userAgents[rand.Intn(len(userAgents))],
		IpAddress:         ips[rand.Intn(len(ips))],
	}
}

func sendEvent(cfg config, event Event) {
	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		return
	}

	req, err := http.NewRequest("POST", cfg.ingestURL, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Request error: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", cfg.apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("HTTP error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		log.Printf("✅ %s %s → %d (%dms)",
			event.Method,
			event.Endpoint,
			event.StatusCode,
			event.ResponseTimeMs,
		)
	} else {
		log.Printf("Failed: HTTP %d", resp.StatusCode)
	}
}
