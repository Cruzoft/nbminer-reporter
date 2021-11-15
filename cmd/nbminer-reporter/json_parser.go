package main

import (
	"fmt"
	"encoding/json"
)

/*
	An struct that represents the JSon object returned by NBMiner Status endpoint
	Based on NBMiner v39.7
*/
type minerStatus struct {
	Miner struct {
		Devices []struct {
			AcceptedShares int `json:"accepted_shares"`
			AcceptedShares2 int `json:"accepted_shares2"`
			CoreClock int `json:"core_clock"`
			CoreUtilization int `json:"core_utilization"`
			Fan int `json:"fan"`
			Fidelity1 float32 `json:"fidelity1"`
			Fidelity2 int `json:"fidelity2"`
			Hashrate string `json:"hashrate"`
			Hashrate2 string `json:"hashrate2"`
			Hashrate2Raw float32 `json:"hashrate2_raw"`
			HashrateRaw float32 `json:"hashrate_raw"`
			Id int `json:"id"`
			Info string `json:"info"`
			InvalidShares int `json:"invalid_shares"`
			MemTemperature int `json:"memTemperature"`
			MemClock int `json:"mem_clock"`
			MemUtilization int `json:"mem_utilization"`
			PCIBusId int `json:"pci_bus_id"`
			Power int `json:"power"`
			RejectedShares int `json:"rejected_shares"`
			RejectedShares2 int `json:"rejected_shares2"`
			Temperature int `json:"temperature"`
		} `json:"devices"`
		TotalHashrate  string `json:"total_hashrate"`
		TotalHashrate2  string `json:"total_hashrate2"`
		TotalHashrateRaw  float32 `json:"total_hashrate_raw"`
		TotalHashrateRaw2  float32 `json:"total_hashrate_raw2"`
		TotalPowerConsume int `json:"total_power_consume"`
	} `json:"miner"`
	RebootTime int `json:"reboot_time"`
	StartTime int `json:"start_time"`
	Stratum struct {
		AcceptedShares int `json:"accepted_shares"`
		AcceptedShares2 int `json:"accepted_shares2"`
		Algorithm string `json:"algorithm"`
		Difficulty string `json:"difficulty"`
		Difficulty2 string `json:"difficulty2"`
		DualMine bool `json:"dual_mine"`
		InvalidShares int `json:"invalid_shares"`
		Latency int `json:"latency"`
		Latency2 int `json:"latency2"`
		PoolHashrate10m string `json:"pool_hashrate_10m"`
		PoolHashrate4h string `json:"pool_hashrate_4h"`
		PoolHashrate24h string `json:"pool_hashrate_24h"`
		RejectedShares int `json:"rejected_shares"`
		RejectedShares2 int `json:"rejected_shares2"`
		URL string `json:"url"`
		URL2 string `json:"url2"`
		UseSSL bool `json:"use_ssl"`
		UseSSL2 bool `json:"use_ssl2"`
		User string `json:"user"`
		User2 string `json:"user2"`
	} `json:"stratum"`
	Version string  `json:"version"`
}

/*
	Parses the json string returned by NBMiner status endpoint into a struct object
*/
func parseStatus(statusData []byte) (minerStatus, error) {

	var status minerStatus

	// Tries to unmarshal the Json data into a struct object
	//   if an error occurs, it returns nil and the error message
	if err := json.Unmarshal(statusData, &status); err != nil {
		return status, fmt.Errorf("Couldn't parse the status json.\n Received Json: %s \n Previous Error: %v",string(statusData), err)
    }

	// Returns the status struct object
	return status, nil
}