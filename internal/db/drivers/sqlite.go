package drivers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/nicolascb/eventdoctor/internal/db/models"
	"github.com/nicolascb/eventdoctor/internal/db/repositories"
)

func MockData(ctx context.Context, db *sql.DB) error {
	db.Exec(`DELETE FROM consumers`)
	db.Exec(`DELETE FROM producers`)
	db.Exec(`DELETE FROM events`)
	db.Exec(`DELETE FROM topics`)

	now := time.Now()

	// Tópicos fictícios
	topics := []models.Topic{
		{Name: "user.events", Owner: "user-service", CreatedAt: now},
		{Name: "order.events", Owner: "order-service", CreatedAt: now},
		{Name: "payment.events", Owner: "payment-service", CreatedAt: now},
	}

	topicIDs := make(map[string]int64)
	for _, t := range topics {
		id, err := repositories.InsertTopic(ctx, db, t)
		if err != nil {
			return err
		}

		topicIDs[t.Name] = id
	}

	schemaVersion := "1.0.0"

	// Eventos fictícios
	eventsSeed := []struct {
		topic string
		event models.Event
	}{
		{"user.events", models.Event{EventType: "UserCreated", SchemaURL: "https://schemas.local/user-created.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"user.events", models.Event{EventType: "UserUpdated", SchemaURL: "https://schemas.local/user-updated.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"order.events", models.Event{EventType: "OrderPlaced", SchemaURL: "https://schemas.local/order-placed.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"order.events", models.Event{EventType: "OrderShipped", SchemaURL: "https://schemas.local/order-shipped.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"payment.events", models.Event{EventType: "PaymentAuthorized", SchemaURL: "https://schemas.local/payment-authorized.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"payment.events", models.Event{EventType: "PaymentCaptured", SchemaURL: "https://schemas.local/payment-captured.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
	}

	eventIDs := map[string]int64{}
	for _, e := range eventsSeed {
		e.event.TopicID = topicIDs[e.topic]
		id, err := repositories.InsertEvent(ctx, db, e.event)
		if err != nil {
			return err
		}
		eventIDs[fmt.Sprintf("%s:%s", e.topic, e.event.EventType)] = id
	}

	fakeRepository := "https://github.com/org/fake-repo"

	// Produtores fictícios
	producersSeed := []struct {
		eventKey string
		prod     models.Producer
	}{
		{"user.events:UserCreated", models.Producer{Service: "user-service", Writes: true, Repository: &fakeRepository, CreatedAt: now}},
		{"user.events:UserUpdated", models.Producer{Service: "user-service", Writes: true, Repository: &fakeRepository, CreatedAt: now}},
		{"order.events:OrderPlaced", models.Producer{Service: "checkout-service", Writes: true, Repository: &fakeRepository, CreatedAt: now}},
		{"order.events:OrderShipped", models.Producer{Service: "logistics-service", Writes: true, Repository: &fakeRepository, CreatedAt: now}},
		{"payment.events:PaymentAuthorized", models.Producer{Service: "payment-service", Writes: true, Repository: &fakeRepository, CreatedAt: now}},
		{"payment.events:PaymentCaptured", models.Producer{Service: "payment-service", Writes: true, Repository: &fakeRepository, CreatedAt: now}},
	}

	for _, p := range producersSeed {
		id, ok := eventIDs[p.eventKey]
		if !ok {
			return fmt.Errorf("evento não encontrado para chave produtor %s", p.eventKey)
		}
		p.prod.EventID = id
		if _, err := repositories.InsertProducer(ctx, db, p.prod); err != nil {
			return err
		}
	}

	// Consumidores fictícios
	consumersSeed := []struct {
		eventKey string
		cons     models.Consumer
	}{
		{"user.events:UserCreated", models.Consumer{ConsumerGroup: "analytics-service-group", Service: "analytics-service", Repository: &fakeRepository, CreatedAt: now}},
		{"user.events:UserCreated", models.Consumer{ConsumerGroup: "email-service-group", Service: "email-service", Repository: &fakeRepository, CreatedAt: now}},
		{"user.events:UserUpdated", models.Consumer{ConsumerGroup: "notification-service-group", Service: "notification-service", Repository: &fakeRepository, CreatedAt: now}},
		{"order.events:OrderPlaced", models.Consumer{ConsumerGroup: "inventory-service-group", Service: "inventory-service", Repository: &fakeRepository, CreatedAt: now}},
		{"order.events:OrderShipped", models.Consumer{ConsumerGroup: "notification-service-group", Service: "notification-service", Repository: &fakeRepository, CreatedAt: now}},
		{"payment.events:PaymentAuthorized", models.Consumer{ConsumerGroup: "fraud-service-group", Service: "fraud-service", Repository: &fakeRepository, CreatedAt: now}},
		{"payment.events:PaymentCaptured", models.Consumer{ConsumerGroup: "billing-service-group", Service: "billing-service", Repository: &fakeRepository, CreatedAt: now}},
	}

	for _, c := range consumersSeed {
		id, ok := eventIDs[c.eventKey]
		if !ok {
			return fmt.Errorf("evento não encontrado para chave consumidor %s", c.eventKey)
		}
		c.cons.EventID = id
		if _, err := repositories.InsertConsumer(ctx, db, c.cons); err != nil {
			return err
		}
	}

	return nil
}
