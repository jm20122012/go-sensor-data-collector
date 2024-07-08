package mqttservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"sensor-data-collection-service/internal/config"
	"sensor-data-collection-service/internal/devices"
	"sensor-data-collection-service/internal/sensordb"
	"sensor-data-collection-service/internal/sensordb/sqlc"
	"sensor-data-collection-service/internal/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TopicData struct {
	DeviceID     int32
	DeviceTypeID int32
	DeviceType   string
}

type MqttService struct {
	Ctx           context.Context
	Wg            *sync.WaitGroup
	Config        config.Config
	Client        mqtt.Client
	TopicMapping  map[string]TopicData
	DeviceListMap map[string]devices.Devices
	Logger        *slog.Logger
	PgPool        *pgxpool.Pool
}

func NewMqttService(wg *sync.WaitGroup, ctx context.Context, cfg config.Config, logger *slog.Logger, pool *pgxpool.Pool, deviceListMap map[string]devices.Devices) MqttService {
	return MqttService{
		Ctx:           ctx,
		Wg:            wg,
		Config:        cfg,
		Logger:        logger,
		PgPool:        pool,
		DeviceListMap: deviceListMap,
	}
}

func PublishConnectSuccessful(client mqtt.Client) {
	connectionMsg := fmt.Sprintf("[%s] Connected to MQTT Broker at %s:%s", time.Now().String(), os.Getenv("MQTT_BROKER_IP"), os.Getenv("MQTT_BROKER_PORT"))
	token := client.Publish("connections", 0, false, connectionMsg)
	token.Wait()
}

func (m *MqttService) MqttSubscribe(topic string) {
	token := m.Client.Subscribe(topic, 1, nil)
	token.Wait()
	m.Logger.Info("Subscribing to MQTT topic", "topic", topic)
}

func (m *MqttService) OnConnectHndlrFactory() mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		m.Logger.Info("MQTT connected")
		PublishConnectSuccessful(client)
	}
}

func (m *MqttService) MqttMsgHndlrFactory() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		m.Logger.Info("Received MQTT message", "msg", msg.Payload(), "topic", msg.Topic())

		switch m.TopicMapping[msg.Topic()].DeviceType {
		case "aqara_temp_sensor":
			m.ProcessAqaraTempSensorMsg(msg)
		case "dht11_sensor":
			m.ProcessDht11TempSensorMsg(msg)
		case "bme280_sensor":
			m.ProcessBme280SensorMsg(msg)
		case "bmp280_sensor":
			m.ProcessBmp280SensorMsg(msg)
		case "sonoff_smart_plug":
			m.ProcessSonoffSmartPlugMsg(msg)
		default:
			m.Logger.Warn("Unknown device type in MQTT message")
		}
	}
}

