package main

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Telemetry struct {
	DeviceID    string `json:"deviceId"`
	CPU         int    `json:"cpu"`
	Memory      int    `json:"memory"`
	Temperature int    `json:"temperature"`
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {

	var telemetry Telemetry

	err := json.Unmarshal(msg.Payload(), &telemetry)

	if err != nil {
		fmt.Println("Failed to parse telemetry:", err)
		return
	}

	fmt.Println("Telemetry Received")

	fmt.Printf("Device ID: %s\n", telemetry.DeviceID)
	fmt.Printf("CPU Usage: %d%%\n", telemetry.CPU)
	fmt.Printf("Memory Usage: %d%%\n", telemetry.Memory)
	fmt.Printf("Temperature: %d°C\n", telemetry.Temperature)

	fmt.Println("--------------------------------")
}

func main() {

	fmt.Println("CloudMesh Telemetry Service Started")

	opts := mqtt.NewClientOptions()

	opts.AddBroker("tcp://localhost:1883")

	client := mqtt.NewClient(opts)

	token := client.Connect()

	token.Wait()

	if token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Connected to MQTT Broker")

	token = client.Subscribe(
		"devices/telemetry",
		0,
		messageHandler,
	)

	token.Wait()

	if token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Subscribed to devices/telemetry")

	select {}
}
