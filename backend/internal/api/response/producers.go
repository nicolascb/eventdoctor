package response

// ProducerTopicItem represents a topic within a service group.
type ProducerTopicItem struct {
	TopicID    int64  `json:"topic_id"`
	Topic      string `json:"topic"`
	EventCount int    `json:"event_count"`
	Owner      bool   `json:"owner"`
	Writes     bool   `json:"writes"`
}

// ProducerServiceItem represents a service with its produced topics.
type ProducerServiceItem struct {
	ServiceID  int64               `json:"service_id"`
	Service    string              `json:"service"`
	Repository string              `json:"repository"`
	Topics     []ProducerTopicItem `json:"topics"`
}

// ProducersListView is the paginated response for the producers listing endpoint.
type ProducersListView struct {
	Producers  []ProducerServiceItem `json:"producers"`
	Pagination *Pagination           `json:"pagination,omitempty"`
}

// ProducerDetailView is the response for a specific producer (service + topic).
type ProducerDetailView struct {
	ServiceID   int64                `json:"service_id"`
	Service     string               `json:"service"`
	Repository  string               `json:"repository"`
	TopicID     int64                `json:"topic_id"`
	Topic       string               `json:"topic"`
	Description string               `json:"description,omitempty"`
	Owner       bool                 `json:"owner"`
	Writes      bool                 `json:"writes"`
	Events      []ProducerEventEntry `json:"events"`
}

// ProducerEventEntry represents an event in the producer detail view.
type ProducerEventEntry struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Version     *string           `json:"version,omitempty"`
	SchemaURL   string            `json:"schema_url,omitempty"`
	Headers     []EventHeaderView `json:"headers,omitempty"`
}