func (m *MqttService) ProcessAqaraTempSensorMsg(msg mqtt.Message) {
	m.Logger.Debug("Processing aqara temp sensor message")

	// Declare the raw message struct
	var rawMsg devices.AqaraTempSensorRawMessage
	err := json.Unmarshal(msg.Payload(), &rawMsg)
	if err != nil {
		m.Logger.Error("Error unmarshalling raw aqara temp sensor message", "error", err)
	}

	m.Logger.Debug("Aqara temp sensor raw message", "msg", rawMsg)

	// Convert data types for DB write
	ts := utils.GetUTCTimestamp()
	tempC := float32(rawMsg.Temperature)
	tempF := utils.CelsiusToFahrenheit(tempC)
	humidity := float32(rawMsg.Humidity)
	pressure := float32(rawMsg.Pressure)
	lq := float32(rawMsg.LinkQuality)
	bp := float32(rawMsg.Battery)
	bv := float32(rawMsg.Voltage)
	poc := int32(rawMsg.PowerOutageCount)
	devId := int32(m.TopicMapping[msg.Topic()].DeviceID)
	devTypeId := int32(m.TopicMapping[msg.Topic()].DeviceTypeID)

	aqaraUniqueParams := sqlc.InsertAqaraUniqueDataParams{
		Timestamp:        ts,
		LinkQuality:      &lq,
		BattPercentage:   &bp,
		BattVoltage:      &bv,
		PowerOutageCount: &poc,
		DeviceID:         &devId,
		DeviceTypeID:     &devTypeId,
	}

	sharedParams := sqlc.InsertReadingParams{
		Timestamp:        ts,
		TempF:            &tempF,
		TempC:            &tempC,
		Humidity:         &humidity,
		AbsolutePressure: &pressure,
		DeviceTypeID:     &devTypeId,
		DeviceID:         &devId,
	}

	db := sensordb.NewDbWrapper(m.PgPool)
	defer db.Conn.Release()

	err = db.DB.InsertAqaraUniqueData(m.Ctx, aqaraUniqueParams)
	if err != nil {
		m.Logger.Error("Error writing aqara temp sensor unique data", "error", err)
	} else {
		m.Logger.Info("Successful write to aqara unique data", "mqttTopic", msg.Topic(), "deviceType", m.TopicMapping[msg.Topic()].DeviceType, "deviceID", m.TopicMapping[msg.Topic()].DeviceID)
	}

	err = db.DB.InsertReading(m.Ctx, sharedParams)
	if err != nil {
		m.Logger.Error("Error writing aqara temp sensor shared data", "error", err)
	} else {
		m.Logger.Info("Successful write to shared data", "mqttTopic", msg.Topic(), "deviceType", m.TopicMapping[msg.Topic()].DeviceType, "deviceID", m.TopicMapping[msg.Topic()].DeviceID)
	}
}

func (m *MqttService) ProcessBme280SensorMsg(msg mqtt.Message) {
	m.Logger.Debug("Processing bme280 sensor message")

	var rawMsg devices.BME280
	err := json.Unmarshal(msg.Payload(), &rawMsg)
	if err != nil {
		m.Logger.Error("Error unmarshalling raw bme280 sensor message", "error", err)
	}

	m.Logger.Debug("BME280 raw message", "msg", rawMsg)

	ts := utils.GetUTCTimestamp()
	tempF := float32(rawMsg.TempF)
	tempC := float32(rawMsg.TempC)
	humidity := float32(rawMsg.Humidity)
	pressure := float32(0)

	deviceName := msg.Topic()
	devId := int32(m.DeviceListMap[deviceName].DeviceID)
	devTypeId := int32(m.DeviceListMap[deviceName].DeviceTypeID)

	sharedParams := sqlc.InsertReadingParams{
		Timestamp:        ts,
		TempF:            &tempF,
		TempC:            &tempC,
		Humidity:         &humidity,
		AbsolutePressure: &pressure,
		DeviceTypeID:     &devTypeId,
		DeviceID:         &devId,
	}

	m.Logger.Debug("Shared params",
		"timestamp", sharedParams.Timestamp,
		"tempF", *sharedParams.TempF,
		"tempC", *sharedParams.TempC,
		"humidity", *sharedParams.Humidity,
		"absolutePressure", *sharedParams.AbsolutePressure,
		"deviceTypeID", *sharedParams.DeviceTypeID,
		"deviceID", *sharedParams.DeviceID,
	)

	db := sensordb.NewDbWrapper(m.PgPool)
	defer db.Conn.Release()

	err = db.DB.InsertReading(m.Ctx, sharedParams)
	if err != nil {
		m.Logger.Error("Error writing bme280 temp sensor shared data", "error", err)
	} else {
		m.Logger.Info("Successful write to shared data", "mqttTopic", msg.Topic(), "deviceType", m.TopicMapping[msg.Topic()].DeviceType, "deviceID", m.TopicMapping[msg.Topic()].DeviceID)
	}
}

