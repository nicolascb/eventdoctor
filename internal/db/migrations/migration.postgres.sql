-- PostgreSQL migration for eventdoctor

CREATE TABLE IF NOT EXISTS topics (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    owner TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    topic_id INTEGER NOT NULL REFERENCES topics(id),
    event_type TEXT NOT NULL,
    schema_url TEXT NOT NULL,
    schema_version TEXT,
    deprecated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS event_headers (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_id, name)
);

CREATE TABLE IF NOT EXISTS producers (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id),
    service TEXT NOT NULL,
    writes BOOLEAN DEFAULT TRUE,
    repository TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS consumers (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id),
    service TEXT NOT NULL,
    consumer_group TEXT NOT NULL,
    event_version TEXT,
    repository TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS missing_topics (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    first_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS missing_events (
    id SERIAL PRIMARY KEY,
    topic TEXT NOT NULL,
    event_type TEXT NOT NULL,
    first_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(topic, event_type)
);

-- Indexes para otimizar queries
CREATE INDEX IF NOT EXISTS idx_events_topic_id ON events(topic_id);
CREATE INDEX IF NOT EXISTS idx_events_topic_event_type ON events(topic_id, event_type);
CREATE INDEX IF NOT EXISTS idx_event_headers_event_id ON event_headers(event_id);
CREATE INDEX IF NOT EXISTS idx_producers_event_id ON producers(event_id);
CREATE INDEX IF NOT EXISTS idx_producers_repository ON producers(repository);
CREATE INDEX IF NOT EXISTS idx_producers_service ON producers(service);
CREATE INDEX IF NOT EXISTS idx_consumers_event_id ON consumers(event_id);
CREATE INDEX IF NOT EXISTS idx_consumers_repository ON consumers(repository);
CREATE INDEX IF NOT EXISTS idx_consumers_service ON consumers(service);
CREATE INDEX IF NOT EXISTS idx_consumers_consumer_group ON consumers(consumer_group);
CREATE INDEX IF NOT EXISTS idx_topics_name ON topics(name);
CREATE INDEX IF NOT EXISTS idx_missing_topics_name ON missing_topics(name);
CREATE INDEX IF NOT EXISTS idx_missing_events_topic_event ON missing_events(topic, event_type);
