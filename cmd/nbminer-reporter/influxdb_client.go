package main

import (
	"strings"
    "time"

    "github.com/influxdata/influxdb-client-go/v2"
    log "github.com/sirupsen/logrus"
)

func writeToInflux(status minerStatus) (error) {

	measurement := "miner-device-status"
	timestamp := time.Now()
		
	//Common tags
	//commonTags := make(map[string]string)
	commonTags := map[string]string{
		"user": status.Stratum.User[strings.LastIndex(status.Stratum.User, ".")+1:], // We don't want to store any wallet addresses 
		"user2": status.Stratum.User2[strings.LastIndex(status.Stratum.User2, ".")+1:], // We don't want to store any wallet addresses 
	}

	//Common fields
	//commonFields := make(map[string]interface{})
	commonFields := map[string]interface{}{
		"reboot_times": status.RebootTime,
	}

    // create new client with default option for server url authenticate by token
    client := influxdb2.NewClient("http://localhost:9999", "my-token")
    // user blocking write client for writes to desired bucket
    writeAPI := client.WriteAPI("my-org", "my-bucket")
	// Get errors channel
    errorsCh := writeAPI.Errors()
    // Create go proc for reading and logging errors
    go func() {
        for err := range errorsCh {
            log.Printf("write error: %s\n", err.Error())
        }
    }()

    // create point using full params constructor
    p := influxdb2.NewPoint(measurement,
		commonTags,
        commonFields,
        timestamp)
    // write point immediately
    writeAPI.WritePoint(p)
    // Force all unwritten data to be sent
    writeAPI.Flush()
    // Ensures background processes finish
    client.Close()

	return nil
}