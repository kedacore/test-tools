
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpRequestsTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of http requests.",
	},
)

func handler(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.Inc()
	msg := "Received a request"
	fmt.Fprint(w, msg)
	fmt.Println(msg)
}

func main() {
	port := "8080"

	prometheus.MustRegister(httpRequestsTotal)

	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Server started on port %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}