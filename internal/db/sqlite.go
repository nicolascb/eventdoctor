package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nicolascb/eventdoctor/internal/api"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {

	fmt.Println("Connecting to SQLite database...", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Ativa enforcement de foreign keys no SQLite
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return nil, fmt.Errorf("falha ao ativar foreign_keys: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.createTables(); err != nil {
		return nil, err
	}

	// Seed de dados fictícios (idempotente)
	if err := store.mockData(); err != nil {
		return nil, err
	}

	fmt.Println("Connected to SQLite database successfully")
	return store, nil
}

func (s *SQLiteStore) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS topics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			owner TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			topic_id INTEGER NOT NULL,
			event_type TEXT NOT NULL,
			schema_url TEXT NOT NULL,
			schema_version TEXT,
			deprecated BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (topic_id) REFERENCES topics (id)
		)`,
		`CREATE TABLE IF NOT EXISTS producers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_id INTEGER NOT NULL,
			producer TEXT NOT NULL,
			repository TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (event_id) REFERENCES events (id)
		)`,
		`CREATE TABLE IF NOT EXISTS consumers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_id INTEGER NOT NULL,
			consumer TEXT NOT NULL,
			repository TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (event_id) REFERENCES events (id)
		)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteStore) mockData() error {
	s.db.Exec(`DELETE FROM consumers`)
	s.db.Exec(`DELETE FROM producers`)
	s.db.Exec(`DELETE FROM events`)
	s.db.Exec(`DELETE FROM topics`)

	now := time.Now()

	// Tópicos fictícios
	topics := []api.Topic{
		{Name: "user.events", Owner: "user-service", CreatedAt: now},
		{Name: "order.events", Owner: "order-service", CreatedAt: now},
		{Name: "payment.events", Owner: "payment-service", CreatedAt: now},
	}

	topicIDs := make(map[string]int)
	for _, t := range topics {
		if err := s.InsertTopic(t); err != nil {
			return err
		}
		stored, err := s.Get(t.Name)
		if err != nil {
			return err
		}
		if stored == nil {
			return fmt.Errorf("falha ao buscar tópico %s após inserção", t.Name)
		}
		topicIDs[t.Name] = stored.ID
	}

	schemaVersion := "1.0.0"

	// Eventos fictícios
	eventsSeed := []struct {
		topic string
		event api.Event
	}{
		{"user.events", api.Event{EventType: "UserCreated", SchemaURL: "https://schemas.local/user-created.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"user.events", api.Event{EventType: "UserUpdated", SchemaURL: "https://schemas.local/user-updated.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"order.events", api.Event{EventType: "OrderPlaced", SchemaURL: "https://schemas.local/order-placed.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"order.events", api.Event{EventType: "OrderShipped", SchemaURL: "https://schemas.local/order-shipped.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"payment.events", api.Event{EventType: "PaymentAuthorized", SchemaURL: "https://schemas.local/payment-authorized.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
		{"payment.events", api.Event{EventType: "PaymentCaptured", SchemaURL: "https://schemas.local/payment-captured.json", SchemaVersion: &schemaVersion, Deprecated: false, CreatedAt: now}},
	}

	for _, e := range eventsSeed {
		e.event.TopicID = topicIDs[e.topic]
		if err := s.InsertEvent(e.event); err != nil {
			return err
		}
	}

	// Mapeia IDs dos eventos para ligação com produtores/consumidores
	eventIDs := map[string]int{}
	for topicName, topicID := range topicIDs {
		evts, err := s.GetEventsByTopic(topicID)
		if err != nil {
			return err
		}
		for _, ev := range evts {
			eventIDs[fmt.Sprintf("%s:%s", topicName, ev.EventType)] = ev.ID
		}
	}

	fakeRepository := "https://github.com/org/fake-repo"

	// Produtores fictícios
	producersSeed := []struct {
		eventKey string
		prod     api.Producer
	}{
		{"user.events:UserCreated", api.Producer{Producer: "user-service", Repository: &fakeRepository, CreatedAt: now}},
		{"user.events:UserUpdated", api.Producer{Producer: "user-service", Repository: &fakeRepository, CreatedAt: now}},
		{"order.events:OrderPlaced", api.Producer{Producer: "checkout-service", Repository: &fakeRepository, CreatedAt: now}},
		{"order.events:OrderShipped", api.Producer{Producer: "logistics-service", Repository: &fakeRepository, CreatedAt: now}},
		{"payment.events:PaymentAuthorized", api.Producer{Producer: "payment-service", Repository: &fakeRepository, CreatedAt: now}},
		{"payment.events:PaymentCaptured", api.Producer{Producer: "payment-service", Repository: &fakeRepository, CreatedAt: now}},
	}

	for _, p := range producersSeed {
		id, ok := eventIDs[p.eventKey]
		if !ok {
			return fmt.Errorf("evento não encontrado para chave produtor %s", p.eventKey)
		}
		p.prod.EventID = id
		if err := s.InsertProducer(p.prod); err != nil {
			return err
		}
	}

	// Consumidores fictícios
	consumersSeed := []struct {
		eventKey string
		cons     api.Consumer
	}{
		{"user.events:UserCreated", api.Consumer{Consumer: "analytics-service", Repository: &fakeRepository, CreatedAt: now}},
		{"user.events:UserCreated", api.Consumer{Consumer: "email-service", Repository: &fakeRepository, CreatedAt: now}},
		{"user.events:UserUpdated", api.Consumer{Consumer: "notification-service", Repository: &fakeRepository, CreatedAt: now}},
		{"order.events:OrderPlaced", api.Consumer{Consumer: "inventory-service", Repository: &fakeRepository, CreatedAt: now}},
		{"order.events:OrderShipped", api.Consumer{Consumer: "notification-service", Repository: &fakeRepository, CreatedAt: now}},
		{"payment.events:PaymentAuthorized", api.Consumer{Consumer: "fraud-service", Repository: &fakeRepository, CreatedAt: now}},
		{"payment.events:PaymentCaptured", api.Consumer{Consumer: "billing-service", Repository: &fakeRepository, CreatedAt: now}},
	}

	for _, c := range consumersSeed {
		id, ok := eventIDs[c.eventKey]
		if !ok {
			return fmt.Errorf("evento não encontrado para chave consumidor %s", c.eventKey)
		}
		c.cons.EventID = id
		if err := s.InsertConsumer(c.cons); err != nil {
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

// TopicStore methods
func (s *SQLiteStore) InsertTopic(topic api.Topic) error {
	ensureCreatedAt(&topic.CreatedAt)
	query := `INSERT INTO topics (name, owner, created_at) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, topic.Name, topic.Owner, topic.CreatedAt)
	return err
}

func (s *SQLiteStore) Get(topicName string) (*api.Topic, error) {
	query := `SELECT id, name, owner, created_at FROM topics WHERE name = ?`
	row := s.db.QueryRow(query, topicName)

	var topic api.Topic
	err := row.Scan(&topic.ID, &topic.Name, &topic.Owner, &topic.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &topic, nil
}

func (s *SQLiteStore) RemoveTopic(topicName string) error {
	query := `DELETE FROM topics WHERE name = ?`
	_, err := s.db.Exec(query, topicName)
	return err
}

func (s *SQLiteStore) ListTopics() ([]api.Topic, error) {
	query := `SELECT id, name, owner, created_at FROM topics ORDER BY name`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []api.Topic
	for rows.Next() {
		var t api.Topic
		if err := rows.Scan(&t.ID, &t.Name, &t.Owner, &t.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}

// EventStore methods
func (s *SQLiteStore) InsertEvent(event api.Event) error {
	ensureCreatedAt(&event.CreatedAt)
	query := `INSERT INTO events (topic_id, event_type, schema_url, schema_version, deprecated, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, event.TopicID, event.EventType, event.SchemaURL, event.SchemaVersion, event.Deprecated, event.CreatedAt)
	return err
}

func (s *SQLiteStore) RemoveEvent(eventID int) error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := s.db.Exec(query, eventID)
	return err
}

func (s *SQLiteStore) GetEventsByTopic(topicID int) ([]api.Event, error) {
	query := `SELECT id, topic_id, event_type, schema_url, schema_version, deprecated, created_at FROM events WHERE topic_id = ?`
	rows, err := s.db.Query(query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []api.Event
	for rows.Next() {
		var event api.Event
		err := rows.Scan(&event.ID, &event.TopicID, &event.EventType, &event.SchemaURL, &event.SchemaVersion, &event.Deprecated, &event.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// ProducerStore methods
func (s *SQLiteStore) InsertProducer(producer api.Producer) error {
	ensureCreatedAt(&producer.CreatedAt)
	query := `INSERT INTO producers (event_id, producer, repository, created_at) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, producer.EventID, producer.Producer, producer.Repository, producer.CreatedAt)
	return err
}

func (s *SQLiteStore) RemoveProducer(producerID int) error {
	query := `DELETE FROM producers WHERE id = ?`
	_, err := s.db.Exec(query, producerID)
	return err
}

func (s *SQLiteStore) GetProducersByEvent(eventID int) ([]api.Producer, error) {
	query := `SELECT id, event_id, producer, repository, created_at FROM producers WHERE event_id = ?`
	rows, err := s.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var producers []api.Producer
	for rows.Next() {
		var producer api.Producer
		err := rows.Scan(&producer.ID, &producer.EventID, &producer.Producer, &producer.Repository, &producer.CreatedAt)
		if err != nil {
			return nil, err
		}
		producers = append(producers, producer)
	}

	return producers, nil
}

// ConsumerStore methods
func (s *SQLiteStore) InsertConsumer(consumer api.Consumer) error {
	ensureCreatedAt(&consumer.CreatedAt)
	query := `INSERT INTO consumers (event_id, consumer, repository, created_at) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, consumer.EventID, consumer.Consumer, consumer.Repository, consumer.CreatedAt)
	return err
}

func (s *SQLiteStore) RemoveConsumer(consumerID int) error {
	query := `DELETE FROM consumers WHERE id = ?`
	_, err := s.db.Exec(query, consumerID)
	return err
}

func (s *SQLiteStore) GetConsumersByEvent(eventID int) ([]api.Consumer, error) {
	query := `SELECT id, event_id, consumer, repository, created_at FROM consumers WHERE event_id = ?`
	rows, err := s.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consumers []api.Consumer
	for rows.Next() {
		var consumer api.Consumer
		err := rows.Scan(&consumer.ID, &consumer.EventID, &consumer.Consumer, &consumer.Repository, &consumer.CreatedAt)
		if err != nil {
			return nil, err
		}
		consumers = append(consumers, consumer)
	}

	return consumers, nil
}
