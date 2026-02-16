package response

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error in the API response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a success message in the API response
type SuccessResponse struct {
	Message string `json:"message"`
}

// Pagination holds pagination metadata for paginated responses.
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func Write(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
