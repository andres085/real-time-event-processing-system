package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andres085/real-time-event-processing-system/internal/api/data"
)

func (app *application) createAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ClientID    int    `json:"client_id"`
		Description string `json:"description"`
		Environment string `json:"environment"`
		RateLimit   int32  `json:"rate_limit_per_minute"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	apiKey := &data.ApiKey{
		ClientId:    input.ClientID,
		Environment: input.Environment,
		Description: input.Description,
		IsActive:    true,
		RateLimit:   input.RateLimit,
		Version:     1,
	}

	plainTextKey, err := apiKey.KeyHash.Generate(input.Environment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.ApiKeys.Insert(apiKey)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/api/api-key/%d", apiKey.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": "api key stored successfully", "apiKey": apiKey, "key": plainTextKey}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAPIKeyByClientIDHandler(w http.ResponseWriter, r *http.Request) {
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
