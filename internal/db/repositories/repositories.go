package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/nicolascb/eventdoctor/internal/db/models"
)

// SQLExecutor é uma interface que abstrai as operações SQL comuns entre *sql.DB e *sql.Tx
type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
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
	query := `INSERT INTO events (topic_id, event_type, schema_url, schema_version, deprecated, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := executor.ExecContext(ctx, query, event.TopicID, event.EventType, event.SchemaURL, event.SchemaVersion, event.Deprecated, event.CreatedAt)
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
	query := `SELECT id, topic_id, event_type, schema_url, schema_version, deprecated, created_at FROM events WHERE topic_id = ?`
	rows, err := executor.QueryContext(ctx, query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.TopicID, &event.EventType, &event.SchemaURL, &event.SchemaVersion, &event.Deprecated, &event.CreatedAt)
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
		return fmt.Errorf("erro ao atualizar owner do tópico: %w", err)
	}
	return nil
}

// GetOrCreateEvent busca um evento por tipo em um tópico ou cria um novo se não existir
func GetOrCreateEvent(ctx context.Context, executor SQLExecutor, topicID int64, eventType string, schemaVersion *string, schemaURL string) (models.Event, error) {
	// Primeiro tenta encontrar o evento existente
	query := `SELECT id, topic_id, event_type, schema_url, schema_version, deprecated, created_at 
              FROM events 
              WHERE topic_id = ? AND event_type = ?`

	row := executor.QueryRowContext(ctx, query, topicID, eventType)

	var event models.Event
	err := row.Scan(
		&event.ID,
		&event.TopicID,
		&event.EventType,
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
		EventType:     eventType,
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
		EventType:     eventType,
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

// DeleteProducersByRepository remove todos os produtores associados a um repositório
func DeleteProducersByRepository(ctx context.Context, tx SQLExecutor, repository string) (int64, error) {
	query := `DELETE FROM producers WHERE repository = $1`
	result, err := tx.ExecContext(ctx, query, repository)
	if err != nil {
		return 0, fmt.Errorf("erro ao excluir produtores: %w", err)
	}

	return result.RowsAffected()
}

// DeleteConsumersByRepository remove todos os consumidores associados a um repositório
func DeleteConsumersByRepository(ctx context.Context, tx SQLExecutor, repository string) (int64, error) {
	query := `DELETE FROM consumers WHERE repository = $1`
	result, err := tx.ExecContext(ctx, query, repository)
	if err != nil {
		return 0, fmt.Errorf("erro ao excluir consumidores: %w", err)
	}

	return result.RowsAffected()
}
