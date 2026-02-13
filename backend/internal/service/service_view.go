package service

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db"
)

// GetServiceView returns a detailed view of a service with what it produces and consumes.
func (s *Service) GetServiceView(ctx context.Context, serviceName string) (*response.ServiceView, error) {
	producerRows, err := db.ListProducersByService(ctx, s.db, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to list producers for service: %w", err)
	}

	consumerRows, err := db.ListConsumersByService(ctx, s.db, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to list consumers for service: %w", err)
	}

	view := &response.ServiceView{
		Service:  serviceName,
		Produces: make([]response.ServiceProduceEntry, 0, len(producerRows)),
		Consumes: make([]response.ServiceConsumeEntry, 0, len(consumerRows)),
	}

	for _, row := range producerRows {
		view.Produces = append(view.Produces, response.ServiceProduceEntry{
			Topic:  row.TopicName,
			Event:  row.EventName,
			Writes: row.Writes,
			Owner:  row.Owner,
		})
	}

	for _, row := range consumerRows {
		view.Consumes = append(view.Consumes, response.ServiceConsumeEntry{
			Topic:   row.TopicName,
			Event:   row.EventName,
			Group:   row.ConsumerGroup,
			Version: row.EventVersion,
		})
	}

	return view, nil
}
