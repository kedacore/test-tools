package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpRequests = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of http requests.",
	},
)

func handler(w http.ResponseWriter, r *http.Request) {
	httpRequests.Inc()
	fmt.Fprint(w, "Hello")
}

func init() {
	prometheus.MustRegister(httpRequests)
}

func main() {
	port := "8080"
	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Server started, listening on port %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
