-- name: GetDevices :many
SELECT 
	dl.device_name, 
	dl.device_location, 
	dl.device_type_id,
	dl.device_id,
	dt.device_type,
    mc.mqtt_topic
from
	device_list dl
inner join device_type_ids dt on dl.device_type_id = dt.device_type_id
left join mqtt_config mc on dl.device_id = mc.device_id;

-- name: GetDeviceIdByName :one
SELECT 
    id
FROM 
    device_list
WHERE 
    device_name = $1;