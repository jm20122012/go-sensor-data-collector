package sensors

import (
	"encoding/json"
	"log"
	"sensor-data-collection-service/internal/db/sqlc"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DHT11Message struct {
	Temp_F         float64 `json:"temp_f"`
	Temp_C         float64 `json:"temp_c"`
	Humidity       float64 `json:"humidity"`
	SensorLocation string  `json:"device_id"`
}

type DHT11Record struct {
	sqlc.PiSensorDatum
}

func (d *DHT11Record) ParseData(data mqtt.Message, sensorIdMap map[string]int64) {
	err := json.Unmarshal(data.Payload(), &d)
	if err != nil {
		log.Println(err)
	}

	log.Println("Sensor data: ", d)

	tempF := float32(*d.TempF)
	tempC := float32(*d.TempC)
	humidity := float32(*d.Humidity)
	sensorID := int64(sensorIdMap[*d.SensorLocation])

}
