package service

import (
	"context"
	"testing"

	"github.com/nicolascb/eventdoctor/internal/db/models"
)

// --- aggregateProducers tests ---

func TestAggregateProducers_Empty(t *testing.T) {
	result := aggregateProducers(nil, nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}

func TestAggregateProducers_SingleService(t *testing.T) {
	services := []models.Service{
		{ID: 1, Name: "svc-a", Repository: "repo-a"},
	}
	topicRows := []models.ProducerListRow{
		{ServiceID: 1, ServiceName: "svc-a", Repository: "repo-a", TopicID: 10, TopicName: "topic-a", EventCount: 5, Owner: true, Writes: true},
	}

	result := aggregateProducers(services, topicRows)
	if len(result) != 1 {
		t.Fatalf("expected 1 producer, got %d", len(result))
	}
	p := result[0]
	if p.Service != "svc-a" {
		t.Errorf("unexpected service name: %s", p.Service)
	}
	if len(p.Topics) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(p.Topics))
	}
	if p.Topics[0].Topic != "topic-a" {
		t.Errorf("unexpected topic name: %s", p.Topics[0].Topic)
	}
	if p.Topics[0].EventCount != 5 {
		t.Errorf("unexpected event count: %d", p.Topics[0].EventCount)
	}
}

func TestAggregateProducers_MultipleServices(t *testing.T) {
	services := []models.Service{
		{ID: 1, Name: "svc-a"},
		{ID: 2, Name: "svc-b"},
	}
	topicRows := []models.ProducerListRow{
		{ServiceID: 1, TopicName: "topic-a"},
		{ServiceID: 2, TopicName: "topic-b"},
		{ServiceID: 1, TopicName: "topic-c"},
	}

	result := aggregateProducers(services, topicRows)
	if len(result) != 2 {
		t.Fatalf("expected 2 producers, got %d", len(result))
	}

	// Result order should match services input order
	if result[0].Service != "svc-a" {
		t.Errorf("expected svc-a first, got %s", result[0].Service)
	}
	if len(result[0].Topics) != 2 {
		t.Errorf("expected 2 topics for svc-a, got %d", len(result[0].Topics))
	}

	if result[1].Service != "svc-b" {
		t.Errorf("expected svc-b second, got %s", result[1].Service)
	}
	if len(result[1].Topics) != 1 {
		t.Errorf("expected 1 topic for svc-b, got %d", len(result[1].Topics))
	}
}

// --- ListProducers integration tests ---

func TestListProducers_Paginated(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListProducers(context.Background(), 1, 10, "")
	if err != nil {
		t.Fatalf("ListProducers() error: %v", err)
	}

	if result.Pagination == nil {
		t.Fatal("expected pagination")
	}
	if len(result.Producers) != 2 {
		t.Fatalf("expected 2 producers, got %d", len(result.Producers))
	}

	// Each producer should have at least 1 topic
	for _, p := range result.Producers {
		if len(p.Topics) == 0 {
			t.Errorf("producer %s has no topics", p.Service)
		}
	}
}

func TestListProducers_Search(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.ListProducers(context.Background(), 1, 10, "service-a")
	if err != nil {
		t.Fatalf("ListProducers() error: %v", err)
	}

	if len(result.Producers) == 0 {
		t.Fatal("expected at least 1 producer matching 'service-a'")
	}
	for _, p := range result.Producers {
		if p.Service != "service-a" {
			t.Errorf("unexpected producer in search results: %s", p.Service)
		}
	}
}

func TestListProducers_RequiresPagination(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)

	_, err := svc.ListProducers(context.Background(), 0, 0, "")
	if err == nil {
		t.Fatal("expected error when page and pageSize are 0")
	}
}

func TestListProducers_EmptyDB(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)

	result, err := svc.ListProducers(context.Background(), 1, 10, "")
	if err != nil {
		t.Fatalf("ListProducers() error: %v", err)
	}

	if len(result.Producers) != 0 {
		t.Errorf("expected 0 producers for empty DB, got %d", len(result.Producers))
	}
}

// --- GetProducerDetail integration tests ---

func TestGetProducerDetail(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	// Get the producer list to find valid IDs
	producers, err := svc.ListProducers(context.Background(), 1, 10, "")
	if err != nil {
		t.Fatalf("ListProducers() error: %v", err)
	}

	var serviceID, topicID int64
	for _, p := range producers.Producers {
		if p.Service == "service-a" && len(p.Topics) > 0 {
			serviceID = p.ServiceID
			topicID = p.Topics[0].TopicID
			break
		}
	}
	if serviceID == 0 {
		t.Fatal("service-a not found in producers")
	}

	result, err := svc.GetProducerDetail(context.Background(), serviceID, topicID)
	if err != nil {
		t.Fatalf("GetProducerDetail() error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Service != "service-a" {
		t.Errorf("expected service-a, got %s", result.Service)
	}
	if len(result.Events) == 0 {
		t.Error("expected at least 1 event")
	}
}

func TestGetProducerDetail_NotFound(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)

	result, err := svc.GetProducerDetail(context.Background(), 99999, 99999)
	if err != nil {
		t.Fatalf("GetProducerDetail() error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil for non-existent producer, got %+v", result)
	}
}
