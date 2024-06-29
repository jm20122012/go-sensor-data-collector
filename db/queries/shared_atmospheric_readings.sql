-- name: InsertReading :exec
INSERT INTO shared_atmospheric_readings
    (
        id, 
        timestamp, 
        temp_f, 
        temp_c,
        humidity,
        absolute_pressure,
        device_type_id,
        device_id 
    )
VALUES
    (
        DEFAULT,
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7
    );