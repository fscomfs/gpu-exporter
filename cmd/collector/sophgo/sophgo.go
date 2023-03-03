package sophgo

import (
	"fmt"
	"github.com/fscomfs/gpu-exporter/cmd/collector"
	"github.com/fscomfs/gpu-exporter/cmd/collector/sophgo/bmctl"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
	"log"
	"runtime"
	"strings"
	"sync"
)

var disabledFlag = false

func init() {
	defer func() {
		if error := recover(); error != nil {
			log.Printf("sophgo init fail %+v", error)
			disabledFlag = true
		}
	}()
	if err := bmctl.InitCtl(); err == nil {
		collector.Register("tpu", new())
		log.Printf("sophgo init success")
	} else {
		log.Printf("sophgo init fail %+v", err)
	}

}

type SophgoExporter struct {
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

func new() *SophgoExporter {
	infoLabels := collector.GetLabels(collector.Fields)
	return &SophgoExporter{
		gpuInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "graphics"),
			fmt.Sprintf("A metric with a constant '1' value labeled by gpu %s.",
				strings.Join(infoLabels, ", ")),
			infoLabels,
			nil,
		),
	}
}

func (e *SophgoExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.gpuInfoDesc
}
func (e *SophgoExporter) Collect(metricCh chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if disabledFlag {
		return
	}
	allInfo := bmctl.GetAllDeviceInfo()
	for deviceIndex, deviceInfo := range allInfo {
		used := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(deviceInfo.MemUsed), cast.ToString(deviceIndex), collector.MemoryUsed, "1684", collector.TPU, runtime.GOARCH)
		metricCh <- used
		free := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(deviceInfo.MemTotal-deviceInfo.MemUsed), cast.ToString(deviceIndex), collector.MemoryFree, "1684", collector.TPU, runtime.GOARCH)
		metricCh <- free
		total := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(deviceInfo.MemTotal), cast.ToString(deviceIndex), collector.MemoryTotal, "1684", collector.TPU, runtime.GOARCH)
		metricCh <- total
	}
}
