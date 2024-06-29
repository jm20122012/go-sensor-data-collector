package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"sync"

	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"sensor-data-collection-service/internal/config"
	"sensor-data-collection-service/internal/mqttservice"
	"sensor-data-collection-service/internal/sensordb"
	"sensor-data-collection-service/internal/sensordb/sqlc"
	"sensor-data-collection-service/internal/utils"
)

type DeviceList struct {
	Devices []*sqlc.GetDevicesRow
}

func (d *DeviceList) GetDeviceList(dbConn *sensordb.SensorDataDB) {
	sensors, err := dbConn.GetDevices(context.Background())
	if err != nil {
		log.Println("Error getting sensors: ", err)
	}
	d.Devices = sensors
}

// func getAvtechData(avtechUrl string) (*datastructs.AvtechResponseData, error) {
// 	resp, err := http.Get(avtechUrl)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resp.Body.Close()

// 	var responseData datastructs.AvtechResponseData

// 	err = json.NewDecoder(resp.Body).Decode(&responseData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &responseData, nil
// }

// func getWeatherStationData(apiUrl string) (datastructs.WeatherStationResponseData, error) {
// 	resp, err := http.Get(apiUrl)
// 	if err != nil {
// 		log.Println("Error getting Ambient Weather Station API info: ", err)
// 		return nil, err
// 	}
// 	// log.Println("Ambient Weather Station API call successful")

// 	defer resp.Body.Close()

// 	var responseData datastructs.WeatherStationResponseData

// 	err = json.NewDecoder(resp.Body).Decode(&responseData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return responseData, nil
// }

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

// func acquirePoolConn(pool *pgxpool.Pool) *pgxpool.Conn {
// 	conn, err := pool.Acquire(context.Background())
// 	if err != nil {
// 		log.Fatalf("Unable to acquire connection: %v\n", err)
// 	}
// 	return conn
// }

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load app config
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		slog.Error("Erorr loading config.yml", "error", err)
	}

	// Setup logger
	var logger *slog.Logger
	switch cfg.DebugLevel {
	case "DEBUG":
		logger = utils.CreateLogger(slog.LevelDebug)
	case "INFO":
		logger = utils.CreateLogger(slog.LevelInfo)
	case "WARNING":
		logger = utils.CreateLogger(slog.LevelWarn)
	case "ERROR":
		logger = utils.CreateLogger(slog.LevelError)
	default:
		logger = utils.CreateLogger(slog.LevelInfo)
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func(cancel context.CancelFunc) {
		<-c
		logger.Info("Ctrl+C pressed, cancelling context...")
		cancel()
	}(cancel)

	connString := sensordb.BuildPgConnectionString(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	pool := sensordb.CreateConnectionPool(connString)
	defer pool.Close()

	conn := utils.AcquirePoolConn(pool)
	db := sensordb.NewSensorDataDB(conn)
	deviceList := DeviceList{}
	deviceList.GetDeviceList(db)
	deviceIdMap := make(map[string]int64)
	for _, device := range deviceList.Devices {
		d := *device
		deviceID, err := db.GetDeviceIdByName(ctx, d.DeviceName)
		if err != nil {
			logger.Error("Error getting sensor ID", "error", err)
		}
		deviceIdMap[d.DeviceName] = deviceID
	}
	conn.Release()

	// var wg *sync.WaitGroup
	wg := &sync.WaitGroup{}

	if os.Getenv("ENABLE_MQTT_LISTENER") == "true" {
		conn := utils.AcquirePoolConn(pool)
		db := sensordb.NewSensorDataDB(conn)
		ms := mqttservice.NewMqttService(wg, ctx, *cfg, logger, pool)
		logger.Info("New MQTT Service created", "service", ms)

		ms.CreateMqttClient()
		logger.Info("MQTT client created")

		topicData, err := db.GetMqttTopicData(ctx)
		if err != nil {
			logger.Error("Could not retrieve topic data", "error", err)
		}

		topicDataMap := make(map[string]mqttservice.TopicData)
		for _, t := range topicData {
			td := mqttservice.TopicData{
				DeviceID:     *t.DeviceID,
				DeviceTypeID: *t.DeviceTypeID,
				DeviceType:   t.DeviceType,
			}
			topicDataMap[*t.MqttTopic] = td
		}

		ms.TopicMapping = topicDataMap

		mqttTopics, err := db.GetUniqueMqttTopics(ctx)
		if err != nil {
			logger.Error("Error getting mqtt topics", "error", err)
		}
		for _, topic := range mqttTopics {
			ms.MqttSubscribe(*topic)
		}

		conn.Release()

		wg.Add(1)
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

	logger.Info("Exiting...")

}
