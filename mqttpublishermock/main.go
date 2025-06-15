package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func getRandomInRange(from, to float64, fixed int) float64 {
	value := rand.Float64()*(to-from) + from
	scale := math.Pow(10, float64(fixed))
	return math.Round(value*scale) / scale
}

func getRandomLocation() (lat, lon float64) {
	lat = getRandomInRange(-6.4, -6.1, 6)
	lon = getRandomInRange(106.6, 107.0, 6)

	return lat, lon
}

type VehicleLocation struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	vehicleID := "B1234XYZ"
	broker := "tcp://localhost:1883"
	topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicleID)
	clientID := "go-publisher"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connected to MQTT broker...")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		lat, lon := getRandomLocation()

		location := VehicleLocation{
			VehicleID: "B1234XYZ",
			Latitude:  lat,
			Longitude: lon,
			Timestamp: t.Unix(),
		}

		message, err := json.Marshal(location)
		if err != nil {
			fmt.Println("Failed to marshal")
			continue
		}

		token := client.Publish(topic, 1, false, message)
		token.Wait()

		fmt.Printf("Publish message to %s topic is success", topic)
	}
}
