# IoT Device Simulator

A professional IoT device simulator developed in Go that manages multiple sensors with NATS communication and MongoDB persistence. It includes a complete Docker Compose setup for both development and production environments.

## 🚀 Quick Start

### 1. Start the Environment
```bash
# Start services (NATS + MongoDB + Client)
make up

# Or manually
docker-compose up -d
```

### 2. Compile and Run the Application
```bash
# Compile
make build

# Run
make run
```

### 3. Test the NATS Endpoints
```bash
# Enter the NATS client shell
make nats-shell

# Query device status
nats req iot.device-001.status ""

# Watch real-time readings
nats sub "iot.device-001.readings.>"
```

## 📖 Full Documentation

All detailed documentation, including architecture diagrams, the NATS command guide, and project requirements, can be found in the [`docs/`](./docs) directory.

- **[System Architecture](./docs/ARCHITECTURE.md)**: Technical design and diagrams.
- **[NATS Commands](./docs/nats-commands.md)**: A complete guide to NATS commands with examples.
- **[Project Requirements](./docs/REQUIREMENTS.md)**: A detailed description of the project requirements.

## 🔧 Main Commands

| Command | Description |
|---|---|
| `make up` | Start the Docker environment |
| `make down` | Stop the Docker environment |
| `make build` | Compile the application |
| `make run` | Run the application |
| `make test` | Execute tests |
| `make nats-shell` | Enter the NATS client shell |

## 🧪 Testing

```bash
# Run all tests
make test
```

## 📄 License

MIT License.