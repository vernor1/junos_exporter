package storage

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vernor1/junos_exporter/collector"
	"github.com/vernor1/junos_exporter/rpc"
)

const prefix string = "junos_storage_"

var (
	totalBlocksDesc     *prometheus.Desc
	usedBlocksDesc      *prometheus.Desc
	availableBlocksDesc *prometheus.Desc
	usedPercentDesc     *prometheus.Desc
)

func init() {
	l := []string{"target", "device", "mountpoint"}
	totalBlocksDesc = prometheus.NewDesc(prefix+"total_blocks_count", "Total number of blocks", l, nil)
	usedBlocksDesc = prometheus.NewDesc(prefix+"used_blocks_count", "Number of used blocks", l, nil)
	availableBlocksDesc = prometheus.NewDesc(prefix+"available_blocks_count", "Number of available blocks", l, nil)
	usedPercentDesc = prometheus.NewDesc(prefix+"used_percent", "Percent of used storage", l, nil)
}

type storageCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &storageCollector{}
}

// Describe describes the metrics
func (*storageCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- totalBlocksDesc
	ch <- usedBlocksDesc
	ch <- availableBlocksDesc
	ch <- usedPercentDesc
}

// Collect collects metrics from JunOS
func (c *storageCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	var x = StorageRpc{}
	err := client.RunCommandAndParse("show system storage", &x)
	if err != nil {
		return err
	}

	for _, f := range x.Information.Filesystems {
		l := append(labelValues, f.FilesystemName, f.MountedOn)

		ch <- prometheus.MustNewConstMetric(totalBlocksDesc, prometheus.GaugeValue, float64(f.TotalBlocks), l...)
		ch <- prometheus.MustNewConstMetric(usedBlocksDesc, prometheus.GaugeValue, float64(f.UsedBlocks), l...)
		ch <- prometheus.MustNewConstMetric(availableBlocksDesc, prometheus.GaugeValue, float64(f.AvailableBlocks), l...)
		percent := strings.TrimSpace(f.UsedPercent)
		value, err := strconv.ParseFloat(percent, 64)
		if err != nil {
			value = 0
		}
		ch <- prometheus.MustNewConstMetric(usedPercentDesc, prometheus.GaugeValue, value, l...)
	}

	return nil
}
