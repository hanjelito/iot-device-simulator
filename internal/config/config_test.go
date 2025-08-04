// Package config_test contains the unit tests for the config package.
package config

import (
	"os"
	"testing"
	"time"
)

// TestLoad tests the Load function to ensure it correctly parses a YAML config file.
func TestLoad(t *testing.T) {
	// Create a temporary config file for testing.
	configContent := `device_id: test-device
nats:
  url: nats://localhost:4222

sensors:
  - id: temp-01
    type: temperature
    enabled: true
    frequency: 5s
    min: 18.0
    max: 25.0
    unit: "°C"
`

	tmpFile, err := os.CreateTemp("", "test-config-*.yml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Load the configuration from the temporary file.
	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	// Assertions to verify the loaded configuration.
	if cfg.DeviceID != "test-device" {
		t.Errorf("Expected device ID 'test-device', got '%s'", cfg.DeviceID)
	}

	if cfg.NATS.URL != "nats://localhost:4222" {
		t.Errorf("Expected NATS URL 'nats://localhost:4222', got '%s'", cfg.NATS.URL)
	}

	if len(cfg.Sensors) != 1 {
		t.Fatalf("Expected 1 sensor, got %d", len(cfg.Sensors))
	}

	sensor := cfg.Sensors[0]
	if sensor.ID != "temp-01" {
		t.Errorf("Expected sensor ID 'temp-01', got '%s'", sensor.ID)
	}

	if sensor.Frequency != 5*time.Second {
		t.Errorf("Expected frequency 5s, got %v", sensor.Frequency)
	}
}
