# Tài liệu Firmware

Firmware chạy trên **ESP32 DevKit V4** sử dụng nền tảng **ESP-IDF v5.5.1**.

## Thành phần phần cứng

| STT | Linh kiện | Giao tiếp |
|-----|-----------|-----------|
| 1 | 4 Quang trở LDR + trở 10k | ADC1 (GPIO 34, 35, 32, 33) |
| 2 | 2 Servo MG90S | PWM LEDC (GPIO 18, 19) |
| 3 | Cảm biến dòng/áp INA219 | I2C (SDA: 21, SCL: 22, Addr: 0x40) |
| 4 | Màn hình OLED SSD1306 | I2C (SDA: 21, SCL: 22, Addr: 0x3C) |

## Các task FreeRTOS

### 1. `solar_tracking_task` (100ms)
- Đọc 4 LDR, lọc nhiễu 32 mẫu
- Tính chênh lệch sáng theo trục ngang và đứng
- Nếu chênh lệch vượt ngưỡng deadband → xoay servo tương ứng

### 2. `monitoring_task` (500ms)
- Đọc INA219: điện áp, dòng, công suất
- Cập nhật OLED dashboard
- In dữ liệu ra Serial theo format `DATA:...`

## Format xuất Serial

```
DATA:volt,curr,power,energy,azi,ele,ldr0,ldr1,ldr2,ldr3,mode
```

> Xem thêm định nghĩa đầy đủ tại [`docs/protocol.md`](protocol.md)

## Build và nạp firmware

```bash
make build    # biên dịch
make flash    # nạp lên ESP32
make monitor  # xem serial log
```
