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
)

type Device struct {
	id      string
	sensors []*sensor.Sensor
	nc      *nats.Conn
}

func NewDivice(cfg *config.Config, nc *nats.Conn) *Device {
	device := &Device{
		id: cfg.DeviceID,
		nc: nc,
	}

	// create sensors
	for _, sensorConfig := range cfg.Sensors {
		s := sensor.New(sensorConfig, nc)
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

func (d *Device) GetID() string {
	return d.id
}