func (m *MqttService) ProcessBmp280SensorMsg(msg mqtt.Message) {
	m.Logger.Debug("Processing bmp280 sensor message")

	var rawMsg devices.BMP280
	err := json.Unmarshal(msg.Payload(), &rawMsg)
	if err != nil {
		m.Logger.Error("Error unmarshalling raw bmp280 sensor message", "error", err)
	}

	m.Logger.Debug("BMP280 raw message", "msg", rawMsg)

	ts := utils.GetUTCTimestamp()
	tempF := float32(rawMsg.TempF)
	tempC := float32(rawMsg.TempC)
	pressure := float32(0)

	deviceName := msg.Topic()
	devId := int32(m.DeviceListMap[deviceName].DeviceID)
	devTypeId := int32(m.DeviceListMap[deviceName].DeviceTypeID)

	sharedParams := sqlc.InsertReadingParams{
		Timestamp:        ts,
		TempF:            &tempF,
		TempC:            &tempC,
		AbsolutePressure: &pressure,
		DeviceTypeID:     &devTypeId,
		DeviceID:         &devId,
	}

	m.Logger.Debug("Shared params",
		"timestamp", sharedParams.Timestamp,
		"tempF", *sharedParams.TempF,
		"tempC", *sharedParams.TempC,
		"absolutePressure", *sharedParams.AbsolutePressure,
		"deviceTypeID", *sharedParams.DeviceTypeID,
		"deviceID", *sharedParams.DeviceID,
	)

	db := sensordb.NewDbWrapper(m.PgPool)
	defer db.Conn.Release()

	err = db.DB.InsertReading(m.Ctx, sharedParams)
	if err != nil {
		m.Logger.Error("Error writing bmp280 temp sensor shared data", "error", err)
	} else {
		m.Logger.Info("Successful write to shared data", "mqttTopic", msg.Topic(), "deviceType", m.TopicMapping[msg.Topic()].DeviceType, "deviceID", m.TopicMapping[msg.Topic()].DeviceID)
	}
}

func (m *MqttService) ProcessDht11TempSensorMsg(msg mqtt.Message) {
	m.Logger.Debug("Processing dht11 sensor message")

	var rawMsg devices.DHT11Message
	err := json.Unmarshal(msg.Payload(), &rawMsg)
	if err != nil {
		m.Logger.Error("Error unmarshalling raw dht11 temp sensor message", "error", err)
	}

	m.Logger.Debug("DHT11 raw message", "msg", rawMsg)

	ts := utils.GetUTCTimestamp()
	tempF := float32(rawMsg.TempF)
	tempC := float32(rawMsg.TempC)
	humidity := float32(rawMsg.Humidity)
	pressure := float32(0)

	// The dht11 sensors don't use the "topic/device_name" mqtt topic
	// convention.  The sensor name is embedded in the message, so for
	// now we're manually building the device name to match the device
	// in the database which is using the convention topic/device_name
	// In the future we need to make sure all sensors use the topic/device_name
	// convention when transmitting mqtt messages
	deviceName := fmt.Sprintf("%s/%s", msg.Topic(), rawMsg.SensorLocation)
	devId := int32(m.DeviceListMap[deviceName].DeviceID)
	devTypeId := int32(m.DeviceListMap[deviceName].DeviceTypeID)

	sharedParams := sqlc.InsertReadingParams{
		Timestamp:        ts,
		TempF:            &tempF,
		TempC:            &tempC,
		Humidity:         &humidity,
		AbsolutePressure: &pressure,
		DeviceTypeID:     &devTypeId,
		DeviceID:         &devId,
	}

	m.Logger.Debug("Shared params",
		"timestamp", sharedParams.Timestamp,
		"tempF", *sharedParams.TempF,
		"tempC", *sharedParams.TempC,
		"humidity", *sharedParams.Humidity,
		"absolutePressure", *sharedParams.AbsolutePressure,
		"deviceTypeID", *sharedParams.DeviceTypeID,
		"deviceID", *sharedParams.DeviceID,
	)

	db := sensordb.NewDbWrapper(m.PgPool)
	defer db.Conn.Release()

	err = db.DB.InsertReading(m.Ctx, sharedParams)
	if err != nil {
		m.Logger.Error("Error writing dht11 temp sensor shared data", "error", err)
	} else {
		m.Logger.Info("Successful write to shared data", "mqttTopic", msg.Topic(), "deviceType", m.TopicMapping[msg.Topic()].DeviceType, "deviceID", m.TopicMapping[msg.Topic()].DeviceID)
	}
}

