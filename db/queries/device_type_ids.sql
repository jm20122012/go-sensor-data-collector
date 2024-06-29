-- name: GetDeviceTypeIDByDeviceType :one
SELECT 
    device_type_id
FROM
    device_type_ids
WHERE
    device_type = $1;