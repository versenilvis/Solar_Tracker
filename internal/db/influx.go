package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/versenilvis/solar-tracker/internal/hub"
)

var client influxdb2.Client

func Init() {
	url := os.Getenv("INFLUX_URL")
	token := os.Getenv("INFLUX_TOKEN")
	if url == "" || token == "" {
		log.Println("[INFLUX] missing INFLUX_URL or INFLUX_TOKEN in .env, storage disabled")
		return
	}
	client = influxdb2.NewClient(url, token)
	log.Println("[INFLUX] connected to InfluxDB")
}

func WriteSensorData(data hub.SensorData) {
	if client == nil {
		return
	}
	org := os.Getenv("INFLUX_ORG")
	bucket := os.Getenv("INFLUX_BUCKET")
	writeAPI := client.WriteAPIBlocking(org, bucket)

	p := influxdb2.NewPointWithMeasurement("solar").
		AddTag("device", "solar_tracker_01").
		AddField("volt", data.Volt).
		AddField("curr", data.Curr).
		AddField("power", data.Power).
		AddField("energy", data.Energy).
		AddField("azi", data.AziAngle).
		AddField("ele", data.EleAngle).
		AddField("ldr_tl", data.LDR[0]).
		AddField("ldr_tr", data.LDR[1]).
		AddField("ldr_bl", data.LDR[2]).
		AddField("ldr_br", data.LDR[3]).
		SetTime(time.Now())

	if err := writeAPI.WritePoint(context.Background(), p); err != nil {
		log.Printf("[INFLUX] write error: %v", err)
	}
}

type DailyStats struct {
	MaxPower  float64 `json:"max_power"`
	TotalEnergy float64 `json:"total_energy"`
}

func GetTodayStats() (DailyStats, error) {
	if client == nil {
		return DailyStats{}, fmt.Errorf("influxdb client not initialized")
	}

	org := os.Getenv("INFLUX_ORG")
	bucket := os.Getenv("INFLUX_BUCKET")
	queryAPI := client.QueryAPI(org)

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startStr := startOfDay.Format(time.RFC3339)

	query := fmt.Sprintf(`
		from(bucket:"%s")
			|> range(start: %s)
			|> filter(fn: (r) => r._measurement == "solar")
			|> filter(fn: (r) => r._field == "power")
			|> max()
	`, bucket, startStr)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return DailyStats{}, err
	}

	var stats DailyStats
	for result.Next() {
		if v, ok := result.Record().Value().(float64); ok {
			stats.MaxPower = v
		}
	}

	return stats, result.Err()
}
