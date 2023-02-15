package disk

import (
	"context"
	"fmt"
	"github.com/fscomfs/gpu-exporter/cmd/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/disk"
	"strings"
	"sync"
)

var disabledFlag = false

func init() {
	collector.Register("disk", new())
}

type DiskExporter struct {
	mutex              sync.RWMutex
	qFields            []collector.QField
	disk_avail_bytes   *prometheus.Desc
	disk_total_bytes   *prometheus.Desc
	disk_used_bytes    *prometheus.Desc
	failedScrapesTotal prometheus.Counter
}
type MetricInfo struct {
	desc            *prometheus.Desc
	MType           prometheus.ValueType
	ValueMultiplier float64
}

func new() *DiskExporter {
	infoLabels := collector.GetLabels(collector.DiskFields)
	return &DiskExporter{
		disk_avail_bytes: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "disk_avail_bytes"),
			fmt.Sprintf("A metric with a constant '1' value labeled by gpu %s.",
				strings.Join(infoLabels, ", ")),
			infoLabels,
			nil,
		),
		disk_total_bytes: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "disk_total_bytes"),
			fmt.Sprintf("A metric with a constant '1' value labeled by gpu %s.",
				strings.Join(infoLabels, ", ")),
			infoLabels,
			nil,
		),
		disk_used_bytes: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "disk_used_bytes"),
			fmt.Sprintf("A metric with a constant '1' value labeled by gpu %s.",
				strings.Join(infoLabels, ", ")),
			infoLabels,
			nil,
		),
	}
}

func (e *DiskExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.disk_used_bytes
	ch <- e.disk_total_bytes
	ch <- e.disk_avail_bytes
}
func (e *DiskExporter) Collect(metricCh chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if disabledFlag {
		return
	}
	if stat, err := disk.UsageWithContext(context.Background(), "/"); err == nil {
		total := prometheus.MustNewConstMetric(e.disk_total_bytes, prometheus.GaugeValue, float64(stat.Total/1024), "ext4")
		metricCh <- total
		used := prometheus.MustNewConstMetric(e.disk_used_bytes, prometheus.GaugeValue, float64(stat.Used/1024), "ext4")
		metricCh <- used
		avail := prometheus.MustNewConstMetric(e.disk_avail_bytes, prometheus.GaugeValue, float64(stat.Free/1024), "ext4")
		metricCh <- avail
	}

}
