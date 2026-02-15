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
	var events []models.EventRow
	var pagination *response.Pagination
	var err error

	if page > 0 && pageSize > 0 {
		var total int
		if search != "" {
			total, err = db.CountEventsSearch(ctx, s.db, search)
		} else {
			total, err = db.CountEvents(ctx, s.db)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to count events: %w", err)
		}

		offset := (page - 1) * pageSize
		if search != "" {
			events, err = db.SearchEventsPaginated(ctx, s.db, search, pageSize, offset)
		} else {
			events, err = db.ListEventsPaginated(ctx, s.db, pageSize, offset)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list events: %w", err)
		}

		totalPages := total / pageSize
		if total%pageSize != 0 {
			totalPages++
		}

		pagination = &response.Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		}
	} else {
		events, err = db.ListAllEvents(ctx, s.db)
		if err != nil {
			return nil, fmt.Errorf("failed to list events: %w", err)
		}
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

func aggregateEvents(
	eventRows []models.EventRow,
	producerRows []models.ProducerRow,
	consumerRows []models.ConsumerRow,
) []response.EventView {
	// Map to store events by key "TopicName/EventName"
	eventMap := make(map[string]*response.EventView)
	// Helper to generate key
	key := func(topic, event string) string {
		return topic + "/" + event
	}

	// 1. Process Event Rows
	for _, row := range eventRows {
		k := key(row.TopicName, row.EventName)
		if _, exists := eventMap[k]; !exists {
			eventMap[k] = &response.EventView{
				ID:          row.TopicName + "/" + row.EventName, // simple ID generation
				Name:        row.EventName,
				Topic:       row.TopicName,
				Description: row.EventDescription,
				Version:     row.SchemaVersion,
				SchemaURL:   row.SchemaURL,
				Headers:     []response.EventHeaderView{},
				Producers:   []response.EventProducer{},
				Consumers:   []response.EventConsumer{},
			}
		}
		event := eventMap[k]
		appendHeader(&event.Headers, row.HeaderName, row.HeaderDescription)
	}

	// 2. Process Producer Rows
	for _, row := range producerRows {
		k := key(row.TopicName, row.EventName)
		if _, exists := eventMap[k]; !exists {
			// Should ideally exist if integrity is maintained, but safety first
			continue
		}
		event := eventMap[k]

		// Avoid duplicates if multiple headers caused multiple rows for same producer
		exists := false
		for _, p := range event.Producers {
			if p.Service == row.ServiceName { // Assuming one producer per service per event
				exists = true
				break
			}
		}
		if !exists {
			event.Producers = append(event.Producers, response.EventProducer{
				ID:         row.ServiceName, // unique ID for producer within event? or just service name
				Service:    row.ServiceName,
				Repository: row.Repository,
				Owner:      fmt.Sprintf("%v", row.Owner), // Owner is bool in DB, string in View?
			})
		}
	}

	// 3. Process Consumer Rows
	for _, row := range consumerRows {
		k := key(row.TopicName, row.EventName)
		if _, exists := eventMap[k]; !exists {
			continue
		}
		event := eventMap[k]

		// Avoid duplicates
		exists := false
		for _, c := range event.Consumers {
			if c.Service == row.ServiceName && c.ID == row.ConsumerGroup { // Consumer Group ID?
				exists = true
				break
			}
		}
		if !exists {
			event.Consumers = append(event.Consumers, response.EventConsumer{
				ID:         row.ConsumerGroup,
				Service:    row.ServiceName,
				Repository: row.Repository,
			})
		}
	}

	// 4. Convert map to slice and sort
	result := make([]response.EventView, 0, len(eventMap))
	for _, ev := range eventMap {
		result = append(result, *ev)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Topic == result[j].Topic {
			return result[i].Name < result[j].Name
		}
		return result[i].Topic < result[j].Topic
	})

	return result
}
