# Análisis de Requisitos - IoT Device Simulator

## Requisitos del Proyecto

### API NATS
- **Registrar y actualizar configuración de sensores**
  - Frecuencia de muestreo
  - Umbrales de alerta
- **Consultar estado actual de configuraciones**
- **Consultar últimos valores leídos**
- **Publicar cambios en lecturas de sensores**

### Persistencia
- **Almacenamiento de datos de lectura de sensores**
  - Formato a elegir por el candidato

### Simulación de Lecturas
- **Componente de lectura periódica**
  - Diferente para cada tipo de sensor
  - Parámetros configurables desde API
- **Simulación de errores de lectura**

### Criterios de Evaluación
- Legibilidad del código y estilo
- Modularidad y organización de código
- Cobertura de código con test
- Buen uso del control de versiones
- Documentación y comentarios

---

## Estado Actual de Implementación

### ✅ Implementado

#### API NATS - Parcial
- **Consultar configuraciones actuales**
  - `iot.{device}.config` - Retorna configuración de todos los sensores
  - Implementado en: `internal/device/device.go:handleConfig()`
  
- **Consultar estado del dispositivo**
  - `iot.{device}.status` - Retorna estado general del dispositivo
  - Implementado en: `internal/device/device.go:handleStatus()`
  
- **Publicar lecturas de sensores**
  - `iot.{device}.readings.{sensor_type}` - Publica lecturas en tiempo real
  - Implementado en: `internal/sensor/sensor.go:publish()`

#### Simulación de Lecturas - Completa
- **Lectura periódica configurable**
  - Frecuencia independiente por sensor
  - Rangos de valores configurables (min/max)
  - Implementado en: `internal/sensor/sensor.go:Start()`
  
- **Simulación de errores**
  - 5% probabilidad de error de comunicación
  - Implementado en: `internal/sensor/sensor.go:generateReading()`

#### Estructura del Código
- **Organización modular**
  - `cmd/iot-device/` - Punto de entrada
  - `internal/config/` - Gestión de configuración
  - `internal/device/` - Lógica del dispositivo
  - `internal/sensor/` - Lógica de sensores

### ❌ Faltante

#### API NATS - Actualización de Configuración
- **Actualizar configuración de sensores**
  - Modificar frecuencia de muestreo
  - Cambiar rangos de valores (min/max)
  - Habilitar/deshabilitar sensores
  - **Sujetos NATS necesarios:**
    - `iot.{device}.config.update.{sensor_id}`
    - `iot.{device}.config.frequency.{sensor_id}`
    - `iot.{device}.config.range.{sensor_id}`

#### API NATS - Consulta de Datos Históricos
- **Consultar últimos valores leídos**
  - **Sujetos NATS necesarios:**
    - `iot.{device}.readings.last.{sensor_id}`
    - `iot.{device}.readings.history.{sensor_id}`

#### Persistencia - Completa
- **Sistema de almacenamiento**
  - Base de datos para lecturas históricas
  - Almacenamiento de configuraciones
  - **Opciones a considerar:**
    - SQLite (local, simple)
    - PostgreSQL (robusto)
    - InfluxDB (time-series)
    - Archivos JSON/CSV (simple)

#### Testing - Completo
- **Tests unitarios**
  - Tests para módulo `device`
  - Tests para módulo `sensor`
  - Tests para módulo `config`
  - **Cobertura objetivo:** >80%

#### Documentación - Mejorar
- **Comentarios en código**
  - Documentación de funciones públicas
  - Explicación de algoritmos complejos
  - **Formato GoDoc**

---

## Funcionalidad Actual por Módulo

### `cmd/iot-device/main.go`
- Carga configuración desde archivo YAML
- Establece conexión con NATS
- Inicializa y arranca el dispositivo
- Maneja señales de interrupción

### `internal/config/config.go`
- Define estructuras de configuración
- Carga configuración desde archivo YAML
- Valida parámetros de configuración

### `internal/device/device.go`
- Gestiona múltiples sensores
- Configura suscripciones NATS
- Maneja peticiones de configuración y estado
- Coordina inicio/parada de sensores

### `internal/sensor/sensor.go`
- Simula lecturas periódicas de sensores
- Genera valores aleatorios en rangos configurados
- Simula errores de comunicación
- Publica lecturas a NATS

### `docker-compose.yml`
- Configura servidor NATS con JetStream
- Incluye cliente NATS para pruebas
- Habilita monitoring web en puerto 8222

---

## Prioridades de Desarrollo

### Alta Prioridad
1. **Implementar persistencia de datos**
2. **API NATS para actualizar configuración**
3. **Tests unitarios básicos**

### Media Prioridad
4. **API NATS para consultar históricos**
5. **Mejorar documentación del código**

### Baja Prioridad
6. **Diagrama de arquitectura**
7. **Métricas y monitoring avanzado**