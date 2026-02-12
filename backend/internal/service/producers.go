package service

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

func (s *Service) ListProducers(ctx context.Context) ([]response.ProducerView, error) {
	rows, err := db.ListAllProducers(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to list producers: %w", err)
	}

	return aggregateProducers(rows), nil
}

// aggregateProducers agrupa as linhas planas em uma estrutura hierárquica
// agrupando por (service, topic) -> events -> headers
func aggregateProducers(rows []models.ProducerRow) []response.ProducerView {
	type producerKey struct {
		Service    string
		Repository string
		Topic      string
	}

	producers := newOrderedMap[producerKey, response.ProducerView]()

	for _, row := range rows {
		key := producerKey{
			Service:    row.ServiceName,
			Repository: row.Repository,
			Topic:      row.TopicName,
		}

		pv := producers.getOrCreate(key, func() response.ProducerView {
			return response.ProducerView{
				Service:    row.ServiceName,
				Repository: row.Repository,
				Topic:      row.TopicName,
				Owner:      row.Owner,
				Writes:     row.Writes,
			}
		})

		event := findOrAppend(&pv.Events,
			func(e *response.ProducerEventView) bool { return e.Name == row.EventName },
			func() response.ProducerEventView {
				return response.ProducerEventView{
					Name:      row.EventName,
					Version:   row.SchemaVersion,
					SchemaURL: row.SchemaURL,
				}
			},
		)

		appendHeader(&event.Headers, row.HeaderName, row.HeaderDescription)
	}

	return producers.collect()
}
