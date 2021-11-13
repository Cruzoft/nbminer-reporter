package main

// Import resty into your code and refer it as `resty`.
import (
    log "github.com/sirupsen/logrus"
)

func main() {
	statusData := requestGet("http://host.docker.internal:8000/api/v1/status")
	status, _ := parseStatus(statusData)

	for _, device := range status.Miner.Devices {
		log.Printf("Found Device %s in pcie %i.", device.Info, device.PCIBusId)
	}
	log.Printf("Total Hashrate: %s.", status.Miner.TotalHashrate)
}