# Giao thức truyền thông

## 1. ESP32 → Broker (MQTT / WiFi)

- **Broker**: `broker.emqx.io:1883`
- **Topic**: `uit/ce224/solar-tracker/data`
- **Payload**:
```json
{
  "volt": 5.18,
  "curr": 182.5,
  "power": 945.3,
  "energy": 12.4,
  "azi": 95,
  "ele": 42,
  "ldr": [2150, 2080, 1950, 1920],
  "mode": "auto"
}
```

## 2. ESP32 → Server (UART / USB)

Format mỗi dòng:

```
DATA:volt,curr,power,energy,azi,ele,ldr0,ldr1,ldr2,ldr3,mode
```

Ví dụ:
```
DATA:5.18,182.5,945.3,12.4,95,42,2150,2080,1950,1920,auto
```

> Baud rate: `115200`

## 3. Server → Trình duyệt (WebSocket)

- **Endpoint**: `ws://localhost:8081/ws`
- **Payload**:

```json
{
  "type": "history | live",
  "data": "SensorData | []SensorData"
}
```

- `history`: Mảng 60 điểm gần nhất, gửi ngay khi client kết nối
- `live`: Một điểm dữ liệu mới, gửi mỗi lần ESP32 bắn lên

## 4. Mô tả trường dữ liệu

| Trường | Đơn vị | Mô tả |
|--------|--------|-------|
| `volt` | V | Điện áp tấm pin mặt trời |
| `curr` | mA | Dòng điện |
| `power` | mW | Công suất tức thời |
| `energy` | Wh | Điện năng tích lũy |
| `azi` | ° | Góc servo Azimuth (trục ngang) |
| `ele` | ° | Góc servo Elevation (trục đứng) |
| `ldr` | raw ADC | 4 giá trị quang trở [TL, TR, BL, BR] |
| `mode` | string | Chế độ: `auto` hoặc `manual` |
