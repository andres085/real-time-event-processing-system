package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/andres085/real-time-event-processing-system/internal/api/data"
	"github.com/andres085/real-time-event-processing-system/internal/api/validator"
	"github.com/julienschmidt/httprouter"
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

func (app *application) showClientHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	client, err := app.models.Clients.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, 200, envelope{"client": client}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listClientsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		// ID   string
		Name string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	// input.ID = app.readString(qs, "id", "")
	input.Name = app.readString(qs, "name", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	clients, metadata, err := app.models.Clients.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"clients": clients, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateClientHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	client, err := app.models.Clients.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name *string `json:"name"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		client.Name = *input.Name
	}

	v := validator.New()

	if data.ValidateClient(v, client); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Clients.Update(client)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, 200, envelope{"client": client}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//
// func (app *application) deleteClientHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
//
// 	err = app.models.Movies.Delete(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, data.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
//
// 	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie sucessfully deleted"}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }
