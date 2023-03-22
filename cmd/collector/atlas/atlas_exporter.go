package atlas

import "C"
import (
	"fmt"
	"github.com/fscomfs/gpu-exporter/cmd/collector"
	"github.com/fscomfs/gpu-exporter/cmd/collector/atlas/common"
	"github.com/fscomfs/gpu-exporter/cmd/collector/atlas/dcmi"
	"github.com/fscomfs/gpu-exporter/cmd/collector/atlas/dsmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
	"log"
	"runtime"
	"strings"
	"sync"
)

var disabledFlag = false
var dc *dcmi.DcManager

func init() {
	defer func() {
		if error := recover(); error != nil {
			log.Printf("dsmi init fail %+v", error)
			disabledFlag = true
		}
	}()
	dc = &dcmi.DcManager{}
	defer dc.DcShutDown()
	err := dc.DcInit()
	if err == nil {
		collector.Register("npu", new(true))
		log.Printf("dcmi init success")
	} else {
		dsmi.Init()
		collector.Register("npu", new(false))
		log.Printf("dsmi init success")
	}

}

type AtlasExporter struct {
	mutex              sync.RWMutex
	qFields            []collector.QField
	gpuInfoDesc        *prometheus.Desc
	failedScrapesTotal prometheus.Counter
	dcmi               bool
}
type MetricInfo struct {
	desc            *prometheus.Desc
	MType           prometheus.ValueType
	ValueMultiplier float64
}

func new(dcmi bool) *AtlasExporter {
	infoLabels := collector.GetLabels(collector.Fields)
	return &AtlasExporter{
		dcmi: dcmi,
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
	if e.dcmi {
		dc.DcInit()
		defer dc.DcShutDown()
		deviceCount, r := dc.DcGetDeviceCount()
		if r == nil && deviceCount > 0 {
			_, carList, _ := dc.DcGetCardList()
			index := 0
			for _, carIndex := range carList {
				deviceIdMax, _ := dc.DcGetDeviceNumInCard(carIndex)
				for deviceId := int32(0); deviceId < deviceIdMax; deviceId++ {
					memoryInfo, _ := dc.DcGetMemoryInfo(carIndex, deviceId)
					chipInfo, _ := dc.DcGetChipInfo(carIndex, deviceId)
					coreRate, _ := dc.DcGetDeviceUtilizationRate(carIndex, deviceId, common.AICore)
					usedMemory := uint64(coreRate) * memoryInfo.MemorySize * uint64(1000*1000) / 100
					used := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(usedMemory), cast.ToString(index), collector.MemoryUsed, chipInfo.Name, collector.NPU, runtime.GOARCH)
					metricCh <- used
					free := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(memoryInfo.MemorySize-usedMemory), cast.ToString(index), collector.MemoryFree, chipInfo.Name, collector.NPU, runtime.GOARCH)
					metricCh <- free
					total := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(memoryInfo.MemorySize*uint64(1000*1000)), cast.ToString(index), collector.MemoryTotal, chipInfo.Name, collector.NPU, runtime.GOARCH)
					metricCh <- total
					index++
				}
			}

		}
	} else {
		deviceCount, r := dsmi.GetDeviceCount()
		if r == nil && deviceCount > 0 {
			allInfo, err := dsmi.AllDeviceInfo()
			if err == nil {
				for deviceIndex, deviceInfo := range allInfo {
					usedMemory := uint64(deviceInfo.CoreRate) * deviceInfo.Total / 100
					used := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(usedMemory), cast.ToString(deviceIndex), collector.MemoryUsed, deviceInfo.ChipName, collector.NPU, runtime.GOARCH)
					metricCh <- used
					free := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(deviceInfo.Total-usedMemory), cast.ToString(deviceIndex), collector.MemoryFree, deviceInfo.ChipName, collector.NPU, runtime.GOARCH)
					metricCh <- free
					total := prometheus.MustNewConstMetric(e.gpuInfoDesc, prometheus.GaugeValue, float64(deviceInfo.Total), cast.ToString(deviceIndex), collector.MemoryTotal, deviceInfo.ChipName, collector.NPU, runtime.GOARCH)
					metricCh <- total
				}
			}
		}
	}

}
