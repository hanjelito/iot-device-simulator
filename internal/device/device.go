// Package device manages the core logic of an IoT device.
// It orchestrates sensors, handles NATS subscriptions, and processes incoming requests.
package device

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"iot-device-simulator/internal/config"
	"iot-device-simulator/internal/sensor"
	"iot-device-simulator/internal/storage"
)

// Device represents a simulated IoT device.
// It holds the device's configuration, sensors, and connections to external services.
// It is the central component for managing the device's state and behavior.
type Device struct {
	id      string
	sensors []*sensor.Sensor
	nc      *nats.Conn
	storage *storage.MongoDB
}

// NewDevice creates and initializes a new Device based on the provided configuration.
// It sets up the device's sensors and establishes connections to NATS and MongoDB.
func NewDevice(cfg *config.Config, nc *nats.Conn, store *storage.MongoDB) *Device {
	device := &Device{
		id:      cfg.DeviceID,
		nc:      nc,
		storage: store,
	}

	// Create sensors from configuration
	for _, sensorConfig := range cfg.Sensors {
		s := sensor.New(sensorConfig, nc, store)
		device.sensors = append(device.sensors, s)
	}

	return device
}

// StartDevice begins the device's operation.
// It sets up NATS subscriptions and starts all enabled sensors in separate goroutines.
func (d *Device) StartDevice(ctx context.Context) {
	// Set up NATS subscriptions
	d.setupSubscriptions()

	// Start sensors
	enabledCount := 0
	for _, s := range d.sensors {
		go s.StartSensor(ctx, d.id)
		if s.GetConfig().Enabled {
			enabledCount++
		}
	}

	log.Printf("Device %s started with %d sensors (%d enabled, %d disabled)", d.id, len(d.sensors), enabledCount, len(d.sensors)-enabledCount)
}

// setupSubscriptions configures the NATS subscriptions for all device endpoints.
func (d *Device) setupSubscriptions() {
	// Get sensor configuration
	d.nc.Subscribe(fmt.Sprintf("iot.%s.config", d.id), d.handleConfig)

	// Get device status
	d.nc.Subscribe(fmt.Sprintf("iot.%s.status", d.id), d.handleStatus)

	// Update sensor configuration
	d.nc.Subscribe(fmt.Sprintf("iot.%s.config.update", d.id), d.handleConfigUpdate)

	// Register a new sensor
	d.nc.Subscribe(fmt.Sprintf("iot.%s.sensor.register", d.id), d.handleSensorRegister)

	// Get the latest readings for a sensor
	d.nc.Subscribe(fmt.Sprintf("iot.%s.readings.latest", d.id), d.handleLatestReadings)
}

// handleConfig responds with the current configuration of all sensors.
func (d *Device) handleConfig(msg *nats.Msg) {
	configs := make(map[string]interface{})
	for _, s := range d.sensors {
		configs[s.GetConfig().ID] = s.GetConfig()
	}

	data, _ := json.Marshal(configs)
	msg.Respond(data)
}

// handleStatus responds with the current operational status of the device.
func (d *Device) handleStatus(msg *nats.Msg) {
	enabledCount := 0
	for _, s := range d.sensors {
		if s.GetConfig().Enabled {
			enabledCount++
		}
	}

	status := map[string]interface{}{
		"device_id":        d.id,
		"total_sensors":    len(d.sensors),
		"enabled_sensors":  enabledCount,
		"disabled_sensors": len(d.sensors) - enabledCount,
		"timestamp":        time.Now(),
	}

	data, _ := json.Marshal(status)
	msg.Respond(data)
}

// handleConfigUpdate processes requests to update a sensor's configuration.
func (d *Device) handleConfigUpdate(msg *nats.Msg) {
	var updateRequest map[string]interface{}
	if err := json.Unmarshal(msg.Data, &updateRequest); err != nil {
		msg.Respond([]byte(`{"error": "invalid JSON"}`))
		return
	}

	sensorID, ok := updateRequest["sensor_id"].(string)
	if !ok {
		msg.Respond([]byte(`{"error": "sensor_id is required"}`))
		return
	}

	// Find the target sensor
	var targetSensor *sensor.Sensor
	for _, s := range d.sensors {
		if s.GetConfig().ID == sensorID {
			targetSensor = s
			break
		}
	}

	if targetSensor == nil {
		msg.Respond([]byte(`{"error": "sensor not found"}`))
		return
	}

	// Update frequency if provided
	if frequency, ok := updateRequest["frequency"]; ok {
		if freqStr, ok := frequency.(string); ok {
			if duration, err := time.ParseDuration(freqStr); err == nil {
				targetSensor.UpdateFrequency(duration)
			}
		}
	}

	// Update thresholds if provided
	thresholdUpdates := make(map[string]interface{})
	if min, ok := updateRequest["min"]; ok {
		thresholdUpdates["min"] = min
	}
	if max, ok := updateRequest["max"]; ok {
		thresholdUpdates["max"] = max
	}
	if len(thresholdUpdates) > 0 {
		targetSensor.UpdateThresholds(thresholdUpdates)
	}

	// Also handle thresholds object for backward compatibility
	if thresholds, ok := updateRequest["thresholds"]; ok {
		if threshMap, ok := thresholds.(map[string]interface{}); ok {
			targetSensor.UpdateThresholds(threshMap)
		}
	}

	// Save updated configuration to MongoDB if storage is available
	if d.storage != nil {
		configs := make(map[string]interface{})
		for _, s := range d.sensors {
			configs[s.GetConfig().ID] = s.GetConfig()
		}
		d.storage.SaveConfig(d.id, configs)
	}

	msg.Respond([]byte(`{"status": "updated"}`))
}

