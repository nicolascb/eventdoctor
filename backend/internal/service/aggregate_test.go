package service

import (
	"testing"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

// orderedMap tests

func TestOrderedMap_GetOrCreate_NewKey(t *testing.T) {
	m := newOrderedMap[string, int]()
	val := m.getOrCreate("a", func() int { return 42 })
	if *val != 42 {
		t.Fatalf("expected 42, got %d", *val)
	}
}

func TestOrderedMap_GetOrCreate_ExistingKey(t *testing.T) {
	m := newOrderedMap[string, int]()
	m.getOrCreate("a", func() int { return 1 })
	val := m.getOrCreate("a", func() int { return 99 })
	if *val != 1 {
		t.Fatalf("expected existing value 1, got %d", *val)
	}
}

func TestOrderedMap_Collect_PreservesOrder(t *testing.T) {
	m := newOrderedMap[string, int]()
	m.getOrCreate("c", func() int { return 3 })
	m.getOrCreate("a", func() int { return 1 })
	m.getOrCreate("b", func() int { return 2 })

	result := m.collect()
	if len(result) != 3 {
		t.Fatalf("expected 3 items, got %d", len(result))
	}
	if result[0] != 3 || result[1] != 1 || result[2] != 2 {
		t.Fatalf("unexpected order: %v", result)
	}
}

func TestOrderedMap_Collect_Empty(t *testing.T) {
	m := newOrderedMap[string, int]()
	result := m.collect()
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %v", result)
	}
}

func TestOrderedMap_MutateViaPointer(t *testing.T) {
	m := newOrderedMap[string, int]()
	p := m.getOrCreate("x", func() int { return 10 })
	*p = 20
	result := m.collect()
	if result[0] != 20 {
		t.Fatalf("expected 20 after mutation, got %d", result[0])
	}
}

// findOrAppend tests

func TestFindOrAppend_AppendsWhenNotFound(t *testing.T) {
	s := []int{}
	p := findOrAppend(&s, func(v *int) bool { return *v == 5 }, func() int { return 5 })
	if *p != 5 {
		t.Fatalf("expected 5, got %d", *p)
	}
	if len(s) != 1 {
		t.Fatalf("expected slice length 1, got %d", len(s))
	}
}

func TestFindOrAppend_ReturnsExisting(t *testing.T) {
	s := []int{3, 7}
	p := findOrAppend(&s, func(v *int) bool { return *v == 7 }, func() int { return 99 })
	if *p != 7 {
		t.Fatalf("expected 7, got %d", *p)
	}
	if len(s) != 2 {
		t.Fatalf("slice should not grow, got length %d", len(s))
	}
}

func TestFindOrAppend_MutateExisting(t *testing.T) {
	s := []int{1, 2, 3}
	p := findOrAppend(&s, func(v *int) bool { return *v == 2 }, func() int { return 99 })
	*p = 42
	if s[1] != 42 {
		t.Fatalf("expected s[1]=42, got %d", s[1])
	}
}

// appendHeader tests

func TestAppendHeader_NilName(t *testing.T) {
	headers := []response.EventHeaderView{}
	appendHeader(&headers, nil, nil)
	if len(headers) != 0 {
		t.Fatal("expected no header appended when name is nil")
	}
}

func TestAppendHeader_WithName(t *testing.T) {
	headers := []response.EventHeaderView{}
	name := "X-Request-ID"
	desc := "unique request identifier"
	appendHeader(&headers, &name, &desc)
	if len(headers) != 1 {
		t.Fatalf("expected 1 header, got %d", len(headers))
	}
	if headers[0].Name != name || headers[0].Description != desc {
		t.Fatalf("unexpected header: %+v", headers[0])
	}
}

