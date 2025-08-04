// Package storage_test contains the unit tests for the storage package.
package storage

import (
	"testing"
	"time"

	"iot-device-simulator/internal/sensor"
)

// TestMongoDB_SaveReading tests the SaveReading method of the MongoDB client.
// It skips the test if a connection to MongoDB cannot be established.
func TestMongoDB_SaveReading(t *testing.T) {
	// Attempt to connect to MongoDB. If it fails, skip the test.
	mongodb, err := NewMongoDB("mongodb://localhost:27017", "test_iot")
	if err != nil {
		t.Skip("MongoDB not available, skipping test")
	}
	defer mongodb.Close()

	reading := sensor.Reading{
		SensorID:  "test-sensor-reading",
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

// TestMongoDB_SaveConfig tests the SaveConfig method of the MongoDB client.
// It skips the test if a connection to MongoDB cannot be established.
func TestMongoDB_SaveConfig(t *testing.T) {
	// Attempt to connect to MongoDB. If it fails, skip the test.
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

	err = mongodb.SaveConfig("test-device-config", configs)
	if err != nil {
		t.Errorf("Error saving config: %v", err)
	}
}
