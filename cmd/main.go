package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sensor-data-collection-service/internal/datastructs"
	"strconv"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"sensor-data-collection-service/internal/db"
	"sensor-data-collection-service/internal/db/sqlc"
	"sensor-data-collection-service/internal/mqttservice"
)

type DeviceList struct {
	Devices []*sqlc.GetDevicesRow
}

func (d *DeviceList) GetDeviceList(dbConn *db.SensorDataDB) {
	sensors, err := dbConn.GetDevices(context.Background())
	if err != nil {
		log.Println("Error getting sensors: ", err)
	}
	d.Devices = sensors
}

func getAvtechData(avtechUrl string) (*datastructs.AvtechResponseData, error) {
	resp, err := http.Get(avtechUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var responseData datastructs.AvtechResponseData

	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, nil
}

func getWeatherStationData(apiUrl string) (datastructs.WeatherStationResponseData, error) {
	resp, err := http.Get(apiUrl)
	if err != nil {
		log.Println("Error getting Ambient Weather Station API info: ", err)
		return nil, err
	}
	// log.Println("Ambient Weather Station API call successful")

	defer resp.Body.Close()

	var responseData datastructs.WeatherStationResponseData

	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

// func avtechWorker(wg *sync.WaitGroup, dbConn *db.SensorDataDB, sensorIdMap map[string]int64) {
// 	defer wg.Done()
// 	// -------------------- Assign ENV Vars -------------------- //
// 	avtechUrl := os.Getenv("AVTECH_URL")

// 	for {
// 		data, err := getAvtechData(avtechUrl)
// 		if err != nil {
// 			log.Println("Error getting Avtech data: ", err)
// 			log.Println("Skipping database write for this iteration")
// 			continue
// 		}
// 		log.Println("Avtech API call successful")

// 		tempF := convertStrToFloat32(data.Sensor[0].TempF)
// 		tempFHigh := convertStrToFloat32(data.Sensor[0].HighF)
// 		tempFLow := convertStrToFloat32(data.Sensor[0].LowF)
// 		tempC := convertStrToFloat32(data.Sensor[0].TempC)
// 		tempCHigh := convertStrToFloat32(data.Sensor[0].HighC)
// 		tempCLow := convertStrToFloat32(data.Sensor[0].LowC)
// 		timestamp := getUTCTimestamp()
// 		sensorID := sensorIdMap["avtech_basement_rack"]

// 		avtechDataParams := sqlc.InsertAvtechDataParams{
// 			Time:      timestamp,
// 			TempF:     &tempF,
// 			TempFHigh: &tempFHigh,
// 			TempFLow:  &tempFLow,
// 			TempC:     &tempC,
// 			TempCHigh: &tempCHigh,
// 			TempCLow:  &tempCLow,
// 			SensorID:  &sensorID,
// 		}

// 		err = dbConn.InsertAvtechData(context.Background(), avtechDataParams)
// 		if err != nil {
// 			log.Println("Error writing Avtech data to database: ", err)
// 		}

// 		log.Println("Avtech data written to database successfully")
// 		log.Println("Data: ", data.Sensor[0])

// 		time.Sleep(1 * time.Minute)
// 	}
// }

// func ambientWeatherStationWorker(wg *sync.WaitGroup, dbConn *db.SensorDataDB, sensorIdMap map[string]int64) {
// 	defer wg.Done()
// 	// -------------------- Assign ENV Vars -------------------- //
// 	ambientUrl := os.Getenv("AMBIENT_FULL_URL")

// 	for {
// 		data, err := getWeatherStationData(ambientUrl)
// 		if err != nil {
// 			log.Println("Error getting data from Ambient Weather Station: ", err)
// 			log.Println("Skipping database write for this iteration")
// 			continue
// 		}
// 		log.Println("Ambient weather station data retrieved successfully")

// 		timestamp := getUTCTimestamp()
// 		dateUTC := int32(data[0].LastData.DateUTC)
// 		insideTempF := float32(data[0].LastData.InsideTempF)
// 		insideFeelsLikeTempF := float32(data[0].LastData.FeelsLikeInside)
// 		outsideTempF := float32(data[0].LastData.OutsideTempF)
// 		outsideFeelsLikeTempF := float32(data[0].LastData.FeelsLikeOutside)
// 		insideHumidity := int32(data[0].LastData.InsideHumidity)
// 		outsideHumidity := int32(data[0].LastData.OutsideHumidity)
// 		insideDewPoint := float32(data[0].LastData.DewPointInside)
// 		outsideDewPoint := float32(data[0].LastData.DewPointOutside)
// 		baroRelative := float32(data[0].LastData.BarometricPressureRelIn)
// 		baroAbsolute := float32(data[0].LastData.BarometricPressureAbsIn)
// 		windDirection := int32(data[0].LastData.WindDirection)
// 		windSpeedMph := float32(data[0].LastData.WindSpeedMPH)
// 		windSpeedGustMph := float32(data[0].LastData.WindGustMPH)
// 		maxDailyGust := float32(data[0].LastData.MaxDailyGust)
// 		hourlyRainInches := float32(data[0].LastData.HourlyRainIn)
// 		eventRainInches := float32(data[0].LastData.EventRainIn)
// 		dailyRainInches := float32(data[0].LastData.DailyRainIn)
// 		weeklyRainInches := float32(data[0].LastData.WeeklyRainIn)
// 		monthlyRainInches := float32(data[0].LastData.MonthlyRainIn)
// 		totalRainInches := float32(data[0].LastData.TotalRainIn)
// 		uvIndex := float32(data[0].LastData.UVIndex)
// 		solarRadiation := float32(data[0].LastData.SolarRadiation)
// 		outsideBattStatus := int32(data[0].LastData.OutsideBattStatus)
// 		battCo2 := int32(data[0].LastData.BattCO2)
// 		sensorID := sensorIdMap["ambient_wx_station"]

// 		ambientStationDataParams := sqlc.InsertAmbientStationDataParams{
// 			Time:                  timestamp,
// 			Date:                  &data[0].LastData.Date,
// 			Timezone:              &data[0].LastData.TZ,
// 			DateUtc:               &dateUTC,
// 			InsideTempF:           &insideTempF,
// 			InsideFeelsLikeTempF:  &insideFeelsLikeTempF,
// 			OutsideTempF:          &outsideTempF,
// 			OutsideFeelsLikeTempF: &outsideFeelsLikeTempF,
// 			InsideHumidity:        &insideHumidity,
// 			OutsideHumidity:       &outsideHumidity,
// 			InsideDewPoint:        &insideDewPoint,
// 			OutsideDewPoint:       &outsideDewPoint,
// 			BaroRelative:          &baroRelative,
// 			BaroAbsolute:          &baroAbsolute,
// 			WindDirection:         &windDirection,
// 			WindSpeedMph:          &windSpeedMph,
// 			WindSpeedGustMph:      &windSpeedGustMph,
// 			MaxDailyGust:          &maxDailyGust,
// 			HourlyRainInches:      &hourlyRainInches,
// 			EventRainInches:       &eventRainInches,
// 			DailyRainInches:       &dailyRainInches,
// 			WeeklyRainInches:      &weeklyRainInches,
// 			MonthlyRainInches:     &monthlyRainInches,
// 			TotalRainInches:       &totalRainInches,
// 			LastRain:              &data[0].LastData.LastRain,
// 			UvIndex:               &uvIndex,
// 			SolarRadiation:        &solarRadiation,
// 			OutsideBattStatus:     &outsideBattStatus,
// 			BattCo2:               &battCo2,
// 			SensorID:              &sensorID,
// 		}

// 		err = dbConn.InsertAmbientStationData(context.Background(), ambientStationDataParams)
// 		if err != nil {
// 			log.Println("Error writing ambient weather station data to database: ", err)
// 		}

// 		log.Println("Ambient weather station data written to database successfully")
// 		log.Println("Data: ", data[0].LastData)
// 		time.Sleep(1 * time.Minute)
// 	}
// }

func mqttMsgHandlerFactory(dbConn *db.SensorDataDB, deviceIdMap map[string]int64) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received MQTT message \" %s \" from topic: %s\n", msg.Payload(), msg.Topic())

		// Convert MQTT json string to struct
		// err := json.Unmarshal(msg.Payload(), &sensorData)
		// if err != nil {
		// 	log.Println("Error unmarshalling MQTT message: ", err)
		// }

		// log.Println("Sensor data: ", sensorData)

	}
}

