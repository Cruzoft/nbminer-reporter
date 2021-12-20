package main

import (
	"os"
	"os/signal"
	"sync"
	"time"
	"fmt"

    log "github.com/sirupsen/logrus"
	getopt "github.com/pborman/getopt"
	cache "github.com/hashicorp/golang-lru"
)

var hostname, _ = os.Hostname()
var optFriendlyName = getopt.StringLong("name", 'n', hostname, "A friendly name for miner. \nDefault: hostname", "string")
var optNBMinerHost = getopt.StringLong("nbhost", 's', "localhost", "NBMiner API Host. \nDefault: localhost", "string")
var optNBMinerPort = getopt.IntLong("nbport", 'r', 8000, "NBMiner API Port. \nDefault: 8000", "strinumberng")
var optInfluxProto = getopt.StringLong("iproto", 'l', "http", "InfluxDB Protocol. \nDefault: http", "string")
var optInfluxHost = getopt.StringLong("ihost", 'h', "localhost", "InfluxDB Host. \nDefault: localhost", "string")
var optInfluxPort = getopt.IntLong("iport", 'p', 8086, "InfluxDB Port. \nDefault: 8086", "number")
var optInfluxToken = getopt.StringLong("token", 't', "", "InfluxDB Access Token.", "string")
var optInfluxUser = getopt.StringLong("username", 'u', "", "InfluxDB Username (For v1.8.x).", "string")
var optInfluxPass = getopt.StringLong("password", 'w', "", "InfluxDB Password (For v1.8.x).", "string")
var optInfluxOrg = getopt.StringLong("org", 'o', "miner-org", "InfluxDB Organization. \nDefault: miner-org", "string")
var optInfluxBucket = getopt.StringLong("bucket", 'b', "miner", "InfluxDB Bucket. \nDefault: miner", "string")
var optCheckFrequency = getopt.IntLong("freq", 'f', 60, "Status check frequency in seconds.\nDefault: 60", "number")
var optCheckFrequencyRound = getopt.IntLong("round", 'd', 1, "Round up the status timestamp seconds.\nDefault: 1", "number")
var optCache = getopt.IntLong("cache", 'c', 60, "Cache size. Set to 0 to disable Cache.\nDefault: 60", "number")
var optVerbose = getopt.Bool('v', "Run in Verbose mode. \nDefault: false", "string")
var optHelp = getopt.BoolLong("help", 0, "Show usage options.")

var token = ""

var localCache *cache.Cache

type cacheItem struct {
	Status minerStatus
	Ping int
}

/*
	A siries of setup actions to run just before the main login gets executed
*/
func init() {
	// Parsing the option flags
	getopt.Parse()

    if *optHelp {
        getopt.Usage()
        os.Exit(0)
    }

	// Checking if shold run in Verbose mode
	if (*optVerbose) {
		log.SetLevel(log.DebugLevel)
        log.Warn("Log level set to DEBUG")
	}

	// Preparing the token for authentication
	token = *optInfluxToken
	if *optInfluxUser != "" {
		log.Debug("Using username as password to authenticate on InfluxDB.")
		token = fmt.Sprintf("%s:%s",*optInfluxUser, *optInfluxPass)
	}

	// Runing a the health and version check before starting.
	healthError := checkInfluxHealth()
	if healthError != nil {
		log.Errorf("Health Error: %s", healthError.Error())
	}
	if *optCache > 0 {
		localCache, _ = cache.New(*optCache)
	}
}
/*
	It does a GET request to NBMiner status endpoint, then parses the response body Json to a Struct object
	and finally it adds it to the local cache, or, if cache is dissabled, it sends the data to InfluxDB

	If anything goes wrong, it's raise an error on the console output, but the process won't stop.
	This is meant to be like this so it can overcome an internet connection issue, or a miner reboot.
*/
func checkMinerStatus() {
	log.Printf("Getting status from miner.")
	// Prepares the ping var
	ping := 1 // If Miner status request succeeds, the ping value will be 1
	// Gets the Miner status data from the endpoint.
	log.Debug("Running GET request to miner status endpoint")
	statusData, err := requestGet(fmt.Sprintf("http://%s:%v/api/v1/status", *optNBMinerHost, *optNBMinerPort))
	if err != nil {
		log.Error("Something occurred while trying to get status from miner.")
		log.Error(err)
		ping = 0 // If Miner status request fails, the ping value is set to 0
	}
	
	log.Debug("Parsing miner status json")
	// Parses the data into a struct object
	status, err := parseStatus(statusData)
	if err != nil {
		log.Error("Something occurred while trying to parse status from miner.")
		log.Error(err)
		ping = 0 // If Miner status parse fails, the ping value is set to 0
	}

	// Creates the timestamp of the datapoint
	timestamp := time.Now().Round(time.Duration(*optCheckFrequencyRound) * time.Second)

	// If cache is enable
	if *optCache > 0 {
		// Storing data on Local Cache
		localCache.Add(time.Time(timestamp), cacheItem{status, ping})
	} else {// If cache is dissable
		// Sends the data to InfluxDB
		log.Printf("Sending data point to Influx.")
		error := writeToInflux(timestamp, status, ping)
		if error != nil {
			log.Errorf("Write error: %s\n", error.Error())
			return
		}
	}

}