func (m *MqttService) ProcessSonoffSmartPlugMsg(msg mqtt.Message) {
	m.Logger.Debug("Processing sonoff smart plug message")

	var rawMsg devices.SonoffSmartPlugRawMsg
	err := json.Unmarshal(msg.Payload(), &rawMsg)
	if err != nil {
		m.Logger.Error("Error unmarshalling raw sonoff smart plug message", "error", err)
	}

	m.Logger.Debug("Sonoff smart plug raw message", "msg", rawMsg)

	timestamp := utils.GetUTCTimestamp()
	lq := int32(rawMsg.LinkQuality)
	devId := int32(m.TopicMapping[msg.Topic()].DeviceID)
	devTypeId := int32(m.TopicMapping[msg.Topic()].DeviceTypeID)

	var state int32
	if rawMsg.State == "ON" {
		state = 1
	} else {
		state = 0
	}

	params := sqlc.InsertSonoffSmartPlugReadingParams{
		Timestamp:    timestamp,
		LinkQuality:  &lq,
		OutletState:  &state,
		DeviceID:     &devId,
		DeviceTypeID: &devTypeId,
	}

	db := sensordb.NewDbWrapper(m.PgPool)
	defer db.Conn.Release()

	err = db.DB.InsertSonoffSmartPlugReading(m.Ctx, params)
	if err != nil {
		m.Logger.Error("Error writing sonoff smart plug update", "error", err)
	} else {
		m.Logger.Info("Successful write to sonoff smart plug update", "mqttTopic", msg.Topic(), "deviceType", m.TopicMapping[msg.Topic()].DeviceType, "deviceID", m.TopicMapping[msg.Topic()].DeviceID)
	}

}

func (m *MqttService) OnConnectLostHndlrFactory() mqtt.ConnectionLostHandler {
	return func(client mqtt.Client, err error) {
		m.Logger.Error("MQTT connection lost", "error", err)
		m.Cleanup()
	}
}

func (m *MqttService) CancelListener() {
	<-m.Ctx.Done()
	m.Logger.Info("MQTT service context done signal detected - cleaning up")
	m.Cleanup()
}

func (m *MqttService) CreateMqttClient() {
	clientID := fmt.Sprintf("goSensorDataCollector-%s", uuid.New().String())
	m.Logger.Debug("Creating MQTT client", "clientID", clientID, "mqttBroker", m.Config.MqttBroker, "mqttPort", m.Config.MqttPort)

	// // -------------------- MQTT Client Setup -------------------- //
	mqttListenerOpts := mqtt.NewClientOptions()
	mqttListenerOpts.AddBroker(fmt.Sprintf("mqtt://%s:%d", m.Config.MqttBroker, m.Config.MqttPort))
	mqttListenerOpts.SetClientID(clientID)
	mqttListenerOpts.SetOrderMatters(false)
	mqttListenerOpts.SetDefaultPublishHandler(m.MqttMsgHndlrFactory())
	mqttListenerOpts.OnConnect = m.OnConnectHndlrFactory()
	mqttListenerOpts.OnConnectionLost = m.OnConnectLostHndlrFactory()
	mqttListenerOpts.SetKeepAlive(60 * time.Second)
	mqttListenerOpts.SetPingTimeout(10 * time.Second)

	newClient := createMqttClient(mqttListenerOpts)
	m.Client = newClient

	go m.CancelListener()
}

func (m *MqttService) Cleanup() {
	m.Logger.Info("MQTT cleanup called")

	for key := range m.TopicMapping {
		m.Client.Unsubscribe(key)
	}

	m.Client.Disconnect(100)

	m.Wg.Done()
}

func createMqttClient(opts *mqtt.ClientOptions) mqtt.Client {
	newClient := mqtt.NewClient(opts)
	if token := newClient.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("Error connecting to MQTT Broker", "error", token.Error())
		panic(token.Error())
	}
	return newClient
}
