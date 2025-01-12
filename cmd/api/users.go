package main

import (
	"errors"
	"github.com/fouched/social/internal/repo"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type userKey string

const userCtx userKey = "user"

// getUser gets a user by id
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	repo.User
//	@Failure		400	{object}	error	"Bad Request"
//	@Failure		404	{object}	error	"Not found"
//	@Failure		500	{object}	error	"Server Error"
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type Follower struct {
	UserID int64 `json:"user_id"`
}

// followUser follows a user based on the post url UserID path
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		500		{object}	error	"Server Error"
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/follow [put]
func (app *application) followUser(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	followedID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
	}

	if err := app.repo.Followers.Follow(followedID, user.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

// unfollowUser unfollows a user based on the post url UserID path
//
//	@Summary		Unfollows a user
//	@Description	Unfollows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int		true	"User ID"
//	@Success		204		{string}	string	"Unfollow user"
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		500		{object}	error	"Server Error"
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/unfollow [put]
func (app *application) unfollowUser(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	followedID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
	}

	if err := app.repo.Followers.Unfollow(followedID, user.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

// activateUser activates a user

//	@Summary		Activates a user
//	@Description	Activates a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error	"Not found"
//	@Failure		500		{object}	error	"Server Error"
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUser(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	err := app.repo.Users.Activate(token)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			app.notFound(w, r, err)
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}

func getUserFromContext(r *http.Request) *repo.User {
	// populated by middleware.go > AuthTokenMiddleware
	user, _ := r.Context().Value(userCtx).(*repo.User)
	return user
}
