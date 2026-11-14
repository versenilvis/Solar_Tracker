package sensor

import (
	"encoding/json"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/versenilvis/solar-tracker/internal/db"
	"github.com/versenilvis/solar-tracker/internal/hub"
)

type MQTTPayload struct {
	Volt      float64 `json:"volt"`
	Curr      float64 `json:"curr"`
	Power     float64 `json:"power"`
	Energy    float64 `json:"energy"`
	AziAngle  int     `json:"azi"`
	EleAngle  int     `json:"ele"`
	LDR       [4]int  `json:"ldr"`
	Mode      string  `json:"mode"`
}

const mqttTopic = "uit/ce224/solar-tracker/data"

func ReadMQTT(h *hub.Hub) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://broker.emqx.io:1883")
	opts.SetClientID("solar_tracker_backend")
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(time.Second * 5)

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("[MQTT] connection lost: %v", err)
	}

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("[MQTT] connected to EMQX broker")
		log.Printf("[MQTT] subscribing to topic: %s", mqttTopic)

		if token := c.Subscribe(mqttTopic, 1, nil); token.Wait() && token.Error() != nil {
			log.Printf("[MQTT] subscribe error: %v", token.Error())
		}
	}

	client := mqtt.NewClient(opts)

	client.AddRoute(mqttTopic, func(c mqtt.Client, m mqtt.Message) {
		var payload MQTTPayload
		if err := json.Unmarshal(m.Payload(), &payload); err != nil {
			log.Printf("[MQTT] payload parse error: %v | raw: %s", err, string(m.Payload()))
			return
		}

		log.Printf("[MQTT] received: Volt=%.2f Curr=%.1f Power=%.1f Azi=%d Ele=%d Mode=%s",
			payload.Volt, payload.Curr, payload.Power, payload.AziAngle, payload.EleAngle, payload.Mode)

		data := hub.SensorData{
			Volt:      payload.Volt,
			Curr:      payload.Curr,
			Power:     payload.Power,
			Energy:    payload.Energy,
			AziAngle:  payload.AziAngle,
			EleAngle:  payload.EleAngle,
			LDR:       payload.LDR,
			Mode:      payload.Mode,
			Timestamp: time.Now(),
		}

		go db.WriteSensorData(data)
		h.Broadcast(data)
	})

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("[MQTT] connect error: %v", token.Error())
	}
}