func TestAppendHeader_NilDescription(t *testing.T) {
	headers := []response.EventHeaderView{}
	name := "Authorization"
	appendHeader(&headers, &name, nil)
	if len(headers) != 1 {
		t.Fatalf("expected 1 header, got %d", len(headers))
	}
	if headers[0].Description != "" {
		t.Fatalf("expected empty description, got %q", headers[0].Description)
	}
}

func strPtr(s string) *string { return &s }

// aggregateConsumers tests

func TestAggregateConsumers_Empty(t *testing.T) {
	result := aggregateConsumers(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}

func TestAggregateConsumers_SingleRow(t *testing.T) {
	rows := []models.ConsumerRow{
		{
			ServiceName:   "svc-a",
			Repository:    "https://github.com/org/svc-a",
			ConsumerGroup: "svc-a-group",
			Description:   "handles events",
			TopicName:     "topic.events",
			EventName:     "EventFired",
			EventVersion:  strPtr("1.0.0"),
		},
	}

	result := aggregateConsumers(rows)
	if len(result) != 1 {
		t.Fatalf("expected 1 consumer, got %d", len(result))
	}
	cv := result[0]
	if cv.Service != "svc-a" || cv.Group != "svc-a-group" {
		t.Fatalf("unexpected consumer view: %+v", cv)
	}
	if len(cv.Topics) != 1 || len(cv.Topics[0].Events) != 1 {
		t.Fatalf("unexpected topics/events structure: %+v", cv.Topics)
	}
}

func TestAggregateConsumers_MultipleTopics(t *testing.T) {
	rows := []models.ConsumerRow{
		{ServiceName: "svc", Repository: "repo", ConsumerGroup: "grp", TopicName: "topic-a", EventName: "EventA"},
		{ServiceName: "svc", Repository: "repo", ConsumerGroup: "grp", TopicName: "topic-b", EventName: "EventB"},
	}

	result := aggregateConsumers(rows)
	if len(result) != 1 {
		t.Fatalf("expected 1 consumer entry, got %d", len(result))
	}
	if len(result[0].Topics) != 2 {
		t.Fatalf("expected 2 topics, got %d", len(result[0].Topics))
	}
}

func TestAggregateConsumers_MultipleGroups(t *testing.T) {
	rows := []models.ConsumerRow{
		{ServiceName: "svc", Repository: "repo", ConsumerGroup: "grp-a", TopicName: "topic", EventName: "Event"},
		{ServiceName: "svc", Repository: "repo", ConsumerGroup: "grp-b", TopicName: "topic", EventName: "Event"},
	}

	result := aggregateConsumers(rows)
	if len(result) != 2 {
		t.Fatalf("expected 2 consumer entries (different groups), got %d", len(result))
	}
}

// aggregateEvents tests

func TestAggregateEvents_Empty(t *testing.T) {
	result := aggregateEvents(nil, nil, nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}

func TestAggregateEvents_SingleRow(t *testing.T) {
	rows := []models.EventRow{
		{
			TopicName:        "topic.events",
			EventName:        "EventFired",
			EventDescription: "something fired",
			SchemaVersion:    strPtr("1.0.0"),
			SchemaURL:        "https://schemas.local/event.json",
		},
	}

	result := aggregateEvents(rows, nil, nil)
	if len(result) != 1 {
		t.Fatalf("expected 1 event, got %d", len(result))
	}
	if result[0].Topic != "topic.events" {
		t.Fatalf("unexpected topic name: %s", result[0].Topic)
	}
	if result[0].Name != "EventFired" {
		t.Fatalf("unexpected event name: %s", result[0].Name)
	}
}

func TestAggregateEvents_MultipleTopics(t *testing.T) {
	rows := []models.EventRow{
		{TopicName: "topic-a", EventName: "EventA"},
		{TopicName: "topic-b", EventName: "EventB"},
		{TopicName: "topic-a", EventName: "EventC"},
	}

	result := aggregateEvents(rows, nil, nil)
	if len(result) != 3 {
		t.Fatalf("expected 3 events, got %d", len(result))
	}
	// Sort order is Topic then EventName
	// Expect: topic-a/EventA, topic-a/EventC, topic-b/EventB
	if result[0].Topic != "topic-a" || result[0].Name != "EventA" {
		t.Errorf("unexpected order at 0: %v", result[0])
	}
	if result[1].Topic != "topic-a" || result[1].Name != "EventC" {
		t.Errorf("unexpected order at 1: %v", result[1])
	}
	if result[2].Topic != "topic-b" || result[2].Name != "EventB" {
		t.Errorf("unexpected order at 2: %v", result[2])
	}
}

func TestAggregateEvents_HeadersDeduplicatedPerEvent(t *testing.T) {
	h1 := "X-Header-A"
	h2 := "X-Header-B"
	rows := []models.EventRow{
		{TopicName: "topic", EventName: "Event", HeaderName: &h1},
		{TopicName: "topic", EventName: "Event", HeaderName: &h2},
	}

	result := aggregateEvents(rows, nil, nil)
	if len(result[0].Headers) != 2 {
		t.Fatalf("expected 2 headers, got %d", len(result[0].Headers))
	}
}

func TestAggregateEvents_WithProducersAndConsumers(t *testing.T) {
	eventRows := []models.EventRow{
		{TopicName: "topic", EventName: "Event"},
	}
	producerRows := []models.ProducerRow{
		{TopicName: "topic", EventName: "Event", ServiceName: "svc-prod", Repository: "repo-prod", Owner: true},
	}
	consumerRows := []models.ConsumerRow{
		{TopicName: "topic", EventName: "Event", ServiceName: "svc-cons", Repository: "repo-cons", ConsumerGroup: "grp"},
	}

	result := aggregateEvents(eventRows, producerRows, consumerRows)
	if len(result) != 1 {
		t.Fatal("expected 1 event")
	}
	ev := result[0]
	if len(ev.Producers) != 1 {
		t.Fatal("expected 1 producer")
	}
	if ev.Producers[0].Service != "svc-prod" {
		t.Errorf("unexpected producer: %v", ev.Producers[0])
	}
	if len(ev.Consumers) != 1 {
		t.Fatal("expected 1 consumer")
	}
	if ev.Consumers[0].Service != "svc-cons" {
		t.Errorf("unexpected consumer: %v", ev.Consumers[0])
	}
}

// aggregateOverview tests

func TestAggregateOverview_Empty(t *testing.T) {
	result := aggregateOverview(nil, nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}

func TestAggregateOverview_ProducerAndConsumer(t *testing.T) {
	producerRows := []models.ProducerRow{
		{ServiceName: "svc-a", Repository: "repo-a", TopicName: "topic", EventName: "Event", Owner: true, Writes: true},
	}
	consumerRows := []models.ConsumerRow{
		{ServiceName: "svc-b", Repository: "repo-b", ConsumerGroup: "grp", TopicName: "topic", EventName: "Event"},
	}

	result := aggregateOverview(producerRows, consumerRows)
	if len(result) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(result))
	}
	topic := result[0]
	if len(topic.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(topic.Events))
	}
	event := topic.Events[0]
	if len(event.Producers) != 1 {
		t.Fatalf("expected 1 producer, got %d", len(event.Producers))
	}
	if len(event.Consumers) != 1 {
		t.Fatalf("expected 1 consumer, got %d", len(event.Consumers))
	}
}

func TestAggregateOverview_ConsumerOnlyTopic(t *testing.T) {
	consumerRows := []models.ConsumerRow{
		{ServiceName: "svc-b", Repository: "repo-b", ConsumerGroup: "grp", TopicName: "orphan-topic", EventName: "Event"},
	}

	result := aggregateOverview(nil, consumerRows)
	if len(result) != 1 {
		t.Fatalf("expected 1 topic from consumer-only data, got %d", len(result))
	}
}
