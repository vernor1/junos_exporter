package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vernor1/junos_exporter/rpc"
)

// RPCCollector collects metrics from JunOS using rpc.Client
type RPCCollector interface {

	// Describe describes the metrics
	Describe(ch chan<- *prometheus.Desc)

	// Collect collects metrics from JunOS
	Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error
}
