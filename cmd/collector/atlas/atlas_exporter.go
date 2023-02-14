package atlas

import (
	"fmt"
	"github.com/fscomfs/gpu-exporter/cmd/collector"
	"github.com/fscomfs/gpu-exporter/cmd/collector/atlas/dsmi"
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
			log.Printf("dcmi init fail %+v", error)
			disabledFlag = true
		}
	}()
	dsmi.Init()
	collector.Register("npu", new())
	log.Printf("dsmi init success")

}

type AtlasExporter struct {
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

func new() *AtlasExporter {
	infoLabels := collector.GetLabels(collector.Fields)
	return &AtlasExporter{
		gpuInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "graphics"),
			fmt.Sprintf("A metric with a constant '1' value labeled by gpu %s.",
				strings.Join(infoLabels, ", ")),
			infoLabels,
			nil,
		),
	}
}

func (e *AtlasExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.gpuInfoDesc
}
func (e *AtlasExporter) Collect(metricCh chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if disabledFlag {
		return
	}
	deviceCount, r := dsmi.GetDeviceCount()
	if r == nil && deviceCount > 0 {
		allInfo, err := dsmi.AllDeviceInfo()
		if err == nil {
			for deviceIndex, deviceInfo := range allInfo {
				usedMemory := uint64(deviceInfo.CoreRate) * deviceInfo.Total / 100
				used := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(usedMemory), cast.ToString(deviceIndex), collector.MemoryFree, deviceInfo.ChipName, collector.NPU, runtime.GOARCH)
				metricCh <- used
				free := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(deviceInfo.Total-usedMemory), cast.ToString(deviceIndex), collector.MemoryUsed, deviceInfo.ChipName, collector.NPU, runtime.GOARCH)
				metricCh <- free
				total := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(deviceInfo.Total), cast.ToString(deviceIndex), collector.MemoryTotal, deviceInfo.ChipName, collector.NPU, runtime.GOARCH)
				metricCh <- total
			}
		}
	}
}
