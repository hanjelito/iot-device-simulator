# NATS Commands for the IoT Device Simulator

## Initial Setup

### Connecting to the NATS client
```bash
# Enter the container with the NATS CLI (recommended)
docker exec -it nats-client sh

# The NATS_URL variable is already configured automatically.
# You don't need to export it; just use the commands directly.
```

---

## 1. Main NATS Endpoints (5/5 Implemented) ✅

### 1.1 Get Configuration of All Sensors
```bash
nats req iot.device-001.config ""
```

**Expected Response:**
```json
{
  "temp-01": {
    "id": "temp-01",
    "type": "temperature",
    "frequency": "5s",
    "min": 15,
    "max": 35,
    "unit": "°C",
    "enabled": true
  },
  "temp-02": {
    "id": "temp-02",
    "type": "temperature",
    "frequency": "3s",
    "min": 18,
    "max": 28,
    "unit": "°C",
    "enabled": true
  }
  // ... more sensors
}
```

### 1.2 Get Device Status
```bash
nats req iot.device-001.status ""
```

**Expected Response:**
```json
{
  "device_id": "device-001",
  "total_sensors": 5,
  "enabled_sensors": 4,
  "disabled_sensors": 1,
  "timestamp": "2025-08-04T10:30:00Z"
}
```

### 1.3 Register a New Sensor ⭐ NEW
```bash
nats req iot.device-001.sensor.register '{
  "sensor_id": "temp-05",
  "type": "temperature",
  "frequency": "10s",
  "min": -10,
  "max": 50,
  "unit": "°C"
}'
```

**Expected Response:**
```json
{
  "status": "registered",
  "sensor_id": "temp-05",
  "config": {
    "id": "temp-05",
    "type": "temperature",
    "frequency": "10s",
    "min": -10,
    "max": 50,
    "unit": "°C",
    "enabled": true
  }
}
```

### 1.4 Update Configuration of an Existing Sensor
```bash
nats req iot.device-001.config.update '{
  "sensor_id": "temp-01",
  "frequency": "2s",
  "min": 10,
  "max": 40
}'
```

**Expected Response:**
```json
{
  "status": "updated"
}
```

### 1.5 Get Latest Readings for a Sensor ⭐ NEW
```bash
nats req iot.device-001.readings.latest '{
  "sensor_id": "temp-01"
}'
```

**Expected Response (if data exists):**
```json
{
  "sensor_id": "temp-01",
  "latest_reading": {
    "sensor_id": "temp-01",
    "type": "temperature",
    "value": 25.4,
    "unit": "°C",
    "timestamp": "2025-08-04T10:30:15Z",
    "error": ""
  }
}
```

**Response if no data exists:**
```json
{
  "error": "no readings found"
}
```

---

## 2. Real-Time Monitoring

### 2.1 Subscribe to All Device Readings
```bash
nats sub "iot.device-001.readings.>"
```

### 2.2 Subscribe to Temperature Readings Only
```bash
nats sub "iot.device-001.readings.temperature"
```

### 2.3 Subscribe to Humidity Readings Only
```bash
nats sub "iot.device-001.readings.humidity"
```

### 2.4 Subscribe to Pressure Readings Only
```bash
nats sub "iot.device-001.readings.pressure"
```

### 2.5 Subscribe to a Specific Sensor
```bash
nats sub "iot.device-001.readings.temperature.temp-01"
```

**Example of a received reading:**
```json
{
  "sensor_id": "temp-02",
  "type": "temperature",
  "value": 27.71,
  "unit": "°C",
  "timestamp": "2025-08-04T10:28:15.428987+02:00"
}
```

---

## 3. Complete Use Cases

### 3.1 Full Flow: Register → Monitor → Query
```bash
# 1. Register a new sensor
nats req iot.device-001.sensor.register '{
  "sensor_id": "light-01",
  "type": "light",
  "frequency": "5s",
  "min": 0,
  "max": 1000,
  "unit": "lux"
}'

# 2. Monitor its readings in real-time
nats sub "iot.device-001.readings.light"

# 3. After a moment, query its latest reading
nats req iot.device-001.readings.latest '{
  "sensor_id": "light-01"
}'
```

### 3.2 Configuration Management
```bash
# 1. View current configuration
nats req iot.device-001.config ""

# 2. Update a sensor's frequency
nats req iot.device-001.config.update '{
  "sensor_id": "humidity-01",
  "frequency": "5s"
}'

# 3. Verify the change was applied
nats req iot.device-001.config ""
```

### 3.3 Device Performance Analysis
```bash
# 1. View general status
nats req iot.device-001.status ""

# 2. Monitor all readings for 30 seconds
timeout 30s nats sub "iot.device-001.readings.>"

# 3. Query latest readings of critical sensors
nats req iot.device-001.readings.latest '{"sensor_id": "temp-01"}'
nats req iot.device-001.readings.latest '{"sensor_id": "temp-02"}'
nats req iot.device-001.readings.latest '{"sensor_id": "humidity-01"}'
```

---

## 4. Error Handling

### 4.1 Common Errors and Solutions

