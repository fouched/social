package main

import (
	"errors"
	"github.com/fouched/social/internal/repo"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type CreatePostPayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Tags    string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &repo.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	if err := app.repo.Posts.Create(post); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Print(err)
		_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := app.repo.Posts.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			_ = writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		_ = writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
