package main

import (
	"github.com/fouched/social/internal/repo"
	"net/http"
)

// getUserFeed search a users geed
//
//	@Summary		Fetches a user feed
//	@Description	Fetches a user feed
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tag		query		string	false	"Tag to search"
//	@Param			search	query		string	false	"Search string"
//	@Param			since	query		string	false	"Start date"
//	@Param			until	query		string	false	"End date"
//	@Success		200		{object}	[]repo.Post
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		500		{object}	error	"Server Error"
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeed(w http.ResponseWriter, r *http.Request) {

	// use defaults since parameters are optional
	pq := repo.PaginatedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	pq, err := pq.Parse(r)
	if err != nil {
		app.badRequest(w, r, err)
		//TODO should we return here - thought we are defaulting...
		return
	}

	if err := Validate.Struct(pq); err != nil {
		app.badRequest(w, r, err)
		return
	}

	userID := int64(2) // hard code for now
	feed, err := app.repo.Posts.GetUserFeed(userID, pq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
