package main

import (
	"net/http"
    "io/ioutil"
)

/*
    Runs a GET HTTP request and returns the response as a byte slice
*/
func requestGet(host string) ([]byte, error) {
    // Execute the HTTP request
    response, err := http.Get(host)

    // If error occurs, stop and report it
    if err != nil {
        return nil, err
    }

    // Read the response body data
    responseData, err := ioutil.ReadAll(response.Body)

    // If error occurs, stop and report it
    if err != nil {
        return nil, err
    }
    
    // Finally, if no error occurs, return the body data
	return responseData, nil
}
