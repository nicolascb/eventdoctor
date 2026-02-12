package commands

import (
	"fmt"
	"io"

	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
	"github.com/nicolascb/eventdoctor/internal/logger"
)

// SpecUploader é a interface para enviar o spec para a API
type SpecUploader interface {
	UploadSpec(serverURL string, spec eventdoctor.EventDoctorSpec) error
}

type ConfigCommand struct {
	uploader SpecUploader
}

func NewConfigCommand(uploader SpecUploader) *ConfigCommand {
	return &ConfigCommand{
		uploader: uploader,
	}
}

// LoadAndValidate carrega e valida o spec a partir de um reader
func (c *ConfigCommand) LoadAndValidate(file io.Reader) (eventdoctor.EventDoctorSpec, error) {
	spec, err := eventdoctor.LoadSpecFromReader(file)
	if err != nil {
		return eventdoctor.EventDoctorSpec{}, fmt.Errorf("fail to parse config: %w", err)
	}

	if err := spec.Validate(); err != nil {
		return eventdoctor.EventDoctorSpec{}, fmt.Errorf("invalid config: %w", err)
	}

	return spec, nil
}

// Validate carrega e valida o arquivo de configuração, exibindo informações
func (c *ConfigCommand) Validate(file io.Reader) error {
	spec, err := c.LoadAndValidate(file)
	if err != nil {
		return err
	}

	logger.Get().Info("Configuration file is valid!")
	logger.Get().Info(fmt.Sprintf("Producers: %d", len(spec.Producers)))
	logger.Get().Info(fmt.Sprintf("Consumers: %d", len(spec.Consumers)))

	return nil
}

// Apply carrega, valida e envia o spec para o servidor do ambiente especificado
func (c *ConfigCommand) Apply(env string, file io.Reader) error {
	spec, err := c.LoadAndValidate(file)
	if err != nil {
		return err
	}

	serverURL, err := spec.GetServerURL(env)
	if err != nil {
		return err
	}

	logger.Get().Info(fmt.Sprintf("Applying configuration to environment: %s (%s)", env, serverURL))

	if err := c.uploader.UploadSpec(serverURL, spec); err != nil {
		return fmt.Errorf("failed to apply configuration: %w", err)
	}

	logger.Get().Info("Configuration applied successfully!")
	return nil
}
