package response

// ProducerView representa um producer agregado com informações de serviço
type ProducerView struct {
	Service    string              `json:"service"`
	Repository string              `json:"repository"`
	Topic      string              `json:"topic"`
	Owner      bool                `json:"owner"`
	Writes     bool                `json:"writes"`
	Events     []ProducerEventView `json:"events"`
}

// ProducerEventView representa um evento produzido
type ProducerEventView struct {
	Name      string            `json:"name"`
	Version   *string           `json:"version,omitempty"`
	SchemaURL string            `json:"schema_url,omitempty"`
	Headers   []EventHeaderView `json:"headers,omitempty"`
}
