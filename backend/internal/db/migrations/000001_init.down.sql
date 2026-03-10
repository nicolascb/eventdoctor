DROP INDEX IF EXISTS idx_topics_owner_service_id;
DROP INDEX IF EXISTS idx_consumers_service_id;
DROP INDEX IF EXISTS idx_producers_service_id;

DROP TABLE IF EXISTS missing_consumers;
DROP TABLE IF EXISTS consumers;
DROP TABLE IF EXISTS producers;
DROP TABLE IF EXISTS event_headers;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS topics;
