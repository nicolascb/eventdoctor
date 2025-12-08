package eventdoctor

type HeaderView struct {
	Name        string
	Description string
}

type ProducerView struct {
	Name       string
	Writes     bool
	Repository *string
}

type ConsumerView struct {
	Name         string
	EventVersion *string
	Repository   *string
}

type EventView struct {
	Type       string
	Version    *string
	SchemaURL  string
	Deprecated bool
	Headers    []HeaderView
	Producers  []ProducerView
	Consumers  []ConsumerView
}

type TopicView struct {
	Name   string
	Owner  string
	Events []EventView
}

type WebData struct {
	Topics []TopicView
}
