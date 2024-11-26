package metric

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricMonitor struct {
	requestsTotal  *prometheus.CounterVec
	requestLatency *prometheus.HistogramVec
	errorCount     *prometheus.CounterVec
}

func New() Monitor {
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_requests_total",
			Help: "Total number of requests",
		},
		[]string{"method", "status"})
	requestLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Latency of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
	errorCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP error responses",
		},
		[]string{"method", "path", "status"},
	)
	return &MetricMonitor{
		requestsTotal:  requestsTotal,
		requestLatency: requestLatency,
		errorCount:     errorCount,
	}
}

func (m *MetricMonitor) IncRequestsTotal(method string, status int) {
	m.requestsTotal.WithLabelValues(method, fmt.Sprint(status)).Inc()
}

func (m *MetricMonitor) IncRequestLatency(method, path string, duration float64) {
	m.requestLatency.WithLabelValues(method, path).Observe(duration)
}

func (m *MetricMonitor) IncErrorCount(method, path string, status int) {
	statusGroup := strconv.Itoa(status % 100)
	m.errorCount.WithLabelValues(method, path, statusGroup).Inc()
}
