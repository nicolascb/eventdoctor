package service

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

func (s *Service) ListConsumers(ctx context.Context, page, pageSize int, search string) (*response.ConsumerView, error) {
	rows, pagination, err := s.retrieveConsumers(ctx, page, pageSize, search)
	if err != nil {
		return nil, err
	}

	return &response.ConsumerView{
		Consumers:  aggregateConsumers(rows),
		Pagination: pagination,
	}, nil
}

func (s *Service) retrieveConsumers(ctx context.Context, page, pageSize int, search string) ([]models.ConsumerRow, *response.Pagination, error) {
	if page > 0 && pageSize > 0 {
		return s.listConsumersPaginated(ctx, page, pageSize, search)
	}

	rows, err := db.ListAllConsumers(ctx, s.db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list consumers: %w", err)
	}
	return rows, nil, nil
}

func (s *Service) listConsumersPaginated(ctx context.Context, page, pageSize int, search string) ([]models.ConsumerRow, *response.Pagination, error) {
	var total int
	var err error

	if search != "" {
		total, err = db.CountConsumerGroupsSearch(ctx, s.db, search)
	} else {
		total, err = db.CountConsumerGroups(ctx, s.db)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count consumer groups: %w", err)
	}

	offset := (page - 1) * pageSize
	var rows []models.ConsumerRow

	if search != "" {
		rows, err = db.SearchConsumersPaginated(ctx, s.db, search, pageSize, offset)
	} else {
		rows, err = db.ListConsumersPaginated(ctx, s.db, pageSize, offset)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list consumers: %w", err)
	}

	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	return rows, &response.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) ListUndocumentedConsumers(ctx context.Context) (*response.UndocumentedConsumerView, error) {
	undocumentedGroups, err := db.ListUndocumentedConsumerGroups(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to list undocumented consumer groups: %w", err)
	}

	res := &response.UndocumentedConsumerView{
		GroupsUndocumented: make([]response.UndocumentedGroup, 0, len(undocumentedGroups)),
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
