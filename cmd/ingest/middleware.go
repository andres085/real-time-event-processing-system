package main

import (
	"net/http"
	"strings"

	"github.com/andres085/real-time-event-processing-system/internal/apikey/data"
	"github.com/andres085/real-time-event-processing-system/internal/apikey/validator"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKeyHeader := r.Header.Get("X-API-KEY")

		if apiKeyHeader == "" {
			app.invalidCredentialsResponse(w, r)
			return
		}

		headerParts := strings.Split(apiKeyHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidCredentialsResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		tokenHash := strings.SplitN(token, "_", 3)

		if len(tokenHash) != 3 {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if data.ValidateTokenPlaintext(v, tokenHash[2]); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		hashedKey := data.HashKey([]byte("pepper"), []byte(token))

		exists := app.apiKeyModels.APIKeys.GetByHash(hashedKey)
		if !exists {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
