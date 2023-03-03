package mlu

import (
	"fmt"
	"github.com/fscomfs/gpu-exporter/cmd/collector"
	"github.com/fscomfs/gpu-exporter/cmd/collector/mlu/cndev"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
	"log"
	"runtime"
	"strings"
	"sync"
)

var disabledFlag = false

var client cndev.Cndev

func init() {
	defer func() {
		if error := recover(); error != nil {
			log.Printf("mlu init fail %+v", error)
			disabledFlag = true
		}
	}()
	client = cndev.NewCndevClient()
	err := client.Init()
	if err != nil {
		disabledFlag = true
	} else {
		collector.Register("mlu", new())
		log.Printf("mlu init success")
	}
}

type MluExporter struct {
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

func new() *MluExporter {
	infoLabels := collector.GetLabels(collector.Fields)
	return &MluExporter{
		gpuInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "graphics"),
			fmt.Sprintf("A metric with a constant '1' value labeled by gpu %s.",
				strings.Join(infoLabels, ", ")),
			infoLabels,
			nil,
		),
	}
}

func (e *MluExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.gpuInfoDesc
}
func (e *MluExporter) Collect(metricCh chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if disabledFlag {
		return
	}
	deviceCount, err := client.GetDeviceCount()
	if err == nil && deviceCount > 0 {
		for i := uint(0); i < deviceCount; i++ {
			pyUsed, pyTotal, _, _, err := client.GetDeviceMemory(i)
			model := client.GetDeviceModel(i)
			if err != nil {
				used := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(pyUsed*1024), cast.ToString(i), collector.MemoryUsed, model, collector.MLU, runtime.GOARCH)
				metricCh <- used
				free := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64((pyTotal-pyUsed)*1024), cast.ToString(i), collector.MemoryFree, model, collector.MLU, runtime.GOARCH)
				metricCh <- free
				total := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(pyTotal*1024), cast.ToString(i), collector.MemoryTotal, model, collector.MLU, runtime.GOARCH)
				metricCh <- total
			}
		}
	}
}
