package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nicolascb/eventdoctor/internal/db/models"
)

func NewSQLiteDB(dbPath string) (*sql.DB, error) {
	if os.Getenv("WITH_MOCK") == "1" {
		os.Remove(dbPath)
	}

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
	now := time.Now()
	fakeRepository := "https://github.com/org/fake-repo"

	serviceNames := []string{
		"user-service", "checkout-service", "logistics-service", "payment-service",
		"analytics-service", "email-service", "notification-service", "inventory-service",
		"fraud-service", "billing-service",
	}

	serviceIDs := make(map[int]int64)
	for i, name := range serviceNames {
		s, err := GetOrCreateService(ctx, db, name, fakeRepository)
		if err != nil {
			return err
		}
		serviceIDs[i] = s.ID
	}

	schemaVersion := "1.0.0"

	// Vocabulary for random names without numbers
	adjectives := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta", "eta", "theta", "iota", "kappa"}
	nouns := []string{"falcon", "eagle", "hawk", "owl", "raven", "wolf", "bear", "lion", "tiger", "shark"}
	actions := []string{
		"Born", "Died", "Awake", "Asleep", "Moving", "Stopped", "Eating", "Drinking",
		"Flying", "Swimming", "Hunting", "Hiding", "Singing", "Running", "Jumping",
		"Climbing", "Falling", "Rising", "Growing", "Shrinking",
	}

	// 100 topics (10 * 10)
	serviceIndex := 0
	for _, adj := range adjectives {
		for _, noun := range nouns {
			topicName := fmt.Sprintf("%s.%s.events", adj, noun)
			ownerID := serviceIDs[serviceIndex%len(serviceIDs)]
			serviceIndex++

			topic := models.Topic{
				Name:           topicName,
				Description:    fmt.Sprintf("Events for the %s %s domain", adj, noun),
				OwnerServiceID: &ownerID,
				CreatedAt:      now,
			}
			topicID, err := InsertTopic(ctx, db, topic)
			if err != nil {
				return err
			}

			// 20 events per topic
			for _, action := range actions {
				eventName := action
				event := models.Event{
					TopicID:       topicID,
					EventName:     eventName,
					Description:   fmt.Sprintf("Fired when %s %s is %s", adj, noun, action),
					SchemaURL:     fmt.Sprintf("https://schemas.local/%s-%s/%s.json", adj, noun, action),
					SchemaVersion: &schemaVersion,
					Deprecated:    false,
					CreatedAt:     now,
				}
				eventID, err := InsertEvent(ctx, db, event)
				if err != nil {
					return err
				}

				// Producer for each event (using the owner service)
				prod := models.Producer{
					EventID:   eventID,
					ServiceID: ownerID,
					Writes:    true,
					CreatedAt: now,
				}
				if _, err := InsertProducer(ctx, db, prod); err != nil {
					return err
				}

				// Add a consumer sometimes to make it realistic
				if serviceIndex%3 == 0 {
					consumerSvcID := serviceIDs[(serviceIndex)%len(serviceIDs)]
					cons := models.Consumer{
						EventID:       eventID,
						ServiceID:     consumerSvcID,
						ConsumerGroup: fmt.Sprintf("%s-group", serviceNames[serviceIndex%len(serviceNames)]),
						Description:   fmt.Sprintf("Consumer for %s events", eventName),
						CreatedAt:     now,
					}
					if _, err := InsertConsumer(ctx, db, cons); err != nil {
						return err
					}
				}
			}
		}
	}

	// 100 undocumented consumers (10 * 10)
	secretAdjs := []string{"secret", "hidden", "ghost", "shadow", "dark", "phantom", "stray", "lost", "void", "null"}
	secretNouns := []string{"watcher", "listener", "sniffer", "prowler", "lurker", "stalker", "scout", "spy", "agent", "wraith"}

	for _, adj := range secretAdjs {
		for _, noun := range secretNouns {
			topicName := fmt.Sprintf("untracked.%s.%s", adj, noun)
			consumerGroup := fmt.Sprintf("%s-%s-group", adj, noun)
			query := `INSERT INTO missing_consumers (topic, consumer_group, created_at, updated_at) VALUES (?, ?, ?, ?)`
			if _, err := db.ExecContext(ctx, query, topicName, consumerGroup, now, now); err != nil {
				return fmt.Errorf("failed to insert missing consumer: %w", err)
			}
		}
	}

	return nil
}

// ensureCreatedAt sets the timestamp to now if it is zero.
func ensureCreatedAt(t *time.Time) {
	if t.IsZero() {
		*t = time.Now()
	}
}

