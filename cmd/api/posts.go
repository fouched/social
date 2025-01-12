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

type CreatePostCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

type postKey string

const postCtx postKey = "post"

// createPost creates a new post and returns the resource
//
//	@Summary		Creates a post
//	@Description	Creates a post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post payload"
//	@Success		201		{object}	repo.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
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

	user := getUserFromContext(r)

	post := &repo.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.ID,
	}

	if err := app.repo.Posts.Create(post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getPost retrieves a post by id
//
//	@Summary		Get a post
//	@Description	Get a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	repo.Post
//	@Failure		500	{object}	error	"Server Error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [get]
func (app *application) getPost(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	comments, err := app.repo.Comments.GetByPostId(post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// updatePost updates a post
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Post payload"
//	@Success		200		{object}	repo.Post
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		404		{object}	error	"Not found"
//	@Failure		409		{object}	error	"Conflict"
//	@Failure		500		{object}	error	"Server Error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
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

	if err := app.repo.Posts.Update(post); err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			app.concurrentModification(w, r, err)
			break
		default:
			app.internalServerError(w, r, err)
			break
		}
		return
	}
	//set the previous comments
	comments, err := app.repo.Comments.GetByPostId(post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// deletePost deletes a post by ID
//
//	@Summary		Delete a post
//	@Description	Delete a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		204
//	@Failure		500	{object}	error	"Server Error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (app *application) deletePost(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	err := app.repo.Posts.Delete(post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// createPostComment creates a post comment
//
//	@Summary		Creates a post comment
//	@Description	Creates a post comment
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Post ID"
//	@Param			payload	body		CreatePostCommentPayload	true	"Post comment payload"
//	@Success		200		{object}	repo.Comment
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		500		{object}	error	"Server Error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{id}/comment [post]
func (app *application) createPostComment(w http.ResponseWriter, r *http.Request) {
	var commentPayload CreatePostCommentPayload
	if err := readJSON(w, r, &commentPayload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(commentPayload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	post := getPostFromContext(r)
	user := getUserFromContext(r)

	comment := repo.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: commentPayload.Content,
	}

	if err := app.repo.Comments.Create(&comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		post, err := app.repo.Posts.GetById(id)
		if err != nil {
			switch {
			case errors.Is(err, repo.ErrNotFound):
				app.notFound(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := r.Context()
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
