package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type telemetry struct {
	Volt     float64 `json:"volt"`
	Curr     float64 `json:"curr"`
	Power    float64 `json:"power"`
	Energy   float64 `json:"energy"`
	AziAngle int     `json:"azi"`
	EleAngle int     `json:"ele"`
	LDR      [4]int  `json:"ldr"`
	Mode     string  `json:"mode"`
}

const topic = "uit/ce224/solar-tracker/data"

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://broker.emqx.io:1883")
	opts.SetClientID(fmt.Sprintf("mock_esp32_%d", time.Now().UnixNano()))
	opts.SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("mqtt connect failed: %v", token.Error())
	}
	defer client.Disconnect(250)

	log.Println("[MOCK] connected to EMQX, broadcasting solar tracker data...")
	log.Printf("[MOCK] topic: %s (press Ctrl+C to stop)", topic)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	t := 0.0
	energy := 0.0

	for range ticker.C {
		t += 0.05
		azi := int(90 + 40*math.Sin(t*0.5))
		ele := int(45 + 25*math.Sin(t*0.3))

		volt := 5.0 + 0.2*math.Sin(t) + (rand.Float64()-0.5)*0.05
		curr := 180.0 + 30.0*math.Sin(t*0.8) + (rand.Float64()-0.5)*5.0
		power := volt * curr
		energy += (power * 1.0) / 3600000.0

		baseLdr := int(2200 + 400*math.Sin(t*0.6))
		ldr := [4]int{
			baseLdr + rand.Intn(100) - 50,
			baseLdr + rand.Intn(100) - 50,
			baseLdr + rand.Intn(100) - 50,
			baseLdr + rand.Intn(100) - 50,
		}

		payload := telemetry{
			Volt:     math.Round(volt*100) / 100,
			Curr:     math.Round(curr*10) / 10,
			Power:    math.Round(power*10) / 10,
			Energy:   math.Round(energy*10000) / 10000,
			AziAngle: azi,
			EleAngle: ele,
			LDR:      ldr,
			Mode:     "auto",
		}

		data, _ := json.Marshal(payload)
		token := client.Publish(topic, 0, false, data)
		token.Wait()

		log.Printf("[MOCK] %s", string(data))
	}
}
