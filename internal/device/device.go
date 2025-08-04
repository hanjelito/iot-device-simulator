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

type Device struct {
	id      string
	sensors []*sensor.Sensor
	nc      *nats.Conn
	storage *storage.MongoDB
}

func NewDivice(cfg *config.Config, nc *nats.Conn, store *storage.MongoDB) *Device {
	device := &Device{
		id:      cfg.DeviceID,
		nc:      nc,
		storage: store,
	}

	// create sensors
	for _, sensorConfig := range cfg.Sensors {
		s := sensor.New(sensorConfig, nc, store)
		device.sensors = append(device.sensors, s)
	}

	return device
}

func (d *Device) StartDivice(ctx context.Context) {
	// Configurar suscripciones NATS
	d.setupSubscriptions()

	// Iniciar sensores
	enabledCount := 0
	for _, s := range d.sensors {
		go s.StartSensor(ctx, d.id)
		if s.GetConfig().Enabled {
			enabledCount++
		}
	}

	log.Printf("Device %s started with %d sensors (%d enabled, %d disabled)", d.id, len(d.sensors), enabledCount, len(d.sensors)-enabledCount)
}

func (d *Device) setupSubscriptions() {
	// Sensor configuration
	d.nc.Subscribe(fmt.Sprintf("iot.%s.config", d.id), d.handleConfig)

	// Device status
	d.nc.Subscribe(fmt.Sprintf("iot.%s.status", d.id), d.handleStatus)

	// Update sensor configuration
	d.nc.Subscribe(fmt.Sprintf("iot.%s.config.update", d.id), d.handleConfigUpdate)

	// Register new sensor configuration
	d.nc.Subscribe(fmt.Sprintf("iot.%s.sensor.register", d.id), d.handleSensorRegister)

	// Get latest readings by sensor
	d.nc.Subscribe(fmt.Sprintf("iot.%s.readings.latest", d.id), d.handleLatestReadings)
}

func (d *Device) handleConfig(msg *nats.Msg) {
	configs := make(map[string]interface{})
	for _, s := range d.sensors {
		configs[s.GetConfig().ID] = s.GetConfig()
	}

	data, _ := json.Marshal(configs)
	msg.Respond(data)
}

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

func (d *Device) handleConfigUpdate(msg *nats.Msg) {
	var updateRequest map[string]interface{}
	if err := json.Unmarshal(msg.Data, &updateRequest); err != nil {
		msg.Respond([]byte(`{"error": "invalid JSON"}`))
		return
	}

	sensorID, ok := updateRequest["sensor_id"].(string)
	if !ok {
		msg.Respond([]byte(`{"error": "sensor_id required"}`))
		return
	}

	// Find sensor
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
	if thresholds, ok := updateRequest["thresholds"]; ok {
		if threshMap, ok := thresholds.(map[string]interface{}); ok {
			targetSensor.UpdateThresholds(threshMap)
		}
	}

	// Save updated configuration to MongoDB
	if d.storage != nil {
		configs := make(map[string]interface{})
		for _, s := range d.sensors {
			configs[s.GetConfig().ID] = s.GetConfig()
		}
		d.storage.SaveConfig(d.id, configs)
	}

	msg.Respond([]byte(`{"status": "updated"}`))
}

func (d *Device) handleSensorRegister(msg *nats.Msg) {
	var registerRequest map[string]interface{}
	if err := json.Unmarshal(msg.Data, &registerRequest); err != nil {
		msg.Respond([]byte(`{"error": "invalid JSON"}`))
		return
	}

	sensorID, ok := registerRequest["sensor_id"].(string)
	if !ok {
		msg.Respond([]byte(`{"error": "sensor_id required"}`))
		return
	}

	sensorType, ok := registerRequest["type"].(string)
	if !ok {
		msg.Respond([]byte(`{"error": "type required"}`))
		return
	}

	// Check if sensor already exists
	for _, s := range d.sensors {
		if s.GetConfig().ID == sensorID {
			msg.Respond([]byte(`{"error": "sensor already exists"}`))
			return
		}
	}

	// Create sensor config
	sensorConfig := config.SensorConfig{
		ID:        sensorID,
		Type:      sensorType,
		Enabled:   true,
		Frequency: 30 * time.Second,
		Min:       0,
		Max:       100,
		Unit:      "",
	}

	// Parse optional parameters
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

	// Create and add new sensor
	newSensor := sensor.New(sensorConfig, d.nc, d.storage)
	d.sensors = append(d.sensors, newSensor)

	// Start the new sensor immediately if enabled
	if sensorConfig.Enabled {
		go newSensor.StartSensor(context.Background(), d.id)
		log.Printf("Started new sensor %s with frequency %v", sensorID, sensorConfig.Frequency)
	}

	// Save updated configuration to MongoDB
	if d.storage != nil {
		configs := make(map[string]interface{})
		for _, s := range d.sensors {
			configs[s.GetConfig().ID] = s.GetConfig()
		}
		d.storage.SaveConfig(d.id, configs)
	}

	response := map[string]interface{}{
		"status": "registered",
		"sensor_id": sensorID,
		"config": sensorConfig,
	}
	
	data, _ := json.Marshal(response)
	msg.Respond(data)
}

func (d *Device) handleLatestReadings(msg *nats.Msg) {
	var request map[string]interface{}
	if err := json.Unmarshal(msg.Data, &request); err != nil {
		msg.Respond([]byte(`{"error": "invalid JSON"}`))
		return
	}

	sensorID, ok := request["sensor_id"].(string)
	if !ok {
		msg.Respond([]byte(`{"error": "sensor_id required"}`))
		return
	}

	// Find sensor
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

	// Get latest readings from MongoDB
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
			"sensor_id": sensorID,
			"latest_reading": readings[0],
		})
		msg.Respond(data)
		return
	}

	msg.Respond([]byte(`{"error": "storage not available"}`))
}

func (d *Device) GetID() string {
	return d.id
}
