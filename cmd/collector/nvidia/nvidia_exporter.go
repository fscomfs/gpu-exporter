package nvidia

import (
	"fmt"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/fscomfs/gpu-exporter/cmd/collector"
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
			disabledFlag = true
			log.Printf("nvml init fail %+v", error)
		}
	}()
	r := nvml.Init()
	if r == nvml.SUCCESS {
		collector.Register("nvidia", new())
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
	log.Printf("collect nvidia")
	deviceCount, r := nvml.DeviceGetCount()
	if r == nvml.SUCCESS {
		for i := 0; i < deviceCount; i++ {
			var deviceHandler nvml.Device
			deviceHandler, r = nvml.DeviceGetHandleByIndex(i)
			if r == nvml.SUCCESS {
				memory, _ := nvml.DeviceGetMemoryInfo(deviceHandler)
				name, _ := nvml.DeviceGetName(deviceHandler)
				free := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(memory.Free), cast.ToString(i), collector.MemoryFree, name, collector.Nvidia, runtime.GOARCH)
				metricCh <- free
				used := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(memory.Used), cast.ToString(i), collector.MemoryUsed, name, collector.Nvidia, runtime.GOARCH)
				metricCh <- used
				total := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(memory.Total), cast.ToString(i), collector.MemoryTotal, name, collector.Nvidia, runtime.GOARCH)
				metricCh <- total
			}
		}
	}
}
