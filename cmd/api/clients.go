package main

import (
	"fmt"
	"net/http"

	"github.com/andres085/real-time-event-processing-system/internal/api/data"
)

func (app *application) createClient(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	client := &data.Client{
		Name: input.Name,
	}

	err = app.models.Clients.Insert(client)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/api/clients/%s", client.ID))

	err = app.writeJSON(
		w,
		http.StatusCreated,
		envelope{"message": "client stored successfully", "client": client},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getClientByClientID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Client By ID Handler")
}
