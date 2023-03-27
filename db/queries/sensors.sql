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