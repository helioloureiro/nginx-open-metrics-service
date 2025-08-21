package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var (
	Version = "development"

	serverAccepts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nginx_server_accepts_total",
		Help: "The total number of server accepted connections",
	})

	serverHandled = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nginx_server_handled_total",
		Help: "The total number of server handled connections",
	})

	serverRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nginx_server_requests_total",
		Help: "The total number of server requests",
	})

	activeConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_active_connections",
		Help: "The number of active connections",
	})

	connectionsReading = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_reading_connections",
		Help: "The number of active reading connections",
	})

	connectionsWriting = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_writing_connections",
		Help: "The number of active writing connections",
	})

	connectionsWaiting = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_waiting_connections",
		Help: "The number of waiting connections",
	})

	service *string
)

func main() {
	service = flag.String("service", "http://localhost:8080/api", "the nginx api service, like http://localhost:8080/api")
	port := flag.Int("port", 9090, "default port to listen the service")
	printVersion := flag.Bool("version", false, "print the version and exit")
	logLevel := flag.String("loglevel", "info", "the logging level (default=info)")
	flag.Parse()

	setUpLogging(*logLevel)

	if *printVersion {
		fmt.Println("nginx-open-metrics-service")
		fmt.Println("version:", Version)
		os.Exit(0)
	}

	if *service == "" {
		logrus.Fatal("missing service information to connecton on nginx")
	}

	logrus.Info("🚚 fetching data from:", *service)
	logrus.Info("🎬 starting service at port:", *port)
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
	// start with data updated
	fetchDataFromNginx()

	go dataUpdater()
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func setUpLogging(logLevel string) {
	logrus.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02T15:04:05",
		LogFormat:       "[%time%] (%lvl%): %msg%\n",
	})

	if logLevel == "info" {
		// default value, so not changes needed
		return
	}

	switch logLevel {
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.Panic("Invalid config.json. Log level must be 'panic' 'fatal' 'error' 'warn' 'info' or 'debug'")
	}

	logrus.Warn(fmt.Sprintf("logging set to %s", logLevel))
}

func dataUpdater() {
	// updated every 15s
	ticker := time.NewTicker(15 * time.Second)
	logrus.Debug("starting time ⌛ ticker - every 15s")
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			logrus.Debug("⌛ ticker after 15s")
			fetchDataFromNginx()
		}
	}
}

func fetchDataFromNginx() {
	logrus.Debug("🚚 fetching data from:", *service)
	resp, err := http.Get(*service)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("failed to fetch data from: %s", *service))
	}
	logrus.Debug(fmt.Sprintf("status_code=%d", resp.StatusCode))
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("failed to fetch data from: %s", *service))
	}
	logrus.Debug(fmt.Sprintf("body: %s", body))
	ac, sa, sh, sr, cr, cw, cwa := parseDataFromNginx(body)
	activeConnections.Set(float64(ac))
	serverAccepts.Add(getDiffValue(serverAccepts, sa))
	serverHandled.Add(getDiffValue(serverHandled, sh))
	serverRequests.Add(getDiffValue(serverRequests, sr))
	connectionsReading.Set(float64(cr))
	connectionsWriting.Set(float64(cw))
	connectionsWaiting.Set(float64(cwa))
}

func parseDataFromNginx(body []byte) (int, int, int, int, int, int, int) {
	bodyStr := string(body)
	lines := strings.Split(bodyStr, "\n")
	logrus.Debug(fmt.Sprintf("lines: %v", lines))

	var tmp string
	tmp = lines[0]
	logrus.Debug(fmt.Sprintf("tmp: %v", tmp))
	ac := convertToInt(strings.Split(tmp, ":")[1])

	tmp = lines[2]
	tmp = sed(tmp, "^ ", "")
	values := strings.Split(tmp, " ")
	sa := convertToInt(values[0])
	sh := convertToInt(values[1])
	sr := convertToInt(values[2])

	tmp = lines[3]
	parameters := strings.Split(tmp, " ")
	cr := convertToInt(parameters[1])
	cw := convertToInt(parameters[3])
	cwa := convertToInt(parameters[5])

	return ac, sa, sh, sr, cr, cw, cwa
}

func convertToInt(value string) int {
	value = sed(value, " ", "")
	v, error := strconv.Atoi(value)
	if error != nil {
		logrus.Fatal(error)
	}
	return v
}

func sed(text, oldPattern, newPattern string) string {
	m := regexp.MustCompile(oldPattern)
	return m.ReplaceAllString(text, newPattern)
}

// src: https://stackoverflow.com/questions/57952695/prometheus-counters-how-to-get-current-value-with-golang-client
func getCounterValue(metric prometheus.Collector) float64 {
	var total float64
	collect(metric, func(m dto.Metric) {
		if h := m.GetHistogram(); h != nil {
			total += float64(h.GetSampleCount())
		} else {
			total += m.GetCounter().GetValue()
		}
	})
	return total
}

// collect calls the function for each metric associated with the Collector
func collect(col prometheus.Collector, do func(dto.Metric)) {
	c := make(chan prometheus.Metric)
	go func(c chan prometheus.Metric) {
		col.Collect(c)
		close(c)
	}(c)
	for x := range c { // eg range across distinct label vector values
		m := dto.Metric{}
		_ = x.Write(&m)
		do(m)
	}
}

func getDiffValue(metric prometheus.Collector, newValue int) float64 {
	currentValue := getCounterValue(metric)
	return float64(newValue) - currentValue
}
