package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	_ "modernc.org/sqlite"
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

var deviceLastSeen = make(map[string]time.Time)
var deviceStatus = make(map[string]string)
var db *sql.DB

func saveTelemetry(telemetry Telemetry) {

	query := `
	INSERT INTO telemetry (
		device_id,
		device_type,
		location,
		cpu,
		memory,
		temperature
	)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(
		query,
		telemetry.DeviceID,
		telemetry.DeviceType,
		telemetry.Location,
		telemetry.CPU,
		telemetry.Memory,
		telemetry.Temperature,
	)

	if err != nil {
		fmt.Println("Failed to save telemetry:", err)
		return
	}

	fmt.Println("Telemetry saved to database")
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {

	var telemetry Telemetry

	err := json.Unmarshal(msg.Payload(), &telemetry)

	if err != nil {
		fmt.Println("Failed to parse telemetry:", err)
		return
	}

	deviceLastSeen[telemetry.DeviceID] = time.Now()
	currentStatus := deviceStatus[telemetry.DeviceID]

	if currentStatus == "OFFLINE" {

		fmt.Println("DEVICE RECOVERED")

		fmt.Printf(
			"Device %s is back ONLINE\n",
			telemetry.DeviceID,
		)

		fmt.Println("++++++++++++++++++++++++++++++++")
	}

	deviceStatus[telemetry.DeviceID] = "ONLINE"
	saveTelemetry(telemetry)
	fmt.Println("Telemetry Received")

	fmt.Printf("Device ID: %s\n", telemetry.DeviceID)
	fmt.Printf("Device Type: %s\n", telemetry.DeviceType)
	fmt.Printf("Location: %s\n", telemetry.Location)
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

func monitorOfflineDevices() {

	for {

		for deviceID, lastSeen := range deviceLastSeen {

			timeSinceLastSeen := time.Since(lastSeen)

			if timeSinceLastSeen > 15*time.Second {

				if deviceStatus[deviceID] != "OFFLINE" {

					deviceStatus[deviceID] = "OFFLINE"

					fmt.Println("DEVICE OFFLINE ALERT")

					fmt.Printf(
						"Device %s is now OFFLINE\n",
						deviceID,
					)

					fmt.Println("################################")
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}
func createTelemetryTable() {

	query := `
	CREATE TABLE IF NOT EXISTS telemetry (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id TEXT,
		device_type TEXT,
		location TEXT,
		cpu INTEGER,
		memory INTEGER,
		temperature INTEGER,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)

	if err != nil {
		panic(err)
	}

	fmt.Println("Telemetry table ready")
}
func main() {

	fmt.Println("CloudMesh Telemetry Service Started")
	var err error

	db, err = sql.Open("sqlite", "cloudmesh.db")

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to SQLite Database")
	createTelemetryTable()

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
	go monitorOfflineDevices()

	select {}
}
