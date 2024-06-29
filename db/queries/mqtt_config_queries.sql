-- name: GetUniqueMqttTopics :many
SELECT DISTINCT mqtt_topic FROM mqtt_config;