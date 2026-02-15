package api

import (
	"net/http"
	"strconv"
)

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func getPaginationFromRequest(r *http.Request) Pagination {
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("page_size")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}

	return Pagination{
		Page:     pageInt,
		PageSize: pageSizeInt,
	}
}
