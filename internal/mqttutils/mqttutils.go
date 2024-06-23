package mqttutils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
	PublishConnectSuccessful(client)
}

func PublishConnectSuccessful(client mqtt.Client) {
	connectionMsg := fmt.Sprintf("[%s] Connected to MQTT Broker at %s:%s", time.Now().String(), os.Getenv("MQTT_BROKER_IP"), os.Getenv("MQTT_BROKER_PORT"))
	token := client.Publish("connections", 0, false, connectionMsg)
	token.Wait()
}

func MqttSubscribe(client mqtt.Client) {
	topic := os.Getenv("MQTT_SUB_TOPIC")
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	log.Println("MQTT Subscribed to topic: ", topic)
}

func CreateMqttClient(connectLostHandler mqtt.ConnectionLostHandler, msgHandler mqtt.MessageHandler) mqtt.Client {
	mqttBroker := os.Getenv("MQTT_BROKER_IP")
	mqttPort, err := strconv.Atoi(os.Getenv("MQTT_BROKER_PORT"))
	if err != nil {
		log.Println("Error converting MQTT port to INT")
	}
	clientID := fmt.Sprintf("goSensorDataCollector-%s", uuid.New().String())
	log.Println("Client ID: ", clientID)
	log.Println("MQTT Broker IP: ", mqttBroker)
	log.Println("MQTT Broker Port: ", mqttPort)

	// // -------------------- MQTT Setup -------------------- //
	mqttListenerOpts := mqtt.NewClientOptions()
	mqttListenerOpts.AddBroker(fmt.Sprintf("mqtt://%s:%d", mqttBroker, mqttPort))
	mqttListenerOpts.SetClientID(clientID)
	mqttListenerOpts.SetOrderMatters(false)
	mqttListenerOpts.SetDefaultPublishHandler(msgHandler)
	mqttListenerOpts.OnConnect = connectHandler
	mqttListenerOpts.OnConnectionLost = connectLostHandler
	mqttListenerOpts.SetKeepAlive(60 * time.Second)
	mqttListenerOpts.SetPingTimeout(10 * time.Second)

	newClient := mqtt.NewClient(mqttListenerOpts)
	if token := newClient.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Error connecting to MQTT Broker: %v", token.Error())
		panic(token.Error())
	}

	return newClient
}
