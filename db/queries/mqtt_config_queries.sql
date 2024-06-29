-- name: GetUniqueMqttTopics :many
SELECT DISTINCT mqtt_topic FROM mqtt_config;

-- name: GetMqttTopicData :many
SELECT
	mc.mqtt_topic,
	mc.device_id,
    mc.device_type_id,
	dti.device_type
FROM 
	mqtt_config mc
INNER JOIN device_type_ids dti ON mc.device_type_id = dti.device_type_id; 