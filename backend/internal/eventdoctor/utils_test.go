package eventdoctor

import (
	"strings"
	"testing"
)

func TestValidateOwnerEventsSchema(t *testing.T) {
	t.Run("owner with schema_url is valid", func(t *testing.T) {
		producers := []Producer{
			{
				Topic: "orders",
				Owner: true,
				Events: []Event{
					{Name: "OrderCreated", SchemaURL: "https://schema.example.com/order"},
				},
			},
		}
		if err := validateOwnerEventsSchema(producers); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("owner without schema_url returns error", func(t *testing.T) {
		producers := []Producer{
			{
				Topic: "orders",
				Owner: true,
				Events: []Event{
					{Name: "OrderCreated", SchemaURL: ""},
				},
			},
		}
		err := validateOwnerEventsSchema(producers)
		if err == nil {
			t.Fatal("expected error for owner event without schema_url")
		}
		if !strings.Contains(err.Error(), "orders") {
			t.Errorf("error should mention topic name, got: %v", err)
		}
		if !strings.Contains(err.Error(), "OrderCreated") {
			t.Errorf("error should mention event name, got: %v", err)
		}
	})

	t.Run("non-owner without schema_url is valid", func(t *testing.T) {
		producers := []Producer{
			{
				Topic: "orders",
				Owner: false,
				Events: []Event{
					{Name: "OrderCreated", SchemaURL: ""},
				},
			},
		}
		if err := validateOwnerEventsSchema(producers); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty producers is valid", func(t *testing.T) {
		if err := validateOwnerEventsSchema(nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("multiple events one missing schema", func(t *testing.T) {
		producers := []Producer{
			{
				Topic: "orders",
				Owner: true,
				Events: []Event{
					{Name: "OrderCreated", SchemaURL: "https://schema.example.com/order"},
					{Name: "OrderUpdated", SchemaURL: ""},
				},
			},
		}
		err := validateOwnerEventsSchema(producers)
		if err == nil {
			t.Fatal("expected error when one owner event lacks schema_url")
		}
		if !strings.Contains(err.Error(), "OrderUpdated") {
			t.Errorf("error should mention the event without schema, got: %v", err)
		}
	})

	t.Run("multiple producers mixed ownership", func(t *testing.T) {
		producers := []Producer{
			{
				Topic: "orders",
				Owner: true,
				Events: []Event{
					{Name: "OrderCreated", SchemaURL: "https://schema.example.com/order"},
				},
			},
			{
				Topic: "payments",
				Owner: false,
				Events: []Event{
					{Name: "PaymentReceived", SchemaURL: ""},
				},
			},
		}
		if err := validateOwnerEventsSchema(producers); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestValidateUniqueEnvironments(t *testing.T) {
	t.Run("unique environments is valid", func(t *testing.T) {
		servers := []Server{
			{Environment: "production", URL: "https://api.example.com"},
			{Environment: "staging", URL: "https://staging.example.com"},
			{Environment: "development", URL: "https://dev.example.com"},
		}
		if err := validateUniqueEnvironments(servers); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("duplicate environments returns error", func(t *testing.T) {
		servers := []Server{
			{Environment: "production", URL: "https://api1.example.com"},
			{Environment: "production", URL: "https://api2.example.com"},
		}
		err := validateUniqueEnvironments(servers)
		if err == nil {
			t.Fatal("expected error for duplicate environments")
		}
		if !strings.Contains(err.Error(), "production") {
			t.Errorf("error should mention duplicate env name, got: %v", err)
		}
	})

	t.Run("single server is valid", func(t *testing.T) {
		servers := []Server{
			{Environment: "production", URL: "https://api.example.com"},
		}
		if err := validateUniqueEnvironments(servers); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty servers is valid", func(t *testing.T) {
		if err := validateUniqueEnvironments(nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("three servers with one duplicate", func(t *testing.T) {
		servers := []Server{
			{Environment: "production", URL: "https://api.example.com"},
			{Environment: "staging", URL: "https://staging.example.com"},
			{Environment: "production", URL: "https://api2.example.com"},
		}
		err := validateUniqueEnvironments(servers)
		if err == nil {
			t.Fatal("expected error for duplicate environments")
		}
	})
}
