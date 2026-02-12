package models

import (
	"time"
)

// Topic representa um tópico no sistema de mensagens
type Topic struct {
	ID             int64
	Name           string
	OwnerServiceID *int64
	CreatedAt      time.Time
}

// Service representa um serviço registrado no sistema
type Service struct {
	ID         int64
	Name       string
	Repository string
	CreatedAt  time.Time
}

// EventHeader representa um header HTTP associado a um evento
type EventHeader struct {
	ID          int64
	EventID     int64
	Name        string
	Description string
	CreatedAt   time.Time
}

// Event representa um evento publicado em um tópico
type Event struct {
	ID            int64
	TopicID       int64
	EventName     string
	SchemaURL     string
	SchemaVersion *string
	Deprecated    bool
	CreatedAt     time.Time
	Headers       []EventHeader
}

// Producer representa um serviço que produz eventos
type Producer struct {
	ID        int64
	EventID   int64
	ServiceID int64
	Writes    bool
	CreatedAt time.Time
}

// Consumer representa um serviço que consome eventos
type Consumer struct {
	ID            int64
	EventID       int64
	ServiceID     int64
	ConsumerGroup string
	EventVersion  *string
	CreatedAt     time.Time
}

// ConsumerRow representa uma linha da query que junta consumers, events, topics e services
type ConsumerRow struct {
	ServiceName   string
	Repository    string
	ConsumerGroup string
	TopicName     string
	EventName     string
	EventVersion  *string
}

// EventRow representa uma linha da query que junta events, topics e headers
type EventRow struct {
	TopicName         string
	EventName         string
	SchemaVersion     *string
	SchemaURL         string
	HeaderName        *string
	HeaderDescription *string
}

// ProducerRow representa uma linha da query que junta producers, services, events, topics e headers
type ProducerRow struct {
	ServiceName       string
	Repository        string
	TopicName         string
	Owner             bool
	Writes            bool
	EventName         string
	SchemaVersion     *string
	SchemaURL         string
	HeaderName        *string
	HeaderDescription *string
}
