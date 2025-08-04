package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
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
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	if cfg.DeviceID != "test-device" {
		t.Errorf("Expected device ID 'test-device', got '%s'", cfg.DeviceID)
	}

	if cfg.NATS.URL != "nats://localhost:4222" {
		t.Errorf("Expected NATS URL 'nats://localhost:4222', got '%s'", cfg.NATS.URL)
	}

	if len(cfg.Sensors) != 1 {
		t.Errorf("Expected 1 sensor, got %d", len(cfg.Sensors))
	}

	sensor := cfg.Sensors[0]
	if sensor.ID != "temp-01" {
		t.Errorf("Expected sensor ID 'temp-01', got '%s'", sensor.ID)
	}

	if sensor.Frequency != 5*time.Second {
		t.Errorf("Expected frequency 5s, got %v", sensor.Frequency)
	}
}