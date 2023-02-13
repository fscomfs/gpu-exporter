package test_1

import (
	"fmt"
	"github.com/fscomfs/gpu-exporter/cmd/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
	"log"
	"strings"
	"sync"
)

var disabledFlag = false

func init() {
	//collector.Register("test", new())
}

type TestExporter struct {
	mutex              sync.RWMutex
	qFields            []collector.QField
	gpuInfoDesc        *prometheus.Desc
	failedScrapesTotal prometheus.Counter
}
type MetricInfo struct {
	desc            *prometheus.Desc
	MType           prometheus.ValueType
	ValueMultiplier float64
}

func new() *TestExporter {
	infoLabels := collector.GetLabels(collector.Fields)
	return &TestExporter{
		gpuInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "graphics"),
			fmt.Sprintf("A metric with a constant '1' value labeled by gpu %s.",
				strings.Join(infoLabels, ", ")),
			infoLabels,
			nil,
		),
	}
}

func (e *TestExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.gpuInfoDesc
}
func (e *TestExporter) Collect(metricCh chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if disabledFlag {
		return
	}
	log.Printf("collect test")
	var deviceCount = 1
	for i := 0; i < deviceCount; i++ {
		free := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(1), cast.ToString(i), collector.MemoryUsed, "test", "test")
		metricCh <- free
		used := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(2), cast.ToString(i), collector.MemoryFree, "test", "test")
		metricCh <- used
		total := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(3), cast.ToString(i), collector.MemoryTotal, "test", "test")
		metricCh <- total
	}
}
