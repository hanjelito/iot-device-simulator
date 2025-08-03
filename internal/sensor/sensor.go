package sensor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/nats-io/nats.go"

	"iot-device-simulator/internal/config"
)

type Reading struct {
	SensorID  string    `json:"sensor_id"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

type Sensor struct {
	config config.SensorConfig
	nc     *nats.Conn
}

func New(sensorConfig config.SensorConfig, nc *nats.Conn) *Sensor {
	return &Sensor{config: sensorConfig, nc: nc}
}

func (s *Sensor) StartSensor(ctx context.Context, deviceID string) {
	if !s.config.Enabled {
		log.Printf("Sensor %s is disabled, not starting", s.config.ID)
		return
	}

	log.Printf("Starting sensor %s with frequency %v", s.config.ID, s.config.Frequency)
	ticker := time.NewTicker(s.config.Frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping sensor %s", s.config.ID)
			return
		case <-ticker.C:
			reading := s.generateReading()
			s.publish(reading, deviceID)
		}
	}
}

// generateReading simulates a sensor reading based on its configuration.
func (s *Sensor) generateReading() Reading {
	reading := Reading{
		SensorID:  s.config.ID,
		Type:      s.config.Type,
		Unit:      s.config.Unit,
		Timestamp: time.Now(),
	}

	// Simular error ocasional (5%)
	if rand.Float64() < 0.05 {
		reading.Error = "sensor communication error"
		return reading
	}

	// Generar valor aleatorio en el rango configurado
	reading.Value = s.config.Min + rand.Float64()*(s.config.Max-s.config.Min)
	return reading
}

func (s *Sensor) publish(reading Reading, deviceID string) {
	data, _ := json.Marshal(reading)
	subject := fmt.Sprintf("iot.%s.readings.%s", deviceID, reading.Type)
	if err := s.nc.Publish(subject, data); err != nil {
		log.Printf("Error publishing reading from %s: %v", reading.SensorID, err)
	}
}

func (s *Sensor) GetConfig() config.SensorConfig {
	return s.config
}
