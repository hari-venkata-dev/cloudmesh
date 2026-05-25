package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Telemetry struct {
	DeviceID    string `json:"deviceId"`
	CPU         int    `json:"cpu"`
	Memory      int    `json:"memory"`
	Temperature int    `json:"temperature"`
}

func main() {

	fmt.Println("CloudMesh Device Simulator Started")

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")

	client := mqtt.NewClient(opts)

	token := client.Connect()
	token.Wait()

	if token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Connected to MQTT Broker")

	deviceID := "AP-101"

	for {

		telemetry := Telemetry{
			DeviceID:    deviceID,
			CPU:         rand.Intn(100),
			Memory:      rand.Intn(100),
			Temperature: rand.Intn(40) + 30,
		}

		jsonData, err := json.Marshal(telemetry)

		if err != nil {
			panic(err)
		}

		fmt.Println(string(jsonData))

		token = client.Publish(
			"devices/telemetry",
			0,
			false,
			jsonData,
		)

		token.Wait()

		fmt.Println("Telemetry published to MQTT")

		time.Sleep(5 * time.Second)
	}
}
