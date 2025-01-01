package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error, method: %s path: %s error: %s", r.Method, r.URL.Path, err)
	_ = writeJSONError(w, http.StatusInternalServerError, "internal server error")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error, method: %s path: %s error: %s", r.Method, r.URL.Path, err)
	_ = writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error, method: %s path: %s error: %s", r.Method, r.URL.Path, err)
	_ = writeJSONError(w, http.StatusNotFound, "resource not found")
}
