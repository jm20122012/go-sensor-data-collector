-- name: InsertSonoffSmartPlugReading :exec
INSERT INTO
    sonoff_smart_plugs (
        id,
        timestamp,
        link_quality,
        outlet_state,
        device_id,
        device_type_id
    )
VALUES (
    DEFAULT,
    $1,
    $2,
    $3,
    $4,
    $5
);