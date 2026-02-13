package response

// ConsumerView represents an aggregated consumer with service and repository information.
type ConsumerView struct {
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
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}
