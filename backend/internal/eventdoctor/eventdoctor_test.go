package eventdoctor

import (
	"strings"
	"testing"
)

func TestValidateSemver(t *testing.T) {
	tests := []struct {
		name    string
		version string
		valid   bool
	}{
		{"valid simple", "1.0.0", true},
		{"valid with zeros", "0.0.0", true},
		{"valid high numbers", "10.20.30", true},
		{"valid with prerelease", "1.0.0-alpha", true},
		{"valid with prerelease dot", "1.0.0-alpha.1", true},
		{"valid with build metadata", "1.0.0+build.123", true},
		{"valid with prerelease and build", "1.0.0-beta.1+build.456", true},
		{"invalid missing patch", "1.0", false},
		{"invalid missing minor and patch", "1", false},
		{"invalid empty", "", false},
		{"invalid letters", "abc", false},
		{"invalid leading zero major", "01.0.0", false},
		{"invalid leading zero minor", "1.01.0", false},
		{"invalid leading zero patch", "1.0.01", false},
		{"invalid prefix v", "v1.0.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := semverRegex.MatchString(tt.version)
			if got != tt.valid {
				t.Errorf("semverRegex.MatchString(%q) = %v, want %v", tt.version, got, tt.valid)
			}
		})
	}
}

func TestLoadSpecFromReader(t *testing.T) {
	t.Run("valid yaml", func(t *testing.T) {
		input := `
version: "1.0.0"
service: my-service
config:
  servers:
    - environment: production
      url: https://api.example.com
  repository: https://github.com/example/repo
producers:
  - topic: orders
    owner: true
    events:
      - name: OrderCreated
        schema_url: https://schema.example.com/order-created
consumers:
  - group: order-processor
    topics:
      - name: payments
        events:
          - name: PaymentReceived
`
		spec, err := LoadSpecFromReader(strings.NewReader(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if spec.Version != "1.0.0" {
			t.Errorf("Version = %q, want %q", spec.Version, "1.0.0")
		}
		if spec.Service != "my-service" {
			t.Errorf("Service = %q, want %q", spec.Service, "my-service")
		}
		if len(spec.Config.Servers) != 1 {
			t.Fatalf("expected 1 server, got %d", len(spec.Config.Servers))
		}
		if spec.Config.Servers[0].Environment != "production" {
			t.Errorf("Server environment = %q, want %q", spec.Config.Servers[0].Environment, "production")
		}
		if spec.Config.Repository != "https://github.com/example/repo" {
			t.Errorf("Repository = %q, want %q", spec.Config.Repository, "https://github.com/example/repo")
		}
		if len(spec.Producers) != 1 {
			t.Fatalf("expected 1 producer, got %d", len(spec.Producers))
		}
		if spec.Producers[0].Topic != "orders" {
			t.Errorf("Producer topic = %q, want %q", spec.Producers[0].Topic, "orders")
		}
		if !spec.Producers[0].Owner {
			t.Error("Producer.Owner = false, want true")
		}
		if len(spec.Consumers) != 1 {
			t.Fatalf("expected 1 consumer, got %d", len(spec.Consumers))
		}
		if spec.Consumers[0].Group != "order-processor" {
			t.Errorf("Consumer group = %q, want %q", spec.Consumers[0].Group, "order-processor")
		}
	})

	t.Run("empty reader", func(t *testing.T) {
		_, err := LoadSpecFromReader(strings.NewReader(""))
		if err == nil {
			t.Error("expected error for empty reader, got nil")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		_, err := LoadSpecFromReader(strings.NewReader("[\x00invalid"))
		if err == nil {
			t.Error("expected error for invalid YAML, got nil")
		}
	})

	t.Run("partial yaml loads without error", func(t *testing.T) {
		input := `
version: "2.0.0"
service: partial-service
`
		spec, err := LoadSpecFromReader(strings.NewReader(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if spec.Version != "2.0.0" {
			t.Errorf("Version = %q, want %q", spec.Version, "2.0.0")
		}
		if spec.Service != "partial-service" {
			t.Errorf("Service = %q, want %q", spec.Service, "partial-service")
		}
	})

	t.Run("parses event headers", func(t *testing.T) {
		input := `
version: "1.0.0"
service: header-service
config:
  servers:
    - environment: dev
      url: https://dev.example.com
  repository: https://github.com/example/repo
producers:
  - topic: events
    owner: true
    events:
      - name: UserCreated
        schema_url: https://schema.example.com/user
        headers:
          - name: X-Correlation-ID
            description: Correlation identifier
          - name: X-Source
`
		spec, err := LoadSpecFromReader(strings.NewReader(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		headers := spec.Producers[0].Events[0].Headers
		if len(headers) != 2 {
			t.Fatalf("expected 2 headers, got %d", len(headers))
		}
		if headers[0].Name != "X-Correlation-ID" {
			t.Errorf("Header[0].Name = %q, want %q", headers[0].Name, "X-Correlation-ID")
		}
		if headers[0].Description != "Correlation identifier" {
			t.Errorf("Header[0].Description = %q, want %q", headers[0].Description, "Correlation identifier")
		}
	})

	t.Run("parses consumer event version", func(t *testing.T) {
		input := `
version: "1.0.0"
service: versioned-consumer
config:
  servers:
    - environment: staging
      url: https://staging.example.com
  repository: https://github.com/example/repo
consumers:
  - group: my-group
    topics:
      - name: orders
        events:
          - name: OrderCreated
            version: "2.1.0"
`
		spec, err := LoadSpecFromReader(strings.NewReader(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		ev := spec.Consumers[0].Topics[0].Events[0]
		if ev.Version == nil {
			t.Fatal("expected event version, got nil")
		}
		if *ev.Version != "2.1.0" {
			t.Errorf("event version = %q, want %q", *ev.Version, "2.1.0")
		}
	})
}
