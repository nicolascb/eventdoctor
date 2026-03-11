package service

import (
"context"
"database/sql"
"os"
"testing"

"github.com/nicolascb/eventdoctor/internal/db"
"github.com/nicolascb/eventdoctor/internal/db/models"
)

func strPtr(s string) *string { return &s }

// setupTestDB initializes an in-memory SQLite database with migrations applied.
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "eventdoctor-test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	testDB, err := db.NewSQLiteDB(tmpFile.Name(), false)
	if err != nil {
		t.Fatalf("failed to create test DB: %v", err)
	}
	t.Cleanup(func() { testDB.Close() })

	return testDB
}

// seedTestData populates the database with a standard set of test data.
func seedTestData(t *testing.T, testDB *sql.DB) {
	t.Helper()
	ctx := context.Background()

	svcA, err := db.GetOrCreateService(ctx, testDB, "service-a", "https://github.com/org/service-a")
	if err != nil {
		t.Fatalf("failed to create service-a: %v", err)
	}

	svcB, err := db.GetOrCreateService(ctx, testDB, "service-b", "https://github.com/org/service-b")
	if err != nil {
		t.Fatalf("failed to create service-b: %v", err)
	}

	ownerID := svcA.ID
	topicOrders, err := db.GetOrCreateTopic(ctx, testDB, "orders.events", &ownerID, "Order domain events")
	if err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}

	topicPayments, err := db.GetOrCreateTopic(ctx, testDB, "payments.events", nil, "Payment domain events")
	if err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}

	v1 := "1.0.0"
	eventCreated, err := db.GetOrCreateEvent(ctx, testDB, topicOrders.ID, "OrderCreated", &v1, "https://schemas.local/order-created.json", "An order was created")
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	eventUpdated, err := db.GetOrCreateEvent(ctx, testDB, topicOrders.ID, "OrderUpdated", nil, "", "An order was updated")
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	eventPaid, err := db.GetOrCreateEvent(ctx, testDB, topicPayments.ID, "PaymentProcessed", nil, "", "A payment was processed")
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	// Insert headers for OrderCreated
	_, err = db.InsertEventHeader(ctx, testDB, models.EventHeader{EventID: eventCreated.ID, Name: "X-Request-ID", Description: "Unique request identifier"})
	if err != nil {
		t.Fatalf("failed to insert header: %v", err)
	}

	// service-a produces OrderCreated and OrderUpdated
	_, err = db.InsertProducer(ctx, testDB, models.Producer{EventID: eventCreated.ID, ServiceID: svcA.ID, Writes: true})
	if err != nil {
		t.Fatalf("failed to insert producer: %v", err)
	}
	_, err = db.InsertProducer(ctx, testDB, models.Producer{EventID: eventUpdated.ID, ServiceID: svcA.ID, Writes: true})
	if err != nil {
		t.Fatalf("failed to insert producer: %v", err)
	}

	// service-b produces PaymentProcessed
	_, err = db.InsertProducer(ctx, testDB, models.Producer{EventID: eventPaid.ID, ServiceID: svcB.ID, Writes: false})
	if err != nil {
		t.Fatalf("failed to insert producer: %v", err)
	}

	// service-b consumes OrderCreated
	_, err = db.InsertConsumer(ctx, testDB, models.Consumer{EventID: eventCreated.ID, ServiceID: svcB.ID, ConsumerGroup: "service-b-orders", Description: "Processes orders", EventVersion: &v1})
	if err != nil {
		t.Fatalf("failed to insert consumer: %v", err)
	}

	// service-a consumes PaymentProcessed
	_, err = db.InsertConsumer(ctx, testDB, models.Consumer{EventID: eventPaid.ID, ServiceID: svcA.ID, ConsumerGroup: "service-a-payments", Description: "Handles payments"})
	if err != nil {
		t.Fatalf("failed to insert consumer: %v", err)
	}
}
