package models

import (
	"time"
)

// Topic represents a topic in the messaging system.
type Topic struct {
	ID             int64
	Name           string
	Description    string
	OwnerServiceID *int64
	CreatedAt      time.Time
}

// Service represents a registered service.
type Service struct {
	ID         int64
	Name       string
	Repository string
	CreatedAt  time.Time
}

// EventHeader represents an HTTP header associated with an event.
type EventHeader struct {
	ID          int64
	EventID     int64
	Name        string
	Description string
	CreatedAt   time.Time
}

// Event represents an event published to a topic.
type Event struct {
	ID            int64
	TopicID       int64
	EventName     string
	Description   string
	SchemaURL     string
	SchemaVersion *string
	Deprecated    bool
	CreatedAt     time.Time
	Headers       []EventHeader
}

// Producer represents a service that produces events.
type Producer struct {
	ID        int64
	EventID   int64
	ServiceID int64
	Writes    bool
	CreatedAt time.Time
}

// Consumer represents a service that consumes events.
type Consumer struct {
	ID            int64
	EventID       int64
	ServiceID     int64
	ConsumerGroup string
	Description   string
	EventVersion  *string
	CreatedAt     time.Time
}

// ConsumerRow represents a query row joining consumers, events, topics, and services.
type ConsumerRow struct {
	ServiceName   string
	Repository    string
	ConsumerGroup string
	Description   string
	TopicName     string
	EventName     string
	EventVersion  *string
}

// EventRow represents a query row joining events, topics, and headers.
type EventRow struct {
	TopicName         string
	EventName         string
	EventDescription  string
	SchemaVersion     *string
	SchemaURL         string
	HeaderName        *string
	HeaderDescription *string
}

// ProducerRow represents a query row joining producers, services, events, topics, and headers.
type ProducerRow struct {
	ServiceName       string
	Repository        string
	TopicName         string
	TopicDescription  string
	Owner             bool
	Writes            bool
	EventName         string
	EventDescription  string
	SchemaVersion     *string
	SchemaURL         string
	HeaderName        *string
	HeaderDescription *string
}

// TopicProducerRow represents a query row for producers of a specific topic.
type TopicProducerRow struct {
	ServiceName  string
	Repository   string
	EventName    string
	Writes       bool
	Owner        bool
	OwnerService *string // name of the service that owns the topic (from topics.owner_service_id)
}

// TopicConsumerRow represents a query row for consumers of a specific topic.
type TopicConsumerRow struct {
	ServiceName   string
	Repository    string
	EventName     string
	ConsumerGroup string
	EventVersion  *string
}

// ServiceProducerRow represents a query row for a service's producer details.
type ServiceProducerRow struct {
	TopicName string
	EventName string
	Writes    bool
	Owner     bool
}

// ServiceConsumerRow represents a query row for a service's consumer details.
type ServiceConsumerRow struct {
	TopicName     string
	EventName     string
	ConsumerGroup string
	EventVersion  *string
}

type UndocumentedConsumerGroupRow struct {
	Topic         string
	ConsumerGroup string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
