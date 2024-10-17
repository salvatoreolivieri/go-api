package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (feedQuery PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	query := r.URL.Query()

	limit := query.Get("limit")

	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return feedQuery, err
		}

		feedQuery.Limit = l

	}

	offset := query.Get("offset")
	if offset != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return feedQuery, err
		}

		feedQuery.Offset = l

	}

	sort := query.Get("sort")
	if sort != "" {
		feedQuery.Sort = sort
	}

	return feedQuery, nil

}
