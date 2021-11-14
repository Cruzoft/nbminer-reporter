package main

import (
	"os"
	"os/signal"
	"sync"
	"time"

    log "github.com/sirupsen/logrus"
)

func checkMinerStatus() {
	statusData := requestGet("http://host.docker.internal:8000/api/v1/status")
	status, _ := parseStatus(statusData)

	for _, device := range status.Miner.Devices {
		log.Printf("Found Device %s in pcie %i.", device.Info, device.PCIBusId)
	}
	log.Printf("Total Hashrate: %s.", status.Miner.TotalHashrate)
	log.Printf("Parsing to Influx line")

	writeToInflux(status)
}

func WaitForCtrlC() {
    var end_waiter sync.WaitGroup
    end_waiter.Add(1)
    var signal_channel chan os.Signal
    signal_channel = make(chan os.Signal, 1)
    signal.Notify(signal_channel, os.Interrupt)
    go func() {
		<-signal_channel
        end_waiter.Done()
    }()
    end_waiter.Wait()
}

func main() {
	log.Printf("NBMiner Status Reporter Initiated")
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			checkMinerStatus()
		}
	}()

	WaitForCtrlC()
	log.Printf("Termination Signal Detected.")
}