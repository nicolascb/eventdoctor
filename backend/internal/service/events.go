package service

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

func (s *Service) ListEvents(ctx context.Context) ([]response.TopicEventsView, error) {
	rows, err := db.ListAllEvents(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	return aggregateEvents(rows), nil
}

func aggregateEvents(rows []models.EventRow) []response.TopicEventsView {
	topics := newOrderedMap[string, response.TopicEventsView]()

	for _, row := range rows {
		topic := topics.getOrCreate(row.TopicName, func() response.TopicEventsView {
			return response.TopicEventsView{Topic: row.TopicName}
		})

		event := findOrAppend(&topic.Events,
			func(e *response.EventView) bool { return e.Name == row.EventName },
			func() response.EventView {
				return response.EventView{
					Name:        row.EventName,
					Description: row.EventDescription,
					Version:     row.SchemaVersion,
					SchemaURL:   row.SchemaURL,
				}
			},
		)

		appendHeader(&event.Headers, row.HeaderName, row.HeaderDescription)
	}

	return topics.collect()
}