/*
	Checks the local cache and if it finds any cache items it'll process them (send its datapoints to InfluxDB)
	till the cache is empty. 
	
	It uses the cache as a FIFO list, so it sends all the datapoint in the order in which they were created.

	If anything goes wrong, it's raise an error on the console output, but the process won't stop.
	This is meant to be like this so it can overcome an internet connection issue, or a miner reboot.
*/
func sendDataToInflux() {
	log.Printf("Sending data point to Influx.")

	// Starts a for loop to process all the items in the cache
	for {
		// Gets the current cache size
		cacheSize := localCache.Len()
		log.Debugf("Items left in cache: %d", cacheSize)
		if (cacheSize < 1) { // If it's empty, returns
			log.Debug("Local cache is empty.")
			break
		}
		// If there are items stored in cache, we get the oldest first
		key, value, ok := localCache.GetOldest()
		timestamp, _ := key.(time.Time) // Data point time
		datapoint, _ := value.(cacheItem) // Data point and ping
		
		if ok {
			log.Debugf("Got oldest data point queue: %v", timestamp)
			
			// Sends the data to InfluxDB
			log.Debug("Writing on InfluxDB")
			error := writeToInflux(timestamp,datapoint.Status, datapoint.Ping)
			if error != nil { // Checking for writing errors
				log.Error("Couldn't write data on InfluxDB.")
				log.Errorf("Write error: %s\n", error.Error())
				return
			}

			// After sending the data to Influx it removes is from the cache
			log.Debug("Data sent")
			removed := localCache.Remove(timestamp)
			if !removed {
				log.Errorf("Something odd occurred, couldn't remove %v from local cache.", timestamp)
				return
			}
			log.Debug("Datapoint removed from local cache")
		} 
	}
}

/*
	A function that will listen for a termination signal to be able to 
	gracefully close all processes.
*/
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

/*
	Cleans up the processes and cache before stoping
*/
func cleanUp(){
	if *optCache > 0 {
		log.Debug("Purging the local cache.")
		localCache.Purge()
	}
}

/*
	Creates a ticker to run a given function periodically on its own go proc.
*/
func scheduleFunction(function func()) *time.Ticker {
	// Creating the ticker
	ticker := time.NewTicker(time.Second * time.Duration(*optCheckFrequency))
	
	// Running initial Check
	function()

	// Creating go proc to run the function loop
    go func() {
        for range ticker.C {
            function()
        }
    }()

    return ticker
}

func main() {
    log.Printf("NBMiner Status Reporter Initiated")
    log.Printf("Using Friendly Name: %s", *optFriendlyName)

	// Starting running loops
	scheduleFunction(checkMinerStatus)
	if *optCache > 0 {
		scheduleFunction(sendDataToInflux)
	}

	waitTerminationSignal()
	log.Printf("Termination Signal Detected.")
	cleanUp()
}