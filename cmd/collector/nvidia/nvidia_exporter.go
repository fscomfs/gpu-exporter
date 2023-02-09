package nvidia

import (
	"fmt"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/fscomfs/gpu-exporter/cmd/collector"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strings"
	"sync"
)

var disabledFlag = false

func init() {
	defer func() {
		if error := recover(); error != nil {
			log.Printf("nvml init fail %+v", error)
		}
	}()
	r := nvml.Init()
	if r == nvml.SUCCESS {
		disabledFlag = true
		c := new()
		collector.Register("nvidia", c)
	} else {
		log.Printf("nvml init fail")
	}
}

type NvidiaExporter struct {
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

func new() *NvidiaExporter {
	infoLabels := collector.GetLabels(collector.Fields)
	return &NvidiaExporter{
		gpuInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "graphics"),
			fmt.Sprintf("A metric with a constant '1' value labeled by gpu %s.",
				strings.Join(infoLabels, ", ")),
			infoLabels,
			nil,
		),
	}
}

func (e *NvidiaExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.gpuInfoDesc
}
func (e *NvidiaExporter) Collect(metricCh chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if disabledFlag {
		return
	}
	deviceCount, r := nvml.DeviceGetCount()
	if r == nvml.SUCCESS {
		for i := 0; i < deviceCount; i++ {
			var deviceHandler nvml.Device
			deviceHandler, r = nvml.DeviceGetHandleByIndex(i)
			if r == nvml.SUCCESS {
				memory, _ := nvml.DeviceGetMemoryInfo(deviceHandler)
				name, _ := nvml.DeviceGetName(deviceHandler)
				free := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(memory.Free), string(i), collector.MemoryUsed, name, collector.Nvidia)
				metricCh <- free
				used := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(memory.Used), string(i), collector.MemoryUsed, name, collector.Nvidia)
				metricCh <- used
				total := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(memory.Total), string(i), collector.MemoryUsed, name, collector.Nvidia)
				metricCh <- total
			}
		}
	}
}
