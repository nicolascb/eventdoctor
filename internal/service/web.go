package service

import (
	"context"
	"fmt"

	"github.com/nicolascb/eventdoctor/internal/db"
	"github.com/nicolascb/eventdoctor/internal/db/models"
	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
)

func (s *Service) WebDocs(ctx context.Context) (eventdoctor.WebData, error) {
	topicsData, err := s.getTopicsData(ctx)
	if err != nil {
		return eventdoctor.WebData{}, err
	}

	data := eventdoctor.WebData{}
	data.Topics = topicsData
	return data, nil
}

func (s *Service) getTopicsData(ctx context.Context) ([]eventdoctor.TopicView, error) {
	topics, err := db.ListTopics(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("error listing topics: %w", err)
	}

	var data []eventdoctor.TopicView
	for _, t := range topics {
		topicView, err := s.processTopicEvents(ctx, t)
		if err != nil {
			return nil, err
		}
		data = append(data, topicView)
	}

	return data, nil
}

func (s *Service) processTopicEvents(ctx context.Context, t models.Topic) (eventdoctor.TopicView, error) {
	tv := eventdoctor.TopicView{Name: t.Name, Owner: t.Owner}

	evs, err := db.GetEventsByTopic(ctx, s.db, t.ID)
	if err != nil {
		return eventdoctor.TopicView{}, fmt.Errorf("error listing events by topic: %w", err)
	}

	for _, ev := range evs {
		eventView, err := s.processEvent(ctx, ev)
		if err != nil {
			return eventdoctor.TopicView{}, err
		}
		tv.Events = append(tv.Events, eventView)
	}

	return tv, nil
}

func (s *Service) processEvent(ctx context.Context, ev models.Event) (eventdoctor.EventView, error) {
	prods, err := db.GetProducersByEvent(ctx, s.db, ev.ID)
	if err != nil {
		return eventdoctor.EventView{}, fmt.Errorf("failed to list producers: %w", err)
	}

	cons, err := db.GetConsumersByEvent(ctx, s.db, ev.ID)
	if err != nil {
		return eventdoctor.EventView{}, fmt.Errorf("failed to list consumers: %w", err)
	}

	headers, err := db.GetEventHeaders(ctx, s.db, ev.ID)
	if err != nil {
		return eventdoctor.EventView{}, fmt.Errorf("failed to list headers: %w", err)
	}

	evv := eventdoctor.EventView{
		Type:       ev.EventType,
		Version:    ev.SchemaVersion,
		SchemaURL:  ev.SchemaURL,
		Deprecated: ev.Deprecated,
	}

	// Adicionar headers
	for _, h := range headers {
		evv.Headers = append(evv.Headers, eventdoctor.HeaderView{
			Name:        h.Name,
			Description: h.Description,
		})
	}

	// Adicionar produtores
	for _, p := range prods {
		evv.Producers = append(evv.Producers, eventdoctor.ProducerView{
			Name:       p.Service,
			Writes:     p.Writes,
			Repository: p.Repository,
		})
	}

	// Adicionar consumidores
	for _, cn := range cons {
		evv.Consumers = append(evv.Consumers, eventdoctor.ConsumerView{
			Name:         cn.Service,
			EventVersion: cn.EventVersion,
			Repository:   cn.Repository,
		})
	}

	return evv, nil
}
