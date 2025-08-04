package sensor

import (
	"testing"
	"time"

	"iot-device-simulator/internal/config"
)

type mockStorage struct{}

func (m *mockStorage) SaveReading(reading Reading) error {
	return nil
}

func TestNew(t *testing.T) {
	cfg := config.SensorConfig{
		ID:        "test-sensor",
		Type:      "temperature",
		Enabled:   true,
		Frequency: time.Second,
		Min:       20.0,
		Max:       30.0,
		Unit:      "°C",
	}

	sensor := New(cfg, nil, &mockStorage{})
	
	if sensor.GetConfig().ID != "test-sensor" {
		t.Errorf("Expected sensor ID 'test-sensor', got '%s'", sensor.GetConfig().ID)
	}
}

func TestGenerateReading(t *testing.T) {
	cfg := config.SensorConfig{
		ID:   "test-sensor",
		Type: "temperature",
		Min:  20.0,
		Max:  30.0,
		Unit: "°C",
	}

	sensor := New(cfg, nil, &mockStorage{})
	reading := sensor.generateReading()

	if reading.SensorID != "test-sensor" {
		t.Errorf("Expected sensor ID 'test-sensor', got '%s'", reading.SensorID)
	}

	if reading.Type != "temperature" {
		t.Errorf("Expected type 'temperature', got '%s'", reading.Type)
	}

	if reading.Unit != "°C" {
		t.Errorf("Expected unit '°C', got '%s'", reading.Unit)
	}

	if reading.Error == "" && (reading.Value < 20.0 || reading.Value > 30.0) {
		t.Errorf("Value %.2f out of range [20.0, 30.0]", reading.Value)
	}
}

func TestUpdateFrequency(t *testing.T) {
	cfg := config.SensorConfig{
		ID:        "test-sensor",
		Frequency: time.Second,
	}

	sensor := New(cfg, nil, &mockStorage{})
	newFreq := 5 * time.Second

	sensor.UpdateFrequency(newFreq)

	if sensor.GetConfig().Frequency != newFreq {
		t.Errorf("Expected frequency %v, got %v", newFreq, sensor.GetConfig().Frequency)
	}
}

func TestUpdateThresholds(t *testing.T) {
	cfg := config.SensorConfig{
		ID:  "test-sensor",
		Min: 20.0,
		Max: 30.0,
	}

	sensor := New(cfg, nil, &mockStorage{})
	thresholds := map[string]interface{}{
		"min": 10.0,
		"max": 40.0,
	}

	sensor.UpdateThresholds(thresholds)
	config := sensor.GetConfig()

	if config.Min != 10.0 {
		t.Errorf("Expected min 10.0, got %.2f", config.Min)
	}

	if config.Max != 40.0 {
		t.Errorf("Expected max 40.0, got %.2f", config.Max)
	}
}