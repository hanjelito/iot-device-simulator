# Variables
BINARY_NAME=iot-device
CONFIG_FILE=cmd/iot-device/config.yml
MAIN_PATH=cmd/iot-device/main.go
DOCKER_COMPOSE=docker compose
DOCKER_PROJECT=iot-device-simulator

# Build y ejecución
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run: build
	./$(BINARY_NAME) $(CONFIG_FILE)

run-dev:
	go run $(MAIN_PATH) $(CONFIG_FILE)

# Docker
up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

restart: down up

# NATS
nats-shell:
	docker exec -it nats-client sh

nats-test:
	docker exec nats-client nats req iot.device-001.status "" --timeout=2s

# Mantenimiento
clean:
	rm -f $(BINARY_NAME)

deps:
	go mod tidy
	go mod download

.DEFAULT_GOAL := help

.PHONY: build run run-dev up down restart nats-shell nats-test nats-monitor mongo-shell mongo-stats mongo-clean \
        test coverage test-integration diagnose app-logs clean-logs clean clean-all deps fmt lint start app-restart \
        run-background stop info
