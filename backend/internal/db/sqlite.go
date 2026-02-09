package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

func NewSQLiteDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable foreign key enforcement in SQLite
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	if os.Getenv("WITH_MOCK") == "1" {
		if err := mockData(context.Background(), db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	migrationBytes, err := os.ReadFile("migration.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute the migration script
	_, err = db.Exec(string(migrationBytes))
	if err != nil {
		return fmt.Errorf("failed to execute migrations: %w", err)
	}

	return nil
}

func mockData(ctx context.Context, db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM consumers`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM producers`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM event_headers`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM events`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM topics`)
	if err != nil {
		return err
	}

	now := time.Now()

	topics := []models.Topic{
		{Name: "user.events", Owner: "user-service", CreatedAt: now},
		{Name: "order.events", Owner: "order-service", CreatedAt: now},
		{Name: "payment.events", Owner: "payment-service", CreatedAt: now},
	}

	topicIDs := make(map[string]int64)
	for _, t := range topics {
		id, err := InsertTopic(ctx, db, t)
		if err != nil {
			return err
		}

		topicIDs[t.Name] = id
	}

	schemaVersion := "1.0.0"

	eventsSeed := []struct {
		topic string
		event models.Event
	}{
		{"user.events", models.Event{EventName: "UserCreated", SchemaURL: "https://schemas.local/user-created.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"user.events", models.Event{EventName: "UserUpdated", SchemaURL: "https://schemas.local/user-updated.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"order.events", models.Event{EventName: "OrderPlaced", SchemaURL: "https://schemas.local/order-placed.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"order.events", models.Event{EventName: "OrderShipped", SchemaURL: "https://schemas.local/order-shipped.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"payment.events", models.Event{EventName: "PaymentAuthorized", SchemaURL: "https://schemas.local/payment-authorized.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"payment.events", models.Event{EventName: "PaymentCaptured", SchemaURL: "https://schemas.local/payment-captured.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
	}

	eventIDs := map[string]int64{}
	for _, e := range eventsSeed {
		e.event.TopicID = topicIDs[e.topic]
		id, err := InsertEvent(ctx, db, e.event)
		if err != nil {
			return err
		}
		eventIDs[fmt.Sprintf("%s:%s", e.topic, e.event.EventName)] = id
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
			return fmt.Errorf("event not found for key %s", p.eventKey)
		}
		p.prod.EventID = id
		if _, err := InsertProducer(ctx, db, p.prod); err != nil {
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
			return fmt.Errorf("event not found for key %s", c.eventKey)
		}
		c.cons.EventID = id
		if _, err := InsertConsumer(ctx, db, c.cons); err != nil {
			return err
		}
	}

	return nil
}

// ensureCreatedAt define now se o campo estiver zerado
func ensureCreatedAt(t *time.Time) {
	if t.IsZero() {
		*t = time.Now()
	}
}

// InsertTopic insere um novo tópico no banco de dados
func InsertTopic(ctx context.Context, executor SQLExecutor, topic models.Topic) (int64, error) {
	ensureCreatedAt(&topic.CreatedAt)
	query := `INSERT INTO topics (name, owner, created_at) VALUES (?, ?, ?)`

	res, err := executor.ExecContext(ctx, query, topic.Name, topic.Owner, topic.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetTopic busca um tópico pelo nome
func GetTopic(ctx context.Context, executor SQLExecutor, topicName string) (*models.Topic, error) {
	query := `SELECT id, name, owner, created_at FROM topics WHERE name = ?`
	row := executor.QueryRowContext(ctx, query, topicName)

	var topic models.Topic
	err := row.Scan(&topic.ID, &topic.Name, &topic.Owner, &topic.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &topic, nil
}

// RemoveTopic remove um tópico pelo nome
func RemoveTopic(ctx context.Context, executor SQLExecutor, topicName string) error {
	query := `DELETE FROM topics WHERE name = ?`
	_, err := executor.ExecContext(ctx, query, topicName)
	return err
}

// ListTopics retorna todos os tópicos ordenados por nome
func ListTopics(ctx context.Context, executor SQLExecutor) ([]models.Topic, error) {
	query := `SELECT id, name, owner, created_at FROM topics ORDER BY name`
	rows, err := executor.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var t models.Topic
		if err := rows.Scan(&t.ID, &t.Name, &t.Owner, &t.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}

// InsertEvent insere um novo evento no banco de dados
func InsertEvent(ctx context.Context, executor SQLExecutor, event models.Event) (int64, error) {
	ensureCreatedAt(&event.CreatedAt)
	query := `INSERT INTO events (topic_id, event_name, schema_url, schema_version, deprecated, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := executor.ExecContext(ctx, query, event.TopicID, event.EventName, event.SchemaURL, event.SchemaVersion, event.Deprecated, event.CreatedAt)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// RemoveEvent remove um evento pelo ID
func RemoveEvent(ctx context.Context, executor SQLExecutor, eventID int64) error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := executor.ExecContext(ctx, query, eventID)
	return err
}

// GetEventsByTopic retorna todos os eventos de um tópico
func GetEventsByTopic(ctx context.Context, executor SQLExecutor, topicID int64) ([]models.Event, error) {
	query := `SELECT id, topic_id, event_name, schema_url, schema_version, deprecated, created_at FROM events WHERE topic_id = ?`
	rows, err := executor.QueryContext(ctx, query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.TopicID, &event.EventName, &event.SchemaURL, &event.SchemaVersion, &event.Deprecated, &event.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// InsertProducer insere um novo produtor no banco de dados
func InsertProducer(ctx context.Context, executor SQLExecutor, producer models.Producer) (int64, error) {
	ensureCreatedAt(&producer.CreatedAt)
	query := `INSERT INTO producers (event_id, service, writes, repository, created_at) VALUES (?, ?, ?, ?, ?)`
	result, err := executor.ExecContext(ctx, query, producer.EventID, producer.Service, producer.Writes, producer.Repository, producer.CreatedAt)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// RemoveProducer remove um produtor pelo ID
func RemoveProducer(ctx context.Context, executor SQLExecutor, producerID int64) error {
	query := `DELETE FROM producers WHERE id = ?`
	_, err := executor.ExecContext(ctx, query, producerID)
	return err
}

// GetProducersByEvent retorna todos os produtores de um evento
func GetProducersByEvent(ctx context.Context, executor SQLExecutor, eventID int64) ([]models.Producer, error) {
	query := `SELECT id, event_id, service, writes, repository, created_at FROM producers WHERE event_id = ?`
	rows, err := executor.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var producers []models.Producer
	for rows.Next() {
		var producer models.Producer
		err := rows.Scan(&producer.ID, &producer.EventID, &producer.Service, &producer.Writes, &producer.Repository, &producer.CreatedAt)
		if err != nil {
			return nil, err
		}
		producers = append(producers, producer)
	}

	return producers, nil
}

// InsertConsumer insere um novo consumidor no banco de dados
func InsertConsumer(ctx context.Context, executor SQLExecutor, consumer models.Consumer) (int64, error) {
	ensureCreatedAt(&consumer.CreatedAt)
	query := `INSERT INTO consumers (event_id, service, consumer_group, event_version, repository, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := executor.ExecContext(ctx, query, consumer.EventID, consumer.Service, consumer.ConsumerGroup, consumer.EventVersion, consumer.Repository, consumer.CreatedAt)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// RemoveConsumer remove um consumidor pelo ID
func RemoveConsumer(ctx context.Context, executor SQLExecutor, consumerID int64) error {
	query := `DELETE FROM consumers WHERE id = ?`
	_, err := executor.ExecContext(ctx, query, consumerID)
	return err
}

// GetConsumersByEvent retorna todos os consumidores de um evento
func GetConsumersByEvent(ctx context.Context, executor SQLExecutor, eventID int64) ([]models.Consumer, error) {
	query := `SELECT id, event_id, service, consumer_group, event_version, repository, created_at FROM consumers WHERE event_id = ?`
	rows, err := executor.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consumers []models.Consumer
	for rows.Next() {
		var consumer models.Consumer
		err := rows.Scan(&consumer.ID, &consumer.EventID, &consumer.Service, &consumer.ConsumerGroup, &consumer.EventVersion, &consumer.Repository, &consumer.CreatedAt)
		if err != nil {
			return nil, err
		}
		consumers = append(consumers, consumer)
	}

	return consumers, nil
}

// GetOrCreateTopic busca um tópico pelo nome ou cria um novo se não existir
func GetOrCreateTopic(ctx context.Context, executor SQLExecutor, topicName string, owner string) (models.Topic, error) {
	// Primeiro tenta encontrar o tópico existente
	topic, err := GetTopic(ctx, executor, topicName)
	if err != nil {
		return models.Topic{}, err
	}

	// Se o tópico já existe
	if topic != nil {
		// Verifica se precisa atualizar o owner
		if owner != "" && owner != topic.Owner {
			if err := updateTopicOwner(ctx, executor, topicName, owner); err != nil {
				return models.Topic{}, err
			}

			updatedTopic := models.Topic{
				ID:        topic.ID,
				Name:      topic.Name,
				Owner:     owner,
				CreatedAt: topic.CreatedAt,
			}

			return updatedTopic, nil
		}
		return *topic, nil
	}

	// Caso contrário, crie um novo tópico
	newTopic := models.Topic{
		Name:      topicName,
		Owner:     owner,
		CreatedAt: time.Now(),
	}

	id, err := InsertTopic(ctx, executor, newTopic)
	if err != nil {
		return models.Topic{}, err
	}

	return models.Topic{
		ID:        id,
		Name:      topicName,
		Owner:     owner,
		CreatedAt: newTopic.CreatedAt,
	}, nil
}

// updateTopicOwner atualiza o owner de um tópico existente
func updateTopicOwner(ctx context.Context, executor SQLExecutor, topicName, newOwner string) error {
	query := `UPDATE topics SET owner = ? WHERE name = ?`
	_, err := executor.ExecContext(ctx, query, newOwner, topicName)
	if err != nil {
		return fmt.Errorf("error updating topic owner: %w", err)
	}
	return nil
}

// GetOrCreateEvent busca um evento por nome em um tópico ou cria um novo se não existir
func GetOrCreateEvent(ctx context.Context, executor SQLExecutor, topicID int64, eventName string, schemaVersion *string, schemaURL string) (models.Event, error) {
	// Primeiro tenta encontrar o evento existente
	query := `SELECT id, topic_id, event_name, schema_url, schema_version, deprecated, created_at 
              FROM events 
              WHERE topic_id = ? AND event_name = ?`

	row := executor.QueryRowContext(ctx, query, topicID, eventName)

	var event models.Event
	err := row.Scan(
		&event.ID,
		&event.TopicID,
		&event.EventName,
		&event.SchemaURL,
		&event.SchemaVersion,
		&event.Deprecated,
		&event.CreatedAt,
	)

	if err == nil {
		// Se o evento existe, retorne-o
		return event, nil
	}

	if err != sql.ErrNoRows {
		// Se houver um erro que não seja ErrNoRows, retorne o erro
		return models.Event{}, err
	}

	// Se o evento não existe, crie um novo
	newEvent := models.Event{
		TopicID:       topicID,
		EventName:     eventName,
		SchemaURL:     schemaURL,
		SchemaVersion: schemaVersion,
		Deprecated:    false,
		CreatedAt:     time.Now(),
	}

	id, err := InsertEvent(ctx, executor, newEvent)
	if err != nil {
		return models.Event{}, err
	}

	return models.Event{
		ID:            id,
		TopicID:       topicID,
		EventName:     eventName,
		SchemaURL:     schemaURL,
		SchemaVersion: schemaVersion,
		Deprecated:    false,
		CreatedAt:     newEvent.CreatedAt,
	}, nil
}

// InsertEventHeader insere um novo header para um evento
func InsertEventHeader(ctx context.Context, executor SQLExecutor, header models.EventHeader) (int64, error) {
	ensureCreatedAt(&header.CreatedAt)
	query := `INSERT INTO event_headers (event_id, name, description, created_at) VALUES (?, ?, ?, ?)`
	result, err := executor.ExecContext(ctx, query, header.EventID, header.Name, header.Description, header.CreatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetEventHeaders retorna todos os headers de um evento
func GetEventHeaders(ctx context.Context, executor SQLExecutor, eventID int64) ([]models.EventHeader, error) {
	query := `SELECT id, event_id, name, description, created_at FROM event_headers WHERE event_id = ? ORDER BY name`
	rows, err := executor.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var headers []models.EventHeader
	for rows.Next() {
		var header models.EventHeader
		err := rows.Scan(&header.ID, &header.EventID, &header.Name, &header.Description, &header.CreatedAt)
		if err != nil {
			return nil, err
		}
		headers = append(headers, header)
	}
	return headers, nil
}

// DeleteEventHeaders remove todos os headers de um evento
func DeleteEventHeaders(ctx context.Context, executor SQLExecutor, eventID int64) error {
	query := `DELETE FROM event_headers WHERE event_id = ?`
	_, err := executor.ExecContext(ctx, query, eventID)
	return err
}
