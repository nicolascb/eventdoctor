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

// Event representa um evento publicado em um tópico
type Event struct {
	ID            int64
	TopicID       int64
	EventType     string
	SchemaURL     string
	SchemaVersion *string
	Deprecated    bool
	CreatedAt     time.Time
}

// Producer representa um serviço que produz eventos
type Producer struct {
	ID         int64
	EventID    int64
	Service    string
	Repository *string
	CreatedAt  time.Time
}

// Consumer representa um serviço que consome eventos
type Consumer struct {
	ID            int64
	EventID       int64
	Service       string
	ConsumerGroup string
	Repository    *string
	CreatedAt     time.Time
}
