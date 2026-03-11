package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

// --- aggregateEvents tests ---

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
		{EventID: 1, TopicName: "topic-a", EventName: "EventA"},
		{EventID: 2, TopicName: "topic-b", EventName: "EventB"},
		{EventID: 3, TopicName: "topic-a", EventName: "EventC"},
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

// --- processEventRows tests ---

func TestProcessEventRows(t *testing.T) {
	events := newOrderedMap[int64, response.EventView]()
	h := "X-Trace"
	rows := []models.EventRow{
		{EventID: 1, TopicName: "topic", EventName: "Event1", EventDescription: "desc1", HeaderName: &h},
		{EventID: 2, TopicName: "topic", EventName: "Event2"},
	}

	processEventRows(&events, rows, func(id int64) int64 { return id })

	result := events.collect()
	if len(result) != 2 {
		t.Fatalf("expected 2 events, got %d", len(result))
	}
	if result[0].Name != "Event1" {
		t.Errorf("expected Event1, got %s", result[0].Name)
	}
	if len(result[0].Headers) != 1 {
		t.Errorf("expected 1 header on Event1, got %d", len(result[0].Headers))
	}
	// Initialized slices should not be nil
	if result[1].Producers == nil {
		t.Error("expected non-nil Producers slice")
	}
	if result[1].Consumers == nil {
		t.Error("expected non-nil Consumers slice")
	}
}

// --- processProducerRows tests ---

func TestProcessProducerRows_AddsToExistingEvents(t *testing.T) {
	events := newOrderedMap[int64, response.EventView]()
	events.getOrCreate(int64(1), func() response.EventView {
		return response.EventView{
			ID: 1, Name: "Event1",
			Producers: []response.EventProducer{},
		}
	})

	producerRows := []models.ProducerRow{
		{EventID: 1, ServiceName: "svc-a", Repository: "repo-a", Owner: true},
	}

	processProducerRows(&events, producerRows, func(id int64) int64 { return id })

	result := events.collect()
	if len(result[0].Producers) != 1 {
		t.Fatalf("expected 1 producer, got %d", len(result[0].Producers))
	}
	if result[0].Producers[0].Service != "svc-a" {
		t.Errorf("expected svc-a, got %s", result[0].Producers[0].Service)
	}
	if result[0].Producers[0].Owner != fmt.Sprintf("%v", true) {
		t.Errorf("expected owner true, got %s", result[0].Producers[0].Owner)
	}
}

func TestProcessProducerRows_IgnoresMissingEvents(t *testing.T) {
	events := newOrderedMap[int64, response.EventView]()

	producerRows := []models.ProducerRow{
		{EventID: 999, ServiceName: "svc-a"},
	}

	processProducerRows(&events, producerRows, func(id int64) int64 { return id })

	result := events.collect()
	if len(result) != 0 {
		t.Errorf("expected 0 events (producer for missing event should be ignored), got %d", len(result))
	}
}

// --- processConsumerRows tests ---

func TestProcessConsumerRows_AddsToExistingEvents(t *testing.T) {
	events := newOrderedMap[int64, response.EventView]()
	events.getOrCreate(int64(1), func() response.EventView {
		return response.EventView{
			ID: 1, Name: "Event1",
			Consumers: []response.EventConsumer{},
		}
	})

	consumerRows := []models.ConsumerRow{
		{EventID: 1, ServiceName: "svc-b", Repository: "repo-b", ConsumerGroup: "grp"},
	}

	processConsumerRows(&events, consumerRows, func(id int64) int64 { return id })

	result := events.collect()
	if len(result[0].Consumers) != 1 {
		t.Fatalf("expected 1 consumer, got %d", len(result[0].Consumers))
	}
	if result[0].Consumers[0].Service != "svc-b" {
		t.Errorf("expected svc-b, got %s", result[0].Consumers[0].Service)
	}
}

func TestProcessConsumerRows_IgnoresMissingEvents(t *testing.T) {
	events := newOrderedMap[int64, response.EventView]()

	consumerRows := []models.ConsumerRow{
		{EventID: 999, ServiceName: "svc-b", ConsumerGroup: "grp"},
	}

	processConsumerRows(&events, consumerRows, func(id int64) int64 { return id })

	result := events.collect()
	if len(result) != 0 {
		t.Errorf("expected 0 events, got %d", len(result))
	}
}

