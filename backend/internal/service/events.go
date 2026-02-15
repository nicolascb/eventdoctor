package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

func (s *Service) ListEvents(ctx context.Context, page, pageSize int, search string) (*response.EventsListView, error) {
	events, pagination, err := s.retrieveEvents(ctx, page, pageSize, search)
	if err != nil {
		return nil, err
	}

	producers, err := db.ListAllProducers(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to list producers: %w", err)
	}

	consumers, err := db.ListAllConsumers(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to list consumers: %w", err)
	}

	return &response.EventsListView{
		Events:     aggregateEvents(events, producers, consumers),
		Pagination: pagination,
	}, nil
}

func (s *Service) retrieveEvents(ctx context.Context, page, pageSize int, search string) ([]models.EventRow, *response.Pagination, error) {
	if page > 0 && pageSize > 0 {
		return s.listEventsPaginated(ctx, page, pageSize, search)
	}
	events, err := db.ListAllEvents(ctx, s.db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list all events: %w", err)
	}
	return events, nil, nil
}

func (s *Service) listEventsPaginated(ctx context.Context, page, pageSize int, search string) ([]models.EventRow, *response.Pagination, error) {
	var total int
	var err error

	if search != "" {
		total, err = db.CountEventsSearch(ctx, s.db, search)
	} else {
		total, err = db.CountEvents(ctx, s.db)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count events: %w", err)
	}

	offset := (page - 1) * pageSize
	var events []models.EventRow

	if search != "" {
		events, err = db.SearchEventsPaginated(ctx, s.db, search, pageSize, offset)
	} else {
		events, err = db.ListEventsPaginated(ctx, s.db, pageSize, offset)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list paginated events: %w", err)
	}

	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	return events, &response.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func aggregateEvents(
	eventRows []models.EventRow,
	producerRows []models.ProducerRow,
	consumerRows []models.ConsumerRow,
) []response.EventView {
	// Use orderedMap to maintain consistent order based on insertion
	events := newOrderedMap[string, response.EventView]()

	// Helper to generate key
	key := func(topic, event string) string {
		return topic + "/" + event
	}

	processEventRows(&events, eventRows, key)
	processProducerRows(&events, producerRows, key)
	processConsumerRows(&events, consumerRows, key)

	result := events.collect()

	// Sort by Topic then Name
	sort.Slice(result, func(i, j int) bool {
		if result[i].Topic == result[j].Topic {
			return result[i].Name < result[j].Name
		}
		return result[i].Topic < result[j].Topic
	})

	return result
}

func processEventRows(events *orderedMap[string, response.EventView], rows []models.EventRow, keyFunc func(string, string) string) {
	for _, row := range rows {
		k := keyFunc(row.TopicName, row.EventName)
		event := events.getOrCreate(k, func() response.EventView {
			return response.EventView{
				ID:          row.TopicName + "/" + row.EventName,
				Name:        row.EventName,
				Topic:       row.TopicName,
				Description: row.EventDescription,
				Version:     row.SchemaVersion,
				SchemaURL:   row.SchemaURL,
				// Initialize slices to avoid nil
				Headers:   []response.EventHeaderView{},
				Producers: []response.EventProducer{},
				Consumers: []response.EventConsumer{},
			}
		})
		appendHeader(&event.Headers, row.HeaderName, row.HeaderDescription)
	}
}

func processProducerRows(events *orderedMap[string, response.EventView], rows []models.ProducerRow, keyFunc func(string, string) string) {
	for _, row := range rows {
		k := keyFunc(row.TopicName, row.EventName)
		// We only add producers to existing events (INNER JOIN logic essentially)
		// If safety is needed, we could use getOrCreate, but typically events should exist first
		if event, ok := events.values[k]; ok {
			findOrAppend(&event.Producers,
				func(p *response.EventProducer) bool { return p.Service == row.ServiceName },
				func() response.EventProducer {
					return response.EventProducer{
						ID:         row.ServiceName,
						Service:    row.ServiceName,
						Repository: row.Repository,
						Owner:      fmt.Sprintf("%v", row.Owner),
					}
				},
			)
		}
	}
}

func processConsumerRows(events *orderedMap[string, response.EventView], rows []models.ConsumerRow, keyFunc func(string, string) string) {
	for _, row := range rows {
		k := keyFunc(row.TopicName, row.EventName)
		if event, ok := events.values[k]; ok {
			findOrAppend(&event.Consumers,
				func(c *response.EventConsumer) bool {
					return c.Service == row.ServiceName && c.ID == row.ConsumerGroup
				},
				func() response.EventConsumer {
					return response.EventConsumer{
						ID:         row.ConsumerGroup,
						Service:    row.ServiceName,
						Repository: row.Repository,
					}
				},
			)
		}
	}
}
