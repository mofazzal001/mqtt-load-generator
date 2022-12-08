package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	MQTTClient "github.com/pablitovicente/mqtt-load-generator/pkg/MQTTClient"
	"github.com/schollz/progressbar/v3"

	"github.com/google/uuid"
)

func main() {
	// Argument parsing
	targetTopic := flag.String("t", "/load", "Target MQTT topic to publish messages to")
	username := flag.String("u", "", "MQTT username")
	password := flag.String("P", "", "MQTT password")
	host := flag.String("h", "localhost", "MQTT host")
	port := flag.Int("p", 1883, "MQTT port")
	qos := flag.Int("q", 1, "MQTT QoS used by all clients")

	flag.Parse()

	if *qos < 0 || *qos > 2 {
		panic("QoS should be any of [0, 1, 2]")
	}

	fmt.Println("press ctrl+c to exit")

	// General Client Config
	mqttClientConfig := MQTTClient.Config{
		TargetTopic: targetTopic,
		Username:    username,
		Password:    password,
		Host:        host,
		Port:        port,
		QoS:         qos,
	}

	rand.Seed(time.Now().UnixNano())
	updates := make(chan int)

	mqttClient := MQTTClient.Client{
		ID:      uuid.NewString(),
		Config:  mqttClientConfig,
		Updates: updates,
	}

	mqttClient.Connect()

	mqttClient.Subscribe(*targetTopic)
	bar := progressbar.Default(-1)
	go func(updates chan int) {
		for update := range updates {
			bar.Add(update)
		}
	}(updates)

	// There's some issue with bar update when traffic is not constant
	// so this go routine updates the bar with 0 just to get the total numbers right
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			// Block until the clock ticks
			<-ticker.C
			// Update bar with 0 to update total
			bar.Add(0)
		}
	}()

	select {}
}
