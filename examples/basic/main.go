package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/loafoe/go-marstek/pkg/marstek"
)

func main() {
	addr := os.Getenv("MARSTEK_IP")
	if addr == "" {
		addr = "192.168.2.189"
	}

	client := marstek.NewClient(addr,
		marstek.WithTimeout(10*time.Second),
	)

	info, err := client.GetDevice("0")
	if err != nil {
		log.Fatalf("Failed to get device info: %v", err)
	}
	fmt.Printf("Connected to %s (firmware v%d)\n", info.Device, info.Version)
	fmt.Printf("IP: %s, WiFi: %s\n\n", info.IP, info.WifiName)

	es, err := client.GetESStatus(0)
	if err != nil {
		log.Fatalf("Failed to get ES status: %v", err)
	}

	fmt.Println("=== Energy System Status ===")
	if es.BatSOC != nil {
		fmt.Printf("Battery SOC: %d%%\n", *es.BatSOC)
	}
	if es.BatCap != nil {
		fmt.Printf("Battery Capacity: %.0f Wh\n", *es.BatCap)
	}
	if es.PVPower != nil {
		fmt.Printf("Solar Power: %.0f W\n", *es.PVPower)
	}
	if es.OngridPower != nil {
		power := *es.OngridPower
		if power < 0 {
			fmt.Printf("Grid: Exporting %.0f W\n", -power)
		} else if power > 0 {
			fmt.Printf("Grid: Importing %.0f W\n", power)
		} else {
			fmt.Println("Grid: Idle")
		}
	}

	mode, err := client.GetESMode(0)
	if err != nil {
		log.Fatalf("Failed to get mode: %v", err)
	}
	fmt.Printf("\nCurrent Mode: %s\n", mode.Mode)

	em, err := client.GetEMStatus(0)
	if err != nil {
		log.Printf("Warning: Failed to get energy meter status: %v", err)
	} else if em.CTState != nil && *em.CTState == 1 {
		fmt.Println("\n=== CT Clamp Readings ===")
		if em.APower != nil {
			fmt.Printf("Phase A: %.0f W\n", *em.APower)
		}
		if em.BPower != nil {
			fmt.Printf("Phase B: %.0f W\n", *em.BPower)
		}
		if em.CPower != nil {
			fmt.Printf("Phase C: %.0f W\n", *em.CPower)
		}
		if em.TotalPower != nil {
			fmt.Printf("Total: %.0f W\n", *em.TotalPower)
		}
	}
}
