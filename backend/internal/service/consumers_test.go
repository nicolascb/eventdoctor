package service

import (
	"context"
	"testing"

	"github.com/nicolascb/eventdoctor/internal/db/models"
)

// --- aggregateConsumers tests ---

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

// --- ListConsumers integration tests ---

func TestListConsumers_All(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListConsumers(context.Background(), 0, 0, "")
	if err != nil {
		t.Fatalf("ListConsumers() error: %v", err)
	}

	if result.Pagination != nil {
		t.Error("expected no pagination when fetching all")
	}
	if len(result.Consumers) != 2 {
		t.Fatalf("expected 2 consumers, got %d", len(result.Consumers))
	}
}

func TestListConsumers_Paginated(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListConsumers(context.Background(), 1, 1, "")
	if err != nil {
		t.Fatalf("ListConsumers() error: %v", err)
	}

	if result.Pagination == nil {
		t.Fatal("expected pagination")
	}
	if result.Pagination.Page != 1 {
		t.Errorf("expected page 1, got %d", result.Pagination.Page)
	}
	if result.Pagination.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Pagination.Total)
	}
	if result.Pagination.TotalPages != 2 {
		t.Errorf("expected 2 total pages, got %d", result.Pagination.TotalPages)
	}
}

func TestListConsumers_Search(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListConsumers(context.Background(), 1, 10, "orders")
	if err != nil {
		t.Fatalf("ListConsumers() error: %v", err)
	}

	if len(result.Consumers) == 0 {
		t.Fatal("expected at least 1 consumer matching 'orders'")
	}
	for _, c := range result.Consumers {
		if c.Group != "service-b-orders" {
			t.Errorf("unexpected consumer group in search results: %s", c.Group)
		}
	}
}

func TestListConsumers_EmptyDB(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)

	result, err := svc.ListConsumers(context.Background(), 0, 0, "")
	if err != nil {
		t.Fatalf("ListConsumers() error: %v", err)
	}

	if len(result.Consumers) != 0 {
		t.Errorf("expected 0 consumers for empty DB, got %d", len(result.Consumers))
	}
}

// --- ListUndocumentedConsumers integration tests ---

func TestListUndocumentedConsumers_Empty(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)

	result, err := svc.ListUndocumentedConsumers(context.Background())
	if err != nil {
		t.Fatalf("ListUndocumentedConsumers() error: %v", err)
	}

	if len(result.GroupsUndocumented) != 0 {
		t.Errorf("expected 0 undocumented groups, got %d", len(result.GroupsUndocumented))
	}
}
