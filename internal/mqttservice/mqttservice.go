package mqttservice

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"sensor-data-collection-service/internal/db"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type MqttService struct {
	Config MqttConfig
	Client mqtt.Client
	DBConn *db.SensorDataDB
}

func NewMqttService() MqttService {
	config, err := readConfig("mqtt_config.yml")
	if err != nil {
		slog.Error("Error loading mqtt yaml file", "error", err)
	}
	return MqttService{
		Config: *config,
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
	PublishConnectSuccessful(client)
}

func PublishConnectSuccessful(client mqtt.Client) {
	connectionMsg := fmt.Sprintf("[%s] Connected to MQTT Broker at %s:%s", time.Now().String(), os.Getenv("MQTT_BROKER_IP"), os.Getenv("MQTT_BROKER_PORT"))
	token := client.Publish("connections", 0, false, connectionMsg)
	token.Wait()
}

func (m *MqttService) MqttSubscribe(topic string) {
	token := m.Client.Subscribe(topic, 1, nil)
	token.Wait()
	log.Println("MQTT Subscribed to topic: ", topic)
}

func (m *MqttService) CreateMqttClient(connectLostHdnlr mqtt.ConnectionLostHandler, msgHndlr mqtt.MessageHandler) {
	log.Println("Creating MQTT client...")
	clientID := fmt.Sprintf("goSensorDataCollector-%s", uuid.New().String())
	log.Println("Client ID: ", clientID)
	log.Println("MQTT Broker IP: ", m.Config.Server)
	log.Println("MQTT Broker Port: ", m.Config.Port)

	// // -------------------- MQTT Client Setup -------------------- //
	mqttListenerOpts := mqtt.NewClientOptions()
	mqttListenerOpts.AddBroker(fmt.Sprintf("mqtt://%s:%d", m.Config.Server, m.Config.Port))
	mqttListenerOpts.SetClientID(clientID)
	mqttListenerOpts.SetOrderMatters(false)
	mqttListenerOpts.SetDefaultPublishHandler(msgHndlr)
	mqttListenerOpts.OnConnect = connectHandler
	mqttListenerOpts.OnConnectionLost = connectLostHdnlr
	mqttListenerOpts.SetKeepAlive(60 * time.Second)
	mqttListenerOpts.SetPingTimeout(10 * time.Second)

	newClient := createMqttClient(mqttListenerOpts)
	m.Client = newClient
}

func createMqttClient(opts *mqtt.ClientOptions) mqtt.Client {
	newClient := mqtt.NewClient(opts)
	if token := newClient.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Error connecting to MQTT Broker: %v", token.Error())
		panic(token.Error())
	}
	return newClient
}

func readConfig(filename string) (*MqttConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	var config MqttConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}

	return &config, nil
}
