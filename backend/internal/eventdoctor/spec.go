package eventdoctor

import (
	"fmt"
)

type EventDoctorSpec struct {
	Version   string     `yaml:"version" json:"version" validate:"required"`
	Service   string     `yaml:"service" json:"service" validate:"required"`
	Config    Config     `yaml:"config" json:"config" validate:"required"`
	Producers []Producer `yaml:"producers" json:"producers" validate:"dive"`
	Consumers []Consumer `yaml:"consumers" json:"consumers" validate:"dive"`
}

type Server struct {
	Environment string `yaml:"environment" json:"environment" validate:"required"`
	URL         string `yaml:"url" json:"url" validate:"required,url"`
}

type Config struct {
	Servers    []Server `yaml:"servers" json:"servers" validate:"required,min=1,dive"`
	Repository string   `yaml:"repository" json:"repository" validate:"required,url"`
}

type EventHeader struct {
	Name        string `yaml:"name" json:"name" validate:"required"`
	Description string `yaml:"description" json:"description"`
}

type Event struct {
	Name        string        `yaml:"name" json:"name" validate:"required"`
	Version     *string       `yaml:"version" json:"version" validate:"omitempty,semver"`
	Description string        `yaml:"description" json:"description"`
	SchemaURL   string        `yaml:"schema_url" json:"schema_url" validate:"required_if=Owner true,omitempty,url"`
	Headers     []EventHeader `yaml:"headers" json:"headers"`
}

type Producer struct {
	Topic       string  `yaml:"topic" json:"topic" validate:"required"`
	Owner       bool    `yaml:"owner" json:"owner"`
	Writes      bool    `yaml:"writes" json:"writes"`
	Description string  `yaml:"description" json:"description"`
	Events      []Event `yaml:"events" json:"events" validate:"required,dive"`
}

type ConsumerEvent struct {
	Name    string  `yaml:"name" json:"name" validate:"required"`
	Version *string `yaml:"version" json:"version" validate:"omitempty,semver"`
}

type Topic struct {
	Name   string          `yaml:"name" json:"name" validate:"required"`
	Events []ConsumerEvent `yaml:"events" json:"events" validate:"required,dive"`
}

type Consumer struct {
	Group       string  `yaml:"group" json:"group" validate:"required"`
	Description string  `yaml:"description" json:"description"`
	Topics      []Topic `yaml:"topics" json:"topics" validate:"required,dive"`
}

func (c *EventDoctorSpec) Validate() error {
	if err := validate.Struct(c); err != nil {
		return err
	}

	// Validação de ambientes únicos por servidor
	if err := validateUniqueEnvironments(c.Config.Servers); err != nil {
		return err
	}

	// Validação para garantir que eventos de produtores owner tenham schema_url
	if err := validateOwnerEventsSchema(c.Producers); err != nil {
		return err
	}

	return nil
}

// GetServerURL retorna a URL do servidor para o ambiente especificado
func (c *EventDoctorSpec) GetServerURL(env string) (string, error) {
	for _, srv := range c.Config.Servers {
		if srv.Environment == env {
			return srv.URL, nil
		}
	}
	return "", fmt.Errorf("no server found for environment %q", env)
}
