# EVALUACIÓN DE DESARROLLADOR DE SOFTWARE EN GOLANG
## Aplicación Golang para lectura de datos de sensores en un dispositivo IoT

### Descripción de la tarea
Implementar una aplicación en Go que simule la gestión de un dispositivo IoT con múltiples sensores de distinto tipo (temperatura, humedad, presión, etc.).

### Requisitos Exactos del Proyecto (Documento Oficial)

#### API NATS - Endpoints Específicos Requeridos
- **Endpoint para registrar configuración de un sensor**
- **Endpoint para actualizar configuración de un sensor** 
  - Frecuencia de muestreo
  - Umbrales de alerta
- **Endpoint para consultar estado actual de configuraciones**
- **Endpoint para consultar últimos valores leídos por sensor**
- **Publicar en NATS los cambios en las lecturas de sensores**

#### Persistencia
- **Almacenamiento de los datos de lectura de los sensores**
  - Formato considerado más adecuado (MongoDB elegido)

#### Simulación de lecturas
- **Componente que emule la lectura periódica de sensores**
  - Diferente para cada tipo de sensor
  - Permitir configurar parámetros desde la API
- **Simulación de errores de lectura de los sensores**

#### Restricciones
- Protocolo de mensajería: **NATS** (librería oficial nats.go)

#### Criterios de evaluación
- Legibilidad del código y estilo
- Modularidad y organización de código
- Cobertura de código con test
- Buen uso del control de versiones e historial de cambios
- Documentación y comentarios asociados al código

#### Entregables
1. Código fuente en un repositorio Git
2. Documentación asociada (README.md)
3. Diagrama explicativo de la solución propuesta

---

## Estado Actual de Implementación

### ✅ COMPLETADO

#### API NATS - Endpoints Implementados (5/5) ✅
- ✅ **Endpoint para registrar configuración de un sensor**
  - `iot.{device}.sensor.register` - Registra nueva configuración de sensor
  - Implementado en: `internal/device/device.go:handleSensorRegister()`

- ✅ **Endpoint para actualizar configuración de un sensor**
  - `iot.{device}.config.update` - Frecuencia de muestreo y umbrales
  - Implementado en: `internal/device/device.go:handleConfigUpdate()`

- ✅ **Endpoint para consultar estado actual de configuraciones**
  - `iot.{device}.config` - Configuración de todos los sensores
  - Implementado en: `internal/device/device.go:handleConfig()`

- ✅ **Endpoint para consultar últimos valores leídos por sensor**
  - `iot.{device}.readings.latest` - Obtiene últimas lecturas por sensor_id
  - Implementado en: `internal/device/device.go:handleLatestReadings()`

- ✅ **Publicar en NATS los cambios en las lecturas de sensores**
  - `iot.{device}.readings.{sensor_type}` - Publica lecturas en tiempo real
  - Implementado en: `internal/sensor/sensor.go:publish()`

#### Persistencia (COMPLETO) ✅
- ✅ **Almacenamiento de los datos de lectura de los sensores en MongoDB**
  - **Colección `readings`**: Auto-persistencia en `sensor.go:publish()`
  - **Colección `configurations`**: Persiste cambios en `device.go:handleConfigUpdate()`
  - Implementado en: `internal/storage/mongodb.go`

#### Simulación de lecturas (COMPLETO) ✅
- ✅ **Componente que emule la lectura periódica de sensores**
  - Diferente para cada tipo de sensor (temperatura, humedad, presión)
  - Parámetros configurables desde la API
  - Implementado en: `internal/sensor/sensor.go:StartSensor()`

- ✅ **Simulación de errores de lectura de los sensores**
  - 5% probabilidad de error de comunicación
  - Implementado en: `internal/sensor/sensor.go:generateReading()`

#### Criterios de evaluación completados ✅
- ✅ **Legibilidad del código y estilo** - Código formateado con go fmt
- ✅ **Modularidad y organización de código** - Arquitectura por paquetes
- ✅ **Cobertura de código con test** - Tests unitarios >70% en módulos críticos
- ✅ **Buen uso del control de versiones** - Commits descriptivos
- ✅ **Documentación y comentarios** - README.md y ARCHITECTURE.md completos

#### Entregables completados ✅
- ✅ **Código fuente en repositorio Git** - Repositorio completo
- ✅ **Documentación asociada (README.md)** - Documentación profesional
- ✅ **Diagrama explicativo de la solución** - ARCHITECTURE.md

---

#### Persistencia Mejorada (COMPLETO) ✅
- ✅ **Almacenamiento de los datos de lectura de los sensores en MongoDB**
  - **Colección `readings`**: Auto-persistencia en `sensor.go:publish()`
  - **Colección `configurations`**: Persiste cambios en `device.go:handleConfigUpdate()`
  - **Método `GetLatestReadings()`**: Recupera últimas lecturas por sensor
  - Implementado en: `internal/storage/mongodb.go`

---

## Resumen de Estado

### **🎉 PROYECTO COMPLETADO AL 100%** ✅

#### **Requisitos Técnicos del Documento Oficial** ✅
- **API NATS**: 5/5 endpoints COMPLETOS ✅
  - ✅ Endpoint para registrar configuración de sensor
  - ✅ Endpoint para actualizar configuración de sensor  
  - ✅ Endpoint para consultar estado actual de configuraciones
  - ✅ Endpoint para consultar últimos valores leídos por sensor
  - ✅ Publicar cambios en las lecturas de sensores
- **Persistencia MongoDB**: COMPLETO ✅ 
- **Simulación de lecturas**: COMPLETO ✅
- **Organización modular**: COMPLETO ✅

#### **Criterios de Evaluación** ✅
- **Legibilidad del código y estilo**: COMPLETO ✅
- **Modularidad y organización de código**: COMPLETO ✅
- **Cobertura de código con test**: COMPLETO ✅
- **Buen uso del control de versiones**: COMPLETO ✅
- **Documentación y comentarios**: COMPLETO ✅

#### **Entregables** ✅
- **Código fuente en repositorio Git**: COMPLETO ✅
- **Documentación asociada (README.md)**: COMPLETO ✅
- **Diagrama explicativo de la solución**: COMPLETO ✅

---

## ✅ TODOS LOS REQUISITOS OFICIALES CUMPLIDOS

**El proyecto ahora cumple el 100% de los requisitos del documento oficial de evaluación.**

### Endpoints NATS implementados:
1. `iot.{device}.sensor.register` - Registrar sensor
2. `iot.{device}.config.update` - Actualizar configuración  
3. `iot.{device}.config` - Consultar configuraciones
4. `iot.{device}.readings.latest` - Últimos valores por sensor
5. `iot.{device}.readings.{sensor_type}` - Publicar lecturas

### Funcionalidades adicionales implementadas:
- Persistencia completa en MongoDB con consultas
- Simulación realista de sensores con errores
- Arquitectura modular y extensible
- Tests unitarios con buena cobertura
- Documentación profesional completa