package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/loafoe/go-marstek/pkg/marstek"
)

func main() {
	var (
		addr     string
		port     int
		timeout  time.Duration
		discover bool
	)

	flag.StringVar(&addr, "addr", "", "Device IP address")
	flag.IntVar(&port, "port", marstek.DefaultPort, "UDP port")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Request timeout")
	flag.BoolVar(&discover, "discover", false, "Discover devices on LAN")
	flag.Parse()

	if discover {
		fmt.Println("Discovering Marstek devices on LAN...")
		devices, err := marstek.DiscoverDevices(port, timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Discovery error: %v\n", err)
			os.Exit(1)
		}
		if len(devices) == 0 {
			fmt.Println("No devices found")
			return
		}
		for _, d := range devices {
			printJSON("Device", d)
		}
		return
	}

	if addr == "" {
		fmt.Fprintln(os.Stderr, "Error: -addr is required")
		flag.Usage()
		os.Exit(1)
	}

	client := marstek.NewClient(addr,
		marstek.WithPort(port),
		marstek.WithTimeout(timeout),
	)

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"status"}
	}

	cmd := args[0]

	switch cmd {
	case "device":
		bleMac := "0"
		if len(args) > 1 {
			bleMac = args[1]
		}
		info, err := client.GetDevice(bleMac)
		handleResult("Device Info", info, err)

	case "wifi":
		status, err := client.GetWifiStatus(0)
		handleResult("WiFi Status", status, err)

	case "ble":
		status, err := client.GetBLEStatus(0)
		handleResult("BLE Status", status, err)

	case "battery", "bat":
		status, err := client.GetBatteryStatus(0)
		handleResult("Battery Status", status, err)

	case "pv", "solar":
		status, err := client.GetPVStatus(0)
		handleResult("PV Status", status, err)

	case "es":
		status, err := client.GetESStatus(0)
		handleResult("Energy System Status", status, err)

	case "mode":
		status, err := client.GetESMode(0)
		handleResult("Energy System Mode", status, err)

	case "em", "meter":
		status, err := client.GetEMStatus(0)
		handleResult("Energy Meter Status", status, err)

	case "set-auto":
		result, err := client.SetAutoMode(0)
		handleResult("Set Auto Mode", result, err)

	case "set-ai":
		result, err := client.SetAIMode(0)
		handleResult("Set AI Mode", result, err)

	case "set-ups":
		result, err := client.SetUPSMode(0)
		handleResult("Set UPS Mode", result, err)

	case "set-passive":
		power := 100
		cdTime := 300
		if len(args) > 1 {
			fmt.Sscanf(args[1], "%d", &power)
		}
		if len(args) > 2 {
			fmt.Sscanf(args[2], "%d", &cdTime)
		}
		result, err := client.SetPassiveMode(0, power, cdTime)
		handleResult("Set Passive Mode", result, err)

	case "led-on":
		result, err := client.SetLED(true)
		handleResult("LED On", result, err)

	case "led-off":
		result, err := client.SetLED(false)
		handleResult("LED Off", result, err)

	case "status":
		fmt.Println("=== Marstek Device Status ===")

		info, err := client.GetDevice("0")
		if err != nil {
			fmt.Printf("Device Info: Error - %v\n\n", err)
		} else {
			printJSON("Device Info", info)
		}

		bat, err := client.GetBatteryStatus(0)
		if err != nil {
			fmt.Printf("Battery Status: Error - %v\n\n", err)
		} else {
			printJSON("Battery Status", bat)
		}

		es, err := client.GetESStatus(0)
		if err != nil {
			fmt.Printf("Energy System: Error - %v\n\n", err)
		} else {
			printJSON("Energy System Status", es)
		}

		mode, err := client.GetESMode(0)
		if err != nil {
			fmt.Printf("Mode: Error - %v\n\n", err)
		} else {
			printJSON("Current Mode", mode)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		fmt.Fprintln(os.Stderr, "Available commands: device, wifi, ble, battery, pv, es, mode, em, set-auto, set-ai, set-ups, set-passive, led-on, led-off, status")
		os.Exit(1)
	}
}

func handleResult(name string, result any, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s Error: %v\n", name, err)
		os.Exit(1)
	}
	printJSON(name, result)
}

func printJSON(name string, v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON error: %v\n", err)
		return
	}
	fmt.Printf("%s:\n%s\n\n", name, string(data))
}
