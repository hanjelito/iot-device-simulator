# Diagrama de Arquitectura - IoT Device Simulator

## Flujo de Datos

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│                 │    │                  │    │                 │
│   NATS Client   │◄──►│   IoT Device     │◄──►│    MongoDB      │
│                 │    │   Simulator      │    │   (Optional)    │
│                 │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
    ┌────▼────┐             ┌────▼────┐            ┌────▼────┐
    │         │             │         │            │         │
    │ Config  │             │ Sensor  │            │Reading  │
    │ Updates │             │Manager  │            │Storage  │
    │         │             │         │            │         │
    └─────────┘             └─────────┘            └─────────┘
```

## Componentes Principales

### 1. Device Manager (`internal/device/`)
```
┌─────────────────────────────────────────┐
│            Device Manager               │
├─────────────────────────────────────────┤
│ • Gestiona múltiples sensores           │
│ • Maneja suscripciones NATS            │
│ • Coordina actualizaciones config      │
│ • Persiste cambios en MongoDB          │
└─────────────────────────────────────────┘
                    │
                    ▼
    ┌─────────────────────────────────────┐
    │          NATS Subscriptions         │
    ├─────────────────────────────────────┤
    │ iot.{device}.config                │
    │ iot.{device}.status                │
    │ iot.{device}.config.update         │
    └─────────────────────────────────────┘
```

### 2. Sensor Simulator (`internal/sensor/`)
```
┌──────────────────────────────────────────────┐
│               Sensor                         │
├──────────────────────────────────────────────┤
│ • Genera lecturas periódicas                 │
│ • Simula errores (5% probabilidad)          │
│ • Thread-safe config updates                │
│ • Publica vía NATS + guarda en MongoDB      │
└──────────────────────────────────────────────┘
                     │
                     ▼
    ┌─────────────────────────────────────────┐
    │          Reading Generation             │
    ├─────────────────────────────────────────┤
    │ • Valores aleatorios en rango          │
    │ • Timestamp automático                 │
    │ • Manejo de errores de comunicación   │
    └─────────────────────────────────────────┘
```

### 3. Data Flow
```
Config.yml ──► Device Manager ──► Sensor[1..N]
                     │                 │
                     │                 ▼
              ┌──────────────┐    Publishing:
              │              │    • NATS: iot.{device}.readings.{type}
              │   MongoDB    │    • MongoDB: readings collection
              │              │
              │ Collections: │         │
              │ • readings   │◄────────┘
              │ • configs    │◄────────────────┐
              └──────────────┘                 │
                     ▲                         │
                     │                         │
              Config Updates ──────────────────┘
```

## Concurrencia y Thread Safety

```
┌─────────────────────────────────────────────────────┐
│                  Main Process                       │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌───────────────┐  ┌───────────────┐              │
│  │   Sensor 1    │  │   Sensor 2    │   ...        │
│  │  (goroutine)  │  │  (goroutine)  │              │
│  │               │  │               │              │
│  │ • ticker      │  │ • ticker      │              │
│  │ • mutex lock  │  │ • mutex lock  │              │
│  │ • publish     │  │ • publish     │              │
│  └───────────────┘  └───────────────┘              │
│                                                     │
│  ┌─────────────────────────────────────────────┐   │
│  │          NATS Subscriptions                │   │
│  │          (separate goroutines)             │   │
│  │                                            │   │
│  │ • iot.{device}.config                     │   │
│  │ • iot.{device}.status                     │   │
│  │ • iot.{device}.config.update              │   │
│  └─────────────────────────────────────────────┘   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

## Persistencia Strategy

```
┌────────────────────────────────────────────────────┐
│                MongoDB Strategy                    │
├────────────────────────────────────────────────────┤
│                                                    │
│  Readings Collection:                              │
│  ┌─────────────────────────────────────────────┐   │
│  │ {                                           │   │
│  │   sensor_id: "temp-01",                     │   │
│  │   type: "temperature",                      │   │
│  │   value: 23.5,                              │   │
│  │   unit: "°C",                               │   │
│  │   timestamp: ISODate(...),                  │   │
│  │   error?: "comm error"                      │   │
│  │ }                                           │   │
│  └─────────────────────────────────────────────┘   │
│                                                    │
│  Configurations Collection:                        │
│  ┌─────────────────────────────────────────────┐   │
│  │ {                                           │   │
│  │   device_id: "device-001",                   │   │
│  │   configs: { ... sensor configs ... },     │   │
│  │   timestamp: ISODate(...)                   │   │
│  │ }                                           │   │
│  └─────────────────────────────────────────────┘   │
│                                                    │
└────────────────────────────────────────────────────┘
```

## Error Handling & Resilience

```
┌─────────────────────────────────────────────────────┐
│                 Error Strategy                      │
├─────────────────────────────────────────────────────┤
│                                                     │
│  MongoDB Connection:                                │
│  ┌─────────────────────────────────────────────┐   │
│  │ • Optional dependency                       │   │
│  │ • Graceful degradation                      │   │
│  │ • App continues without persistence         │   │
│  └─────────────────────────────────────────────┘   │
│                                                     │
│  Sensor Communication:                              │
│  ┌─────────────────────────────────────────────┐   │
│  │ • 5% error simulation                       │   │
│  │ • Error field in reading                    │   │
│  │ • Continue operation                        │   │
│  └─────────────────────────────────────────────┘   │
│                                                     │
│  NATS Connection:                                   │
│  ┌─────────────────────────────────────────────┐   │
│  │ • Critical dependency                       │   │
│  │ • App exits if unavailable                  │   │
│  │ • Required for core functionality           │   │
│  └─────────────────────────────────────────────┘   │
│                                                     │
└─────────────────────────────────────────────────────┘
```