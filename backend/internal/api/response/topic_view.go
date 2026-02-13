package response

// TopicView represents a detailed view of a topic with its producers and consumers.
type TopicView struct {
	Topic        string               `json:"topic"`
	OwnerService *string              `json:"owner_service,omitempty"`
	Producers    []TopicProducerEntry `json:"producers"`
	Consumers    []TopicConsumerEntry `json:"consumers"`
}

// TopicProducerEntry represents a producer entry in the topic view.
type TopicProducerEntry struct {
	Service    string `json:"service"`
	Repository string `json:"repository"`
	Event      string `json:"event"`
	Writes     bool   `json:"writes"`
	Owner      bool   `json:"owner"`
}

// TopicConsumerEntry represents a consumer entry in the topic view.
type TopicConsumerEntry struct {
	Service    string  `json:"service"`
	Repository string  `json:"repository"`
	Event      string  `json:"event"`
	Group      string  `json:"group"`
	Version    *string `json:"version,omitempty"`
}
