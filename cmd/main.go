package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"sensor-data-collection-service/internal/ambientstationservice"
	"sensor-data-collection-service/internal/avtechservice"
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
	deviceListMap := make(map[string]sqlc.GetDevicesRow)
	for _, device := range deviceList.Devices {
		d := *device
		deviceListMap[d.DeviceName] = sqlc.GetDevicesRow{
			DeviceName:     d.DeviceName,
			DeviceLocation: d.DeviceLocation,
			DeviceTypeID:   d.DeviceTypeID,
			DeviceID:       d.DeviceID,
			DeviceType:     d.DeviceType,
			MqttTopic:      d.MqttTopic,
		}
	}
	conn.Release()

	wg := &sync.WaitGroup{}

	if cfg.EnableMqttWorker {
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

	if cfg.EnableAmbientStationWorker {
		wg.Add(1)
		apiUrl := fmt.Sprintf("%s?applicationKey=%s&apiKey=%s",
			os.Getenv("AMBIENT_API_URL_BASE"),
			os.Getenv("AMBIENT_APPLICATION_KEY"),
			os.Getenv("AMBIENT_API_KEY"),
		)

		ambientService := ambientstationservice.NewAmbientStationService(
			ctx,
			wg,
			logger,
			pool,
			apiUrl,
			deviceListMap,
		)

		ambientService.Run()

	}

	if cfg.EnableAvtechWorker {
		wg.Add(1)
		apiUrl := os.Getenv("AVTECH_API_URL")

		avtechService := avtechservice.NewAvtechService(
			ctx,
			wg,
			logger,
			pool,
			apiUrl,
			deviceListMap,
		)

		avtechService.Run()

	}

	wg.Wait()

	logger.Info("Exiting...")

}
