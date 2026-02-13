package response

// OverviewTopicView represents a topic in the overview with its events.
type OverviewTopicView struct {
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Events      []OverviewEventView `json:"events"`
}

// OverviewEventView represents an event in the overview with producers and consumers.
type OverviewEventView struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	SchemaURL   string                 `json:"schema_url,omitempty"`
	Headers     []EventHeaderView      `json:"headers,omitempty"`
	Producers   []OverviewProducerView `json:"producers"`
	Consumers   []OverviewConsumerView `json:"consumers"`
}

// OverviewProducerView represents a summarized producer in the overview.
type OverviewProducerView struct {
	Service    string `json:"service"`
	Repository string `json:"repository"`
	Owner      bool   `json:"owner"`
	Writes     bool   `json:"writes"`
}

// OverviewConsumerView represents a summarized consumer in the overview.
type OverviewConsumerView struct {
	Service    string `json:"service"`
	Repository string `json:"repository"`
	Group      string `json:"group"`
}

// OverviewResponse is the response envelope for the overview with consolidated statistics.
type OverviewResponse struct {
	TotalTopics    int `json:"total_topics"`
	TotalEvents    int `json:"total_events"`
	TotalProducers int `json:"total_producers"`
	TotalConsumers int `json:"total_consumers"`
}