// handleSensorRegister processes requests to register a new sensor with the device.
func (d *Device) handleSensorRegister(msg *nats.Msg) {
	var registerRequest map[string]interface{}
	if err := json.Unmarshal(msg.Data, &registerRequest); err != nil {
		msg.Respond([]byte(`{"error": "invalid JSON"}`))
		return
	}

	sensorID, ok := registerRequest["sensor_id"].(string)
	if !ok {
		msg.Respond([]byte(`{"error": "sensor_id is required"}`))
		return
	}

	sensorType, ok := registerRequest["type"].(string)
	if !ok {
		msg.Respond([]byte(`{"error": "type is required"}`))
		return
	}

	// Check if sensor already exists
	for _, s := range d.sensors {
		if s.GetConfig().ID == sensorID {
			msg.Respond([]byte(`{"error": "sensor already exists"}`))
			return
		}
	}

	// Create a new sensor configuration from the request
	sensorConfig := config.SensorConfig{
		ID:        sensorID,
		Type:      sensorType,
		Enabled:   true,
		Frequency: 30 * time.Second, // Default frequency
		Min:       0,                // Default min
		Max:       100,              // Default max
		Unit:      "",
	}

	// Parse optional parameters from the request
	if frequency, ok := registerRequest["frequency"].(string); ok {
		if duration, err := time.ParseDuration(frequency); err == nil {
			sensorConfig.Frequency = duration
		}
	}
	if min, ok := registerRequest["min"].(float64); ok {
		sensorConfig.Min = min
	}
	if max, ok := registerRequest["max"].(float64); ok {
		sensorConfig.Max = max
	}
	if unit, ok := registerRequest["unit"].(string); ok {
		sensorConfig.Unit = unit
	}

	// Create and add the new sensor
	newSensor := sensor.New(sensorConfig, d.nc, d.storage)
	d.sensors = append(d.sensors, newSensor)

	// Start the new sensor immediately in a new goroutine
	if sensorConfig.Enabled {
		go newSensor.StartSensor(context.Background(), d.id)
		log.Printf("Started new sensor %s with frequency %v", sensorID, sensorConfig.Frequency)
	}

	// Save the updated device configuration to MongoDB
	if d.storage != nil {
		configs := make(map[string]interface{})
		for _, s := range d.sensors {
			configs[s.GetConfig().ID] = s.GetConfig()
		}
		d.storage.SaveConfig(d.id, configs)
	}

	response := map[string]interface{}{
		"status":    "registered",
		"sensor_id": sensorID,
		"config":    sensorConfig,
	}

	data, _ := json.Marshal(response)
	msg.Respond(data)
}

// handleLatestReadings responds with the most recent reading for a given sensor.
func (d *Device) handleLatestReadings(msg *nats.Msg) {
	var request map[string]interface{}
	if err := json.Unmarshal(msg.Data, &request); err != nil {
		msg.Respond([]byte(`{"error": "invalid JSON"}`))
		return
	}

	sensorID, ok := request["sensor_id"].(string)
	if !ok {
		msg.Respond([]byte(`{"error": "sensor_id is required"}`))
		return
	}

	// Find the target sensor
	var targetSensor *sensor.Sensor
	for _, s := range d.sensors {
		if s.GetConfig().ID == sensorID {
			targetSensor = s
			break
		}
	}

	if targetSensor == nil {
		msg.Respond([]byte(`{"error": "sensor not found"}`))
		return
	}

	// Get latest readings from MongoDB if storage is available
	if d.storage != nil {
		readings, err := d.storage.GetLatestReadings(sensorID, 1)
		if err != nil {
			msg.Respond([]byte(`{"error": "failed to retrieve readings"}`))
			return
		}

		if len(readings) == 0 {
			msg.Respond([]byte(`{"error": "no readings found"}`))
			return
		}

		data, _ := json.Marshal(map[string]interface{}{
			"sensor_id":      sensorID,
			"latest_reading": readings[0],
		})
		msg.Respond(data)
		return
	}

	msg.Respond([]byte(`{"error": "storage not available"}`))
}

// GetID returns the unique identifier of the device.
func (d *Device) GetID() string {
	return d.id
}
