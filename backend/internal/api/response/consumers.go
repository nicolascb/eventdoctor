package response

// ConsumerView represents an aggregated consumer with service and repository information.
type ConsumerView struct {
	Consumers  []Consumer  `json:"consumers"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type UndocumentedConsumerView struct {
	GroupsUndocumented []UndocumentedGroup `json:"groups_undocumented"`
}

type UndocumentedGroup struct {
	Topic     string `json:"topic"`
	Group     string `json:"group"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Consumer struct {
	Service     string              `json:"service"`
	Repository  string              `json:"repository"`
	Group       string              `json:"group"`
	Description string              `json:"description,omitempty"`
	Topics      []ConsumerTopicView `json:"topics"`
}

// ConsumerTopicView represents a consumed topic.
type ConsumerTopicView struct {
	Name   string              `json:"name"`
	Events []ConsumerEventView `json:"events"`
}

// ConsumerEventView represents a consumed event.
type ConsumerEventView struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}
