package main

// Import resty into your code and refer it as `resty`.
import (
	"net/http"
    "io/ioutil"

    log "github.com/sirupsen/logrus"
)

/*

*/
func requestGet(host string) (result []byte)/*(result string, statusCode int, status string)*/ {
    response, err := http.Get(host)

    if err != nil {
        log.Println(err)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatalln(err)
    }
    
	return responseData
}
