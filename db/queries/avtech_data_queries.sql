-- name: InsertAvtechData :exec
INSERT INTO avtech_data 
(
    id, 
    "time", 
    temp_f, 
    temp_f_low, 
    temp_f_high, 
    temp_c, 
    temp_c_low, 
    temp_c_high, 
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
    $8
);