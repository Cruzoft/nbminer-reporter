package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/status", simulateNBMinerHandler)
	http.HandleFunc("/summary", simulateTRexMinerHandler)
    log.Println("[INFO] NBMiner Status Rest API Simulator")
    log.Println("[INFO] Listening for requests at:")
    log.Println("[INFO] NBMiner Status http://localhost:8000/api/v1/status")
    log.Println("[INFO] TRex Miner Status http://localhost:8000/summary")
	log.Fatal(http.ListenAndServe(":8000", nil))
}