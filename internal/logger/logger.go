// Package logger provides logging utilities for the application.
// It offers functions to log messages at various levels (info, warning, error) using structured logging.
// Uses the standard library slog package with a singleton pattern for global access.
package logger

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/lmittmann/tint"
)

// Level represents the logging level.
type Level = slog.Level

// Log levels.
const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// Format represents the output format for logs.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
	FormatTint Format = "tint"
)

// Config holds the configuration for the logger.
type Config struct {
	Component string
	Version   string
	Level     Level
	Format    Format
	Output    io.Writer
}

// DefaultConfig returns the default logger configuration.
func DefaultConfig() Config {
	return Config{
		Component: os.Getenv("COMPONENT"),
		Version:   os.Getenv("VERSION"),
		Level:     LevelInfo,
		Format:    FormatJSON,
		Output:    os.Stdout,
	}
}

var (
	instance *slog.Logger
	once     sync.Once
)

// Init initializes the global logger with the given configuration.
// This function is idempotent - it will only initialize the logger once.
// For testing or reconfiguration, use Reset() followed by Init().
func Init(cfg Config) {
	once.Do(func() {
		instance = newLogger(cfg)
	})
}

// Reset resets the singleton instance, allowing reinitialization.
// This is primarily useful for testing purposes.
func Reset() {
	once = sync.Once{}
	instance = nil
}

// newLogger creates a new slog.Logger with the given configuration.
func newLogger(cfg Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: cfg.Level,
	}

	var handler slog.Handler
	switch cfg.Format {
	case FormatText:
		handler = slog.NewTextHandler(cfg.Output, opts)
	case FormatTint:
		handler = tint.NewHandler(cfg.Output, &tint.Options{
			Level: cfg.Level,
		})
	default:
		handler = slog.NewJSONHandler(cfg.Output, opts)
	}

	l := slog.New(handler)

	if cfg.Component != "" || cfg.Version != "" {
		attrs := make([]any, 0, 4)
		if cfg.Component != "" {
			attrs = append(attrs, slog.String("component", cfg.Component))
		}
		if cfg.Version != "" {
			attrs = append(attrs, slog.String("version", cfg.Version))
		}
		l = l.With(attrs...)
	}

	return l
}

// Get returns the global logger instance.
// If Init() has not been called, it initializes the logger with default configuration.
func Get() *slog.Logger {
	if instance == nil {
		Init(DefaultConfig())
	}
	return instance
}
