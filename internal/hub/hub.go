package hub

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/v3/websocket"
)

type SensorData struct {
	Volt      float64   `json:"volt"`
	Curr      float64   `json:"curr"`
	Power     float64   `json:"power"`
	Energy    float64   `json:"energy"`
	AziAngle  int       `json:"azi"`
	EleAngle  int       `json:"ele"`
	LDR       [4]int    `json:"ldr"`
	Mode      string    `json:"mode"`
	Timestamp time.Time `json:"timestamp"`
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Hub struct {
	mu         sync.RWMutex
	clients    map[*websocket.Conn]struct{}
	history    []SensorData
	maxHistory int
}

func New(maxHistory int) *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]struct{}),
		history:    make([]SensorData, 0, maxHistory),
		maxHistory: maxHistory,
	}
}

func (h *Hub) Register(c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[c] = struct{}{}

	msg := Message{
		Type: "history",
		Data: h.history,
	}
	payload, _ := json.Marshal(msg)
	_ = c.WriteMessage(websocket.TextMessage, payload)
}

func (h *Hub) Unregister(c *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
}

func (h *Hub) Broadcast(data SensorData) {
	h.mu.Lock()
	if len(h.history) >= h.maxHistory {
		h.history = h.history[1:]
	}
	h.history = append(h.history, data)
	h.mu.Unlock()

	log.Printf("[SENSOR] Volt: %.2fV | Curr: %.1fmA | Power: %.1fmW | Azi: %d° | Ele: %d° | Mode: %s",
		data.Volt, data.Curr, data.Power, data.AziAngle, data.EleAngle, data.Mode)

	msg := Message{
		Type: "live",
		Data: data,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.clients {
		_ = c.WriteMessage(websocket.TextMessage, payload)
	}
}
