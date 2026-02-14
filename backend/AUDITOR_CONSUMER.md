# EventDoctor Auditor

The `eventdoctor-auditor` is a command line tool that acts as a Kafka event consumer.

Planned features roadmap:

| Feature                         | Description                                                                    | Status |
| ------------------------------- | ------------------------------------------------------------------------------ | ------ |
| Undocumented events database    | Enrichment of undocumented events database in EventDoctor                      | -      |
| Undocumented consumers database | Enrichment of undocumented consumers database in EventDoctor                   | -      |
| Kafka event validation          | Consumes events from Kafka topics and validates against schemas in EventDoctor | -      |


## Undocumented Consumers/Events/Topics

Considered as undocumented:

- Topic without `owner` defined in the EventDoctor configuration file.
- Event published to a topic that is not documented in EventDoctor.
- Consumer that consumes from a topic that is not documented in EventDoctor.

When the auditor identifies any of these situations, it will report them to the EventDoctor API to enrich its database.

Undocumented producers:

- We consume events from a topic that is not documented in EventDoctor.
- Save it in the events table.

Undocumented consumers:

- We consume events from a topic that does not have a consumer documented in EventDoctor.
- Save it in the missing_consumers table.

When saving a spec, it must check whether the consumer exists in the missing_consumers table. If it exists, it should remove it from that table and save it in the consumers table.