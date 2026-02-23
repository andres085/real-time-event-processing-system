package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andres085/real-time-event-processing-system/internal/api/data"
)

func (app *application) createApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ClientId    int       `json:"client_id"`
		Description string    `json:"description"`
		IsActive    bool      `json:"is_active"`
		RateLimit   int32     `json:"rate_limit_per_minute"`
		CreatedAt   time.Time `json:"created_at"`
		Version     int32     `json:"version"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	keyHash := "asd123456"

	apiKey := &data.ApiKey{
		ClientId:    input.ClientId,
		KeyHash:     keyHash,
		Description: input.Description,
		IsActive:    input.IsActive,
		RateLimit:   input.RateLimit,
		Version:     input.Version,
	}

	err = app.models.ApiKeys.Insert(apiKey)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/api/api-key/%d", apiKey.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": "api key stored successfully", "apiKey": input}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getApiKeyByClientIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	apiKey := data.ApiKey{
		ID:          1,
		ClientId:    id,
		Description: "Test Api Key",
		IsActive:    true,
		RateLimit:   2000,
		CreatedAt:   time.Now(),
		Version:     2,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"apiKey": apiKey}, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
