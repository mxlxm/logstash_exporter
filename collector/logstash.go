package collector

import "github.com/prometheus/client_golang/prometheus"

const (
	// Namespace const string
	Namespace = "logstash"
)

// Collector interface implement Collect function
type Collector interface {
	Describe(ch chan<- *prometheus.Desc)
	Collect(ch chan<- prometheus.Metric)
}
