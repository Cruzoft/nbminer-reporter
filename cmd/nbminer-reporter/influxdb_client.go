package main

import (
	"strconv"
	"strings"
    "time"

    "github.com/influxdata/influxdb-client-go/v2"
    log "github.com/sirupsen/logrus"
)

func writeToInflux(status minerStatus) (error) {

	measurement := "miner-device-status"
	timestamp := time.Now()
	
    // create new client with default option for server url authenticate by token
    client := influxdb2.NewClient("http://host.docker.internal:8086", "shhh-secret-token")
    // user blocking write client for writes to desired bucket
    writeAPI := client.WriteAPI("nlsrig", "mainr")
	// Get errors channel
    errorsCh := writeAPI.Errors()
    // Create go proc for reading and logging errors
    go func() {
        for err := range errorsCh {
            log.Printf("write error: %s\n", err.Error())
        }
    }()

	for _, device := range status.Miner.Devices {
		tags := map[string]string{
			// Common Tags
			"friendlyName": "rig03",
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
			"latency": status.Stratum.Latency,
			"latency2": status.Stratum.Latency2,
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
			"mem_clock": device.MemClock,
			"mem_utilization": device.MemUtilization,
			"power": device.Power,
			"device_rejected_shares": device.RejectedShares,
			"device_rejected_shares2": device.RejectedShares2,
			"temperature": device.Temperature,
		}
		// create point using full params constructor
		p := influxdb2.NewPoint(measurement,
			tags,
			fields,
			timestamp)
		// write point immediately
		writeAPI.WritePoint(p)
	}

    // Force all unwritten data to be sent
    writeAPI.Flush()
    // Ensures background processes finish
    client.Close()

	return nil
}