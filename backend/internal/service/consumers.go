package service

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

func (s *Service) ListConsumers(ctx context.Context) (*response.ConsumerView, error) {
	rows, err := db.ListAllConsumers(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to list consumers: %w", err)
	}

	undocumentedGroups, err := db.ListUndocumentedConsumerGroups(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to list undocumented consumer groups: %w", err)
	}

	res := &response.ConsumerView{
		GroupsUndocumented: make([]response.UndocumentedGroup, 0),
		Consumers:          aggregateConsumers(rows),
	}

	for _, group := range undocumentedGroups {
		res.GroupsUndocumented = append(res.GroupsUndocumented, response.UndocumentedGroup{
			Topic:     group.Topic,
			Group:     group.ConsumerGroup,
			CreatedAt: group.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: group.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return res, nil
}

func aggregateConsumers(rows []models.ConsumerRow) []response.Consumer {
	type consumerKey struct {
		Service    string
		Repository string
		Group      string
	}

	consumers := newOrderedMap[consumerKey, response.Consumer]()

	for _, row := range rows {
		key := consumerKey{
			Service:    row.ServiceName,
			Repository: row.Repository,
			Group:      row.ConsumerGroup,
		}

		cv := consumers.getOrCreate(key, func() response.Consumer {
			return response.Consumer{
				Service:     row.ServiceName,
				Repository:  row.Repository,
				Group:       row.ConsumerGroup,
				Description: row.Description,
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
