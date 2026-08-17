package common

type Pagination struct {
	CurrentPage int   `json:"currentPage"`
	PerPage     int   `json:"perPage"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"totalPage"`
}
