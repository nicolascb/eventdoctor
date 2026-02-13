package response

// ServiceView represents a detailed view of a service with what it produces and consumes.
type ServiceView struct {
	Service   string                `json:"service"`
	Produces  []ServiceProduceEntry `json:"produces"`
	Consumes  []ServiceConsumeEntry `json:"consumes"`
}

// ServiceProduceEntry represents a topic/event that a service produces.
type ServiceProduceEntry struct {
	Topic  string `json:"topic"`
	Event  string `json:"event"`
	Writes bool   `json:"writes"`
	Owner  bool   `json:"owner"`
}

// ServiceConsumeEntry represents a topic/event that a service consumes.
type ServiceConsumeEntry struct {
	Topic   string  `json:"topic"`
	Event   string  `json:"event"`
	Group   string  `json:"group"`
	Version *string `json:"version,omitempty"`
}
