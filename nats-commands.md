# Comandos NATS para IoT Device Simulator

## Configuración Inicial

### Conexión al cliente NATS
```bash
# Entrar al contenedor con CLI NATS (recomendado)
docker exec -it nats-client sh

# La variable NATS_URL ya está configurada automáticamente
# No necesitas export, solo usar los comandos directamente
```

---

## 1. Endpoints NATS Principales (5/5 Implementados) ✅

### 1.1 Consultar configuración de todos los sensores
```bash
nats req iot.device-001.config ""
```

**Respuesta esperada:**
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
  // ... más sensores
}
```

### 1.2 Consultar estado del dispositivo
```bash
nats req iot.device-001.status ""
```

**Respuesta esperada:**
```json
{
  "device_id": "device-001",
  "total_sensors": 5,
  "enabled_sensors": 4,
  "disabled_sensors": 1,
  "timestamp": "2025-08-04T10:30:00Z"
}
```

### 1.3 Registrar nuevo sensor ⭐ NUEVO
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

**Respuesta esperada:**
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

### 1.4 Actualizar configuración de sensor existente
```bash
nats req iot.device-001.config.update '{
  "sensor_id": "temp-01",
  "frequency": "2s",
  "min": 10,
  "max": 40
}'
```

**Respuesta esperada:**
```json
{
  "status": "updated"
}
```

### 1.5 Consultar últimas lecturas por sensor ⭐ NUEVO
```bash
nats req iot.device-001.readings.latest '{
  "sensor_id": "temp-01"
}'
```

**Respuesta esperada (si hay datos):**
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

**Respuesta si no hay datos:**
```json
{
  "error": "no readings found"
}
```

---

## 2. Monitoreo en Tiempo Real

### 2.1 Suscribirse a todas las lecturas del dispositivo
```bash
nats sub "iot.device-001.readings.>"
```

### 2.2 Suscribirse solo a lecturas de temperatura
```bash
nats sub "iot.device-001.readings.temperature"
```

### 2.3 Suscribirse solo a lecturas de humedad
```bash
nats sub "iot.device-001.readings.humidity"
```

### 2.4 Suscribirse solo a lecturas de presión
```bash
nats sub "iot.device-001.readings.pressure"
```

**Ejemplo de lectura recibida:**
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

## 3. Casos de Uso Completos

### 3.1 Flujo completo: Registrar → Monitorear → Consultar
```bash
# 1. Registrar nuevo sensor
nats req iot.device-001.sensor.register '{
  "sensor_id": "light-01",
  "type": "light",
  "frequency": "5s",
  "min": 0,
  "max": 1000,
  "unit": "lux"
}'

# 2. Monitorear sus lecturas en tiempo real
nats sub "iot.device-001.readings.light"

# 3. Después de un momento, consultar su última lectura
nats req iot.device-001.readings.latest '{
  "sensor_id": "light-01"
}'
```

### 3.2 Gestión de configuración
```bash
# 1. Ver configuración actual
nats req iot.device-001.config ""

# 2. Actualizar frecuencia de un sensor
nats req iot.device-001.config.update '{
  "sensor_id": "humidity-01",
  "frequency": "5s"
}'

# 3. Verificar que se aplicó el cambio
nats req iot.device-001.config ""
```

### 3.3 Análisis de rendimiento del dispositivo
```bash
# 1. Ver estado general
nats req iot.device-001.status ""

# 2. Monitorear todas las lecturas por 30 segundos
timeout 30s nats sub "iot.device-001.readings.>"

