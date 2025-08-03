# iot-device-simulator

docker exec -it nats-client sh

nats sub "iot.device-001.readings.*"

nats sub "iot.device-001.readings.temperature"

nats sub "iot.device-001.readings.pressure"

nats sub "iot.device-001.readings.humidity"
