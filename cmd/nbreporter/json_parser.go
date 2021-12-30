package main

import (
	"fmt"
	"encoding/json"
)

/*
	An struct that represents the JSon object returned by NBMiner Status endpoint
	Based on NBMiner v39.7
*/
type minerStatusNBminer struct {
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
	An struct that represents the JSon object returned by TRex Summary endpoint
	Based on T-Rex v0.24.8
*/
type minerStatusTrex struct {
	AcceptedCount int `json:"accepted_count"`
	ActivePool    struct {
		Difficulty int    `json:"difficulty"`
		Ping       int    `json:"ping"`
		Retries    int    `json:"retries"`
		URL        string `json:"url"`
		User       string `json:"user"`
	} `json:"active_pool"`
	Algorithm   string  `json:"algorithm"`
	API         string  `json:"api"`
	Cuda        string  `json:"cuda"`
	Description string  `json:"description"`
	Difficulty  float64 `json:"difficulty"`
	GpuTotal    int     `json:"gpu_total"`
	Gpus        []struct {
		DeviceID              int     `json:"device_id"`
		FanSpeed              int     `json:"fan_speed"`
		GpuUserID             int     `json:"gpu_user_id"`
		Hashrate              float64 `json:"hashrate"`
		HashrateDay           int     `json:"hashrate_day"`
		HashrateHour          int     `json:"hashrate_hour"`
		HashrateMinute        int     `json:"hashrate_minute"`
		Intensity             float64 `json:"intensity"`
		Name                  string  `json:"name"`
		Temperature           int     `json:"temperature"`
		Vendor                string  `json:"vendor"`
		Disabled              bool    `json:"disabled"`
		DisabledAtTemperature int     `json:"disabled_at_temperature"`
		Shares                struct {
			AcceptedCount int `json:"accepted_count"`
			InvalidCount  int `json:"invalid_count"`
			RejectedCount int `json:"rejected_count"`
			SolvedCount   int `json:"solved_count"`
		} `json:"shares"`
	} `json:"gpus"`
	Hashrate       float64 `json:"hashrate"`
	HashrateDay    int     `json:"hashrate_day"`
	HashrateHour   int     `json:"hashrate_hour"`
	HashrateMinute int     `json:"hashrate_minute"`
	Name           string  `json:"name"`
	Os             string  `json:"os"`
	RejectedCount  int     `json:"rejected_count"`
	SolvedCount    int     `json:"solved_count"`
	Ts             int     `json:"ts"`
	Uptime         int     `json:"uptime"`
	Version        string  `json:"version"`
	Updates        struct {
		URL            string `json:"url"`
		Md5Sum         string `json:"md5sum"`
		Version        string `json:"version"`
		NotesFull      string `json:"notes_full"`
		DownloadStatus struct {
			DownloadedBytes  int     `json:"downloaded_bytes"`
			TotalBytes       int     `json:"total_bytes"`
			LastError        string  `json:"last_error"`
			TimeElapsedSec   float64 `json:"time_elapsed_sec"`
			UpdateInProgress bool    `json:"update_in_progress"`
			UpdateState      string  `json:"update_state"`
			URL              string  `json:"url"`
		} `json:"download_status"`
	} `json:"updates"`
}

/*
	Parses the json string returned by NBMiner status endpoint into a struct object
*/
func parseStatus (statusData []byte) (interface{}, error) {

	var status interface{}
	
	switch minerName {
	case "trex":
		status = &minerStatusTrex{}
	default:
		status = &minerStatusNBminer{}
	}
	// Tries to unmarshal the Json data into a struct object
	//   if an error occurs, it returns nil and the error message
	if err := json.Unmarshal(statusData, &status); err != nil {
		return status, fmt.Errorf("Couldn't parse the status json.\n Received Json: %s \n Previous Error: %v",string(statusData), err)
    }

	// Returns the status struct object
	return status, nil
}