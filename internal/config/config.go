package config

import (
	"io"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Version   string      `yaml:"version"`
	Metadata  Metadata    `yaml:"metadata"`
	Producers []Producers `yaml:"producers"`
	Consumers []Consumers `yaml:"consumers"`
}
type Metadata struct {
	ServerURL  string `yaml:"server_url"`
	Repository string `yaml:"repository"`
}
type Events struct {
	Type        string `yaml:"type"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	SchemaURL   string `yaml:"schema_url"`
}
type Producers struct {
	Name        string   `yaml:"name"`
	Topic       string   `yaml:"topic"`
	Owner       bool     `yaml:"owner"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Events      []Events `yaml:"events"`
}
type ConsumerEvents struct {
	Type string `yaml:"type"`
}

type Topics struct {
	Name   string           `yaml:"name"`
	Events []ConsumerEvents `yaml:"events"`
}
type Consumers struct {
	Service     string   `yaml:"service"`
	Group       string   `yaml:"group"`
	Description string   `yaml:"description"`
	Topics      []Topics `yaml:"topics"`
}

func LoadConfigFromReader(f io.Reader) (Config, error) {
	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	// TODO: Implementar validações conforme necessário
	return nil
}
