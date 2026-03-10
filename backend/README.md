# EventDoctor Backend

Backend API, CLI, and consumer for EventDoctor written in Go.

## Quick Start

### Prerequisites

- Go 1.25+
- SQLite3

### Installation

```bash
go mod download
```

### Running the API

```bash
go run ./cmd/api/main.go
```

The API will be available at `http://localhost:8087`.

### Running the Consumer

```bash
go run ./cmd/consumer/main.go
```

### Using the CLI

```bash
go run ./cmd/cli/main.go --help
```

## Building

```bash
go build -o eventdoctor ./cmd/api/main.go
go build -o eventdoctor-cli ./cmd/cli/main.go
go build -o eventdoctor-consumer ./cmd/consumer/main.go
```

## Documentation

- [CLI Documentation](./CLI.md)
- [Consumer Documentation](./CONSUMER.md)

## Testing

```bash
go test ./...
```

## Project Structure

- `cmd/` - Command entry points (api, cli, consumer)
- `internal/` - Internal packages
  - `api/` - HTTP API handlers
  - `client/` - Client logic
  - `commands/` - CLI commands
  - `db/` - Database layer
  - `eventdoctor/` - Core domain logic
  - `logger/` - Logging utilities
  - `service/` - Business logic

## API Endpoints

- `POST /v1/config` - Upload configuration
- `GET /v1/producers` - List producers
- `GET /v1/events` - List events
- `GET /v1/consumers` - List consumers
