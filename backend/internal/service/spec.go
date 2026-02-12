package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/db/models"
	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) SaveSpec(ctx context.Context, spec eventdoctor.EventDoctorSpec) error {
	return db.WithTransaction(ctx, s.db, func(ctx context.Context, tx db.SQLExecutor) error {
		repository := spec.Config.Repository
		service := spec.Service

		if err := s.cleanExistingData(ctx, tx, repository); err != nil {
			return err
		}

		if err := s.insertProducers(ctx, tx, spec.Producers, repository, service); err != nil {
			return err
		}

		if err := s.insertConsumers(ctx, tx, spec.Consumers, repository, service); err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) cleanExistingData(ctx context.Context, tx db.SQLExecutor, repository string) error {
	_, err := db.DeleteProducersByRepository(ctx, tx, repository)
	if err != nil {
		return fmt.Errorf("failed to delete existing producers: %w", err)
	}

	_, err = db.DeleteConsumersByRepository(ctx, tx, repository)
	if err != nil {
		return fmt.Errorf("failed to delete existing consumers: %w", err)
	}

	return nil
}

func (s *Service) insertProducers(ctx context.Context, tx db.SQLExecutor, producers []eventdoctor.Producer, repository, service string) error {
	for _, p := range producers {
		owner := ""
		if p.Owner {
			owner = service
		}

		topic, err := db.GetOrCreateTopic(ctx, tx, p.Topic, owner)
		if err != nil {
			return fmt.Errorf("insertProducers: failed to get/create topic: %w", err)
		}

		for _, e := range p.Events {
			ev, err := db.GetOrCreateEvent(ctx, tx, topic.ID, e.Name, e.Version, e.SchemaURL)
			if err != nil {
				return fmt.Errorf("insertProducers: failed to get/create event: %w", err)
			}

			// Inserir headers do evento se houver
			for _, h := range e.Headers {
				header := models.EventHeader{
					EventID:     ev.ID,
					Name:        h.Name,
					Description: h.Description,
					CreatedAt:   time.Now(),
				}
				if _, err := db.InsertEventHeader(ctx, tx, header); err != nil {
					// Ignorar erro de constraint única (header já existe)
					// Seria melhor fazer um upsert, mas por enquanto ignoramos
				}
			}

			producer := models.Producer{
				EventID:    ev.ID,
				Service:    service,
				Writes:     p.Writes,
				Repository: &repository,
			}

			if _, err := db.InsertProducer(ctx, tx, producer); err != nil {
				return fmt.Errorf("insertProducers: failed to insert: %w", err)
			}
		}
	}
	return nil
}

func (s *Service) insertConsumers(ctx context.Context, tx db.SQLExecutor, consumers []eventdoctor.Consumer, repository string, service string) error {
	for _, c := range consumers {
		for _, t := range c.Topics {
			topic, err := db.GetOrCreateTopic(ctx, tx, t.Name, "")
			if err != nil {
				return fmt.Errorf("insertConsumers: failed to get/create topic: %w", err)
			}

			for _, e := range t.Events {
				ev, err := db.GetOrCreateEvent(ctx, tx, topic.ID, e.Name, nil, "")
				if err != nil {
					return fmt.Errorf("insertConsumers: failed to get/create event: %w", err)
				}

				consumer := models.Consumer{
					EventID:       ev.ID,
					Service:       service,
					ConsumerGroup: c.Group,
					EventVersion:  e.Version,
					Repository:    &repository,
				}

				if _, err := db.InsertConsumer(ctx, tx, consumer); err != nil {
					return fmt.Errorf("insertConsumers: failed to insert consumer: %w", err)
				}
			}
		}
	}
	return nil
}

func (s *Service) ListProducers(ctx context.Context) ([]eventdoctor.Producer, error) {
	// For now, return empty slice. In a real implementation,
	// this would query the database and aggregate the data
	return []eventdoctor.Producer{}, nil
}

func (s *Service) ListConsumers(ctx context.Context) ([]eventdoctor.Consumer, error) {
	// For now, return empty slice. In a real implementation,
	// this would query the database and aggregate the data
	return []eventdoctor.Consumer{}, nil
}

func (s *Service) ListEvents(ctx context.Context) ([]eventdoctor.Event, error) {
	// For now, return empty slice. In a real implementation,
	// this would query the database and aggregate the data
	return []eventdoctor.Event{}, nil
}
