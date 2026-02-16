package service

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

func (s *Service) Overview(ctx context.Context) (*response.OverviewResponse, error) {
	producerRows, consumerRows, err := s.retrieveOverviewData(ctx)
	if err != nil {
		return nil, err
	}

	topics := aggregateOverview(producerRows, consumerRows)

	return calculateOverviewMetrics(topics), nil
}

func (s *Service) retrieveOverviewData(ctx context.Context) ([]models.ProducerRow, []models.ConsumerRow, error) {
	producerRows, err := db.ListAllProducers(ctx, s.db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list producers for response.Overview: %w", err)
	}

	consumerRows, err := db.ListAllConsumers(ctx, s.db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list consumers for response.Overview: %w", err)
	}

	return producerRows, consumerRows, nil
}

func calculateOverviewMetrics(topics []response.OverviewTopicView) *response.OverviewResponse {
	producers := make(map[string]struct{})
	consumers := make(map[string]struct{})
	totalEvents := 0

	for _, t := range topics {
		totalEvents += len(t.Events)
		for _, e := range t.Events {
			trackProducers(producers, e.Producers)
			trackConsumers(consumers, e.Consumers)
		}
	}

	return &response.OverviewResponse{
		TotalTopics:    len(topics),
		TotalEvents:    totalEvents,
		TotalProducers: len(producers),
		TotalConsumers: len(consumers),
	}
}

func trackProducers(producers map[string]struct{}, rows []response.OverviewProducerView) {
	for _, row := range rows {
		key := row.Service + ":" + row.Repository
		producers[key] = struct{}{}
	}
}

func trackConsumers(consumers map[string]struct{}, rows []response.OverviewConsumerView) {
	for _, c := range rows {
		key := c.Service + ":" + c.Repository + ":" + c.Group
		consumers[key] = struct{}{}
	}
}

func aggregateOverview(producerRows []models.ProducerRow, consumerRows []models.ConsumerRow) []response.OverviewTopicView {
	topics := newOrderedMap[string, response.OverviewTopicView]()

	processOverviewProducerRows(&topics, producerRows)
	processOverviewConsumerRows(&topics, consumerRows)

	return topics.collect()
}

func processOverviewProducerRows(topics *orderedMap[string, response.OverviewTopicView], rows []models.ProducerRow) {
	for _, row := range rows {
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
}

func processOverviewConsumerRows(topics *orderedMap[string, response.OverviewTopicView], rows []models.ConsumerRow) {
	for _, row := range rows {
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
}
