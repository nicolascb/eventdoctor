---
title: Producers & Consumers
description: Understanding the role of producers and consumers in EventDoctor.
---

## Producers

A **producer** is a service that writes events to a topic. In the `eventdoctor.yaml`, producers are declared under the `producers` section:

```yaml
producers:
  - topic: "order.events"
    owner: true
    writes: true
    description: "Events related to the order system"
    events:
      - name: "OrderCreated"
        version: "1.0.0"
        description: "Fired when a new order is created"
        schema_url: "https://schemas.example.com/order_created.json"
```

### Ownership

The `owner` flag indicates that this service **owns** the topic — it defines the schema and is the authority on the event format. A topic can have multiple producers, but typically only one owner.

When a producer has `owner: true`:

- A `schema_url` is **required** for every event.
- The service is shown as the **authoritative source** in the Web UI.

### Writes Flag

The `writes` flag indicates whether the service actively writes to the topic. A service might own a topic (define the schema) but delegate writing to another service.

## Consumers

A **consumer** is a service that reads events from one or more topics. Consumers are grouped by their **consumer group**:

```yaml
consumers:
  - group: "notification-group"
    description: "Service responsible for push notifications"
    topics:
      - name: "order.events"
        events:
          - name: "OrderCreated"
            version: "1.0.0"
          - name: "OrderCancelled"
            version: "1.0.0"
```

### Consumer Groups

The `group` field maps to the consumer group concept in message brokers (e.g., Kafka consumer groups). Events are load-balanced across instances within the same group.

### Cross-Topic Consumers

A single consumer group can subscribe to multiple topics. This is common for aggregation services (e.g., a notification service that reacts to events from multiple domains).

## Relationships in the UI

The Web UI provides several views to understand the relationships:

- **Topics View** — shows all topics with their producers and consumers.
- **Producers View** — lists all producing services and the topics they write to.
- **Consumers View** — lists all consumer groups and the topics they subscribe to.
- **Event Detail** — for any event, shows its schema, producers, and consumers side by side.
