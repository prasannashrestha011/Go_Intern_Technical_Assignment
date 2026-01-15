package utils

import "github.com/prometheus/client_golang/prometheus"

var HttpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total http requests",
	},
	[]string{"method", "path", "status"},
)

var ActiveSessionsGuage = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "active_sessions",
	Help: "shows number of active sessions in the server.",
})

var HttpRequestDuration = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "shows observation of no. of request within interval of seconds.",
		Buckets: prometheus.LinearBuckets(0.01, 0.05, 10),
	},
)
