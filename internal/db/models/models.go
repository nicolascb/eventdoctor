package models

import (
	"time"
)

// Topic representa um tópico no sistema de mensagens
type Topic struct {
	ID        int64
	Name      string
	Owner     string
	CreatedAt time.Time
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
	EventType     string
	SchemaURL     string
	SchemaVersion *string
	Deprecated    bool
	CreatedAt     time.Time
	Headers       []EventHeader
}

// Producer representa um serviço que produz eventos
type Producer struct {
	ID         int64
	EventID    int64
	Service    string
	Writes     bool
	Repository *string
	CreatedAt  time.Time
}

// Consumer representa um serviço que consome eventos
type Consumer struct {
	ID            int64
	EventID       int64
	Service       string
	ConsumerGroup string
	EventVersion  *string
	Repository    *string
	CreatedAt     time.Time
}
