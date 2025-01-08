package repo

import (
	"net/http"
	"strconv"
	"time"
)

type PaginatedQuery struct {
	Limit  int       `json:"limit" validate:"gte=1,lte=20"`
	Offset int       `json:"offset" validate:"gte=0"`
	Sort   string    `json:"sort" validate:"oneof=asc desc"`
	Tag    string    `json:"tag" validate:"max=5"`
	Search string    `json:"search" validate:"max=100"`
	Since  time.Time `json:"since"`
	Until  time.Time `json:"until"`
}

func (pq PaginatedQuery) Parse(r *http.Request) (PaginatedQuery, error) {
	qs := r.URL.Query()

	// runs before validation so do checks as well
	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return pq, nil
		}

		pq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return pq, nil
		}

		pq.Offset = o
	}

	sort := qs.Get("sort")
	if sort != "" {
		pq.Sort = sort
	}

	tag := qs.Get("tag")
	if tag != "" {
		pq.Tag = tag
	}

	search := qs.Get("search")
	if search != "" {
		pq.Search = search
	}

	since := qs.Get("since")
	isDateSet := false
	if since != "" {
		sinceDate, err := time.Parse(time.DateOnly, since)
		if err != nil {
			isDateSet = false
		}
		isDateSet = true
		pq.Since = sinceDate
	}
	// default to last year
	if !isDateSet {
		pq.Since = time.Now().Add(time.Hour * 24 * 365 * -1)
	}

	until := qs.Get("until")
	untilDate, err := time.Parse(time.DateOnly, until)
	if err != nil {
		pq.Until = untilDate
	}

	return pq, nil
}