# 3. Consultar últimas lecturas de sensores críticos
nats req iot.device-001.readings.latest '{"sensor_id": "temp-01"}'
nats req iot.device-001.readings.latest '{"sensor_id": "temp-02"}'
nats req iot.device-001.readings.latest '{"sensor_id": "humidity-01"}'
```

---

## 4. Manejo de Errores

### 4.1 Errores comunes y soluciones

**Error: "No responders are available"**
```bash
# Causa: La aplicación Go no está ejecutándose
# Solución: Verificar que ./iot-device esté corriendo
ps aux | grep iot-device | grep -v grep
```

**Error: "sensor not found"**
```bash
# Causa: El sensor_id no existe
# Solución: Ver lista de sensores disponibles
nats req iot.device-001.config ""
```

**Error: "no readings found"**
```bash
# Causa: El sensor es nuevo y no ha generado lecturas aún
# Solución: Esperar unos segundos o verificar que está habilitado
nats req iot.device-001.config ""
```

### 4.2 Simulación de errores
```bash
# Los sensores simulan 5% de errores automáticamente
# Para ver errores, monitoreea las lecturas:
nats sub "iot.device-001.readings.>" | grep "error"
```

---

## 5. Persistencia en MongoDB

### 5.1 Verificar datos en MongoDB
```bash
# Conectar a MongoDB
docker exec -it mongodb mongosh iot_simulator

# Ver lecturas más recientes
db.readings.find().sort({timestamp: -1}).limit(5)

# Contar lecturas por sensor
db.readings.aggregate([
  {$group: {_id: "$sensor_id", count: {$sum: 1}}},
  {$sort: {count: -1}}
])

# Ver configuraciones guardadas
db.configurations.find()
```

### 5.2 Limpiar datos de prueba
```bash
# Conectar a MongoDB
docker exec -it mongodb mongosh iot_simulator

# Eliminar lecturas de sensores de prueba
db.readings.deleteMany({"sensor_id": /test/})

# Limpiar todas las lecturas (¡CUIDADO!)
db.readings.deleteMany({})
```

---

## 6. Monitoreo Web

### 6.1 NATS Server Monitoring
```bash
# Abrir monitoring de NATS en navegador
open http://localhost:8222

# O verificar via curl
curl -s http://localhost:8222/varz | jq
```

### 6.2 Estadísticas via API
```bash
# Información general del servidor
curl -s http://localhost:8222/varz | jq '{
  server_id: .server_id,
  version: .version,
  connections: .connections,
  subscriptions: .subscriptions
}'

# Conexiones activas
curl -s http://localhost:8222/connz | jq '.connections | length'

# Subscripciones activas
curl -s http://localhost:8222/subsz | jq '.subscriptions | length'
```

---

## 7. Resumen de Subjects NATS

### Request/Response (síncronos)
- `iot.device-001.config` - Obtener configuración
- `iot.device-001.status` - Obtener estado
- `iot.device-001.sensor.register` - Registrar sensor
- `iot.device-001.config.update` - Actualizar configuración
- `iot.device-001.readings.latest` - Últimas lecturas

### Publish/Subscribe (asíncronos)
- `iot.device-001.readings.temperature` - Lecturas de temperatura
- `iot.device-001.readings.humidity` - Lecturas de humedad
- `iot.device-001.readings.pressure` - Lecturas de presión
- `iot.device-001.readings.>` - Todas las lecturas (wildcard)

---

## 8. Notas Importantes

### Device ID
- El dispositivo usa ID: `device-001` (según config.yml)
- Todos los subjects incluyen este ID
- Para múltiples dispositivos, cambiar el ID en config.yml

### Tipos de Sensores Soportados
- `temperature` - Temperatura en °C
- `humidity` - Humedad en %
- `pressure` - Presión en hPa
- `light` - Luz en lux (personalizable)
- Cualquier tipo personalizado

### Persistencia
- Las lecturas se guardan automáticamente en MongoDB
- Las configuraciones se persisten al actualizarse
- Los datos sobreviven reinicios del sistema

### Rendimiento
- El sistema maneja múltiples sensores simultáneamente
- Frecuencias configurables por sensor (mínimo 1s)
- Simulación automática de errores (~5%)
- Respuesta en tiempo real (<5ms típico)