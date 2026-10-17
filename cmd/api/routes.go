package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/api/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/api/api-key/:id", app.getAPIKeyByClientIDHandler)
	router.HandlerFunc(http.MethodPost, "/v1/api/api-key", app.createAPIKeyHandler)

	router.HandlerFunc(http.MethodPost, "/v1/api/clients", app.createClient)
	router.HandlerFunc(http.MethodGet, "/v1/api/clients/:id", app.showClientHandler)
	router.HandlerFunc(http.MethodGet, "/v1/api/clients", app.listClientsHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/api/clients/:id", app.updateClientHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/api/clients/:id", app.deleteClientHandler)

	router.HandlerFunc(http.MethodGet, "/v1/api/metrics/volume/:apiKeyId", app.getMetricsVolumeHandler)

	return router
}
