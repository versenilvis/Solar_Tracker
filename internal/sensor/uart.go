package sensor

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/versenilvis/solar-tracker/internal/hub"
	"go.bug.st/serial"
)

func ReadUart(h *hub.Hub, portName string) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatal("cannot open UART serial port: ", err)
	}

	log.Printf("[UART] listening on %s @ 115200", portName)

	scanner := bufio.NewScanner(port)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "DATA:") {
			continue
		}

		// format: DATA:volt,curr,power,energy,azi,ele,ldr0,ldr1,ldr2,ldr3,mode
		var volt, curr, power, energy float64
		var azi, ele int
		var ldr0, ldr1, ldr2, ldr3 int
		var mode string

		_, err := fmt.Sscanf(line[5:], "%f,%f,%f,%f,%d,%d,%d,%d,%d,%d,%s",
			&volt, &curr, &power, &energy,
			&azi, &ele,
			&ldr0, &ldr1, &ldr2, &ldr3,
			&mode)
		if err != nil {
			log.Printf("[UART] parse error: %v | raw: %s", err, line)
			continue
		}

		h.Broadcast(hub.SensorData{
			Volt:      volt,
			Curr:      curr,
			Power:     power,
			Energy:    energy,
			AziAngle:  azi,
			EleAngle:  ele,
			LDR:       [4]int{ldr0, ldr1, ldr2, ldr3},
			Mode:      mode,
			Timestamp: time.Now(),
		})
	}
}
