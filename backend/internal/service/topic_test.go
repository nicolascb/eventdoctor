package service

import (
	"context"
	"testing"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

// --- aggregateTopicEvents tests ---

func TestAggregateTopicEvents_Empty(t *testing.T) {
	result := aggregateTopicEvents(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d items", len(result))
	}
}

func TestAggregateTopicEvents_SingleEvent(t *testing.T) {
	v := "1.0.0"
	rows := []models.EventRow{
		{
			EventID:          1,
			TopicName:        "orders",
			EventName:        "OrderCreated",
			EventDescription: "Order was created",
			SchemaVersion:    &v,
			SchemaURL:        "https://schemas.local/order.json",
		},
	}

	result := aggregateTopicEvents(rows)
	if len(result) != 1 {
		t.Fatalf("expected 1 event, got %d", len(result))
	}
	if result[0].Name != "OrderCreated" {
		t.Errorf("expected OrderCreated, got %s", result[0].Name)
	}
	if *result[0].Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %v", result[0].Version)
	}
}

func TestAggregateTopicEvents_DeduplicatesWithHeaders(t *testing.T) {
	h1 := "X-Header-A"
	h2 := "X-Header-B"
	rows := []models.EventRow{
		{EventID: 1, EventName: "Event1", HeaderName: &h1},
		{EventID: 1, EventName: "Event1", HeaderName: &h2},
		{EventID: 2, EventName: "Event2"},
	}

	result := aggregateTopicEvents(rows)
	if len(result) != 2 {
		t.Fatalf("expected 2 events, got %d", len(result))
	}
	if len(result[0].Headers) != 2 {
		t.Errorf("expected 2 headers for Event1, got %d", len(result[0].Headers))
	}
}

// --- distinctEvents tests ---

func TestDistinctEvents_Empty(t *testing.T) {
	view := &response.TopicView{}
	result := distinctEvents(view)
	if len(result) != 0 {
		t.Fatalf("expected empty set, got %d items", len(result))
	}
}

func TestDistinctEvents_FromAllSources(t *testing.T) {
	view := &response.TopicView{
		Events: []response.EventView{
			{Name: "EventA"},
			{Name: "EventB"},
		},
		Producers: []response.TopicProducerEntry{
			{Event: "EventA"},
			{Event: "EventC"},
		},
		Consumers: []response.TopicConsumerEntry{
			{Event: "EventB"},
			{Event: "EventD"},
		},
	}

	result := distinctEvents(view)
	if len(result) != 4 {
		t.Errorf("expected 4 distinct events, got %d", len(result))
	}

	for _, expected := range []string{"EventA", "EventB", "EventC", "EventD"} {
		if _, ok := result[expected]; !ok {
			t.Errorf("expected %s in distinct events", expected)
		}
	}
}

func TestDistinctEvents_Deduplicates(t *testing.T) {
	view := &response.TopicView{
		Events:    []response.EventView{{Name: "Same"}},
		Producers: []response.TopicProducerEntry{{Event: "Same"}},
		Consumers: []response.TopicConsumerEntry{{Event: "Same"}},
	}

	result := distinctEvents(view)
	if len(result) != 1 {
		t.Errorf("expected 1 distinct event, got %d", len(result))
	}
}

// --- ListTopics integration tests ---

func TestListTopics_All(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListTopics(context.Background(), 0, 0, "")
	if err != nil {
		t.Fatalf("ListTopics() error: %v", err)
	}

	if len(result.Topics) != 2 {
		t.Fatalf("expected 2 topics, got %d", len(result.Topics))
	}
	if result.Pagination != nil {
		t.Error("expected no pagination when fetching all")
	}
	if result.CountEvents == 0 {
		t.Error("expected non-zero event count")
	}
}

func TestListTopics_Paginated(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListTopics(context.Background(), 1, 1, "")
	if err != nil {
		t.Fatalf("ListTopics() error: %v", err)
	}

	if result.Pagination == nil {
		t.Fatal("expected pagination")
	}
	if result.Pagination.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Pagination.Total)
	}
	if len(result.Topics) != 1 {
		t.Errorf("expected 1 topic on page 1, got %d", len(result.Topics))
	}
}

func TestListTopics_Search(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListTopics(context.Background(), 1, 10, "orders")
	if err != nil {
		t.Fatalf("ListTopics() error: %v", err)
	}

	if len(result.Topics) != 1 {
		t.Fatalf("expected 1 topic matching 'orders', got %d", len(result.Topics))
	}
	if result.Topics[0].Topic != "orders.events" {
		t.Errorf("expected orders.events, got %s", result.Topics[0].Topic)
	}
}

func TestListTopics_CountsUnconsumedAndOrphaned(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListTopics(context.Background(), 0, 0, "")
	if err != nil {
		t.Fatalf("ListTopics() error: %v", err)
	}

	// OrderUpdated has no consumer → unconsumed
	if result.CountUnconsumed == 0 {
		t.Error("expected at least 1 unconsumed event (OrderUpdated)")
	}
}

// --- GetTopicView integration tests ---

func TestGetTopicView(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.GetTopicView(context.Background(), "orders.events")
	if err != nil {
		t.Fatalf("GetTopicView() error: %v", err)
	}

	if result.Topic != "orders.events" {
		t.Errorf("expected topic orders.events, got %s", result.Topic)
	}
	if len(result.Events) != 2 {
		t.Errorf("expected 2 events (OrderCreated, OrderUpdated), got %d", len(result.Events))
	}
	if len(result.Producers) == 0 {
		t.Error("expected at least 1 producer")
	}
	if len(result.Consumers) == 0 {
		t.Error("expected at least 1 consumer")
	}
	if result.OwnerService == nil {
		t.Error("expected owner service to be set")
	} else if *result.OwnerService != "service-a" {
		t.Errorf("expected owner service-a, got %s", *result.OwnerService)
	}
}

func TestGetTopicView_EmptyTopic(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)

	result, err := svc.GetTopicView(context.Background(), "nonexistent.topic")
	if err != nil {
		t.Fatalf("GetTopicView() error: %v", err)
	}

	if result.Topic != "nonexistent.topic" {
		t.Errorf("expected topic name, got %s", result.Topic)
	}
	if len(result.Events) != 0 {
		t.Errorf("expected 0 events, got %d", len(result.Events))
	}
}
