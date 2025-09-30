package api

import "time"

type Topic struct {
	ID        int
	Name      string
	Owner     string
	CreatedAt time.Time
}

type Event struct {
	ID            int
	TopicID       int
	EventType     string
	SchemaURL     string
	SchemaVersion *string
	Deprecated    bool
	CreatedAt     time.Time
}

type Producer struct {
	ID         int
	EventID    int
	Producer   string
	Repository *string
	CreatedAt  time.Time
}

type Consumer struct {
	ID         int
	EventID    int
	Consumer   string
	Repository *string
	CreatedAt  time.Time
}

type Store interface {
	TopicStore
	EventStore
	ProducerStore
	ConsumerStore
}

type TopicStore interface {
	InsertTopic(topic Topic) error
	Get(topicName string) (*Topic, error)
	RemoveTopic(topicName string) error
}

type EventStore interface {
	InsertEvent(event Event) error
	RemoveEvent(eventID int) error
	GetEventsByTopic(topicID int) ([]Event, error)
}

type ProducerStore interface {
	InsertProducer(producer Producer) error
	RemoveProducer(producerID int) error
	GetProducersByEvent(eventID int) ([]Producer, error)
}

type ConsumerStore interface {
	InsertConsumer(consumer Consumer) error
	RemoveConsumer(consumerID int) error
	GetConsumersByEvent(eventID int) ([]Consumer, error)
}
