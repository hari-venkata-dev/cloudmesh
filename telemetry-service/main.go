package main

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Telemetry struct {
	DeviceID    string `json:"deviceId"`
	CPU         int    `json:"cpu"`
	Memory      int    `json:"memory"`
	Temperature int    `json:"temperature"`
}

var deviceLastSeen = make(map[string]time.Time)

func messageHandler(client mqtt.Client, msg mqtt.Message) {

	var telemetry Telemetry

	err := json.Unmarshal(msg.Payload(), &telemetry)

	if err != nil {
		fmt.Println("Failed to parse telemetry:", err)
		return
	}

	deviceLastSeen[telemetry.DeviceID] = time.Now()
	fmt.Println("Telemetry Received")

	fmt.Printf("Device ID: %s\n", telemetry.DeviceID)
	fmt.Printf("CPU Usage: %d%%\n", telemetry.CPU)
	fmt.Printf("Memory Usage: %d%%\n", telemetry.Memory)
	fmt.Printf("Temperature: %d°C\n", telemetry.Temperature)

	fmt.Println("--------------------------------")
	fmt.Printf(
		"Last Seen: %s\n",
		deviceLastSeen[telemetry.DeviceID].Format(time.RFC3339),
	)
	if telemetry.CPU > 80 {
		fmt.Println("HIGH CPU ALERT")
		fmt.Printf(
			"Device %s CPU usage is critically high: %d%%\n",
			telemetry.DeviceID,
			telemetry.CPU,
		)

		fmt.Println("********************************")
	}
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
