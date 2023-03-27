-- Create "device_type_ids" table
CREATE TABLE "device_type_ids" (
  "id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  "device_type" character varying NOT NULL,
  "device_type_id" integer NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "device_type_ids_unique" UNIQUE ("device_type"),
  CONSTRAINT "device_type_ids_unique_1" UNIQUE ("device_type_id")
);
-- Create "sensors" table
CREATE TABLE "sensors" (
  "id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  "sensor_name" character varying NOT NULL,
  "sensor_location" character varying NOT NULL,
  "device_type_id" integer NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "sensors_device_type_ids_fk" FOREIGN KEY ("device_type_id") REFERENCES "device_type_ids" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "ambient_station_data" table
CREATE TABLE "ambient_station_data" (
  "id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  "time" timestamptz NULL,
  "date" text NULL,
  "timezone" text NULL,
  "date_utc" integer NULL,
  "inside_temp_f" real NULL,
  "inside_feels_like_temp_f" real NULL,
  "outside_temp_f" real NULL,
  "outside_feels_like_temp_f" real NULL,
  "inside_humidity" integer NULL,
  "outside_humidity" integer NULL,
  "inside_dew_point" real NULL,
  "outside_dew_point" real NULL,
  "baro_relative" real NULL,
  "baro_absolute" real NULL,
  "wind_direction" integer NULL,
  "wind_speed_mph" real NULL,
  "wind_speed_gust_mph" real NULL,
  "max_daily_gust" real NULL,
  "hourly_rain_inches" real NULL,
  "event_rain_inches" real NULL,
  "daily_rain_inches" real NULL,
  "weekly_rain_inches" real NULL,
  "monthly_rain_inches" real NULL,
  "total_rain_inches" real NULL,
  "last_rain" text NULL,
  "uv_index" real NULL,
  "solar_radiation" real NULL,
  "outside_batt_status" integer NULL,
  "batt_co2" integer NULL,
  "sensor_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "ambient_station_data_sensor_id_fkey" FOREIGN KEY ("sensor_id") REFERENCES "sensors" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "avtech_data" table
CREATE TABLE "avtech_data" (
  "id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  "time" timestamptz NULL,
  "temp_f" real NULL,
  "temp_f_low" real NULL,
  "temp_f_high" real NULL,
  "temp_c" real NULL,
  "temp_c_low" real NULL,
  "temp_c_high" real NULL,
  "sensor_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "avtech_data_sensors_fk" FOREIGN KEY ("sensor_id") REFERENCES "sensors" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "pi_sensor_data" table
CREATE TABLE "pi_sensor_data" (
  "id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
  "time" timestamptz NULL,
  "temp_f" real NULL,
  "temp_c" real NULL,
  "humidity" real NULL,
  "sensor_location" text NULL,
  "sensor_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "pi_sensor_data_sensors_fk" FOREIGN KEY ("sensor_id") REFERENCES "sensors" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
