package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	defaultBuckets = []float64{10, 20, 30, 50, 80, 100, 200, 300, 500, 1000, 2000, 3000}
)

type Monitor struct {
	h           http.Handler
	reqCounter  *prometheus.CounterVec
	errCounter  *prometheus.CounterVec
	respLatency *prometheus.HistogramVec
}

func NewMonitor(application string, port int, buckets ...float64) *Monitor {
	var m Monitor
	process := strconv.Itoa(port)

	m.reqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   application,
			Name:        "requests_total",
			Help:        "Total request counts",
			ConstLabels: prometheus.Labels{"process": process},
		},
		[]string{
			"method",
			"endpoint",
		},
	)
	prometheus.MustRegister(m.reqCounter)

	m.errCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   application,
			Name:        "error_total",
			Help:        "Total error counts",
			ConstLabels: prometheus.Labels{"process": process},
		},
		[]string{
			"method",
			"endpoint",
		},
	)
	prometheus.MustRegister(m.errCounter)

	if len(buckets) == 0 {
		buckets = defaultBuckets
	}
	m.respLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   application,
			Name:        "response_latency_millisecond",
			Help:        "Response latency (millisecond)",
			ConstLabels: prometheus.Labels{"process": process},
			Buckets:     buckets,
		},
		[]string{
			"method",
			"endpoint",
		},
	)
	prometheus.MustRegister(m.respLatency)

	return &m
}

func Monitoring(m *Monitor) func(http.Handler) http.Handler {
	fn := func(h http.Handler) http.Handler {
		m.h = h
		return m
	}
	return fn
}

func (m *Monitor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()
	m.h.ServeHTTP(rw, r)
	m.reqCounter.WithLabelValues(r.Method, r.URL.Path).Inc()
	m.respLatency.WithLabelValues(r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
}
