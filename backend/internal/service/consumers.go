package service

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

func (s *Service) ListConsumers(ctx context.Context) ([]response.ConsumerView, error) {
	rows, err := db.ListAllConsumers(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to list consumers: %w", err)
	}

	return aggregateConsumers(rows), nil
}

func aggregateConsumers(rows []models.ConsumerRow) []response.ConsumerView {
	type consumerKey struct {
		Service    string
		Repository string
		Group      string
	}

	consumers := newOrderedMap[consumerKey, response.ConsumerView]()

	for _, row := range rows {
		key := consumerKey{
			Service:    row.ServiceName,
			Repository: row.Repository,
			Group:      row.ConsumerGroup,
		}

		cv := consumers.getOrCreate(key, func() response.ConsumerView {
			return response.ConsumerView{
				Service:    row.ServiceName,
				Repository: row.Repository,
				Group:      row.ConsumerGroup,
			}
		})

		topic := findOrAppend(&cv.Topics,
			func(t *response.ConsumerTopicView) bool { return t.Name == row.TopicName },
			func() response.ConsumerTopicView { return response.ConsumerTopicView{Name: row.TopicName} },
		)

		topic.Events = append(topic.Events, response.ConsumerEventView{
			Name:    row.EventName,
			Version: row.EventVersion,
		})
	}

	return consumers.collect()
}