**Error: "No responders are available"**
```bash
# Cause: The Go application is not running.
# Solution: Verify that ./iot-device is running.
ps aux | grep iot-device | grep -v grep
```

**Error: "sensor not found"**
```bash
# Cause: The sensor_id does not exist.
# Solution: View the list of available sensors.
nats req iot.device-001.config ""
```

**Error: "no readings found"**
```bash
# Cause: The sensor is new and has not generated any readings yet.
# Solution: Wait a few seconds or verify that it is enabled.
nats req iot.device-001.config ""
```

### 4.2 Error Simulation
```bash
# Sensors automatically simulate a 5% error rate.
# To see errors, monitor the readings:
nats sub "iot.device-001.readings.>" | grep "error"
```

---

## 5. MongoDB Persistence

### 5.1 Verifying Data in MongoDB
```bash
# Connect to MongoDB
docker exec -it mongodb mongosh iot_simulator

# View the most recent readings
db.readings.find().sort({timestamp: -1}).limit(5)

# Count readings by sensor
db.readings.aggregate([
  {$group: {_id: "$sensor_id", count: {$sum: 1}}},
  {$sort: {count: -1}}
])

# View saved configurations
db.configurations.find()
```

### 5.2 Clearing Test Data
```bash
# Connect to MongoDB
docker exec -it mongodb mongosh iot_simulator

# Delete readings from test sensors
db.readings.deleteMany({"sensor_id": /test/})

# Clear all readings (CAUTION!)
db.readings.deleteMany({})
```

---

## 6. Web Monitoring

### 6.1 NATS Server Monitoring
```bash
# Open NATS monitoring in a browser
open http://localhost:8222

# Or verify via curl
curl -s http://localhost:8222/varz | jq
```

### 6.2 Statistics via API
```bash
# General server information
curl -s http://localhost:8222/varz | jq '{
  server_id: .server_id,
  version: .version,
  connections: .connections,
  subscriptions: .subscriptions
}'

# Active connections
curl -s http://localhost:8222/connz | jq '.connections | length'

# Active subscriptions
curl -s http://localhost:8222/subsz | jq '.subscriptions | length'
```

---

## 7. NATS Subject Summary

### Request/Response (Synchronous)
- `iot.device-001.config` - Get configuration
- `iot.device-001.status` - Get status
- `iot.device-001.sensor.register` - Register a sensor
- `iot.device-001.config.update` - Update configuration
- `iot.device-001.readings.latest` - Get latest readings

### Publish/Subscribe (Asynchronous)
- `iot.device-001.readings.temperature` - Temperature readings
- `iot.device-001.readings.humidity` - Humidity readings
- `iot.device-001.readings.pressure` - Pressure readings
- `iot.device-001.readings.>` - All readings (wildcard)

---

## 8. Important Notes

### Device ID
- The device uses the ID: `device-001` (from config.yml)
- All subjects include this ID.
- For multiple devices, change the ID in config.yml.

### Supported Sensor Types
- `temperature` - Temperature in °C
- `humidity` - Humidity in %
- `pressure` - Pressure in hPa
- `light` - Light in lux (customizable)
- Any custom type

### Persistence
- Readings are automatically saved to MongoDB.
- Configurations are persisted upon update.
- Data survives system restarts.

### Performance
- The system handles multiple sensors simultaneously.
- Configurable frequencies per sensor (1s minimum).
- Automatic error simulation (~5%).
- Real-time response (<5ms typical).

---

## 9. Makefile Commands

### Useful Commands
```bash
# Start the complete environment (NATS + MongoDB + Client)
make up

# Stop the environment
make down

# Enter the NATS client shell
make nats-shell

# Compile and run the application
make run

# Run tests
make test
```

---

## 10. Next Steps

### Future Improvements
- [ ] **Alerts**: Publish to `iot.device-001.alerts`
- [ ] **Metrics**: Publish to `iot.device-001.metrics`
- [ ] **Commands**: `iot.device-001.command.{action}`
- [ ] **Security**: Authentication with NKEYs
- [ ] **Clustering**: Support for multiple devices

---

## 11. Advanced Troubleshooting

### Connection Debugging
```bash
# Check application logs
tail -f app.log

# Check NATS server logs
docker logs nats-server

# Check MongoDB logs
docker logs mongodb
```

### Forcing Reconnection
```bash
# Restart the Go application
make restart

# Restart the NATS server
docker restart nats-server
```

### Cleaning the Environment
```bash
# Stop and remove containers
make down

# Clean build artifacts
make clean

# Start again
make up
```

---

## 12. Command Summary

| Command | Description |
|---------|-------------|
| `nats req iot.device-001.config ""` | Get configuration |
| `nats req iot.device-001.status ""` | Get status |
| `nats req iot.device-001.sensor.register '{'...'}'` | Register sensor |
| `nats req iot.device-001.config.update '{'...'}'` | Update sensor |
| `nats req iot.device-001.readings.latest '{'...'}'` | Get latest readings |
| `nats sub "iot.device-001.readings.>" ` | Monitor readings |

---