// GetOrCreateService looks up a service by name and repository, or creates it if not found.
func GetOrCreateService(ctx context.Context, executor SQLExecutor, name, repository string) (models.Service, error) {
	query := `SELECT id, name, repository, created_at FROM services WHERE name = ? AND repository = ?`
	row := executor.QueryRowContext(ctx, query, name, repository)

	var svc models.Service
	err := row.Scan(&svc.ID, &svc.Name, &svc.Repository, &svc.CreatedAt)
	if err == nil {
		return svc, nil
	}
	if err != sql.ErrNoRows {
		return models.Service{}, err
	}

	now := time.Now()
	insertQuery := `INSERT INTO services (name, repository, created_at) VALUES (?, ?, ?)`
	result, err := executor.ExecContext(ctx, insertQuery, name, repository, now)
	if err != nil {
		return models.Service{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Service{}, err
	}

	return models.Service{ID: id, Name: name, Repository: repository, CreatedAt: now}, nil
}

// InsertTopic inserts a new topic into the database.
func InsertTopic(ctx context.Context, executor SQLExecutor, topic models.Topic) (int64, error) {
	ensureCreatedAt(&topic.CreatedAt)
	query := `INSERT INTO topics (name, description, owner_service_id, created_at) VALUES (?, ?, ?, ?)`

	res, err := executor.ExecContext(ctx, query, topic.Name, topic.Description, topic.OwnerServiceID, topic.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetTopic looks up a topic by name.
func GetTopic(ctx context.Context, executor SQLExecutor, topicName string) (*models.Topic, error) {
	query := `SELECT id, name, description, owner_service_id, created_at FROM topics WHERE name = ?`
	row := executor.QueryRowContext(ctx, query, topicName)

	var topic models.Topic
	err := row.Scan(&topic.ID, &topic.Name, &topic.Description, &topic.OwnerServiceID, &topic.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &topic, nil
}

// RemoveTopic deletes a topic by name.
func RemoveTopic(ctx context.Context, executor SQLExecutor, topicName string) error {
	query := `DELETE FROM topics WHERE name = ?`
	_, err := executor.ExecContext(ctx, query, topicName)
	return err
}

// CountTopics returns the total number of topics.
func CountTopics(ctx context.Context, executor SQLExecutor) (int, error) {
	var count int
	if err := executor.QueryRowContext(ctx, `SELECT COUNT(*) FROM topics`).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// ListTopicNamesPaginated returns topic names for a specific page, ordered by name.
func ListTopicNamesPaginated(ctx context.Context, executor SQLExecutor, limit, offset int) ([]string, error) {
	rows, err := executor.QueryContext(ctx, `SELECT name FROM topics ORDER BY name LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

// ListAllTopicNames returns all topic names ordered by name.
func ListAllTopicNames(ctx context.Context, executor SQLExecutor) ([]string, error) {
	rows, err := executor.QueryContext(ctx, `SELECT name FROM topics ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

// ListTopics returns all topics ordered by name.
func ListTopics(ctx context.Context, executor SQLExecutor) ([]models.Topic, error) {
	query := `SELECT id, name, description, owner_service_id, created_at FROM topics ORDER BY name`
	rows, err := executor.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var t models.Topic
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.OwnerServiceID, &t.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}

// InsertEvent inserts a new event into the database.
func InsertEvent(ctx context.Context, executor SQLExecutor, event models.Event) (int64, error) {
	ensureCreatedAt(&event.CreatedAt)
	query := `INSERT INTO events (topic_id, event_name, description, schema_url, schema_version, deprecated, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := executor.ExecContext(ctx, query, event.TopicID, event.EventName, event.Description, event.SchemaURL, event.SchemaVersion, event.Deprecated, event.CreatedAt)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// RemoveEvent deletes an event by ID.
func RemoveEvent(ctx context.Context, executor SQLExecutor, eventID int64) error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := executor.ExecContext(ctx, query, eventID)
	return err
}

// GetEventsByTopic returns all events for a given topic.
func GetEventsByTopic(ctx context.Context, executor SQLExecutor, topicID int64) ([]models.Event, error) {
	query := `SELECT id, topic_id, event_name, description, schema_url, schema_version, deprecated, created_at FROM events WHERE topic_id = ?`
	rows, err := executor.QueryContext(ctx, query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.TopicID, &event.EventName, &event.Description, &event.SchemaURL, &event.SchemaVersion, &event.Deprecated, &event.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// InsertProducer inserts a new producer into the database.
func InsertProducer(ctx context.Context, executor SQLExecutor, producer models.Producer) (int64, error) {
	ensureCreatedAt(&producer.CreatedAt)
	query := `INSERT INTO producers (event_id, service_id, writes, created_at) VALUES (?, ?, ?, ?)`
	result, err := executor.ExecContext(ctx, query, producer.EventID, producer.ServiceID, producer.Writes, producer.CreatedAt)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// RemoveProducer deletes a producer by ID.
func RemoveProducer(ctx context.Context, executor SQLExecutor, producerID int64) error {
	query := `DELETE FROM producers WHERE id = ?`
	_, err := executor.ExecContext(ctx, query, producerID)
	return err
}

// GetProducersByEvent returns all producers for a given event.
func GetProducersByEvent(ctx context.Context, executor SQLExecutor, eventID int64) ([]models.Producer, error) {
	query := `SELECT id, event_id, service_id, writes, created_at FROM producers WHERE event_id = ?`
	rows, err := executor.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var producers []models.Producer
	for rows.Next() {
		var producer models.Producer
		err := rows.Scan(&producer.ID, &producer.EventID, &producer.ServiceID, &producer.Writes, &producer.CreatedAt)
		if err != nil {
			return nil, err
		}
		producers = append(producers, producer)
	}

	return producers, nil
}

// InsertConsumer inserts a new consumer into the database.
func InsertConsumer(ctx context.Context, executor SQLExecutor, consumer models.Consumer) (int64, error) {
	ensureCreatedAt(&consumer.CreatedAt)
	query := `INSERT INTO consumers (event_id, service_id, consumer_group, description, event_version, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := executor.ExecContext(ctx, query, consumer.EventID, consumer.ServiceID, consumer.ConsumerGroup, consumer.Description, consumer.EventVersion, consumer.CreatedAt)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// RemoveConsumer deletes a consumer by ID.
func RemoveConsumer(ctx context.Context, executor SQLExecutor, consumerID int64) error {
	query := `DELETE FROM consumers WHERE id = ?`
	_, err := executor.ExecContext(ctx, query, consumerID)
	return err
}

// GetConsumersByEvent returns all consumers for a given event.
func GetConsumersByEvent(ctx context.Context, executor SQLExecutor, eventID int64) ([]models.Consumer, error) {
	query := `SELECT id, event_id, service_id, consumer_group, description, event_version, created_at FROM consumers WHERE event_id = ?`
	rows, err := executor.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consumers []models.Consumer
	for rows.Next() {
		var consumer models.Consumer
		err := rows.Scan(&consumer.ID, &consumer.EventID, &consumer.ServiceID, &consumer.ConsumerGroup, &consumer.Description, &consumer.EventVersion, &consumer.CreatedAt)
		if err != nil {
			return nil, err
		}
		consumers = append(consumers, consumer)
	}

	return consumers, nil
}

// GetOrCreateTopic looks up a topic by name, or creates it if not found.
func GetOrCreateTopic(ctx context.Context, executor SQLExecutor, topicName string, ownerServiceID *int64, description string) (models.Topic, error) {
	topic, err := GetTopic(ctx, executor, topicName)
	if err != nil {
		return models.Topic{}, err
	}

	if topic != nil {
		needsUpdate := false
		if ownerServiceID != nil && (topic.OwnerServiceID == nil || *ownerServiceID != *topic.OwnerServiceID) {
			topic.OwnerServiceID = ownerServiceID
			needsUpdate = true
		}
		if description != "" && description != topic.Description {
			topic.Description = description
			needsUpdate = true
		}
		if needsUpdate {
			if err := updateTopic(ctx, executor, topic.Name, topic.OwnerServiceID, topic.Description); err != nil {
				return models.Topic{}, err
			}
		}
		return *topic, nil
	}

	newTopic := models.Topic{
		Name:           topicName,
		Description:    description,
		OwnerServiceID: ownerServiceID,
		CreatedAt:      time.Now(),
	}

	id, err := InsertTopic(ctx, executor, newTopic)
	if err != nil {
		return models.Topic{}, err
	}

	newTopic.ID = id
	return newTopic, nil
}

// updateTopic updates the owner and description of an existing topic.
func updateTopic(ctx context.Context, executor SQLExecutor, topicName string, ownerServiceID *int64, description string) error {
	query := `UPDATE topics SET owner_service_id = ?, description = ? WHERE name = ?`
	_, err := executor.ExecContext(ctx, query, ownerServiceID, description, topicName)
	if err != nil {
		return fmt.Errorf("error updating topic: %w", err)
	}
	return nil
}

// GetOrCreateEvent looks up an event by name within a topic, or creates it if not found.
// If the event already exists, non-empty fields are updated (upsert).
func GetOrCreateEvent(ctx context.Context, executor SQLExecutor, topicID int64, eventName string, schemaVersion *string, schemaURL string, description string) (models.Event, error) {
	query := `SELECT id, topic_id, event_name, description, schema_url, schema_version, deprecated, created_at
              FROM events
              WHERE topic_id = ? AND event_name = ?`

	row := executor.QueryRowContext(ctx, query, topicID, eventName)

	var event models.Event
	err := row.Scan(
		&event.ID,
		&event.TopicID,
		&event.EventName,
		&event.Description,
		&event.SchemaURL,
		&event.SchemaVersion,
		&event.Deprecated,
		&event.CreatedAt,
	)

	if err == nil {
		needsUpdate := false
		if schemaURL != "" && schemaURL != event.SchemaURL {
			event.SchemaURL = schemaURL
			needsUpdate = true
		}
		if schemaVersion != nil && (event.SchemaVersion == nil || *schemaVersion != *event.SchemaVersion) {
			event.SchemaVersion = schemaVersion
			needsUpdate = true
		}
		if description != "" && description != event.Description {
			event.Description = description
			needsUpdate = true
		}
		if needsUpdate {
			updateQuery := `UPDATE events SET schema_url = ?, schema_version = ?, description = ? WHERE id = ?`
			if _, err := executor.ExecContext(ctx, updateQuery, event.SchemaURL, event.SchemaVersion, event.Description, event.ID); err != nil {
				return models.Event{}, fmt.Errorf("error updating event: %w", err)
			}
		}
		return event, nil
	}

	if err != sql.ErrNoRows {
		return models.Event{}, err
	}

	newEvent := models.Event{
		TopicID:       topicID,
		EventName:     eventName,
		Description:   description,
		SchemaURL:     schemaURL,
		SchemaVersion: schemaVersion,
		Deprecated:    false,
		CreatedAt:     time.Now(),
	}

	id, err := InsertEvent(ctx, executor, newEvent)
	if err != nil {
		return models.Event{}, err
	}

	newEvent.ID = id
	return newEvent, nil
}

// InsertEventHeader inserts a new header for an event.
func InsertEventHeader(ctx context.Context, executor SQLExecutor, header models.EventHeader) (int64, error) {
	ensureCreatedAt(&header.CreatedAt)
	query := `INSERT INTO event_headers (event_id, name, description, created_at) VALUES (?, ?, ?, ?)`
	result, err := executor.ExecContext(ctx, query, header.EventID, header.Name, header.Description, header.CreatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpsertEventHeader inserts a new event header or updates the existing one with the same event_id and name
func UpsertEventHeader(ctx context.Context, executor SQLExecutor, header models.EventHeader) error {
	ensureCreatedAt(&header.CreatedAt)
	query := `
		INSERT INTO event_headers (event_id, name, description, created_at) 
		VALUES (?, ?, ?, ?)
		ON CONFLICT(event_id, name) 
		DO UPDATE SET description = excluded.description
	`
	_, err := executor.ExecContext(ctx, query, header.EventID, header.Name, header.Description, header.CreatedAt)
	return err
}

// GetEventHeaders returns all headers for a given event.
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

// DeleteEventHeaders deletes all headers for a given event.
func DeleteEventHeaders(ctx context.Context, executor SQLExecutor, eventID int64) error {
	query := `DELETE FROM event_headers WHERE event_id = ?`
	_, err := executor.ExecContext(ctx, query, eventID)
	return err
}

// ListAllEvents returns all events with topic and header information.
func ListAllEvents(ctx context.Context, executor SQLExecutor) ([]models.EventRow, error) {
	query := `
		SELECT
			t.name AS topic_name,
			e.event_name,
			e.description AS event_description,
			e.schema_version,
			e.schema_url,
			eh.name AS header_name,
			eh.description AS header_description
		FROM events e
		JOIN topics t ON e.topic_id = t.id
		LEFT JOIN event_headers eh ON eh.event_id = e.id
		ORDER BY t.name, e.event_name, eh.name
	`
	rows, err := executor.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.EventRow
	for rows.Next() {
		var row models.EventRow
		if err := rows.Scan(&row.TopicName, &row.EventName, &row.EventDescription, &row.SchemaVersion, &row.SchemaURL, &row.HeaderName, &row.HeaderDescription); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// CountEvents returns the total number of events.
func CountEvents(ctx context.Context, executor SQLExecutor) (int, error) {
	var count int
	if err := executor.QueryRowContext(ctx, `SELECT COUNT(*) FROM events`).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// ListEventsPaginated returns events for a specific page.
func ListEventsPaginated(ctx context.Context, executor SQLExecutor, limit, offset int) ([]models.EventRow, error) {
	// We need to paginate on the EVENTS table first, then join headers.
	// 1. Select the event IDs for the current page, filtering/ordering as needed.
	// Since the main list is ordered by TopicName then EventName, we need to join topics for ordering.
	query := `
		SELECT
			t.name AS topic_name,
			e.event_name,
			e.description AS event_description,
			e.schema_version,
			e.schema_url,
			eh.name AS header_name,
			eh.description AS header_description
		FROM events e
		JOIN topics t ON e.topic_id = t.id
		LEFT JOIN event_headers eh ON eh.event_id = e.id
		WHERE e.id IN (
			SELECT e2.id
			FROM events e2
			JOIN topics t2 ON e2.topic_id = t2.id
			ORDER BY t2.name, e2.event_name
			LIMIT ? OFFSET ?
		)
		ORDER BY t.name, e.event_name, eh.name
	`
	rows, err := executor.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.EventRow
	for rows.Next() {
		var row models.EventRow
		if err := rows.Scan(&row.TopicName, &row.EventName, &row.EventDescription, &row.SchemaVersion, &row.SchemaURL, &row.HeaderName, &row.HeaderDescription); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// CountEventsSearch returns the total number of events matching a search term.
func CountEventsSearch(ctx context.Context, executor SQLExecutor, search string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM events e
		JOIN topics t ON e.topic_id = t.id
		WHERE e.event_name LIKE ? OR t.name LIKE ?
	`
	pattern := "%" + search + "%"
	var count int
	if err := executor.QueryRowContext(ctx, query, pattern, pattern).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// SearchEventsPaginated returns events matching a search term for a specific page.
func SearchEventsPaginated(ctx context.Context, executor SQLExecutor, search string, limit, offset int) ([]models.EventRow, error) {
	query := `
		SELECT
			t.name AS topic_name,
			e.event_name,
			e.description AS event_description,
			e.schema_version,
			e.schema_url,
			eh.name AS header_name,
			eh.description AS header_description
		FROM events e
		JOIN topics t ON e.topic_id = t.id
		LEFT JOIN event_headers eh ON eh.event_id = e.id
		WHERE e.id IN (
			SELECT e2.id
			FROM events e2
			JOIN topics t2 ON e2.topic_id = t2.id
			WHERE e2.event_name LIKE ? OR t2.name LIKE ?
			ORDER BY t2.name, e2.event_name
			LIMIT ? OFFSET ?
		)
		ORDER BY t.name, e.event_name, eh.name
	`
	pattern := "%" + search + "%"
	rows, err := executor.QueryContext(ctx, query, pattern, pattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.EventRow
	for rows.Next() {
		var row models.EventRow
		if err := rows.Scan(&row.TopicName, &row.EventName, &row.EventDescription, &row.SchemaVersion, &row.SchemaURL, &row.HeaderName, &row.HeaderDescription); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// ListEventsByTopicName returns all events (with headers) for a specific topic.
func ListEventsByTopicName(ctx context.Context, executor SQLExecutor, topicName string) ([]models.EventRow, error) {
	query := `
		SELECT
			t.name AS topic_name,
			e.event_name,
			e.description AS event_description,
			e.schema_version,
			e.schema_url,
			eh.name AS header_name,
			eh.description AS header_description
		FROM events e
		JOIN topics t ON e.topic_id = t.id
		LEFT JOIN event_headers eh ON eh.event_id = e.id
		WHERE t.name = ?
		ORDER BY e.event_name, eh.name
	`
	rows, err := executor.QueryContext(ctx, query, topicName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.EventRow
	for rows.Next() {
		var row models.EventRow
		if err := rows.Scan(&row.TopicName, &row.EventName, &row.EventDescription, &row.SchemaVersion, &row.SchemaURL, &row.HeaderName, &row.HeaderDescription); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// ListAllProducers returns all producers with service, topic, event, and header information.
func ListAllProducers(ctx context.Context, executor SQLExecutor) ([]models.ProducerRow, error) {
	query := `
		SELECT
			s.name AS service_name,
			s.repository,
			t.name AS topic_name,
			t.description AS topic_description,
			CASE WHEN t.owner_service_id = p.service_id THEN 1 ELSE 0 END AS is_owner,
			p.writes,
			e.event_name,
			e.description AS event_description,
			e.schema_version,
			e.schema_url,
			eh.name AS header_name,
			eh.description AS header_description
		FROM producers p
		JOIN services s ON p.service_id = s.id
		JOIN events e ON p.event_id = e.id
		JOIN topics t ON e.topic_id = t.id
		LEFT JOIN event_headers eh ON eh.event_id = e.id
		ORDER BY s.name, t.name, e.event_name, eh.name
	`
	rows, err := executor.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ProducerRow
	for rows.Next() {
		var row models.ProducerRow
		if err := rows.Scan(&row.ServiceName, &row.Repository, &row.TopicName, &row.TopicDescription, &row.Owner, &row.Writes, &row.EventName, &row.EventDescription, &row.SchemaVersion, &row.SchemaURL, &row.HeaderName, &row.HeaderDescription); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// CountServices returns the total number of services.
func CountServices(ctx context.Context, executor SQLExecutor) (int, error) {
	var count int
	if err := executor.QueryRowContext(ctx, "SELECT COUNT(*) FROM services").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// CountServicesSearch returns the total number of services matching a search term.
func CountServicesSearch(ctx context.Context, executor SQLExecutor, search string) (int, error) {
	query := "SELECT COUNT(*) FROM services WHERE name LIKE ?"
	pattern := "%" + search + "%"
	var count int
	if err := executor.QueryRowContext(ctx, query, pattern).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// ListServicesPaginated returns services for a specific page.
func ListServicesPaginated(ctx context.Context, executor SQLExecutor, limit, offset int) ([]models.Service, error) {
	query := `
		SELECT id, name, repository
		FROM services
		ORDER BY name
		LIMIT ? OFFSET ?
	`
	rows, err := executor.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Service
	for rows.Next() {
		var row models.Service
		if err := rows.Scan(&row.ID, &row.Name, &row.Repository); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// SearchServicesPaginated returns services matching a search term for a specific page.
func SearchServicesPaginated(ctx context.Context, executor SQLExecutor, search string, limit, offset int) ([]models.Service, error) {
	query := `
		SELECT id, name, repository
		FROM services
		WHERE name LIKE ?
		ORDER BY name
		LIMIT ? OFFSET ?
	`
	pattern := "%" + search + "%"
	rows, err := executor.QueryContext(ctx, query, pattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Service
	for rows.Next() {
		var row models.Service
		if err := rows.Scan(&row.ID, &row.Name, &row.Repository); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// ListTopicsForServices returns all topics produced by the given services.
func ListTopicsForServices(ctx context.Context, executor SQLExecutor, serviceIDs []int64) ([]models.ProducerListRow, error) {
	if len(serviceIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(serviceIDs))
	args := make([]interface{}, len(serviceIDs))
	for i, id := range serviceIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT
			p.service_id,
			s.name AS service_name,
			s.repository,
			e.topic_id,
			t.name AS topic_name,
			COUNT(DISTINCT e.id) AS event_count,
			MAX(CASE WHEN t.owner_service_id = p.service_id THEN 1 ELSE 0 END) AS is_owner,
			MAX(CASE WHEN p.writes THEN 1 ELSE 0 END) AS writes
		FROM producers p
		JOIN services s ON p.service_id = s.id
		JOIN events e ON p.event_id = e.id
		JOIN topics t ON e.topic_id = t.id
		WHERE p.service_id IN (%s)
		GROUP BY p.service_id, e.topic_id
		ORDER BY s.name, t.name
	`, strings.Join(placeholders, ","))

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ProducerListRow
	for rows.Next() {
		var row models.ProducerListRow
		if err := rows.Scan(&row.ServiceID, &row.ServiceName, &row.Repository, &row.TopicID, &row.TopicName, &row.EventCount, &row.Owner, &row.Writes); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// GetProducerDetail returns full event/header data for a specific (service, topic) producer pair.
func GetProducerDetail(ctx context.Context, executor SQLExecutor, serviceID, topicID int64) ([]models.ProducerRow, error) {
	query := `
		SELECT
			p.service_id,
			s.name AS service_name,
			s.repository,
			e.topic_id,
			t.name AS topic_name,
			t.description AS topic_description,
			CASE WHEN t.owner_service_id = p.service_id THEN 1 ELSE 0 END AS is_owner,
			p.writes,
			e.event_name,
			e.description AS event_description,
			e.schema_version,
			e.schema_url,
			eh.name AS header_name,
			eh.description AS header_description
		FROM producers p
		JOIN services s ON p.service_id = s.id
		JOIN events e ON p.event_id = e.id
		JOIN topics t ON e.topic_id = t.id
		LEFT JOIN event_headers eh ON eh.event_id = e.id
		WHERE p.service_id = ? AND e.topic_id = ?
		ORDER BY e.event_name, eh.name
	`
	rows, err := executor.QueryContext(ctx, query, serviceID, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ProducerRow
	for rows.Next() {
		var row models.ProducerRow
		if err := rows.Scan(&row.ServiceID, &row.ServiceName, &row.Repository, &row.TopicID, &row.TopicName, &row.TopicDescription, &row.Owner, &row.Writes, &row.EventName, &row.EventDescription, &row.SchemaVersion, &row.SchemaURL, &row.HeaderName, &row.HeaderDescription); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// ListAllConsumers returns all consumers with topic and event information.
func ListAllConsumers(ctx context.Context, executor SQLExecutor) ([]models.ConsumerRow, error) {
	query := `
		SELECT
			s.name AS service_name,
			s.repository,
			c.consumer_group,
			c.description,
			t.name AS topic_name,
			e.event_name,
			c.event_version
		FROM consumers c
		JOIN services s ON c.service_id = s.id
		JOIN events e ON c.event_id = e.id
		JOIN topics t ON e.topic_id = t.id
		ORDER BY s.name, c.consumer_group, t.name, e.event_name
	`
	rows, err := executor.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ConsumerRow
	for rows.Next() {
		var row models.ConsumerRow
		if err := rows.Scan(&row.ServiceName, &row.Repository, &row.ConsumerGroup, &row.Description, &row.TopicName, &row.EventName, &row.EventVersion); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// CountConsumerGroups returns the total number of distinct consumer groups.
func CountConsumerGroups(ctx context.Context, executor SQLExecutor) (int, error) {
	query := `SELECT COUNT(*) FROM (SELECT DISTINCT c.service_id, c.consumer_group FROM consumers c)`
	var count int
	if err := executor.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// ListConsumersPaginated returns consumers for a specific page of consumer groups.
// It paginates on distinct (service_id, consumer_group) pairs, then fetches all
// rows (events/topics) for those groups.
func ListConsumersPaginated(ctx context.Context, executor SQLExecutor, limit, offset int) ([]models.ConsumerRow, error) {
	query := `
		SELECT
			s.name AS service_name,
			s.repository,
			c.consumer_group,
			c.description,
			t.name AS topic_name,
			e.event_name,
			c.event_version
		FROM consumers c
		JOIN services s ON c.service_id = s.id
		JOIN events e ON c.event_id = e.id
		JOIN topics t ON e.topic_id = t.id
		WHERE (c.service_id, c.consumer_group) IN (
			SELECT DISTINCT c2.service_id, c2.consumer_group
			FROM consumers c2
			JOIN services s2 ON c2.service_id = s2.id
			ORDER BY s2.name, c2.consumer_group
			LIMIT ? OFFSET ?
		)
		ORDER BY s.name, c.consumer_group, t.name, e.event_name
	`
	rows, err := executor.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ConsumerRow
	for rows.Next() {
		var row models.ConsumerRow
		if err := rows.Scan(&row.ServiceName, &row.Repository, &row.ConsumerGroup, &row.Description, &row.TopicName, &row.EventName, &row.EventVersion); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// CountConsumerGroupsSearch returns the total number of distinct consumer groups matching a search term.
func CountConsumerGroupsSearch(ctx context.Context, executor SQLExecutor, search string) (int, error) {
	query := `
		SELECT COUNT(*) FROM (
			SELECT DISTINCT c.service_id, c.consumer_group
			FROM consumers c
			JOIN services s ON c.service_id = s.id
			JOIN events e ON c.event_id = e.id
			JOIN topics t ON e.topic_id = t.id
			WHERE c.consumer_group LIKE ? OR s.name LIKE ? OR t.name LIKE ? OR e.event_name LIKE ?
		)
	`
	pattern := "%" + search + "%"
	var count int
	if err := executor.QueryRowContext(ctx, query, pattern, pattern, pattern, pattern).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// SearchConsumersPaginated returns consumers matching a search term for a specific page.
func SearchConsumersPaginated(ctx context.Context, executor SQLExecutor, search string, limit, offset int) ([]models.ConsumerRow, error) {
	query := `
		SELECT
			s.name AS service_name,
			s.repository,
			c.consumer_group,
			c.description,
			t.name AS topic_name,
			e.event_name,
			c.event_version
		FROM consumers c
		JOIN services s ON c.service_id = s.id
		JOIN events e ON c.event_id = e.id
		JOIN topics t ON e.topic_id = t.id
		WHERE (c.service_id, c.consumer_group) IN (
			SELECT DISTINCT c2.service_id, c2.consumer_group
			FROM consumers c2
			JOIN services s2 ON c2.service_id = s2.id
			JOIN events e2 ON c2.event_id = e2.id
			JOIN topics t2 ON e2.topic_id = t2.id
			WHERE c2.consumer_group LIKE ? OR s2.name LIKE ? OR t2.name LIKE ? OR e2.event_name LIKE ?
			ORDER BY s2.name, c2.consumer_group
			LIMIT ? OFFSET ?
		)
		ORDER BY s.name, c.consumer_group, t.name, e.event_name
	`
	pattern := "%" + search + "%"
	rows, err := executor.QueryContext(ctx, query, pattern, pattern, pattern, pattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ConsumerRow
	for rows.Next() {
		var row models.ConsumerRow
		if err := rows.Scan(&row.ServiceName, &row.Repository, &row.ConsumerGroup, &row.Description, &row.TopicName, &row.EventName, &row.EventVersion); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// ListProducersByTopic returns all producers for a specific topic,
// including whether each producer owns the topic and the owner service name.
func ListProducersByTopic(ctx context.Context, executor SQLExecutor, topicName string) ([]models.TopicProducerRow, error) {
	query := `
		SELECT
			s.name AS service_name,
			s.repository,
			e.event_name,
			p.writes,
			CASE WHEN t.owner_service_id = p.service_id THEN 1 ELSE 0 END AS is_owner,
			os.name AS owner_service
		FROM producers p
		JOIN services s ON p.service_id = s.id
		JOIN events e ON p.event_id = e.id
		JOIN topics t ON e.topic_id = t.id
		LEFT JOIN services os ON os.id = t.owner_service_id
		WHERE t.name = ?
		ORDER BY s.name, e.event_name
	`
	rows, err := executor.QueryContext(ctx, query, topicName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TopicProducerRow
	for rows.Next() {
		var row models.TopicProducerRow
		if err := rows.Scan(&row.ServiceName, &row.Repository, &row.EventName, &row.Writes, &row.Owner, &row.OwnerService); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// ListConsumersByTopic returns all consumers for a specific topic.
func ListConsumersByTopic(ctx context.Context, executor SQLExecutor, topicName string) ([]models.TopicConsumerRow, error) {
	query := `
		SELECT
			s.name AS service_name,
			s.repository,
			e.event_name,
			c.consumer_group,
			c.event_version
		FROM consumers c
		JOIN services s ON c.service_id = s.id
		JOIN events e ON c.event_id = e.id
		JOIN topics t ON e.topic_id = t.id
		WHERE t.name = ?
		ORDER BY s.name, c.consumer_group, e.event_name
	`
	rows, err := executor.QueryContext(ctx, query, topicName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TopicConsumerRow
	for rows.Next() {
		var row models.TopicConsumerRow
		if err := rows.Scan(&row.ServiceName, &row.Repository, &row.EventName, &row.ConsumerGroup, &row.EventVersion); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// ListProducersByService returns all events produced by a specific service.
func ListProducersByService(ctx context.Context, executor SQLExecutor, serviceName string) ([]models.ServiceProducerRow, error) {
	query := `
		SELECT
			t.name AS topic_name,
			e.event_name,
			p.writes,
			CASE WHEN t.owner_service_id = p.service_id THEN 1 ELSE 0 END AS is_owner
		FROM producers p
		JOIN services s ON p.service_id = s.id
		JOIN events e ON p.event_id = e.id
		JOIN topics t ON e.topic_id = t.id
		WHERE s.name = ?
		ORDER BY t.name, e.event_name
	`
	rows, err := executor.QueryContext(ctx, query, serviceName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ServiceProducerRow
	for rows.Next() {
		var row models.ServiceProducerRow
		if err := rows.Scan(&row.TopicName, &row.EventName, &row.Writes, &row.Owner); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// ListConsumersByService returns all events consumed by a specific service.
func ListConsumersByService(ctx context.Context, executor SQLExecutor, serviceName string) ([]models.ServiceConsumerRow, error) {
	query := `
		SELECT
			t.name AS topic_name,
			e.event_name,
			c.consumer_group,
			c.event_version
		FROM consumers c
		JOIN services s ON c.service_id = s.id
		JOIN events e ON c.event_id = e.id
		JOIN topics t ON e.topic_id = t.id
		WHERE s.name = ?
		ORDER BY t.name, e.event_name
	`
	rows, err := executor.QueryContext(ctx, query, serviceName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ServiceConsumerRow
	for rows.Next() {
		var row models.ServiceConsumerRow
		if err := rows.Scan(&row.TopicName, &row.EventName, &row.ConsumerGroup, &row.EventVersion); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// ListUndocumentedConsumerGroups returns all consumer groups that are consuming events without documentation.
func ListUndocumentedConsumerGroups(ctx context.Context, executor SQLExecutor) ([]models.UndocumentedConsumerGroupRow, error) {
	query := `select topic, consumer_group, created_at, updated_at from missing_consumers`

	rows, err := executor.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.UndocumentedConsumerGroupRow
	for rows.Next() {
		var row models.UndocumentedConsumerGroupRow
		if err := rows.Scan(&row.Topic, &row.ConsumerGroup, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}
