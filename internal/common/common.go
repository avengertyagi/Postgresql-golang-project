package common

type Pagination struct {
	CurrentPage int   `json:"currentPage"`
	PerPage     int   `json:"perPage"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"totalPage"`
}

type ApiSuccessResponse struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

type ApiErrorResponse struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}
