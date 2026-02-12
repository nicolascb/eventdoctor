package response

// TopicEventsView representa um tópico com seus eventos agregados
type TopicEventsView struct {
	Topic  string      `json:"topic"`
	Events []EventView `json:"events"`
}

// EventView representa um evento com seus headers
type EventView struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Version     *string           `json:"version,omitempty"`
	SchemaURL   string            `json:"schema_url,omitempty"`
	Headers     []EventHeaderView `json:"headers,omitempty"`
}

// EventHeaderView representa um header de evento
type EventHeaderView struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
