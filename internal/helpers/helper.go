package helpers

type ApiResponse struct {
	Status     bool        `json:"status"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
