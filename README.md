# go-marstek

Go client library for Marstek battery systems (Venus E3, Venus C, Venus D, Venus A) local API.

## Disclaimer

This project is not affiliated with, endorsed by, or connected to Marstek in any way. This is an independent, community-developed library based on publicly available API documentation.

**USE AT YOUR OWN RISK.** This software interacts with battery hardware and energy systems. Improper use could potentially affect your battery system's operation. The authors and contributors are not responsible for any damage, data loss, or other issues that may arise from using this software. Always ensure you understand the commands you are sending to your device.

## Installation

```bash
go get github.com/loafoe/go-marstek
```

## Usage

### Library

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/loafoe/go-marstek/pkg/marstek"
)

func main() {
    client := marstek.NewClient("192.168.1.100",
        marstek.WithTimeout(10*time.Second),
    )

    // Get device info
    info, err := client.GetDevice("0")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Device: %s (v%d)\n", info.Device, info.Version)

    // Get energy system status
    es, err := client.GetESStatus(0)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Battery: %d%%, Grid: %.0fW\n", *es.BatSOC, *es.OngridPower)

    // Get current operating mode
    mode, err := client.GetESMode(0)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Mode: %s\n", mode.Mode)
}
```

### CLI

```bash
# Build the CLI
go build -o marstek ./cmd/marstek

# Get full status
./marstek -addr 192.168.1.100 status

# Get specific information
./marstek -addr 192.168.1.100 device
./marstek -addr 192.168.1.100 battery
./marstek -addr 192.168.1.100 es
./marstek -addr 192.168.1.100 mode
./marstek -addr 192.168.1.100 em

# Set operating mode
./marstek -addr 192.168.1.100 set-auto
./marstek -addr 192.168.1.100 set-ai
./marstek -addr 192.168.1.100 set-ups
./marstek -addr 192.168.1.100 set-passive 100 300  # power(W) countdown(s)

# Control LED
./marstek -addr 192.168.1.100 led-on
./marstek -addr 192.168.1.100 led-off

# Discover devices on LAN
./marstek -discover
```

## API Methods

### Device & Network

| Method | Description |
|--------|-------------|
| `GetDevice(bleMac)` | Get device information |
| `GetWifiStatus(id)` | Get WiFi connection status |
| `GetBLEStatus(id)` | Get Bluetooth status |

### Energy System

| Method | Description |
|--------|-------------|
| `GetBatteryStatus(id)` | Get battery SOC, temperature, capacity |
| `GetPVStatus(id)` | Get solar panel power/voltage/current |
| `GetESStatus(id)` | Get energy system status (power, energy totals) |
| `GetESMode(id)` | Get current operating mode |
| `GetEMStatus(id)` | Get energy meter / CT clamp readings |

### Configuration

| Method | Description |
|--------|-------------|
| `SetAutoMode(id)` | Switch to Auto mode |
| `SetAIMode(id)` | Switch to AI mode |
| `SetManualMode(id, cfg)` | Set manual schedule |
| `SetPassiveMode(id, power, cdTime)` | Set passive control mode |
| `SetUPSMode(id)` | Switch to UPS mode |
| `SetDOD(value)` | Set Depth of Discharge (30-88%) |
| `SetBLEAdvertising(enable)` | Enable/disable Bluetooth |
| `SetLED(on)` | Control panel LED |

## Operating Modes

- **Auto**: Automatic optimization based on solar production and consumption
- **AI**: AI-driven optimization (requires cloud connection)
- **Manual**: User-defined time-based schedules
- **Passive**: External control via power setpoint with countdown
- **UPS**: Uninterruptible power supply mode

## Device Rate Limiting

The Marstek devices require approximately 2 seconds between API requests. The client handles this automatically with built-in rate limiting.

## Prerequisites

1. Device must be connected to your local network (WiFi or Ethernet)
2. Open API feature must be enabled in the Marstek mobile app
3. Default UDP port is 30000 (configurable in the app, recommended range: 49152-65535)

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.
