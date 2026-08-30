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

	router.HandlerFunc(http.MethodGet, "/v1/api/metrics/volume/:apiKeyId", app.getMetricsVolumeHandler)

	return router
}
