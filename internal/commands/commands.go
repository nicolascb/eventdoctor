package commands

import (
	"fmt"
	"os"

	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
)

func ValidateConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer file.Close()

	spec, err := eventdoctor.LoadSpecFromReader(file)
	if err != nil {
		return fmt.Errorf("erro ao carregar especificação: %w", err)
	}

	if err := spec.Validate(); err != nil {
		return fmt.Errorf("configuração inválida: %w", err)
	}

	return nil
}

func ApplyConfig(filePath string, serverURL string) error {
	return nil
}
