package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func aSingleRequest(url string) {
	client := http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Println("Something went wrong")
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Got Status Code: " + string(resp.StatusCode))
	}
}

func main() {
	serverAddress := fmt.Sprintf("%s:%s/", os.Getenv("SERVER"), os.Getenv("PORT"))

	url := "http://" + serverAddress
	fmt.Println("Will talk to " + url)
	// Let the program run forever...
	// send one request per second
	for {
		go aSingleRequest(url)
		time.Sleep(1 * time.Second)
	}
}
