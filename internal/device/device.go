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

func New(cfg *config.Config, nc *nats.Conn) *Device {
	device := &Device{
		id: cfg.DeviceID,
		nc: nc,
	}

	// Crear sensores
	for _, sensorConfig := range cfg.Sensors {
		s := sensor.New(sensorConfig, nc)
		device.sensors = append(device.sensors, s)
	}

	return device
}

func (d *Device) Start(ctx context.Context) {
	// Configurar suscripciones NATS
	d.setupSubscriptions()

	// Iniciar sensores
	for _, s := range d.sensors {
		go s.Start(ctx, d.id)
	}

	log.Printf("Device %s started with %d sensors", d.id, len(d.sensors))
}

func (d *Device) setupSubscriptions() {
	// Configuración de sensores
	d.nc.Subscribe(fmt.Sprintf("iot.%s.config", d.id), d.handleConfig)
	
	// Estado del dispositivo
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
	status := map[string]interface{}{
		"device_id": d.id,
		"sensors":   len(d.sensors),
		"timestamp": time.Now(),
	}
	
	data, _ := json.Marshal(status)
	msg.Respond(data)
}

func (d *Device) GetID() string {
	return d.id
}