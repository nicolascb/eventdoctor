package service

import (
	"context"
	"testing"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

// --- calculateOverviewMetrics tests ---

func TestCalculateOverviewMetrics_Empty(t *testing.T) {
	result := calculateOverviewMetrics(nil)

	if result.TotalTopics != 0 || result.TotalEvents != 0 || result.TotalProducers != 0 || result.TotalConsumers != 0 {
		t.Errorf("expected all zeros, got %+v", result)
	}
}

func TestCalculateOverviewMetrics_SingleTopic(t *testing.T) {
	topics := []response.OverviewTopicView{
		{
			Name: "topic-a",
			Events: []response.OverviewEventView{
				{
					Name: "EventA",
					Producers: []response.OverviewProducerView{
						{Service: "svc-a", Repository: "repo-a"},
					},
					Consumers: []response.OverviewConsumerView{
						{Service: "svc-b", Repository: "repo-b", Group: "grp-b"},
					},
				},
				{
					Name: "EventB",
					Producers: []response.OverviewProducerView{
						{Service: "svc-a", Repository: "repo-a"},
					},
				},
			},
		},
	}

	result := calculateOverviewMetrics(topics)

	if result.TotalTopics != 1 {
		t.Errorf("expected 1 topic, got %d", result.TotalTopics)
	}
	if result.TotalEvents != 2 {
		t.Errorf("expected 2 events, got %d", result.TotalEvents)
	}
	if result.TotalProducers != 1 {
		t.Errorf("expected 1 unique producer, got %d", result.TotalProducers)
	}
	if result.TotalConsumers != 1 {
		t.Errorf("expected 1 unique consumer, got %d", result.TotalConsumers)
	}
}

func TestCalculateOverviewMetrics_MultipleTopicsDeduplicateProducers(t *testing.T) {
	topics := []response.OverviewTopicView{
		{
			Name: "topic-a",
			Events: []response.OverviewEventView{
				{
					Name: "EventA",
					Producers: []response.OverviewProducerView{
						{Service: "svc-shared", Repository: "repo"},
					},
				},
			},
		},
		{
			Name: "topic-b",
			Events: []response.OverviewEventView{
				{
					Name: "EventB",
					Producers: []response.OverviewProducerView{
						{Service: "svc-shared", Repository: "repo"},
					},
				},
			},
		},
	}

	result := calculateOverviewMetrics(topics)

	if result.TotalProducers != 1 {
		t.Errorf("expected 1 unique producer (deduplicated), got %d", result.TotalProducers)
	}
}

// --- aggregateOverview tests ---

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

// --- processOverviewProducerRows tests ---

func TestProcessOverviewProducerRows(t *testing.T) {
	topics := newOrderedMap[string, response.OverviewTopicView]()
	h := "X-Trace"
	rows := []models.ProducerRow{
		{
			ServiceName: "svc-a",
			Repository:  "repo",
			TopicName:   "topic-a",
			EventName:   "EventA",
			Owner:       true,
			Writes:      true,
			HeaderName:  &h,
		},
	}

	processOverviewProducerRows(&topics, rows)

	result := topics.collect()
	if len(result) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(result))
	}
	if result[0].Name != "topic-a" {
		t.Errorf("expected topic-a, got %s", result[0].Name)
	}
	if len(result[0].Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(result[0].Events))
	}
	if len(result[0].Events[0].Producers) != 1 {
		t.Errorf("expected 1 producer, got %d", len(result[0].Events[0].Producers))
	}
	if len(result[0].Events[0].Headers) != 1 {
		t.Errorf("expected 1 header, got %d", len(result[0].Events[0].Headers))
	}
}

// --- processOverviewConsumerRows tests ---

func TestProcessOverviewConsumerRows(t *testing.T) {
	topics := newOrderedMap[string, response.OverviewTopicView]()
	rows := []models.ConsumerRow{
		{
			ServiceName:   "svc-b",
			Repository:    "repo",
			ConsumerGroup: "grp",
			TopicName:     "topic-a",
			EventName:     "EventA",
		},
	}

	processOverviewConsumerRows(&topics, rows)

	result := topics.collect()
	if len(result) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(result))
	}
	if len(result[0].Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(result[0].Events))
	}
	if len(result[0].Events[0].Consumers) != 1 {
		t.Errorf("expected 1 consumer, got %d", len(result[0].Events[0].Consumers))
	}
}

func TestProcessOverviewConsumerRows_DeduplicatesConsumers(t *testing.T) {
	topics := newOrderedMap[string, response.OverviewTopicView]()
	rows := []models.ConsumerRow{
		{ServiceName: "svc-b", Repository: "repo", ConsumerGroup: "grp", TopicName: "topic", EventName: "Event"},
		{ServiceName: "svc-b", Repository: "repo", ConsumerGroup: "grp", TopicName: "topic", EventName: "Event"},
	}

	processOverviewConsumerRows(&topics, rows)

	result := topics.collect()
	if len(result[0].Events[0].Consumers) != 1 {
		t.Errorf("expected 1 deduplicated consumer, got %d", len(result[0].Events[0].Consumers))
	}
}

// --- trackProducers / trackConsumers tests ---

func TestTrackProducers(t *testing.T) {
	m := make(map[string]struct{})
	trackProducers(m, []response.OverviewProducerView{
		{Service: "svc-a", Repository: "repo-a"},
		{Service: "svc-b", Repository: "repo-b"},
		{Service: "svc-a", Repository: "repo-a"}, // duplicate
	})

	if len(m) != 2 {
		t.Errorf("expected 2 unique producer keys, got %d", len(m))
	}
}

func TestTrackConsumers(t *testing.T) {
	m := make(map[string]struct{})
	trackConsumers(m, []response.OverviewConsumerView{
		{Service: "svc-a", Repository: "repo-a", Group: "grp-1"},
		{Service: "svc-a", Repository: "repo-a", Group: "grp-2"},
		{Service: "svc-a", Repository: "repo-a", Group: "grp-1"}, // duplicate
	})

	if len(m) != 2 {
		t.Errorf("expected 2 unique consumer keys, got %d", len(m))
	}
}

// --- Overview integration tests ---

func TestOverview(t *testing.T) {
	testDB := setupTestDB(t)
	seedTestData(t, testDB)
	svc := NewService(testDB)

	result, err := svc.Overview(context.Background())
	if err != nil {
		t.Fatalf("Overview() error: %v", err)
	}

	if result.TotalTopics != 2 {
		t.Errorf("expected 2 topics, got %d", result.TotalTopics)
	}
	if result.TotalEvents != 3 {
		t.Errorf("expected 3 events, got %d", result.TotalEvents)
	}
	if result.TotalProducers != 2 {
		t.Errorf("expected 2 producers, got %d", result.TotalProducers)
	}
	if result.TotalConsumers != 2 {
		t.Errorf("expected 2 consumers, got %d", result.TotalConsumers)
	}
}

func TestOverview_EmptyDB(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)

	result, err := svc.Overview(context.Background())
	if err != nil {
		t.Fatalf("Overview() error: %v", err)
	}

	if result.TotalTopics != 0 || result.TotalEvents != 0 || result.TotalProducers != 0 || result.TotalConsumers != 0 {
		t.Errorf("expected all zeros for empty DB, got topics=%d events=%d producers=%d consumers=%d",
			result.TotalTopics, result.TotalEvents, result.TotalProducers, result.TotalConsumers)
	}
}