func onConnectionLostHdnlrFactory(wg *sync.WaitGroup, conn *pgxpool.Conn) mqtt.ConnectionLostHandler {
	return func(client mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v", err)
		conn.Release()
		wg.Done()
	}
}

func getUTCTimestamp() pgtype.Timestamptz {
	// Get the current time in UTC
	currentTime := time.Now().UTC()

	// Create a Timestamptz and set its value
	var timestamptz pgtype.Timestamptz
	timestamptz.Time = currentTime
	timestamptz.Valid = true

	return timestamptz
}

func convertStrToFloat32(value string) float32 {
	floatValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		log.Println("Error converting value to float32: ", err)
	}
	return float32(floatValue)
}

func acquirePoolConn(pool *pgxpool.Pool) *pgxpool.Conn {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire connection: %v\n", err)
	}
	return conn
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	connString := db.BuildPgConnectionString(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	pool := db.CreateConnectionPool(connString)
	defer pool.Close()

	conn := acquirePoolConn(pool)
	dbConn := db.NewSensorDataDB(conn)
	deviceList := DeviceList{}
	deviceList.GetDeviceList(dbConn)
	deviceIdMap := make(map[string]int64)
	for _, device := range deviceList.Devices {
		d := *device
		deviceID, err := dbConn.GetDeviceIdByName(ctx, d.DeviceName)
		if err != nil {
			log.Println("Error getting sensor ID: ", err)
		}
		deviceIdMap[d.DeviceName] = deviceID
	}
	conn.Release()

	// var wg *sync.WaitGroup
	wg := &sync.WaitGroup{}

	if os.Getenv("ENABLE_MQTT_LISTENER") == "true" {
		mqttService := mqttservice.NewMqttService()
		log.Println("New MQTT Service created", mqttService)

		conn := acquirePoolConn(pool)
		dbConn := db.NewSensorDataDB(conn)
		mqttService.DBConn = dbConn
		mqttService.CreateMqttClient(
			onConnectionLostHdnlrFactory(wg, conn),
			mqttMsgHandlerFactory(dbConn, deviceIdMap),
		)

		log.Println("Getting MQTT topics...")
		mqttTopics, err := dbConn.GetUniqueMqttTopics(ctx)
		if err != nil {
			log.Println("Error getting mqtt topics", err)
		}
		log.Println("MQTT topics: ", mqttTopics)
		for _, topic := range mqttTopics {
			log.Printf("Attempting to subscribe to %s\n", *topic)
			mqttService.MqttSubscribe(*topic)
		}
		wg.Add(1)
		// for _, device := range deviceList.Devices {
		// 	if *device.MqttTopic != "" {
		// 		mqttService.MqttSubscribe(*device.MqttTopic)
		// 	}
		// }
		// wg.Add(1)
	}
	// if os.Getenv("ENABLE_MQTT_LISTENER") == "true" {
	// 	log.Println("Adding 1 to WaitGroup for MQTT Listener...")
	// 	wg.Add(1)

	// 	log.Println("Getting postgres connection for MQTT Listener...")
	// 	conn := acquirePoolConn(pool)
	// 	defer conn.Release()

	// 	log.Println("Creating DB connection for MQTT Listener...")
	// 	dbConn := db.NewSensorDataDB(conn)

	// 	log.Println("Creating MQTT message handler for MQTT Listener...")
	// 	msgHndlr := mqttMsgHandlerFactory(dbConn, sensorIdMap)

	// 	log.Println("Creating connection lost handler for MQTT Listener...")
	// 	connectLostHdnlr := onConnectionLostHdnlrFactory(wg, conn)

	// 	log.Println("Creating MQTT client for MQTT Listener...")
	// 	mqttClient := mqttutils.CreateMqttClient(connectLostHdnlr, msgHndlr)

	// 	log.Println("Subscribing to MQTT topic...")
	// 	mqttutils.MqttSubscribe(mqttClient)

	// }

	// if os.Getenv("ENABLE_AVTECH_WORKER") == "true" {
	// 	conn := acquirePoolConn(pool)
	// 	defer conn.Release()

	// 	dbConn := db.NewSensorDataDB(conn)

	// 	log.Println("Starting Avtech Worker")

	// 	wg.Add(1)
	// 	go avtechWorker(wg, dbConn, sensorIdMap)
	// }

	// if os.Getenv("ENABLE_AMBIENT_WORKER") == "true" {
	// 	conn := acquirePoolConn(pool)
	// 	defer conn.Release()

	// 	dbConn := db.NewSensorDataDB(conn)

	// 	log.Println("Starting Ambient Weather Station Worker")

	// 	wg.Add(1)
	// 	go ambientWeatherStationWorker(wg, dbConn, sensorIdMap)
	// }

	wg.Wait()

}
