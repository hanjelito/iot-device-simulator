# Makefile - IoT Device Simulator

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

nats-monitor:
	timeout 30s docker exec nats-client nats sub "iot.device-001.readings.>" || true

# MongoDB
mongo-shell:
	docker exec -it mongodb mongosh iot_simulator

mongo-stats:
	docker exec mongodb mongosh iot_simulator --eval "db.stats()" --quiet

mongo-clean:
	@read -p "Are you sure? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	docker exec mongodb mongosh iot_simulator --eval "db.readings.deleteMany({sensor_id: /test/})" --quiet

# Tests
test:
	go test ./... -v

coverage:
	go test ./... -cover -v

test-integration: up
	@sleep 5
	$(MAKE) nats-test

# Logs y monitoreo
diagnose:
	$(DOCKER_COMPOSE) ps
	@docker exec nats-client nats req iot.device-001.status "" --timeout=2s || echo "NATS not responding"
	@docker exec mongodb mongosh iot_simulator --eval "db.runCommand('ping')" --quiet || echo "MongoDB not responding"
	@ps aux | grep "$(BINARY_NAME)" | grep -v grep || echo "Application not running"

app-logs:
	tail -f app.log 2>/dev/null || echo "No app.log file found"

clean-logs:
	rm -f app.log *.log

# Mantenimiento
clean:
	rm -f $(BINARY_NAME)

clean-all: clean clean-logs
	$(DOCKER_COMPOSE) down -v --remove-orphans
	docker system prune -f

deps:
	go mod tidy
	go mod download

fmt:
	go fmt ./...

lint:
	golangci-lint run

# Operaciones
start: deps up build
	@echo "System ready. Run with 'make run'"

app-restart:
	@pkill -f "$(BINARY_NAME)" 2>/dev/null || true
	@sleep 2
	$(MAKE) run-background

run-background: build
	@pkill -f "$(BINARY_NAME)" 2>/dev/null || true
	nohup ./$(BINARY_NAME) $(CONFIG_FILE) > app.log 2>&1 &

stop:
	@pkill -f "$(BINARY_NAME)" 2>/dev/null || echo "Application not running"

# Información
info:
	@echo "Project: $(DOCKER_PROJECT)"
	@echo "Binary: $(BINARY_NAME)"
	@echo "Config: $(CONFIG_FILE)"
	@go version
	@docker compose version

.DEFAULT_GOAL := help

.PHONY: build run run-dev up down restart nats-shell nats-test nats-monitor mongo-shell mongo-stats mongo-clean \
        test coverage test-integration diagnose app-logs clean-logs clean clean-all deps fmt lint start app-restart \
        run-background stop info
