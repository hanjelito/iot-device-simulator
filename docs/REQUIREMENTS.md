# GOLANG SOFTWARE DEVELOPER EVALUATION
## Go Application for Reading Sensor Data from an IoT Device

### Task Description
Implement an application in Go that simulates the management of an IoT device with multiple sensors of different types (temperature, humidity, pressure, etc.).

### Exact Project Requirements (Official Document)

#### NATS API - Specific Required Endpoints
- **Endpoint to register a sensor's configuration**
- **Endpoint to update a sensor's configuration**
  - Sampling frequency
  - Alert thresholds
- **Endpoint to query the current state of configurations**
- **Endpoint to query the latest values read by a sensor**
- **Publish changes in sensor readings to NATS**

#### Persistence
- **Storage of sensor reading data**
  - In a format deemed most appropriate (MongoDB was chosen)

#### Reading Simulation
- **A component that emulates the periodic reading of sensors**
  - Different for each sensor type
  - Allow configuration of parameters via the API
- **Simulation of sensor reading errors**

#### Constraints
- Messaging protocol: **NATS** (official nats.go library)

#### Evaluation Criteria
- Code readability and style
- Modularity and code organization
- Code coverage with tests
- Good use of version control and change history
- Documentation and comments associated with the code

#### Deliverables
1. Source code in a Git repository
2. Associated documentation (README.md)
3. Explanatory diagram of the proposed solution

---

## Current Implementation Status

### ✅ COMPLETE

#### NATS API - Implemented Endpoints (5/5) ✅
- ✅ **Endpoint to register a sensor's configuration**
  - `iot.{device}.sensor.register` - Registers a new sensor configuration
  - Implemented in: `internal/device/device.go:handleSensorRegister()`

- ✅ **Endpoint to update a sensor's configuration**
  - `iot.{device}.config.update` - Updates sampling frequency and thresholds
  - Implemented in: `internal/device/device.go:handleConfigUpdate()`

- ✅ **Endpoint to query the current state of configurations**
  - `iot.{device}.config` - Configuration of all sensors
  - Implemented in: `internal/device/device.go:handleConfig()`

- ✅ **Endpoint to query the latest values read by a sensor**
  - `iot.{device}.readings.latest` - Gets the latest readings by sensor_id
  - Implemented in: `internal/device/device.go:handleLatestReadings()`

- ✅ **Publish changes in sensor readings to NATS**
  - `iot.{device}.readings.{sensor_type}` - Publishes readings in real-time
  - Implemented in: `internal/sensor/sensor.go:publish()`

#### Persistence (COMPLETE) ✅
- ✅ **Storage of sensor reading data in MongoDB**
  - **`readings` collection**: Auto-persistence in `sensor.go:publish()`
  - **`configurations` collection**: Persists changes in `device.go:handleConfigUpdate()`
  - Implemented in: `internal/storage/mongodb.go`

#### Reading Simulation (COMPLETE) ✅
- ✅ **Component that emulates the periodic reading of sensors**
  - Different for each sensor type (temperature, humidity, pressure)
  - Configurable parameters via the API
  - Implemented in: `internal/sensor/sensor.go:StartSensor()`

- ✅ **Simulation of sensor reading errors**
  - 5% probability of communication error
  - Implemented in: `internal/sensor/sensor.go:generateReading()`

#### Evaluation Criteria Met ✅
- ✅ **Code readability and style** - Code formatted with `go fmt`
- ✅ **Modularity and code organization** - Package-based architecture
- ✅ **Code coverage with tests** - Unit tests >70% in critical modules
- ✅ **Good use of version control** - Descriptive commits
- ✅ **Documentation and comments** - README.md and ARCHITECTURE.md are complete

#### Deliverables Met ✅
- ✅ **Source code in Git repository** - Complete repository
- ✅ **Associated documentation (README.md)** - Professional documentation
- ✅ **Explanatory diagram of the solution** - ARCHITECTURE.md

---

## Status Summary

### **🎉 PROJECT 100% COMPLETE** ✅

#### **Technical Requirements from the Official Document** ✅
- **NATS API**: 5/5 endpoints COMPLETE ✅
  - ✅ Endpoint to register sensor configuration
  - ✅ Endpoint to update sensor configuration
  - ✅ Endpoint to query current configuration status
  - ✅ Endpoint to query latest values read by sensor
  - ✅ Publish changes in sensor readings
- **MongoDB Persistence**: COMPLETE ✅
- **Reading Simulation**: COMPLETE ✅
- **Modular Organization**: COMPLETE ✅

#### **Evaluation Criteria** ✅
- **Code readability and style**: COMPLETE ✅
- **Modularity and code organization**: COMPLETE ✅
- **Code coverage with tests**: COMPLETE ✅
- **Good use of version control**: COMPLETE ✅
- **Documentation and comments**: COMPLETE ✅

#### **Deliverables** ✅
- **Source code in Git repository**: COMPLETE ✅
- **Associated documentation (README.md)**: COMPLETE ✅
- **Explanatory diagram of the solution**: COMPLETE ✅

---

## ✅ ALL OFFICIAL REQUIREMENTS MET

**The project now meets 100% of the requirements from the official evaluation document.**

### Implemented NATS Endpoints:
1. `iot.{device}.sensor.register` - Register sensor
2. `iot.{device}.config.update` - Update configuration
3. `iot.{device}.config` - Query configurations
4. `iot.{device}.readings.latest` - Latest values per sensor
5. `iot.{device}.readings.{sensor_type}` - Publish readings

### Additional Implemented Features:
- Complete persistence in MongoDB with queries
- Realistic sensor simulation with errors
- Modular and extensible architecture
- Unit tests with good coverage
- Professional and complete documentation
