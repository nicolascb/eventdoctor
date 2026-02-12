package response

// ErrorResponse represents an error in the API response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a success message in the API response
type SuccessResponse struct {
	Message string `json:"message"`
}
