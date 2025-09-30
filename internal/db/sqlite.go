package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nicolascb/eventdoctor/internal/api"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.createTables(); err != nil {
		return nil, err
	}

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

// TopicStore methods
func (s *SQLiteStore) InsertTopic(topic api.Topic) error {
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

// EventStore methods
func (s *SQLiteStore) InsertEvent(event api.Event) error {
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
