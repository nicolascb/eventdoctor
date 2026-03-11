package eventdoctor

import (
	"testing"
)

func validSpec() EventDoctorSpec {
	return EventDoctorSpec{
		Version: "1.0.0",
		Service: "test-service",
		Config: Config{
			Servers: []Server{
				{Environment: "production", URL: "https://api.example.com"},
			},
			Repository: "https://github.com/example/repo",
		},
		Producers: []Producer{
			{
				Topic: "orders",
				Owner: true,
				Events: []Event{
					{
						Name:      "OrderCreated",
						SchemaURL: "https://schema.example.com/order-created",
					},
				},
			},
		},
		Consumers: []Consumer{
			{
				Group: "order-processor",
				Topics: []Topic{
					{
						Name: "payments",
						Events: []ConsumerEvent{
							{Name: "PaymentReceived"},
						},
					},
				},
			},
		},
	}
}

func TestEventDoctorSpec_Validate(t *testing.T) {
	t.Run("valid spec passes", func(t *testing.T) {
		spec := validSpec()
		if err := spec.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing version", func(t *testing.T) {
		spec := validSpec()
		spec.Version = ""
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing version")
		}
	})

	t.Run("missing service", func(t *testing.T) {
		spec := validSpec()
		spec.Service = ""
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing service")
		}
	})

	t.Run("missing servers", func(t *testing.T) {
		spec := validSpec()
		spec.Config.Servers = nil
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing servers")
		}
	})

	t.Run("empty servers slice", func(t *testing.T) {
		spec := validSpec()
		spec.Config.Servers = []Server{}
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for empty servers")
		}
	})

	t.Run("server missing environment", func(t *testing.T) {
		spec := validSpec()
		spec.Config.Servers = []Server{
			{Environment: "", URL: "https://api.example.com"},
		}
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing server environment")
		}
	})

	t.Run("server missing url", func(t *testing.T) {
		spec := validSpec()
		spec.Config.Servers = []Server{
			{Environment: "prod", URL: ""},
		}
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing server URL")
		}
	})

	t.Run("server invalid url", func(t *testing.T) {
		spec := validSpec()
		spec.Config.Servers = []Server{
			{Environment: "prod", URL: "not-a-url"},
		}
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for invalid server URL")
		}
	})

	t.Run("missing repository", func(t *testing.T) {
		spec := validSpec()
		spec.Config.Repository = ""
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing repository")
		}
	})

	t.Run("invalid repository url", func(t *testing.T) {
		spec := validSpec()
		spec.Config.Repository = "not-a-url"
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for invalid repository URL")
		}
	})

	t.Run("producer missing topic", func(t *testing.T) {
		spec := validSpec()
		spec.Producers[0].Topic = ""
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing producer topic")
		}
	})

	t.Run("producer missing events", func(t *testing.T) {
		spec := validSpec()
		spec.Producers[0].Events = nil
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing producer events")
		}
	})

	t.Run("producer event missing name", func(t *testing.T) {
		spec := validSpec()
		spec.Producers[0].Events = []Event{
			{Name: "", SchemaURL: "https://schema.example.com/ev"},
		}
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing event name")
		}
	})

	t.Run("producer event invalid semver", func(t *testing.T) {
		spec := validSpec()
		invalidVersion := "not-semver"
		spec.Producers[0].Events[0].Version = &invalidVersion
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for invalid semver")
		}
	})

	t.Run("producer event valid semver", func(t *testing.T) {
		spec := validSpec()
		validVersion := "2.0.0"
		spec.Producers[0].Events[0].Version = &validVersion
		if err := spec.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("owner producer event missing schema_url", func(t *testing.T) {
		spec := validSpec()
		spec.Producers[0].Owner = true
		spec.Producers[0].Events[0].SchemaURL = ""
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for owner producer event without schema_url")
		}
	})

	t.Run("non-owner producer event without schema_url is valid", func(t *testing.T) {
		spec := validSpec()
		spec.Producers[0].Owner = false
		spec.Producers[0].Events[0].SchemaURL = ""
		if err := spec.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("consumer missing group", func(t *testing.T) {
		spec := validSpec()
		spec.Consumers[0].Group = ""
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing consumer group")
		}
	})

	t.Run("consumer topic missing name", func(t *testing.T) {
		spec := validSpec()
		spec.Consumers[0].Topics[0].Name = ""
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing consumer topic name")
		}
	})

	t.Run("consumer topic missing events", func(t *testing.T) {
		spec := validSpec()
		spec.Consumers[0].Topics[0].Events = nil
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing consumer topic events")
		}
	})

	t.Run("consumer event missing name", func(t *testing.T) {
		spec := validSpec()
		spec.Consumers[0].Topics[0].Events = []ConsumerEvent{
			{Name: ""},
		}
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for missing consumer event name")
		}
	})

	t.Run("consumer event invalid semver", func(t *testing.T) {
		spec := validSpec()
		badVersion := "abc"
		spec.Consumers[0].Topics[0].Events = []ConsumerEvent{
			{Name: "SomeEvent", Version: &badVersion},
		}
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for invalid consumer event semver")
		}
	})

	t.Run("duplicate environments", func(t *testing.T) {
		spec := validSpec()
		spec.Config.Servers = []Server{
			{Environment: "production", URL: "https://api1.example.com"},
			{Environment: "production", URL: "https://api2.example.com"},
		}
		if err := spec.Validate(); err == nil {
			t.Error("expected validation error for duplicate environments")
		}
	})

	t.Run("multiple unique environments", func(t *testing.T) {
		spec := validSpec()
		spec.Config.Servers = []Server{
			{Environment: "production", URL: "https://api.example.com"},
			{Environment: "staging", URL: "https://staging.example.com"},
		}
		if err := spec.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("no producers is valid", func(t *testing.T) {
		spec := validSpec()
		spec.Producers = nil
		if err := spec.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("no consumers is valid", func(t *testing.T) {
		spec := validSpec()
		spec.Consumers = nil
		if err := spec.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestEventDoctorSpec_GetServerURL(t *testing.T) {
	spec := validSpec()
	spec.Config.Servers = []Server{
		{Environment: "production", URL: "https://api.example.com"},
		{Environment: "staging", URL: "https://staging.example.com"},
	}

	t.Run("found production", func(t *testing.T) {
		url, err := spec.GetServerURL("production")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "https://api.example.com" {
			t.Errorf("URL = %q, want %q", url, "https://api.example.com")
		}
	})

	t.Run("found staging", func(t *testing.T) {
		url, err := spec.GetServerURL("staging")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url != "https://staging.example.com" {
			t.Errorf("URL = %q, want %q", url, "https://staging.example.com")
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := spec.GetServerURL("development")
		if err == nil {
			t.Error("expected error for unknown environment")
		}
	})

	t.Run("empty environment", func(t *testing.T) {
		_, err := spec.GetServerURL("")
		if err == nil {
			t.Error("expected error for empty environment")
		}
	})
}
