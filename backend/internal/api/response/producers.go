package response

// ProducerView represents an aggregated producer with service information.
type ProducerView struct {
	Service     string              `json:"service"`
	Repository  string              `json:"repository"`
	Topic       string              `json:"topic"`
	Description string              `json:"description,omitempty"`
	Owner       bool                `json:"owner"`
	Writes      bool                `json:"writes"`
	Events      []ProducerEventView `json:"events"`
}

// ProducerEventView represents a produced event.
type ProducerEventView struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Version     *string           `json:"version,omitempty"`
	SchemaURL   string            `json:"schema_url,omitempty"`
	Headers     []EventHeaderView `json:"headers,omitempty"`
}
