package service

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

func (s *Service) ListProducers(ctx context.Context, page, pageSize int, search string) (*response.ProducersListView, error) {
	var services []models.Service
	var pagination *response.Pagination
	var err error

	if page > 0 && pageSize > 0 {
		var total int
		if search != "" {
			total, err = db.CountServicesSearch(ctx, s.db, search)
		} else {
			total, err = db.CountServices(ctx, s.db)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to count services: %w", err)
		}

		offset := (page - 1) * pageSize
		if search != "" {
			services, err = db.SearchServicesPaginated(ctx, s.db, search, pageSize, offset)
		} else {
			services, err = db.ListServicesPaginated(ctx, s.db, pageSize, offset)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list services: %w", err)
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
	}

	if len(services) == 0 {
		return &response.ProducersListView{
			Producers:  []response.ProducerServiceItem{},
			Pagination: pagination,
		}, nil
	}

	// Extract service IDs
	serviceIDs := make([]int64, len(services))
	serviceMap := make(map[int64]*response.ProducerServiceItem, len(services))

	// Prepare response items in order
	producers := make([]response.ProducerServiceItem, len(services))
	for i, svc := range services {
		serviceIDs[i] = svc.ID
		producers[i] = response.ProducerServiceItem{
			ServiceID:  svc.ID,
			Service:    svc.Name,
			Repository: svc.Repository,
			Topics:     []response.ProducerTopicItem{},
		}
		serviceMap[svc.ID] = &producers[i]
	}

	// Fetch topics for these services
	topicRows, err := db.ListTopicsForServices(ctx, s.db, serviceIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to list topics for services: %w", err)
	}

	// Group topics by service
	for _, row := range topicRows {
		if item, exists := serviceMap[row.ServiceID]; exists {
			item.Topics = append(item.Topics, response.ProducerTopicItem{
				TopicID:    row.TopicID,
				Topic:      row.TopicName,
				EventCount: row.EventCount,
				Owner:      row.Owner,
				Writes:     row.Writes,
			})
		}
	}

	return &response.ProducersListView{
		Producers:  producers,
		Pagination: pagination,
	}, nil
}

func (s *Service) GetProducerDetail(ctx context.Context, serviceID, topicID int64) (*response.ProducerDetailView, error) {
	rows, err := db.GetProducerDetail(ctx, s.db, serviceID, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get producer detail: %w", err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	first := rows[0]
	detail := &response.ProducerDetailView{
		ServiceID:   first.ServiceID,
		Service:     first.ServiceName,
		Repository:  first.Repository,
		TopicID:     first.TopicID,
		Topic:       first.TopicName,
		Description: first.TopicDescription,
		Owner:       first.Owner,
		Writes:      first.Writes,
	}

	// Aggregate events and headers from flat rows
	for _, row := range rows {
		event := findOrAppend(&detail.Events,
			func(e *response.ProducerEventEntry) bool { return e.Name == row.EventName },
			func() response.ProducerEventEntry {
				return response.ProducerEventEntry{
					Name:        row.EventName,
					Description: row.EventDescription,
					Version:     row.SchemaVersion,
					SchemaURL:   row.SchemaURL,
				}
			},
		)

		appendHeader(&event.Headers, row.HeaderName, row.HeaderDescription)
	}

	return detail, nil
}
