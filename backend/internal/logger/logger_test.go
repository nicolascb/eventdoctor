package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	Reset()

	var buf bytes.Buffer
	cfg := Config{
		Level:  LevelInfo,
		Format: FormatJSON,
		Output: &buf,
	}

	Init(cfg)

	Get().Info("test message", "key", "value")

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("failed to parse log entry: %v", err)
	}

	if logEntry["msg"] != "test message" {
		t.Errorf("expected msg 'test message', got %v", logEntry["msg"])
	}

	if logEntry["key"] != "value" {
		t.Errorf("expected key 'value', got %v", logEntry["key"])
	}
}

func TestSingleton(t *testing.T) {
	Reset()

	var buf bytes.Buffer
	cfg := Config{
		Level:  LevelInfo,
		Format: FormatJSON,
		Output: &buf,
	}

	Init(cfg)
	logger1 := Get()

	// Try to init again with different config
	Init(Config{
		Level:  LevelDebug,
		Format: FormatText,
		Output: &buf,
	})
	logger2 := Get()

	if logger1 != logger2 {
		t.Error("expected singleton to return the same instance")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Level != LevelInfo {
		t.Errorf("expected default level Info, got %v", cfg.Level)
	}

	if cfg.Format != FormatJSON {
		t.Errorf("expected default format JSON, got %v", cfg.Format)
	}

	if cfg.Output == nil {
		t.Error("expected default output to be non-nil")
	}
}

func TestTextFormat(t *testing.T) {
	Reset()

	var buf bytes.Buffer
	cfg := Config{
		Level:  LevelInfo,
		Format: FormatText,
		Output: &buf,
	}

	Init(cfg)
	Get().Info("text format test")

	output := buf.String()
	if !strings.Contains(output, "text format test") {
		t.Errorf("expected output to contain 'text format test', got %s", output)
	}

	// Text format should not be valid JSON
	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err == nil {
		t.Error("expected text format to not be valid JSON")
	}
}

func TestLogLevels(t *testing.T) {
	Reset()

	var buf bytes.Buffer
	cfg := Config{
		Level:  LevelWarn,
		Format: FormatJSON,
		Output: &buf,
	}

	Init(cfg)

	log := Get()

	// Debug and Info should not be logged
	log.Debug("debug message")
	log.Info("info message")

	if buf.Len() > 0 {
		t.Error("expected no output for debug/info when level is warn")
	}

	// Warn and Error should be logged
	log.Warn("warn message")
	output := buf.String()
	if !strings.Contains(output, "warn message") {
		t.Errorf("expected warn message in output, got %s", output)
	}

	buf.Reset()
	log.Error("error message")
	output = buf.String()
	if !strings.Contains(output, "error message") {
		t.Errorf("expected error message in output, got %s", output)
	}
}

func TestWith(t *testing.T) {
	Reset()

	var buf bytes.Buffer
	cfg := Config{
		Level:  LevelInfo,
		Format: FormatJSON,
		Output: &buf,
	}

	Init(cfg)

	childLogger := Get().With("component", "test")
	childLogger.Info("child logger message")

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("failed to parse log entry: %v", err)
	}

	if logEntry["component"] != "test" {
		t.Errorf("expected component 'test', got %v", logEntry["component"])
	}
}

func TestWithGroup(t *testing.T) {
	Reset()

	var buf bytes.Buffer
	cfg := Config{
		Level:  LevelInfo,
		Format: FormatJSON,
		Output: &buf,
	}

	Init(cfg)

	groupLogger := Get().WithGroup("request")
	groupLogger.Info("grouped message", "id", "123")

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("failed to parse log entry: %v", err)
	}

	request, ok := logEntry["request"].(map[string]any)
	if !ok {
		t.Fatalf("expected request group in output, got %v", logEntry)
	}

	if request["id"] != "123" {
		t.Errorf("expected id '123' in request group, got %v", request["id"])
	}
}

func TestAttributes(t *testing.T) {
	Reset()

	var buf bytes.Buffer
	cfg := Config{
		Level:  LevelInfo,
		Format: FormatJSON,
		Output: &buf,
	}

	Init(cfg)

	Get().Info("attribute test",
		slog.String("name", "john"),
		slog.Int("age", 30),
		slog.Int64("id", 12345678901234),
		slog.Bool("active", true),
	)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("failed to parse log entry: %v", err)
	}

	if logEntry["name"] != "john" {
		t.Errorf("expected name 'john', got %v", logEntry["name"])
	}

	if int(logEntry["age"].(float64)) != 30 {
		t.Errorf("expected age 30, got %v", logEntry["age"])
	}

	if logEntry["active"] != true {
		t.Errorf("expected active true, got %v", logEntry["active"])
	}
}

func TestGetInitializesDefault(t *testing.T) {
	Reset()

	// Get() should initialize with defaults if not initialized
	log := Get()
	if log == nil {
		t.Error("expected Get() to return a non-nil logger")
	}
}
