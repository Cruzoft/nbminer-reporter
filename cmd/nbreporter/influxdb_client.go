package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
    "time"

    log "github.com/sirupsen/logrus"
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

/*
	Checks the connection with the remote InfluxDB and shows the version detected
*/
func checkInfluxHealth() (error) {
	// Creates the influx client
	influxClient := influxdb2.NewClient(fmt.Sprintf("%s://%s:%v", *optInfluxProto, *optInfluxHost, *optInfluxPort), token)
	// Runs the healthCheck 
	health, err := influxClient.Health(context.Background())
	// Closes the client connection since it won't be used anymore.
	influxClient.Close()
	
	// Checking if any errors occured while running InfluxDB healthCheck
	if err != nil {
		log.Error("Couldn't check InfluxDB Health.")
		return err
	} 
	
	// Shows the detected version
	log.Infof("Detected InfluxDB version: %s", *domain.HealthCheck(*health).Version)
	return nil
}

/*
	It takes a NBMiner status object, creates an InfluxDB data point
		and writes it into a remote InfluxDB service.
*/
func writeToInflux(timestamp time.Time, statusRaw interface{}, ping int) (error) {

	measurement := "miner-device-status"
	
    // create new client with default option for server url authenticate by token
	influxClient := influxdb2.NewClient(fmt.Sprintf("%s://%s:%v", *optInfluxProto, *optInfluxHost, *optInfluxPort), token)
    // user blocking write client for writes to desired bucket

    writeAPI := influxClient.WriteAPIBlocking(*optInfluxOrg, *optInfluxBucket)

	// Creating context for the writing API
	ctx := context.Background()

	// Create the tags and fields bars to store all the datapoint values
	tags := map[string]string{"friendlyName": *optFriendlyName}
	fields := map[string]interface{}{"ping": ping}

	var err error
	
	if (ping == 1) {
		switch minerName {
		case "trex":
			status, casted := statusRaw.(*minerStatusTrex)
			if !casted {
				return fmt.Errorf("Miner status couldn't be casted into TRex Miner status struct: %v", statusRaw)
			}
			// Common Tags
			tags["user"] = status.Name
			// Common Fields
			fields["total_hashrate_raw"] = status.Hashrate
			fields["accepted_shares"] = status.AcceptedCount
			fields["start_time"] = status.Ts
			fields["version"] = status.Version

			for _, device := range status.Gpus {
				// Device Tags
				tags["id"] = strconv.Itoa(device.DeviceID)
				tags["pci_bus_id"] = strconv.Itoa(device.GpuUserID)
				tags["info"] = device.Name
		
				// Device Fields
				fields["device_accepted_shares"] = device.Shares.AcceptedCount
				fields["device_invalid_shares"] = device.Shares.InvalidCount
				fields["device_rejected_shares"] = device.Shares.RejectedCount
				fields["fan"] = device.FanSpeed
				fields["hashrate_raw"] = device.Hashrate
				fields["temperature"] = device.Temperature
				fields["power"] = 0
				//// create point using full params constructor
				//dataPoint := influxdb2.NewPoint(measurement,
				//	tags,
				//	fields,
				//	timestamp)
				//// write point immediately
				//err = writeAPI.WritePoint(ctx, dataPoint)
			}
		default:
			status, casted := statusRaw.(*minerStatusNBminer)
			if !casted {
				return fmt.Errorf("Miner status couldn't be casted into NBMiner status struct: %v", statusRaw)
			}
			// Common Tags
			tags["user"] = status.Stratum.User[strings.LastIndex(status.Stratum.User, ".")+1:] // We don't want to store any wallet addresses 
			tags["user2"] = status.Stratum.User2[strings.LastIndex(status.Stratum.User2, ".")+1:] // We don't want to store any wallet addresses 
			// Common Fields
			fields["total_hashrate"] = status.Miner.TotalHashrate
			fields["total_hashrate2"] = status.Miner.TotalHashrate2
			fields["total_hashrate_raw"] = status.Miner.TotalHashrateRaw
			fields["total_hashrate_raw2"] = status.Miner.TotalHashrateRaw2
			fields["total_power_consume"] = status.Miner.TotalPowerConsume
			fields["reboot_times"] = status.RebootTime
			fields["start_time"] = status.StartTime
			fields["accepted_shares"] = status.Stratum.AcceptedShares
			fields["accepted_shares2"] = status.Stratum.AcceptedShares2
			fields["algorithm"] = status.Stratum.Algorithm
			fields["difficulty"] = status.Stratum.Difficulty
			fields["difficulty2"] = status.Stratum.Difficulty2
			fields["dual_mine"] = status.Stratum.DualMine
			fields["invalid_shares"] = status.Stratum.InvalidShares
			fields["latency"] = status.Stratum.Latency
			fields["latency2"] = status.Stratum.Latency2
			fields["pool_hashrate_10m"] = status.Stratum.PoolHashrate10m
			fields["pool_hashrate_4h"] = status.Stratum.PoolHashrate4h
			fields["pool_hashrate_24h"] = status.Stratum.PoolHashrate24h
			fields["rejected_shares"] = status.Stratum.RejectedShares
			fields["rejected_shares2"] = status.Stratum.RejectedShares2
			fields["url"] = status.Stratum.URL
			fields["url2"] = status.Stratum.URL2
			fields["use_ssl"] = status.Stratum.UseSSL
			fields["use_ssl2"] = status.Stratum.UseSSL2
			fields["version"] = status.Version
	
			for _, device := range status.Miner.Devices {
				// Device Tags
				tags["id"] = strconv.Itoa(device.Id)
				tags["pci_bus_id"] = strconv.Itoa(device.PCIBusId)
				tags["info"] = device.Info
		
				// Device Fields
				fields["device_accepted_shares"] = device.AcceptedShares
				fields["device_accepted_shares2"] = device.AcceptedShares2
				fields["core_clock"] = device.CoreClock
				fields["core_utilization"] = device.CoreUtilization
				fields["fan"] = device.Fan
				fields["fidelity1"] = device.Fidelity1
				fields["fidelity2"] = device.Fidelity2
				fields["hashrate"] = device.Hashrate
				fields["hashrate2"] = device.Hashrate2
				fields["hashrate_raw"] = device.HashrateRaw
				fields["hashrate2_raw"] = device.Hashrate2Raw
				fields["device_invalid_shares"] = device.InvalidShares
				fields["mem_temperature"] = device.MemTemperature
				fields["mem_clock"] = device.MemClock
				fields["mem_utilization"] = device.MemUtilization
				fields["power"] = device.Power
				fields["device_rejected_shares"] = device.RejectedShares
				fields["device_rejected_shares2"] = device.RejectedShares2
				fields["temperature"] = device.Temperature
				
				//// create point using full params constructor
				//dataPoint := influxdb2.NewPoint(measurement,
				//	tags,
				//	fields,
				//	timestamp)
				//// write point immediately
				//err = writeAPI.WritePoint(ctx, dataPoint)
			}
		}
	} else {
		log.Warn("There is no Miner status data, an empty datapoint will be sent to Influx.")
	}
	// create point using full params constructor
	dataPoint := influxdb2.NewPoint(measurement,
		tags,
		fields,
		timestamp)
	// write point immediately
	err = writeAPI.WritePoint(ctx, dataPoint)

    // Ensures background processes finish
    influxClient.Close()

	// We return no error, which should be nil if everything went well
	return err
}