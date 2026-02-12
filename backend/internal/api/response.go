package api

// TopicWithEvents represents the API response for events grouped by topic
type TopicWithEvents struct {
	TopicName string          `json:"topic"`
	Events    []EventResponse `json:"events"`
}

// EventHeaderResponse represents an event header in the API response
type EventHeaderResponse struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// EventResponse represents an event in the API response
type EventResponse struct {
	Name        string                `json:"name"`
	Version     *string               `json:"version,omitempty"`
	Description string                `json:"description"`
	SchemaURL   string                `json:"schema_url,omitempty"`
	Headers     []EventHeaderResponse `json:"headers,omitempty"`
}

// ProducerResponse represents a producer in the API response
type ProducerResponse struct {
	Topic       string          `json:"topic"`
	Name        string          `json:"name"`
	Owner       bool            `json:"owner"`
	Writes      bool            `json:"writes"`
	Description string          `json:"description"`
	Events      []EventResponse `json:"events"`
}

// ConsumerEventResponse represents a consumed event in the API response
type ConsumerEventResponse struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}

// ConsumerTopicResponse represents a consumed topic in the API response
type ConsumerTopicResponse struct {
	Name   string                  `json:"name"`
	Events []ConsumerEventResponse `json:"events"`
}

// ConsumerResponse represents a consumer in the API response
type ConsumerResponse struct {
	Group       string                  `json:"group"`
	Description string                  `json:"description"`
	Topics      []ConsumerTopicResponse `json:"topics"`
}

// ErrorResponse represents an error in the API response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a success message in the API response
type SuccessResponse struct {
	Message string `json:"message"`
}
