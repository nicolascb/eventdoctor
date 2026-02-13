package response

// OverviewTopicView representa um tópico no overview com seus eventos
type OverviewTopicView struct {
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Events      []OverviewEventView `json:"events"`
}

// OverviewEventView representa um evento no overview com producers e consumers
type OverviewEventView struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	SchemaURL   string                 `json:"schema_url,omitempty"`
	Headers     []EventHeaderView      `json:"headers,omitempty"`
	Producers   []OverviewProducerView `json:"producers"`
	Consumers   []OverviewConsumerView `json:"consumers"`
}

// OverviewProducerView representa um producer resumido no overview
type OverviewProducerView struct {
	Service    string `json:"service"`
	Repository string `json:"repository"`
	Owner      bool   `json:"owner"`
	Writes     bool   `json:"writes"`
}

// OverviewConsumerView representa um consumer resumido no overview
type OverviewConsumerView struct {
	Service    string `json:"service"`
	Repository string `json:"repository"`
	Group      string `json:"group"`
}

// OverviewResponse é o envelope da resposta do overview com estatísticas consolidadas
type OverviewResponse struct {
	TotalTopics    int `json:"total_topics"`
	TotalEvents    int `json:"total_events"`
	TotalProducers int `json:"total_producers"`
	TotalConsumers int `json:"total_consumers"`
}
