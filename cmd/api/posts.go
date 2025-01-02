package main

import (
	"errors"
	"github.com/fouched/social/internal/repo"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// CreatePostPayload defines a struct for the post payload. Note the optional validate syntax and keep to it
type CreatePostPayload struct {
	Title   string `json:"title" validate:"required,max=100"`
	Content string `json:"content" validate:"required,max=2000"`
	Tags    string `json:"tags"`
}

// createPost creates a new post and returns the resource
func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	post := &repo.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	if err := app.repo.Posts.CreatePost(post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getPost retrieves a post by id
func (app *application) getPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	post, err := app.repo.Posts.GetPostById(id)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			app.notFound(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	comments, err := app.repo.Comments.GetCommentsByPostId(id)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	post.Comments = comments

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
