# Tài liệu Backend Server

Backend viết bằng **Go**, sử dụng **Fiber v3**. Nhận dữ liệu từ ESP32 qua MQTT hoặc UART, phát lại cho frontend qua WebSocket, lưu lịch sử vào InfluxDB.

## Yêu cầu hạ tầng

```bash
docker compose up -d   # khởi động InfluxDB 2.7
```

## Các module

### `internal/hub` — WebSocket Hub
- Quản lý các kết nối trình duyệt
- Duy trì sliding window 60 điểm trong RAM
- Gửi lịch sử ngay khi client kết nối

### `internal/sensor/mqtt.go` — MQTT receiver
- Kết nối tới `broker.emqx.io:1883`
- Subscribe topic `uit/ce224/solar-tracker/data`
- Parse JSON payload → `hub.SensorData` → broadcast

### `internal/sensor/uart.go` — UART receiver
- Mở cổng serial (mặc định `/dev/ttyUSB0`, baud 115200)
- Parse format `DATA:...`→ `hub.SensorData` → broadcast

> Chọn chế độ nhận bằng cách comment/uncomment trong `cmd/main.go`

### `internal/db/influx.go` — InfluxDB client
- Measurement: `solar`, tag: `device=solar_tracker_01`
- Lưu tất cả trường: volt, curr, power, energy, azi, ele, ldr

## API

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | `/ws` | WebSocket stream dữ liệu live |
| GET | `/health` | Kiểm tra server còn sống |
| GET | `/api/stats` | Thống kê trong ngày (max power, total energy) |

## Chạy server

```bash
make server   # chỉ chạy Go server
make dev      # chỉ chạy Svelte dev server
make app      # chạy cả hai + docker influxdb
```
