package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
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

	search := query.Get("search")
	if search != "" {
		feedQuery.Search = search
	}

	// TODO: implement tags, since and until
	tags := query.Get("tags")
	if tags != "" {
		feedQuery.Tags = strings.Split(tags, ",")
	}

	since := query.Get("since")
	if since != "" {
		feedQuery.Since = parseTime(since)
	}

	until := query.Get("until")
	if until != "" {
		feedQuery.Until = parseTime(until)
	}

	return feedQuery, nil

}

func parseTime(s string) string {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}
	return t.Format(time.DateTime)
}
