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

var hostname, _ = os.Hostname()
var optFriendlyName = getopt.StringLong("name", 'n', hostname, "A friendly name for miner. \nDefault: hostname", "string")
var optNBMinerHost = getopt.StringLong("nbhost", 's', "localhost", "NBMiner API Host. \nDefault: localhost", "string")
var optNBMinerPort = getopt.IntLong("nbport", 'r', 8000, "NBMiner API Port. \nDefault: 8000", "strinumberng")
var optInfluxProto = getopt.StringLong("iproto", 'l', "http", "InfluxDB Protocol. \nDefault: http", "string")
var optInfluxHost = getopt.StringLong("ihost", 'h', "localhost", "InfluxDB Host. \nDefault: localhost", "string")
var optInfluxPort = getopt.IntLong("iport", 'p', 8086, "InfluxDB Port. \nDefault: 8086", "number")
var optInfluxToken = getopt.StringLong("itoken", 't', "", "InfluxDB Access Token.", "string")
var optInfluxOrg = getopt.StringLong("iorg", 'o', "miner-org", "InfluxDB Organization. \nDefault: miner-org", "string")
var optInfluxBucket = getopt.StringLong("ibucket", 'b', "miner", "InfluxDB Bucket. \nDefault: miner", "string")
var optCheckFrequency = getopt.IntLong("freq", 'f', 60, "Status check frequency in seconds.\nDefault: 60", "number")
var optVerbose = getopt.Bool('v', "Run in Verbose mode. \nDefault: false", "string")
var optHelp = getopt.BoolLong("help", 0, "Show usage options.")


func init() {
	getopt.Parse()

    if *optHelp {
        getopt.Usage()
        os.Exit(0)
    }

	if (*optVerbose) {
		log.SetLevel(log.DebugLevel)
        log.Warn("Log level set to DEBUG")
	}
}
/*
	This is the process main logic
	It does a GET request to NBMiner status endpoint, then parses the response body Json to a Struct object
	and finally it sends its data to InfluxDB

	If anything goes wrong, it's raise an error on the console output, but the process won't stop.
	This is meant to be like this so it can overcome an internet connection issue, or a miner reboot.
*/
func checkMinerStatus() {
	log.Printf("Checking Status.")
	// Gets the Miner status data from the endpoint.
	log.Debug("Running GET request to miner status endpoint")
	statusData, err := requestGet(fmt.Sprintf("http://%s:%v/api/v1/status", *optNBMinerHost, *optNBMinerPort))
	if err != nil {
		log.Error("Something occurred while trying to get status from miner.")
		log.Error(err)
		return
	}
	
	log.Debug("Parsing miner status json")
	// Parses the data into a struct object
	status, err := parseStatus(statusData)
	if err != nil {
		log.Error("Something occurred while trying to parse status from miner.")
		log.Error(err)
		return
	}
	// Sends the data to InfluxDB
	log.Debug("Sending data to InfluxDB")
	writeToInflux(status)
	log.Debug("Data sent")
}

func waitTerminationSignal() {
    var endWaiter sync.WaitGroup
    var signalChannel chan os.Signal
    
	endWaiter.Add(1)
    signalChannel = make(chan os.Signal, 1)
    signal.Notify(signalChannel, os.Interrupt)
    
	go func() {
		<-signalChannel
        endWaiter.Done()
    }()
    
	endWaiter.Wait()
}

func main() {
    log.Printf("NBMiner Status Reporter Initiated")
    log.Printf("Using Friendly Name: %s", *optFriendlyName)

	// Starting running loop
	ticker := time.NewTicker(time.Second * time.Duration(*optCheckFrequency))
	go func() {
		for range ticker.C {
			checkMinerStatus()
		}
	}()

	waitTerminationSignal()
	log.Printf("Termination Signal Detected.")
}