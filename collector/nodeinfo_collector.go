package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// NodeInfoCollector type
type NodeInfoCollector struct {
	endpoint string

	Error     prometheus.Gauge
	NodeInfos *prometheus.Desc
	OsInfos   *prometheus.Desc
	JvmInfos  *prometheus.Desc
}

// NewNodeInfoCollector function
func NewNodeInfoCollector(logstashEndpoint string) (Collector, error) {
	const subsystem = "info"

	return &NodeInfoCollector{
		endpoint: logstashEndpoint,

		Error: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: subsystem,
			Name:      "last_scrape_error",
			Help:      "Whether the last scrape of metrics from logstash resulted in an error (1 for error, 0 for success).",
		}),

		NodeInfos: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "node"),
			"A metric with a constant '1' value labeled by Logstash version.",
			[]string{"version"},
			nil,
		),

		OsInfos: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "os"),
			"A metric with a constant '1' value labeled by name, arch, version and available_processors to the OS running Logstash.",
			[]string{"name", "arch", "version", "available_processors"},
			nil,
		),

		JvmInfos: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "jvm"),
			"A metric with a constant '1' value labeled by name, version and vendor of the JVM running Logstash.",
			[]string{"name", "version", "vendor"},
			nil,
		),
	}, nil
}

func (c *NodeInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Error.Desc()
}

// Collect function implements nodestats_collector collector
func (c *NodeInfoCollector) Collect(ch chan<- prometheus.Metric) {
	stats, err := NodeInfo(c.endpoint)
	if err != nil {
		c.Error.Set(1)
		return
	}

	ch <- c.Error
	ch <- prometheus.MustNewConstMetric(
		c.NodeInfos,
		prometheus.CounterValue,
		float64(1),
		stats.Version,
	)

	ch <- prometheus.MustNewConstMetric(
		c.OsInfos,
		prometheus.CounterValue,
		float64(1),
		stats.Os.Name,
		stats.Os.Arch,
		stats.Os.Version,
		strconv.Itoa(stats.Os.AvailableProcessors),
	)

	ch <- prometheus.MustNewConstMetric(
		c.JvmInfos,
		prometheus.CounterValue,
		float64(1),
		stats.Jvm.VMName,
		stats.Jvm.VMVersion,
		stats.Jvm.VMVendor,
	)

}
