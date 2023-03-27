-- name: InsertAmbientStationData :exec
INSERT INTO ambient_station_data
    (
        id, 
        "time", 
        "date", 
        timezone, 
        date_utc, 
        inside_temp_f, 
        inside_feels_like_temp_f, 
        outside_temp_f, 
        outside_feels_like_temp_f, 
        inside_humidity, 
        outside_humidity, 
        inside_dew_point, 
        outside_dew_point, 
        baro_relative, 
        baro_absolute, 
        wind_direction, 
        wind_speed_mph, 
        wind_speed_gust_mph, 
        max_daily_gust, 
        hourly_rain_inches, 
        event_rain_inches, 
        daily_rain_inches, 
        weekly_rain_inches, 
        monthly_rain_inches, 
        total_rain_inches, 
        last_rain, 
        uv_index, 
        solar_radiation, 
        outside_batt_status, 
        batt_co2, 
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
    $6, 
    $7, 
    $8, 
    $9, 
    $10, 
    $11, 
    $12, 
    $13, 
    $14, 
    $15, 
    $16, 
    $17, 
    $18, 
    $19, 
    $20, 
    $21, 
    $22, 
    $23, 
    $24, 
    $25, 
    $26, 
    $27, 
    $28, 
    $29,
    $30
);