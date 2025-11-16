package main

import "github.com/prometheus/client_golang/prometheus"

var (
	TotalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "lb_total_requests",
			Help: "Total number of requests received",
		},
	)

	ActiveBackends = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "lb_active_backends",
			Help: "Number of alive backends",
		},
	)
)

func InitMetrics() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(ActiveBackends)
}
