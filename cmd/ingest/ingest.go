package main

import (
	"net/http"
	"time"

	"github.com/andres085/real-time-event-processing-system/internal/ingest/data"
	"github.com/andres085/real-time-event-processing-system/internal/ingest/validator"
)

func (app *application) ingestHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Timestamp         time.Time `json:"timestamp"`
		Source            string    `json:"source"`
		Method            string    `json:"method"`
		Endpoint          string    `json:"endpoint"`
		StatusCode        int32     `json:"status_code"`
		ResponseTimeMs    int64     `json:"response_time_ms"`
		RequestSizeBytes  int64     `json:"request_size_bytes"`
		ResponseSizeBytes int64     `json:"response_size_bytes"`
		UserAgent         string    `json:"user_agent"`
		IPAddress         string    `json:"ip_address"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	v.Check(!input.Timestamp.IsZero(), "timestamp", "must be provided")
	v.Check(input.Timestamp.Before(time.Now()), "timestamp", "cannot be in the future")
	v.Check(input.Source != "", "source", "must be provided")
	v.Check(input.Method != "", "method", "must be provided")
	v.Check(input.Endpoint != "", "endpoint", "must be provided")
	v.Check(input.StatusCode >= 100 && input.StatusCode <= 599, "status_code", "must be between 100 and 599")
	v.Check(input.ResponseTimeMs >= 0, "response_time_ms", "must be greater than 0")
	v.Check(input.RequestSizeBytes >= 0, "request_size_bytes", "must be greater than 0")
	v.Check(input.ResponseSizeBytes >= 0, "response_size_bytes", "must be greater than 0")
	v.Check(input.UserAgent != "", "user_agent", "must be provided")
	v.Check(input.IPAddress != "", "ip_address", "must be provided")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	rawEvent := &data.RawEvent{
		Timestamp:         input.Timestamp,
		Source:            input.Source,
		Method:            input.Method,
		Endpoint:          input.Endpoint,
		StatusCode:        input.StatusCode,
		ResponseTimeMs:    input.ResponseTimeMs,
		RequestSizeBytes:  input.RequestSizeBytes,
		ResponseSizeBytes: input.ResponseSizeBytes,
		UserAgent:         input.UserAgent,
		IPAddress:         input.IPAddress,
	}

	err = app.models.RawEvents.Insert(rawEvent)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "metrics stored successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
