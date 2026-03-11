package service

import (
	"context"
	"testing"

	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
)

// --- SaveSpec integration tests ---

func TestSaveSpec(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)
	ctx := context.Background()

	v1 := "1.0.0"
	spec := eventdoctor.EventDoctorSpec{
		Version: "1",
		Service: "my-service",
		Config: eventdoctor.Config{
			Repository: "https://github.com/org/my-service",
		},
		Producers: []eventdoctor.Producer{
			{
				Topic:       "orders.events",
				Owner:       true,
				Writes:      true,
				Description: "Order events topic",
				Events: []eventdoctor.Event{
					{
						Name:        "OrderCreated",
						Version:     &v1,
						Description: "Order was created",
						SchemaURL:   "https://schemas.local/order-created.json",
						Headers: []eventdoctor.EventHeader{
							{Name: "X-Request-ID", Description: "Request ID"},
						},
					},
				},
			},
		},
		Consumers: []eventdoctor.Consumer{
			{
				Group:       "my-service-payments",
				Description: "Consumes payment events",
				Topics: []eventdoctor.Topic{
					{
						Name: "payments.events",
						Events: []eventdoctor.ConsumerEvent{
							{Name: "PaymentProcessed", Version: &v1},
						},
					},
				},
			},
		},
	}

	err := svc.SaveSpec(ctx, spec)
	if err != nil {
		t.Fatalf("SaveSpec() error: %v", err)
	}

	// Verify the data was saved correctly
	overview, err := svc.Overview(ctx)
	if err != nil {
		t.Fatalf("Overview() error: %v", err)
	}

	if overview.TotalTopics != 2 {
		t.Errorf("expected 2 topics, got %d", overview.TotalTopics)
	}
	if overview.TotalEvents != 2 {
		t.Errorf("expected 2 events, got %d", overview.TotalEvents)
	}

	// Verify service view
	serviceView, err := svc.GetServiceView(ctx, "my-service")
	if err != nil {
		t.Fatalf("GetServiceView() error: %v", err)
	}
	if len(serviceView.Produces) != 1 {
		t.Errorf("expected 1 produces entry, got %d", len(serviceView.Produces))
	}
	if len(serviceView.Consumes) != 1 {
		t.Errorf("expected 1 consumes entry, got %d", len(serviceView.Consumes))
	}
}

func TestSaveSpec_Idempotent(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)
	ctx := context.Background()

	spec := eventdoctor.EventDoctorSpec{
		Version: "1",
		Service: "idempotent-svc",
		Config: eventdoctor.Config{
			Repository: "https://github.com/org/idempotent-svc",
		},
		Producers: []eventdoctor.Producer{
			{
				Topic:  "test.topic",
				Owner:  true,
				Writes: true,
				Events: []eventdoctor.Event{
					{Name: "TestEvent", Description: "A test event"},
				},
			},
		},
	}

	// Apply twice
	if err := svc.SaveSpec(ctx, spec); err != nil {
		t.Fatalf("first SaveSpec() error: %v", err)
	}
	if err := svc.SaveSpec(ctx, spec); err != nil {
		t.Fatalf("second SaveSpec() error: %v", err)
	}

	// Should still have only 1 producer
	producers, err := svc.ListProducers(ctx, 1, 10, "")
	if err != nil {
		t.Fatalf("ListProducers() error: %v", err)
	}

	count := 0
	for _, p := range producers.Producers {
		if p.Service == "idempotent-svc" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 producer entry for idempotent-svc after double apply, got %d", count)
	}
}

func TestSaveSpec_ReplacesOldData(t *testing.T) {
	testDB := setupTestDB(t)
	svc := NewService(testDB)
	ctx := context.Background()

	repo := "https://github.com/org/replace-svc"

	// Apply initial spec
	spec1 := eventdoctor.EventDoctorSpec{
		Version: "1",
		Service: "replace-svc",
		Config:  eventdoctor.Config{Repository: repo},
		Producers: []eventdoctor.Producer{
			{
				Topic: "old.topic", Owner: true, Writes: true,
				Events: []eventdoctor.Event{
					{Name: "OldEvent"},
				},
			},
		},
	}
	if err := svc.SaveSpec(ctx, spec1); err != nil {
		t.Fatalf("first SaveSpec() error: %v", err)
	}

	// Apply updated spec with different events
	spec2 := eventdoctor.EventDoctorSpec{
		Version: "1",
		Service: "replace-svc",
		Config:  eventdoctor.Config{Repository: repo},
		Producers: []eventdoctor.Producer{
			{
				Topic: "new.topic", Owner: true, Writes: true,
				Events: []eventdoctor.Event{
					{Name: "NewEvent"},
				},
			},
		},
	}
	if err := svc.SaveSpec(ctx, spec2); err != nil {
		t.Fatalf("second SaveSpec() error: %v", err)
	}

	// Verify the old producer data was replaced
	serviceView, err := svc.GetServiceView(ctx, "replace-svc")
	if err != nil {
		t.Fatalf("GetServiceView() error: %v", err)
	}

	if len(serviceView.Produces) != 1 {
		t.Fatalf("expected 1 produces entry, got %d", len(serviceView.Produces))
	}
	if serviceView.Produces[0].Topic != "new.topic" {
		t.Errorf("expected new.topic, got %s", serviceView.Produces[0].Topic)
	}
	if serviceView.Produces[0].Event != "NewEvent" {
		t.Errorf("expected NewEvent, got %s", serviceView.Produces[0].Event)
	}
}
