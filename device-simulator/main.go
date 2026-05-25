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
	DeviceType  string `json:"deviceType"`
	Location    string `json:"location"`
	CPU         int    `json:"cpu"`
	Memory      int    `json:"memory"`
	Temperature int    `json:"temperature"`
}

type Device struct {
	ID       string
	Type     string
	Location string
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

	devices := []Device{
		{
			ID:       "AP-101",
			Type:     "Access Point",
			Location: "Building A",
		},
		{
			ID:       "AP-102",
			Type:     "Access Point",
			Location: "Building B",
		},
		{
			ID:       "SW-201",
			Type:     "Switch",
			Location: "Floor 1",
		},
		{
			ID:       "RTR-301",
			Type:     "Router",
			Location: "Datacenter",
		},
	}

	for {

		device := devices[rand.Intn(len(devices))]

		telemetry := Telemetry{
			DeviceID:    device.ID,
			DeviceType:  device.Type,
			Location:    device.Location,
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
