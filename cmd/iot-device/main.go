package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/nats-io/nats.go"

	"iot-device-simulator/internal/config"
	"iot-device-simulator/internal/device"
)

func main() {
	// Cargar configuración - obligatorio
	configFile := ""
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	} else {
		programName := filepath.Base(os.Args[0])
		log.Fatalf("Configuration file is required. Usage: %s <config.yml>", programName)
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("Error loading config from %s: %v", configFile, err)
	}

	// Log view sensors active
	log.Printf("Loaded %d sensors from config file", len(cfg.Sensors))
	for _, sensor := range cfg.Sensors {
		log.Printf("  - %s (%s): %v enabled=%v", sensor.ID, sensor.Type, sensor.Frequency, sensor.Enabled)
	}

	// connect to NATS
	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		log.Fatal("Error connecting to NATS:", err)
	}
	defer nc.Close()

	// Crear y iniciar dispositivo
	dev := device.New(cfg, nc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dev.Start(ctx)

	log.Printf("NATS subjects:")
	log.Printf("  - iot.%s.config (get sensor configs)", dev.GetID())
	log.Printf("  - iot.%s.status (get device status)", dev.GetID())
	log.Printf("  - iot.%s.readings.* (sensor readings)", dev.GetID())

	// Esperar señal de interrupción
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	cancel()
}
