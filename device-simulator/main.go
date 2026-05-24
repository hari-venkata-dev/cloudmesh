package main

import (
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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

		cpuUsage := rand.Intn(100)
		memoryUsage := rand.Intn(100)
		temperature := rand.Intn(40) + 30

		telemetry := fmt.Sprintf(
			"Device: %s | CPU: %d%% | Memory: %d%% | Temperature: %d°C",
			deviceID,
			cpuUsage,
			memoryUsage,
			temperature,
		)

		fmt.Println(telemetry)

		token = client.Publish(
			"devices/telemetry",
			0,
			false,
			telemetry,
		)

		token.Wait()

		fmt.Println("Telemetry published to MQTT")

		time.Sleep(5 * time.Second)
	}
}
