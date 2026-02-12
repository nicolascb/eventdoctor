package response

// ConsumerView representa um consumer agregado com informações de serviço e repositório
type ConsumerView struct {
	Service    string              `json:"service"`
	Repository string              `json:"repository"`
	Group      string              `json:"group"`
	Topics     []ConsumerTopicView `json:"topics"`
}

// ConsumerTopicView representa um tópico consumido
type ConsumerTopicView struct {
	Name   string              `json:"name"`
	Events []ConsumerEventView `json:"events"`
}

// ConsumerEventView representa um evento consumido
type ConsumerEventView struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}
