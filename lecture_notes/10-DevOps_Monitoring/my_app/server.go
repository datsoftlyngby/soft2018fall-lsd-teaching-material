package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	cpuLoad = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_load_percent",
		Help: "Current load of the CPU in percent.",
	})
	totalRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_responses_total",
		Help: "The count of http responses sent.",
	})
	responseMetric = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "request_duration_milliseconds",
			Help:    "Request latency distribution",
			Buckets: prometheus.ExponentialBuckets(10.0, 1.13, 40),
		})

	page = `<!DOCTYPE html>
<html>
<title>Hej!</title>
<body>
	<h1>Welcome</h1>
	<p>To this server.</p>
</body>
</html>`
)

func init() {
	prometheus.MustRegister(cpuLoad)
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(responseMetric)
}

func main() {
	// Call it for example as env IP=0.0.0.0 PORT=8080 server
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	if ip == "" {
		ip = "127.0.0.1"
	}
	if port == "" {
		port = "8080"
	}
	serverAddress := fmt.Sprintf("%s:%s", ip, port)
	log.SetPrefix("INFO: ")
	log.Println("Listening to http://" + serverAddress)
	http.HandleFunc("/", landingPage)
	http.Handle("/metrics", prometheus.Handler())
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

func logCPUUsage() {
	// This code will only run on a system with top, grep, sed, and awk installed
	cmd := `top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1}'`
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.SetPrefix("ERROR: ")
		log.Println(err)
	}
	load, err := strconv.ParseFloat(strings.Trim(string(out), "\n"), 64)
	if err != nil {
		log.SetPrefix("ERROR: ")
		log.Println(err)
	}
	cpuLoad.Set(load)
}

func landingPage(res http.ResponseWriter, req *http.Request) {
	totalRequests.Inc()
	start := time.Now()

	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	// for simulation purposes, make the server return the page after a random
	// break
	myPause := rand.Intn(3) // 0 <= n <= 2
	log.SetPrefix("INFO: ")
	log.Println(fmt.Sprintf("Serving request after %ds", myPause))
	time.Sleep(time.Duration(myPause) * time.Second)
	fmt.Fprintf(res, page)
	msElapsed := time.Since(start) / time.Millisecond
	responseMetric.Observe(float64(msElapsed))
	logCPUUsage()
}

// https://github.com/the-hobbes/mockingProduction/blob/f0019173f00ad6c631f21269bc7409001bc5e536/http_server/http.go
// https://stackoverflow.com/questions/37611754/how-to-push-metrics-to-prometheus-using-client-golang