// --- ListEvents integration tests ---

func TestListEvents_All(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListEvents(context.Background(), 0, 0, "")
	if err != nil {
		t.Fatalf("ListEvents() error: %v", err)
	}

	if result.Pagination != nil {
		t.Error("expected no pagination when fetching all")
	}
	if len(result.Events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(result.Events))
	}
}

func TestListEvents_Paginated(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListEvents(context.Background(), 1, 2, "")
	if err != nil {
		t.Fatalf("ListEvents() error: %v", err)
	}

	if result.Pagination == nil {
		t.Fatal("expected pagination")
	}
	if result.Pagination.Total != 3 {
		t.Errorf("expected total 3, got %d", result.Pagination.Total)
	}
	if len(result.Events) != 2 {
		t.Errorf("expected 2 events on page 1, got %d", len(result.Events))
	}
}

func TestListEvents_Search(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListEvents(context.Background(), 1, 10, "Payment")
	if err != nil {
		t.Fatalf("ListEvents() error: %v", err)
	}

	if len(result.Events) == 0 {
		t.Fatal("expected at least 1 event matching 'Payment'")
	}
	for _, e := range result.Events {
		if e.Name != "PaymentProcessed" {
			t.Errorf("unexpected event in search results: %s", e.Name)
		}
	}
}

func TestListEvents_EventsHaveProducersAndConsumers(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListEvents(context.Background(), 0, 0, "")
	if err != nil {
		t.Fatalf("ListEvents() error: %v", err)
	}

	// Find OrderCreated - it should have 1 producer (service-a) and 1 consumer (service-b)
	for _, e := range result.Events {
		if e.Name == "OrderCreated" {
			if len(e.Producers) != 1 {
				t.Errorf("OrderCreated: expected 1 producer, got %d", len(e.Producers))
			}
			if len(e.Consumers) != 1 {
				t.Errorf("OrderCreated: expected 1 consumer, got %d", len(e.Consumers))
			}
			if len(e.Headers) != 1 {
				t.Errorf("OrderCreated: expected 1 header, got %d", len(e.Headers))
			}
			return
		}
	}
	t.Error("OrderCreated event not found")
}

func TestListEvents_EmptyDB(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)

	result, err := svc.ListEvents(context.Background(), 0, 0, "")
	if err != nil {
		t.Fatalf("ListEvents() error: %v", err)
	}

	if len(result.Events) != 0 {
		t.Errorf("expected 0 events for empty DB, got %d", len(result.Events))
	}
}

// --- GetEvent integration tests ---

func TestGetEvent(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	// First, list events to get an ID
	events, err := svc.ListEvents(context.Background(), 0, 0, "")
	if err != nil {
		t.Fatalf("ListEvents() error: %v", err)
	}

	var orderCreatedID int64
	for _, e := range events.Events {
		if e.Name == "OrderCreated" {
			orderCreatedID = e.ID
			break
		}
	}
	if orderCreatedID == 0 {
		t.Fatal("OrderCreated event not found")
	}

	result, err := svc.GetEvent(context.Background(), orderCreatedID)
	if err != nil {
		t.Fatalf("GetEvent() error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Name != "OrderCreated" {
		t.Errorf("expected OrderCreated, got %s", result.Name)
	}
	if result.Topic != "orders.events" {
		t.Errorf("expected topic orders.events, got %s", result.Topic)
	}
	if len(result.Producers) != 1 {
		t.Errorf("expected 1 producer, got %d", len(result.Producers))
	}
	if len(result.Consumers) != 1 {
		t.Errorf("expected 1 consumer, got %d", len(result.Consumers))
	}
	if len(result.Headers) != 1 {
		t.Errorf("expected 1 header, got %d", len(result.Headers))
	}
}

func TestGetEvent_NotFound(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)

	result, err := svc.GetEvent(context.Background(), 99999)
	if err != nil {
		t.Fatalf("GetEvent() error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil for non-existent event, got %+v", result)
	}
}
