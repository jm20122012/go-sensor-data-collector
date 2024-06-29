-- name: GetUniqueMqttTopics :many
SELECT DISTINCT mqtt_topic FROM mqtt_config;

-- name: GetMqttTopicData :many
SELECT mqtt_topic, device_id, device_type_id FROM mqtt_config;