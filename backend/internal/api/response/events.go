package response

type EventsListView struct {
	Events     []EventView `json:"events"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type EventView struct {
	ID          int64             `json:"id"`
	Topic       string            `json:"topic"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Version     *string           `json:"version,omitempty"`
	SchemaURL   string            `json:"schema_url,omitempty"`
	Headers     []EventHeaderView `json:"headers,omitempty"`
	Producers   []EventProducer   `json:"producers"`
	Consumers   []EventConsumer   `json:"consumers"`
}

type EventProducer struct {
	ID         string `json:"id"`
	Service    string `json:"service"`
	Repository string `json:"repository"`
	Owner      string `json:"owner"`
}

type EventConsumer struct {
	ID         string `json:"id"`
	Service    string `json:"service"`
	Repository string `json:"repository"`
}

type EventHeaderView struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
