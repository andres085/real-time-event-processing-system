package main

import (
	"net/http"
	"time"
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
		IpAddress         string    `json:"ip_address"`
		Processed         bool      `json:"processed"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "metrics stored successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
