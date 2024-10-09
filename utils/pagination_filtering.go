package utils

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func ApplyPaginationAndFiltering(r *http.Request, db *gorm.DB) *gorm.DB {

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	releaseDate := r.URL.Query().Get("releaseDate")

	query := db.Offset(offset).Limit(limit)
	if group != "" {
		query = query.Where("group = ?", group)
	}
	if song != "" {
		query = query.Where("song = ?", song)
	}
	if releaseDate != "" {
		query = query.Where("release_date = ?", releaseDate)
	}

	return query
}
