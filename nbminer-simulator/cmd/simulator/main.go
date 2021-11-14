package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"math/rand"
)

func main() {
	
	statusHandler := func(w http.ResponseWriter, req *http.Request) {
		log.Println("[INFO] Status report")
		io.WriteString(w, fmt.Sprintf(`{
	"miner": {
		"devices": [
			{
				"accepted_shares": 2,
				"accepted_shares2": 0,
				"core_clock": 1620,
				"core_utilization": 100,
				"fan": 47,
				"fidelity1": 5.859799716605649,
				"fidelity2": 0,
				"hashrate": "217.1 M",
				"hashrate2": "36.19 M",
				"hashrate_raw": %v,
				"hashrate2_raw": 36190716.266428046,
				"id": 0,
				"info": "GeForce RTX 2070",
				"mem_clock": 6801,
				"mem_utilization": 86,
				"pci_bus_id": 1,
				"power": 188,
				"rejected_shares": 0,
				"rejected_shares2": 0,
				"temperature": 72
			},
			{
				"accepted_shares": 0,
				"accepted_shares2": 0,
				"core_clock": 1607,
				"core_utilization": 100,
				"fan": 0,
				"fidelity1": 0,
				"fidelity2": 0,
				"hashrate": "168.5 M",
				"hashrate2": "42.11 M",
				"hashrate_raw": %v,
				"hashrate2_raw": 42113955.19774488,
				"id": 1,
				"info": "P102-100",
				"mem_clock": 5508,
				"mem_utilization": 100,
				"pci_bus_id": 4,
				"power": 232,
				"rejected_shares": 0,
				"rejected_shares2": 0,
				"temperature": 57
			}
		],
		"total_hashrate": "708 M",
		"total_hashrate2": "164.4 M",
		"total_hashrate2_raw": 164395439.13815895,
		"total_hashrate_raw": 708044466.8349969,
		"total_power_consume": 839
	},
	"reboot_times": 0,
	"start_time": 1586944619,
	"stratum": {
		"accepted_shares": 2,
		"accepted_shares2": 0,
		"algorithm": "hns_ethash",
		"difficulty": "8.59 G",
		"difficulty2": "8.59 G",
		"dual_mine": true,
		"latency": 221,
		"latency2": 0,
		"rejected_shares": 0,
		"rejected_shares2": 0,
		"url": "handshake.hk.nicehash.com:3384",
		"url2": "daggerhashimoto.hk.nicehash.com:3353",
		"use_ssl": false,
		"use_ssl2": false,
		"user": "3QHNv52ahdCyeYTGVYDPGjRzMpkknjjfAf.test",
		"user2": "3QHNv52ahdCyeYTGVYDPGjRzMpkknjjfAf.test"
	},
	"version": "30.0"
}
`,
(rand.Float32() * 5) + 180,
(rand.Float32() * 5) + 180))
	}

	http.HandleFunc("/api/v1/status", statusHandler)
    log.Println("[INFO] NBMiner Status Rest API Simulator")
    log.Println("[INFO] Listening for requests at http://localhost:8000/api/v1/status")
	log.Fatal(http.ListenAndServe(":8000", nil))
}