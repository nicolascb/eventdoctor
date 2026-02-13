package service

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

func (s *Service) Overview(ctx context.Context) (*response.OverviewResponse, error) {
	producerRows, err := db.ListAllProducers(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to list producers for response.Overview: %w", err)
	}

	consumerRows, err := db.ListAllConsumers(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to list consumers for response.Overview: %w", err)
	}

	topics := aggregateOverview(producerRows, consumerRows)

	producers := make(map[string]struct{})
	consumers := make(map[string]struct{})
	totalEvents := 0

	for _, t := range topics {
		totalEvents += len(t.Events)
		for _, e := range t.Events {
			for _, p := range e.Producers {
				producerKey := p.Service + ":" + p.Repository
				producers[producerKey] = struct{}{}
			}
			for _, c := range e.Consumers {
				consumerKey := c.Service + ":" + c.Repository + ":" + c.Group
				consumers[consumerKey] = struct{}{}
			}
		}
	}

	return &response.OverviewResponse{
		TotalTopics:    len(topics),
		TotalEvents:    totalEvents,
		TotalProducers: len(producers),
		TotalConsumers: len(consumers),
	}, nil
}

func aggregateOverview(producerRows []models.ProducerRow, consumerRows []models.ConsumerRow) []response.OverviewTopicView {
	topics := newOrderedMap[string, response.OverviewTopicView]()

	// Process producers (includes headers).
	for _, row := range producerRows {
		topic := topics.getOrCreate(row.TopicName, func() response.OverviewTopicView {
			return response.OverviewTopicView{Name: row.TopicName, Description: row.TopicDescription}
		})

		event := findOrAppend(&topic.Events,
			func(e *response.OverviewEventView) bool { return e.Name == row.EventName },
			func() response.OverviewEventView {
				return response.OverviewEventView{
					Name:        row.EventName,
					Description: row.EventDescription,
					SchemaURL:   row.SchemaURL,
				}
			},
		)

		appendHeader(&event.Headers, row.HeaderName, row.HeaderDescription)

		findOrAppend(&event.Producers,
			func(p *response.OverviewProducerView) bool { return p.Service == row.ServiceName },
			func() response.OverviewProducerView {
				return response.OverviewProducerView{
					Service:    row.ServiceName,
					Repository: row.Repository,
					Owner:      row.Owner,
					Writes:     row.Writes,
				}
			},
		)
	}

	// Process consumers (no headers).
	for _, row := range consumerRows {
		topic := topics.getOrCreate(row.TopicName, func() response.OverviewTopicView {
			return response.OverviewTopicView{Name: row.TopicName}
		})

		event := findOrAppend(&topic.Events,
			func(e *response.OverviewEventView) bool { return e.Name == row.EventName },
			func() response.OverviewEventView {
				return response.OverviewEventView{Name: row.EventName}
			},
		)

		findOrAppend(&event.Consumers,
			func(c *response.OverviewConsumerView) bool {
				return c.Service == row.ServiceName && c.Group == row.ConsumerGroup
			},
			func() response.OverviewConsumerView {
				return response.OverviewConsumerView{
					Service:    row.ServiceName,
					Repository: row.Repository,
					Group:      row.ConsumerGroup,
				}
			},
		)
	}

	return topics.collect()
}
