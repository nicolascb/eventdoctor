package service

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

func (s *Service) ListTopics(ctx context.Context, page, pageSize int, search string) (*response.TopicListView, error) {
	var names []string
	var pagination *response.Pagination
	var err error

	if page > 0 && pageSize > 0 {
		var total int
		if search != "" {
			total, err = db.CountTopicsSearch(ctx, s.db, search)
		} else {
			total, err = db.CountTopics(ctx, s.db)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to count topics: %w", err)
		}

		offset := (page - 1) * pageSize
		if search != "" {
			names, err = db.SearchTopicNamesPaginated(ctx, s.db, search, pageSize, offset)
		} else {
			names, err = db.ListTopicNamesPaginated(ctx, s.db, pageSize, offset)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list topics: %w", err)
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
		names, err = db.ListAllTopicNames(ctx, s.db)
		if err != nil {
			return nil, fmt.Errorf("failed to list topics: %w", err)
		}
	}

	topics := make([]response.TopicView, 0, len(names))
	var countEvents, countUnconsumed, countOrphaned int

	for _, name := range names {
		view, err := s.GetTopicView(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get topic view for %q: %w", name, err)
		}
		topics = append(topics, *view)

		// Collect distinct events and compute counts.
		events := distinctEvents(view)
		countEvents += len(events)
		for eventName := range events {
			hasProducer := false
			for _, p := range view.Producers {
				if p.Event == eventName {
					hasProducer = true
					break
				}
			}
			hasConsumer := false
			for _, c := range view.Consumers {
				if c.Event == eventName {
					hasConsumer = true
					break
				}
			}
			if !hasProducer {
				countOrphaned++
			}
			if !hasConsumer {
				countUnconsumed++
			}
		}
	}

	return &response.TopicListView{
		Topics:          topics,
		CountEvents:     countEvents,
		CountUnconsumed: countUnconsumed,
		CountOrphaned:   countOrphaned,
		Pagination:      pagination,
	}, nil
}

// aggregateTopicEvents converts flat EventRow results into deduplicated EventView slices,
// accumulating headers per event (same pattern as aggregateEvents but for a single topic).
func aggregateTopicEvents(rows []models.EventRow) []response.EventView {
	var events []response.EventView
	for _, row := range rows {
		ev := findOrAppend(&events,
			func(e *response.EventView) bool { return e.Name == row.EventName },
			func() response.EventView {
				return response.EventView{
					ID:          row.EventID,
					Name:        row.EventName,
					Description: row.EventDescription,
					Version:     row.SchemaVersion,
					SchemaURL:   row.SchemaURL,
				}
			},
		)
		appendHeader(&ev.Headers, row.HeaderName, row.HeaderDescription)
	}
	return events
}

// distinctEvents returns the set of unique event names from a topic view.
func distinctEvents(view *response.TopicView) map[string]struct{} {
	events := make(map[string]struct{})
	for _, e := range view.Events {
		events[e.Name] = struct{}{}
	}
	for _, p := range view.Producers {
		events[p.Event] = struct{}{}
	}
	for _, c := range view.Consumers {
		events[c.Event] = struct{}{}
	}
	return events
}

// GetTopicView returns a detailed view of a topic with its events, producers and consumers.
func (s *Service) GetTopicView(ctx context.Context, topicName string) (*response.TopicView, error) {
	eventRows, err := db.ListEventsByTopicName(ctx, s.db, topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to list events for topic: %w", err)
	}

	producerRows, err := db.ListProducersByTopic(ctx, s.db, topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to list producers for topic: %w", err)
	}

	consumerRows, err := db.ListConsumersByTopic(ctx, s.db, topicName)
	if err != nil {
		return nil, fmt.Errorf("failed to list consumers for topic: %w", err)
	}

	// Aggregate events (reuse the same pattern from aggregateEvents).
	eventViews := aggregateTopicEvents(eventRows)

	view := &response.TopicView{
		Topic:     topicName,
		Events:    eventViews,
		Producers: make([]response.TopicProducerEntry, 0, len(producerRows)),
		Consumers: make([]response.TopicConsumerEntry, 0, len(consumerRows)),
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
