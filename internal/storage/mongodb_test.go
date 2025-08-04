package storage

import (
	"testing"
	"time"

	"iot-device-simulator/internal/sensor"
)

func TestMongoDB_SaveReading(t *testing.T) {
	// Skip if no MongoDB available
	mongodb, err := NewMongoDB("mongodb://localhost:27017", "test_iot")
	if err != nil {
		t.Skip("MongoDB not available, skipping test")
	}
	defer mongodb.Close()

	reading := sensor.Reading{
		SensorID:  "test-sensor",
		Type:      "temperature",
		Value:     25.5,
		Unit:      "°C",
		Timestamp: time.Now(),
	}

	err = mongodb.SaveReading(reading)
	if err != nil {
		t.Errorf("Error saving reading: %v", err)
	}
}

func TestMongoDB_SaveConfig(t *testing.T) {
	// Skip if no MongoDB available
	mongodb, err := NewMongoDB("mongodb://localhost:27017", "test_iot")
	if err != nil {
		t.Skip("MongoDB not available, skipping test")
	}
	defer mongodb.Close()

	configs := map[string]any{
		"temp-01": map[string]any{
			"id":        "temp-01",
			"type":      "temperature",
			"frequency": "5s",
			"min":       20.0,
			"max":       30.0,
		},
	}

	err = mongodb.SaveConfig("test-device", configs)
	if err != nil {
		t.Errorf("Error saving config: %v", err)
	}
}