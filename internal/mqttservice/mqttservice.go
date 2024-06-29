package mqttservice

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"sensor-data-collection-service/internal/config"
	"sensor-data-collection-service/internal/db"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TopicData struct {
	DeviceID     int32
	DeviceTypeID int32
}

type MqttService struct {
	Ctx          context.Context
	Wg           *sync.WaitGroup
	Config       config.Config
	Client       mqtt.Client
	DB           *db.SensorDataDB
	Conn         *pgxpool.Conn
	TopicMapping map[string]TopicData
	Logger       *slog.Logger
}

func NewMqttService(wg *sync.WaitGroup, ctx context.Context, db *db.SensorDataDB, conn *pgxpool.Conn, cfg config.Config, logger *slog.Logger) MqttService {
	return MqttService{
		Ctx:    ctx,
		Wg:     wg,
		Config: cfg,
		DB:     db,
		Conn:   conn,
		Logger: logger,
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
		m.Logger.Info("Receiving MQTT message", "msg", msg.Payload(), "topic", msg.Topic())
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

	m.Conn.Release()

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
