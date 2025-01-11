package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//read auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorized(w, r, errors.New("authorization header missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorized(w, r, errors.New("authorization header malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorized(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorized(w, r, err)
			return
		}

		user, err := app.repo.Users.GetById(userID)
		if err != nil {
			app.unauthorized(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// BasicAuthMiddleware returning a func seems overly complex, let's see...
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
