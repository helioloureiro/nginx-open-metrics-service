package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Version = "development"

	activeConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_connections",
		Help: "The number of active connections",
	})

	serverAccepts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "server_accepts_total",
		Help: "The total number of server accepted connections",
	})

	serverHandled = promauto.NewCounter(prometheus.CounterOpts{
		Name: "server_handled_total",
		Help: "The total number of server handled connections",
	})

	serverRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "server_requests_total",
		Help: "The total number of server requests",
	})

	connectionsReading = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "reading_connections",
		Help: "The number of active reading connections",
	})

	connectionsWriting = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "writing_connections",
		Help: "The number of active writing connections",
	})

	connectionsWaiting = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "waiting_connections",
		Help: "The number of waiting connections",
	})
)

func main() {
	service := flag.String("service", "http://localhost:8080/api", "the nginx api service, like http://localhost:8080/api")
	port := flag.Int("port", 9090, "default port to listen the service")
	printVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *printVersion {
		fmt.Println("nginx-open-metrics-service")
		fmt.Println("version:", Version)
		os.Exit(0)
	}

	if *service == "" {
		log.Fatal("missing service service")
	}

	fmt.Println("fetching data from:", *service)
	fmt.Println("starting service at port:", *port)
	/*
	 * Active connections: 39
	 * server accepts handled requests
	 * 286479 286479 1417563
	 * Reading: 0 Writing: 64 Waiting: 10
	 */
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		activeConnections,
		serverAccepts,
		serverHandled,
		serverRequests,
		connectionsReading,
		connectionsWriting,
		connectionsWaiting,
	)

	http.Handle(
		"/metrics", promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			}),
	)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
