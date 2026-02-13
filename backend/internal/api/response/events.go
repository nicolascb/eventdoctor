package response

// TopicEventsView represents a topic with its aggregated events.
type TopicEventsView struct {
	Topic  string      `json:"topic"`
	Events []EventView `json:"events"`
}

// EventView represents an event with its headers.
type EventView struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Version     *string           `json:"version,omitempty"`
	SchemaURL   string            `json:"schema_url,omitempty"`
	Headers     []EventHeaderView `json:"headers,omitempty"`
}

// EventHeaderView represents an event header.
type EventHeaderView struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
