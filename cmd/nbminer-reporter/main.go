package main

import (
	"os"
	"os/signal"
	"sync"
	"time"
	"fmt"

    log "github.com/sirupsen/logrus"
	getopt "github.com/pborman/getopt"
)

var optFriendlyName = getopt.StringLong("name", 'n', "a-miner", "A friendly name for miner.")
var optNBMinerHost = getopt.StringLong("nbhost", 's', "localhost", "NBMiner API Host. Default: localhost")
var optNBMinerPort = getopt.IntLong("nbport", 'r', 8000, "NBMiner API Port. Default: 8000")
var optInfluxHost = getopt.StringLong("ihost", 'h', "localhost", "InfluxDB Host. Default: localhost")
var optInfluxPort = getopt.IntLong("iport", 'p', 8086, "InfluxDB Port. Default: 8086")
var optInfluxToken = getopt.StringLong("itoken", 't', "", "InfluxDB Access Token.")
var optInfluxOrg = getopt.StringLong("iorg", 'o', "miner-org", "InfluxDB Organization. Default: miner-org")
var optInfluxBucket = getopt.StringLong("ibucket", 'b', "miner", "InfluxDB Bucket. Default: miner")
var optCheckFrequency = getopt.IntLong("freq", 'f', 300, "Status check frequency in seconds. Default: 300")
var optHelp = getopt.BoolLong("help", 0, "Help")

func checkMinerStatus() {
	//statusData := requestGet("http://host.docker.internal:8000/api/v1/status")
	statusData := requestGet(fmt.Sprintf("http://%s:%v/api/v1/status", *optNBMinerHost, *optNBMinerPort))
	status, _ := parseStatus(statusData)

	for _, device := range status.Miner.Devices {
		log.Printf("Found Device %s in pcie %v.", device.Info, device.PCIBusId)
	}
	log.Printf("Total Hashrate: %s.", status.Miner.TotalHashrate)
	
	writeToInflux(status)
	log.Printf("Influx line sent")
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
    getopt.Parse()

    if *optHelp {
        getopt.Usage()
        os.Exit(0)
    }

    log.Printf("NBMiner Status Reporter Initiated")
    log.Printf("Using Friendly Name: %s", *optFriendlyName)
	ticker := time.NewTicker(time.Second * time.Duration(*optCheckFrequency))
	go func() {
		for range ticker.C {
			checkMinerStatus()
		}
	}()

	WaitForCtrlC()
	log.Printf("Termination Signal Detected.")
}