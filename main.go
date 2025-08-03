package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"gopkg.in/yaml.v3"
)

// Config representa la configuración completa del dispositivo
type Config struct {
	DeviceID string         `yaml:"device_id"`
	NATS     NATSConfig     `yaml:"nats"`
	Sensors  []SensorConfig `yaml:"sensors"`
}

type NATSConfig struct {
	URL string `yaml:"url"`
}

type SensorConfig struct {
	ID        string        `yaml:"id"`
	Type      string        `yaml:"type"`
	Frequency time.Duration `yaml:"frequency"`
	Min       float64       `yaml:"min"`
	Max       float64       `yaml:"max"`
	Unit      string        `yaml:"unit"`
	Enabled   bool          `yaml:"enabled"`
}

type Reading struct {
	SensorID  string    `json:"sensor_id"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

// Sensor simula un sensor IoT
type Sensor struct {
	config SensorConfig
	nc     *nats.Conn
}

func NewSensor(config SensorConfig, nc *nats.Conn) *Sensor {
	return &Sensor{config: config, nc: nc}
}

func (s *Sensor) Start(ctx context.Context, deviceID string) {
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

// Device representa el dispositivo IoT principal
type Device struct {
	id      string
	sensors []*Sensor
	nc      *nats.Conn
}

func NewDevice(config Config, nc *nats.Conn) *Device {
	device := &Device{
		id: config.DeviceID,
		nc: nc,
	}

	// Crear sensores
	for _, sensorConfig := range config.Sensors {
		sensor := NewSensor(sensorConfig, nc)
		device.sensors = append(device.sensors, sensor)
	}

	return device
}

func (d *Device) Start(ctx context.Context) {
	// Configurar suscripciones NATS
	d.setupSubscriptions()

	// Iniciar sensores
	for _, sensor := range d.sensors {
		go sensor.Start(ctx, d.id)
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
	for _, sensor := range d.sensors {
		configs[sensor.config.ID] = sensor.config
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

// loadConfig carga la configuración desde archivo YAML
func loadConfig(filename string) (Config, error) {
	var config Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	// Log para debug
	log.Printf("Loaded %d sensors from config file", len(config.Sensors))
	for _, sensor := range config.Sensors {
		log.Printf("  - %s (%s): %v enabled=%v", sensor.ID, sensor.Type, sensor.Frequency, sensor.Enabled)
	}

	return config, nil
}

func main() {
	// Cargar configuración - obligatorio
	configFile := "config.yml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading config from %s: %v", configFile, err)
	}

	// Conectar a NATS
	nc, err := nats.Connect(config.NATS.URL)
	if err != nil {
		log.Fatal("Error connecting to NATS:", err)
	}
	defer nc.Close()

	// Crear y iniciar dispositivo
	device := NewDevice(config, nc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	device.Start(ctx)

	log.Printf("NATS subjects:")
	log.Printf("  - iot.%s.config (get sensor configs)", config.DeviceID)
	log.Printf("  - iot.%s.status (get device status)", config.DeviceID)
	log.Printf("  - iot.%s.readings.* (sensor readings)", config.DeviceID)

	// Esperar señal de interrupción
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	cancel()
}
