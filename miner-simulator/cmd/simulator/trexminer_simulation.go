package main

import (
	"fmt"
	"log"
	"io"
	"math/rand"
	"net/http"
)

func simulateTRexMinerHandler (w http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] TRex Miner Status report")
	hashrate1 := (rand.Float32() * 5000000) + 30000000
	hashrate2 := (rand.Float32() * 5000000) + 30000000
	hashrateTotal := hashrate1 + hashrate2
	io.WriteString(w, fmt.Sprintf(`{
		"accepted_count": 6,
		"active_pool":
		{
		  "difficulty": 5,
		  "ping": 97,
		  "retries": 0,
		  "url": "stratum+tcp://...",
		  "user": "..."
		},
		"algorithm": "x16r",
		"api": "1.2",
		"cuda": "9.10",
		"description": "T-Rex NVIDIA GPU miner",
		"difficulty": 31968.245093004043,
		"gpu_total": 1,
		"gpus": [{
		  "device_id": 0,                        
		  "fan_speed": 66,                       
		  "gpu_user_id": 0,                        
		  "hashrate": %f,                   
		  "hashrate_day": 5023728,    
		  "hashrate_hour": 0,          
		  "hashrate_minute": 4671930,    
		  "intensity": 21.5,        
		  "name": "GeForce GTX 1050",
		  "temperature": 80,            
		  "vendor": "Gigabyte", 
		  "disabled":true,                       
		  "disabled_at_temperature": 77,
		  "shares": {
			  "accepted_count": 3,
			  "invalid_count": 0,
			  "rejected_count": 0,
			  "solved_count": 0
		  	}
		  },
		  {
			"device_id": 1,                        
			"fan_speed": 66,                       
			"gpu_user_id": 1,                        
			"hashrate": %f,                   
			"hashrate_day": 5023728,    
			"hashrate_hour": 0,          
			"hashrate_minute": 4671930,    
			"intensity": 21.5,        
			"name": "GeForce GTX 1050",
			"temperature": 75,            
			"vendor": "Gigabyte", 
			"disabled":false,                       
			"disabled_at_temperature": 77,
			"shares": {
				"accepted_count": 6,
				"invalid_count": 0,
				"rejected_count": 0,
				"solved_count": 0
				}
			}
		],
		"hashrate": %f,                       
		"hashrate_day": 5023728,                   
		"hashrate_hour": 0,                        
		"hashrate_minute": 4671930,                
		"name": "t-rex",
		"os": "linux",
		"rejected_count": 0,                       
		"solved_count": 0, 
		"ts": 1537095257,    
		"uptime": 108,     
		"version": "0.6.5",
		"updates":{
		  "url": "https://fileurl",  
		  "md5sum": "md5...",  
		  "version": "0.8.0",  
		  "notes_full": "full update info",  
		  "download_status": {
			"downloaded_bytes": 1775165,
			"total_bytes": 5245345,
			"last_error":"",
			"time_elapsed_sec": 2.887111,
			"update_in_progress": true,
			"update_state": "downloading",
			"url": "https://fileurl"
		  }
		}
	  }
    	`,
    	hashrate1,
    	hashrate2,
    	hashrateTotal))
    }