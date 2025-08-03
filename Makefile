# Variables
BINARY_NAME=iot-device
CONFIG_FILE=cmd/iot-device/config.yml
MAIN_PATH=cmd/iot-device/main.go

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application with go run
run:
	@echo "Running IoT Device Simulator..."
	go run $(MAIN_PATH) $(CONFIG_FILE)

# Run the compiled binary
run-binary: build
	@echo "Running compiled binary..."
	./$(BINARY_NAME) $(CONFIG_FILE)

# Clean up built files
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

# Install dependencies
deps:
	go mod tidy

# Default target
.DEFAULT_GOAL := build

.PHONY: build run run-binary clean deps