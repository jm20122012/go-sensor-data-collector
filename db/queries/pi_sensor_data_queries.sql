-- name: InsertPiSensorData :exec
INSERT INTO pi_sensor_data
(
    id,
    "time",
    temp_f,
    temp_c,
    humidity,
    sensor_location,
    sensor_id
)
VALUES
(
    DEFAULT,
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);