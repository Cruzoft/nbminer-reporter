package main

import (
	"fmt"
	"log"
	"io"
	"math/rand"
	"net/http"
)

func simulateNBMinerHandler (w http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] NBMiner Status report")
	hashrate1 := (rand.Float32() * 5000000) + 30000000
	hashrate2 := (rand.Float32() * 5000000) + 30000000
	hashrate3 := (rand.Float32() * 5000000) + 30000000
	hashrateTotal := hashrate1 + hashrate2 + hashrate3
	io.WriteString(w, fmt.Sprintf(`{
    			"miner": {
    				"devices": [{
    					"accepted_shares": 0,
    					"core_clock": 32,
    					"core_utilization": 100,
    					"fan": 60,
    					"hashrate": "30.99 M",
    					"hashrate2": "0.000 ",
    					"hashrate2_raw": 0,
    					"hashrate_raw": %f,
    					"id": 0,
    					"info": "CMP 30HX",
    					"invalid_shares": 0,
    					"memTemperature": -8,
    					"mem_clock": 0,
    					"mem_utilization": 85,
    					"pci_bus_id": 1,
    					"power": 78,
    					"rejected_shares": 0,
    					"temperature": 67
    				}, {
    					"accepted_shares": 1,
    					"core_clock": 32,
    					"core_utilization": 99,
    					"fan": 60,
    					"hashrate": "31.31 M",
    					"hashrate2": "0.000 ",
    					"hashrate2_raw": 0,
    					"hashrate_raw": %f,
    					"id": 1,
    					"info": "CMP 30HX",
    					"invalid_shares": 0,
    					"memTemperature": -8,
    					"mem_clock": 0,
    					"mem_utilization": 85,
    					"pci_bus_id": 2,
    					"power": 78,
    					"rejected_shares": 0,
    					"temperature": 62
    				}, {
    					"accepted_shares": 0,
    					"core_clock": 32,
    					"core_utilization": 100,
    					"fan": 60,
    					"hashrate": "31.27 M",
    					"hashrate2": "0.000 ",
    					"hashrate2_raw": 0,
    					"hashrate_raw": %f,
    					"id": 2,
    					"info": "CMP 30HX",
    					"invalid_shares": 0,
    					"memTemperature": -8,
    					"mem_clock": 0,
    					"mem_utilization": 86,
    					"pci_bus_id": 3,
    					"power": 78,
    					"rejected_shares": 0,
    					"temperature": 60
    				}],
    				"total_hashrate": "186.9 M",
    				"total_hashrate2": "0.000 ",
    				"total_hashrate2_raw": 0,
    				"total_hashrate_raw": %f,
    				"total_power_consume": 467
    			},
    			"reboot_times": 0,
    			"start_time": 1636943904,
    			"stratum": {
    				"accepted_shares": 1,
    				"algorithm": "ethash",
    				"difficulty": "8.726 G",
    				"dual_mine": false,
    				"invalid_shares": 0,
    				"latency": 186,
    				"pool_hashrate_10m": "50.15 M",
    				"pool_hashrate_24h": "50.15 M",
    				"pool_hashrate_4h": "50.15 M",
    				"rejected_shares": 0,
    				"url": "what.a-miner.com:2121",
    				"use_ssl": false,
    				"user": "hfa7erw80fh43hti5r9pgwh8594wjgirtesj498tj94wjptje.rig03"
    			},
    			"version": "39.7"
    		}
    	`,
    	hashrate1,
    	hashrate2,
    	hashrate3,
    	hashrateTotal))
    }