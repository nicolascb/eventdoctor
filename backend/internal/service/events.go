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

	var eventIDs []int64
	for _, e := range events {
		eventIDs = append(eventIDs, e.EventID)
	}

	producers, err := db.ListProducersByEventIDs(ctx, s.db, eventIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to list producers: %w", err)
	}

	consumers, err := db.ListConsumersByEventIDs(ctx, s.db, eventIDs)
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

// GetEvent returns a single event by ID.
func (s *Service) GetEvent(ctx context.Context, id int64) (*response.EventView, error) {
	eventRows, err := db.GetEventByID(ctx, s.db, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event by id: %w", err)
	}
	if len(eventRows) == 0 {
		return nil, nil // Not found
	}

	producers, err := db.ListProducersByEventIDs(ctx, s.db, []int64{id})
	if err != nil {
		return nil, fmt.Errorf("failed to list producers: %w", err)
	}

	consumers, err := db.ListConsumersByEventIDs(ctx, s.db, []int64{id})
	if err != nil {
		return nil, fmt.Errorf("failed to list consumers: %w", err)
	}

	// No need to filter in memory anymore as the DB returns only relevant data
	events := aggregateEvents(eventRows, producers, consumers)
	if len(events) == 0 {
		return nil, nil
	}
	// aggregateEvents returns a slice, we want the first one (should be only one)
	result := events[0]
	return &result, nil
}

func aggregateEvents(
	eventRows []models.EventRow,
	producerRows []models.ProducerRow,
	consumerRows []models.ConsumerRow,
) []response.EventView {
	// Use orderedMap to maintain consistent order based on insertion
	events := newOrderedMap[int64, response.EventView]()

	// Helper to generate key
	key := func(id int64) int64 {
		return id
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

func processEventRows(events *orderedMap[int64, response.EventView], rows []models.EventRow, keyFunc func(int64) int64) {
	for _, row := range rows {
		k := keyFunc(row.EventID)
		event := events.getOrCreate(k, func() response.EventView {
			return response.EventView{
				ID:          row.EventID,
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

func processProducerRows(events *orderedMap[int64, response.EventView], rows []models.ProducerRow, keyFunc func(int64) int64) {
	for _, row := range rows {
		k := keyFunc(row.EventID)
		// We only add producers to existing events (INNER JOIN logic essentially)
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

func processConsumerRows(events *orderedMap[int64, response.EventView], rows []models.ConsumerRow, keyFunc func(int64) int64) {
	for _, row := range rows {
		k := keyFunc(row.EventID)
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
