package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err)
	_ = writeJSONError(w, http.StatusInternalServerError, "internal server error")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Infow("bad request", "method", r.Method, "path", r.URL.Path, "error", err)
	_ = writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Infow("not found", "method", r.Method, "path", r.URL.Path, "error", err)
	_ = writeJSONError(w, http.StatusNotFound, "resource not found")
}

func (app *application) concurrentModification(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("concurrent modification", "method", r.Method, "path", r.URL.Path, "error", err)
	_ = writeJSONError(w, http.StatusConflict, "concurrent modification error")
}

func (app *application) unauthorizedBasicAuth(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Infow("Unauthorized basic auth", "method", r.Method, "path", r.URL.Path, "error", err)
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	_ = writeJSONError(w, http.StatusUnauthorized, "Unauthorized basic auth")
}
