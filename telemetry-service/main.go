package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func messageHandler(client mqtt.Client, msg mqtt.Message) {

	fmt.Printf(
		"Received telemetry -> Topic: %s | Message: %s\n",
		msg.Topic(),
		msg.Payload(),
	)
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
