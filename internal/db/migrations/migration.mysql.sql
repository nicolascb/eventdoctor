-- MySQL migration for eventdoctor

CREATE TABLE IF NOT EXISTS topics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    owner VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    topic_id INT NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    schema_url VARCHAR(500) NOT NULL,
    schema_version VARCHAR(50),
    deprecated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (topic_id) REFERENCES topics(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS event_headers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id),
    UNIQUE KEY unique_event_header (event_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS producers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT NOT NULL,
    service VARCHAR(255) NOT NULL,
    writes BOOLEAN DEFAULT TRUE,
    repository VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS consumers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT NOT NULL,
    service VARCHAR(255) NOT NULL,
    consumer_group VARCHAR(255) NOT NULL,
    event_version VARCHAR(50),
    repository VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS missing_topics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    first_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS missing_events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    topic VARCHAR(255) NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    first_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_missing_event (topic, event_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Indexes para otimizar queries
CREATE INDEX idx_events_topic_id ON events(topic_id);
CREATE INDEX idx_events_topic_event_type ON events(topic_id, event_type);
CREATE INDEX idx_event_headers_event_id ON event_headers(event_id);
CREATE INDEX idx_producers_event_id ON producers(event_id);
CREATE INDEX idx_producers_repository ON producers(repository);
CREATE INDEX idx_producers_service ON producers(service);
CREATE INDEX idx_consumers_event_id ON consumers(event_id);
CREATE INDEX idx_consumers_repository ON consumers(repository);
CREATE INDEX idx_consumers_service ON consumers(service);
CREATE INDEX idx_consumers_consumer_group ON consumers(consumer_group);
CREATE INDEX idx_topics_name ON topics(name);
CREATE INDEX idx_missing_topics_name ON missing_topics(name);
CREATE INDEX idx_missing_events_topic_event ON missing_events(topic, event_type);
