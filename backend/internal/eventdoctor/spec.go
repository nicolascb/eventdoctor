package eventdoctor

import (
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

var (
	validate           *validator.Validate
	semverRegex        = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	ErrNoOwnerForTopic = errors.New("deve existir um produtor owner para cada tópico")
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

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("semver", validateSemver)
}

func validateSemver(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return semverRegex.MatchString(value)
}

func LoadSpecFromReader(f io.Reader) (EventDoctorSpec, error) {
	var cfg EventDoctorSpec
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return EventDoctorSpec{}, err
	}

	return cfg, nil
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

// validateOwnerEventsSchema verifica se todos os eventos de produtores owner têm schema_url definido
func validateOwnerEventsSchema(producers []Producer) error {
	for _, producer := range producers {
		if producer.Owner {
			for _, event := range producer.Events {
				if event.SchemaURL == "" {
					return fmt.Errorf("o produtor é owner do tópico '%s', mas o evento '%s' não possui schema_url definido",
						producer.Topic, event.Name)
				}
			}
		}
	}
	return nil
}

// validateUniqueEnvironments verifica se cada ambiente é único
func validateUniqueEnvironments(servers []Server) error {
	environments := make(map[string]bool)

	for _, server := range servers {
		if environments[server.Environment] {
			return fmt.Errorf("env duplicate: '%s'", server.Environment)
		}
		environments[server.Environment] = true
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
