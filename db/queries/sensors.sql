-- name: GetSensors :many
SELECT 
    id, sensor_name, sensor_location, device_type_id
FROM 
    sensors;

-- name: GetSensorIdBySensorName :one
SELECT 
    id
FROM 
    sensors
WHERE 
    sensor_name = $1;

-- name: InsertDeviceIfNotExists :exec
INSERT INTO sensors (sensor_name, sensor_location, device_type_id)
SELECT $1, $2, dt.id
FROM device_type_ids dt
WHERE dt.device_type = $3
  AND NOT EXISTS (
    SELECT 1 FROM sensors WHERE sensor_name = $1
  );