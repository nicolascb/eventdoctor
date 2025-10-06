package eventdoctor

type ProducerView struct {
	Name       string
	Repository *string
}
type ConsumerView struct {
	Name       string
	Repository *string
}
type EventView struct {
	Type       string
	Version    *string
	SchemaURL  string
	Deprecated bool
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
