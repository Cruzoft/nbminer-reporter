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

func checkInfluxHealth() (error) {
	influxClient := influxdb2.NewClient(fmt.Sprintf("%s://%s:%v", *optInfluxProto, *optInfluxHost, *optInfluxPort), token)
	health, err := influxClient.Health(context.Background())
	
	if err != nil {
		log.Error("Couldn't check InfluxDB Health.")
		return err
	} 
	
	log.Infof("InfluxDB Version: %s", *domain.HealthCheck(*health).Version)
	return nil
}

/*
	It takes a NBMiner status object, creates an InfluxDB data point
		and writes it into a remote InfluxDB service.
*/
func writeToInflux(status minerStatus) (error) {

	measurement := "   "
	timestamp := time.Now().Round(time.Duration(*optCheckFrequencyRound) * time.Second)
	
    // create new client with default option for server url authenticate by token
	influxClient := influxdb2.NewClient(fmt.Sprintf("%s://%s:%v", *optInfluxProto, *optInfluxHost, *optInfluxPort), token)
    // user blocking write client for writes to desired bucket
    writeAPI := influxClient.WriteAPI(*optInfluxOrg, *optInfluxBucket)
	// Get errors channel
    errorsCh := writeAPI.Errors()
    // Create go proc for reading and logging errors
    go func() {
        for err := range errorsCh {
            log.Error("Something when wrong while trying to send data to InfluxDB.")
            log.Errorf("Write error: %s\n", err.Error())
        }
    }()

	for _, device := range status.Miner.Devices {
		tags := map[string]string{
			// Common Tags
			"friendlyName": *optFriendlyName,
			"user": status.Stratum.User[strings.LastIndex(status.Stratum.User, ".")+1:], // We don't want to store any wallet addresses 
			"user2": status.Stratum.User2[strings.LastIndex(status.Stratum.User2, ".")+1:], // We don't want to store any wallet addresses 
			// Device Tags
			"id": strconv.Itoa(device.Id),
			"pci_bus_id": strconv.Itoa(device.PCIBusId),
			"info": device.Info,
		}

		fields := map[string]interface{}{
			// Common Fields
			"total_hashrate": status.Miner.TotalHashrate,
			"total_hashrate2": status.Miner.TotalHashrate2,
			"total_hashrate_raw": status.Miner.TotalHashrateRaw,
			"total_hashrate_raw2": status.Miner.TotalHashrateRaw2,
			"total_power_consume": status.Miner.TotalPowerConsume,
			"reboot_times": status.RebootTime,
			"start_time": status.StartTime,
			"accepted_shares": status.Stratum.AcceptedShares,
			"accepted_shares2": status.Stratum.AcceptedShares2,
			"algorithm": status.Stratum.Algorithm,
			"difficulty": status.Stratum.Difficulty,
			"difficulty2": status.Stratum.Difficulty2,
			"dual_mine": status.Stratum.DualMine,
			"invalid_shares": status.Stratum.InvalidShares,
			"latency": status.Stratum.Latency,
			"latency2": status.Stratum.Latency2,
			"pool_hashrate_10m": status.Stratum.PoolHashrate10m,
			"pool_hashrate_4h": status.Stratum.PoolHashrate4h,
			"pool_hashrate_24h": status.Stratum.PoolHashrate24h,
			"rejected_shares": status.Stratum.RejectedShares,
			"rejected_shares2": status.Stratum.RejectedShares2,
			"url": status.Stratum.URL,
			"url2": status.Stratum.URL2,
			"use_ssl": status.Stratum.UseSSL,
			"use_ssl2": status.Stratum.UseSSL2,
			"version": status.Version,
			// Device Fields
			"device_accepted_shares": device.AcceptedShares,
			"device_accepted_shares2": device.AcceptedShares2,
			"core_clock": device.CoreClock,
			"core_utilization": device.CoreUtilization,
			"fan": device.Fan,
			"fidelity1": device.Fidelity1,
			"fidelity2": device.Fidelity2,
			"hashrate": device.Hashrate,
			"hashrate2": device.Hashrate2,
			"hashrate_raw": device.HashrateRaw,
			"hashrate2_raw": device.Hashrate2Raw,
			"device_invalid_shares": device.InvalidShares,
			"mem_temperature": device.MemTemperature,
			"mem_clock": device.MemClock,
			"mem_utilization": device.MemUtilization,
			"power": device.Power,
			"device_rejected_shares": device.RejectedShares,
			"device_rejected_shares2": device.RejectedShares2,
			"temperature": device.Temperature,
		}
		// create point using full params constructor
		dataPoint := influxdb2.NewPoint(measurement,
			tags,
			fields,
			timestamp)
		// write point immediately
		writeAPI.WritePoint(dataPoint)
	}

    // Force all unwritten data to be sent
    writeAPI.Flush()
    // Ensures background processes finish
    influxClient.Close()

	// Since everything when will, we return no error
	return nil
}