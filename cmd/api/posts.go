package main

import (
	"context"
	"errors"
	"github.com/fouched/social/internal/repo"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// CreatePostPayload defines a struct for the post payload. Note the optional validate syntax and keep to it
type CreatePostPayload struct {
	Title   string `json:"title" validate:"required,max=128"`
	Content string `json:"content" validate:"required,max=2000"`
	Tags    string `json:"tags,max=128"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=128"`
	Content *string `json:"content" validate:"omitempty,max=2000"`
	Tags    *string `json:"tags" validate:"omitempty,max=128"`
}

type postKey string

const postCtx postKey = "post"

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
	post := getPostFromContext(r)

	comments, err := app.repo.Comments.GetCommentsByPostId(post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	post.Comments = comments

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// updatePost retrieves a post by id
func (app *application) updatePost(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	// read payload, validate and update
	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	// partial update - only specified fields
	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Tags != nil {
		post.Tags = *payload.Tags
	}

	if err := app.repo.Posts.UpdatePost(post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	//set the previous comments
	comments, err := app.repo.Comments.GetCommentsByPostId(post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) deletePost(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	err := app.repo.Posts.DeletePost(post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		ctx := r.Context()

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

		// never mutate a context always create a new one
		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromContext(r *http.Request) *repo.Post {
	// cast to appropriate type
	post, _ := r.Context().Value(postCtx).(*repo.Post)
	return post
}
