package service

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db"
)

// GetTopicView returns a detailed view of a topic with its producers and consumers.
func (s *Service) GetTopicView(ctx context.Context, topicName string) (*response.TopicView, error) {
	producerRows, err := db.ListProducersByTopic(ctx, s.db, topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to list producers for topic: %w", err)
	}

	consumerRows, err := db.ListConsumersByTopic(ctx, s.db, topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to list consumers for topic: %w", err)
	}

	view := &response.TopicView{
		Topic: topicName,
	}

	for _, row := range producerRows {
		// Populate the owner service name from the first row that has it set.
		if view.OwnerService == nil && row.OwnerService != nil {
			view.OwnerService = row.OwnerService
		}

		view.Producers = append(view.Producers, response.TopicProducerEntry{
			Service:    row.ServiceName,
			Repository: row.Repository,
			Event:      row.EventName,
			Writes:     row.Writes,
			Owner:      row.Owner,
		})
	}

	for _, row := range consumerRows {
		view.Consumers = append(view.Consumers, response.TopicConsumerEntry{
			Service:    row.ServiceName,
			Repository: row.Repository,
			Event:      row.EventName,
			Group:      row.ConsumerGroup,
			Version:    row.EventVersion,
		})
	}

	return view, nil
}
