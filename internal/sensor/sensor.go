// Package sensor proporciona la lógica para simular un sensor de IoT.
package sensor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/nats-io/nats.go"

	"iot-device-simulator/internal/config"
)

// Reading representa una única lectura de un sensor.
// Contiene el valor medido, timestamp, y posible información de error.
type Reading struct {
	SensorID  string    `json:"sensor_id" bson:"sensor_id"`
	Type      string    `json:"type" bson:"type"`
	Value     float64   `json:"value" bson:"value"`
	Unit      string    `json:"unit" bson:"unit"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Error     string    `json:"error,omitempty" bson:"error,omitempty"`
}

// Storage define la interfaz para el almacenamiento persistente de lecturas.
// Esto permite desacoplar el sensor de una implementación de base de datos específica.
type Storage interface {
	SaveReading(reading Reading) error
}

// Sensor simula un sensor de IoT. Es responsable de generar lecturas periódicas,
// publicarlas en NATS y guardarlas en un almacenamiento persistente.
// Es seguro para uso concurrente.
type Sensor struct {
	config  config.SensorConfig
	nc      *nats.Conn
	storage Storage
	mu      sync.RWMutex
}

// New crea y devuelve una nueva instancia de Sensor.
func New(sensorConfig config.SensorConfig, nc *nats.Conn, storage Storage) *Sensor {
	return &Sensor{config: sensorConfig, nc: nc, storage: storage}
}

// StartSensor inicia el ciclo de vida del sensor en una nueva goroutine.
// Genera lecturas a la frecuencia especificada en su configuración.
// Se detiene cuando el contexto es cancelado.
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

// generateReading simula una lectura de sensor basada en su configuración.
// Incluye una probabilidad del 5% de simular un error de comunicación.
func (s *Sensor) generateReading() Reading {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

// publish envía una lectura a través de NATS y la guarda en el almacenamiento si está configurado.
func (s *Sensor) publish(reading Reading, deviceID string) {
	// Publish to NATS first
	data, _ := json.Marshal(reading)
	// subject := fmt.Sprintf("iot.%s.readings.%s", deviceID, reading.Type)
	subject := fmt.Sprintf("iot.%s.readings.%s.%s", deviceID, reading.Type, reading.SensorID)
	if err := s.nc.Publish(subject, data); err != nil {
		log.Printf("Error publishing reading from %s: %v", reading.SensorID, err)
	}

	// Save to MongoDB if available
	if s.storage != nil {
		if err := s.storage.SaveReading(reading); err != nil {
			log.Printf("Error saving reading to storage: %v", err)
		}
	}
}

// GetConfig devuelve una copia de la configuración actual del sensor de forma segura.
func (s *Sensor) GetConfig() config.SensorConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

// UpdateFrequency actualiza la frecuencia de lectura del sensor de forma segura.
func (s *Sensor) UpdateFrequency(frequency time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.Frequency = frequency
	log.Printf("Sensor %s frequency updated to %v", s.config.ID, frequency)
}

// UpdateThresholds actualiza los umbrales (min/max) del sensor de forma segura.
func (s *Sensor) UpdateThresholds(thresholds map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if min, ok := thresholds["min"].(float64); ok {
		s.config.Min = min
	}
	if max, ok := thresholds["max"].(float64); ok {
		s.config.Max = max
	}

	log.Printf("Sensor %s thresholds updated: min=%.2f, max=%.2f", s.config.ID, s.config.Min, s.config.Max)
}
