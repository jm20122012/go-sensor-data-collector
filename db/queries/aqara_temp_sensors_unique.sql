-- name: InsertAqaraUniqueData :exec 
INSERT INTO aqara_temp_sensors_unique
    (
        id,
        timestamp,
        link_quality,
        batt_percentage,
        power_outage_count,
        device_id,
        device_type_id
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