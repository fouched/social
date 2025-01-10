package main

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

// BasicAuthMiddleware returning a func seems overly complex, let us see...
func (app *application) BasicAuthMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//read auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicAuth(w, r, errors.New("authorization header missing"))
				return
			}

			//parse it > get base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicAuth(w, r, errors.New("authorization header malformed"))
				return
			}

			//decode
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicAuth(w, r, err)
				return
			}

			//check credentials
			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			credentials := strings.SplitN(string(decoded), ":", 2)
			if len(credentials) != 2 || credentials[0] != username || credentials[1] != pass {
				app.unauthorizedBasicAuth(w, r, errors.New("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
