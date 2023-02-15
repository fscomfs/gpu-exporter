package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var Collectors = make(map[string]prometheus.Collector)

const (
	MemoryFree               = "memory.free"
	MemoryTotal              = "memory.total"
	MemoryUsed               = "memory.used"
	UtilizationGpu           = "utilization.memory"
	UtilizationMemory        = "utilization.memory"
	Nvidia                   = "nvidia"
	NPU                      = "npu"
	TPU                      = "tpu"
	Index             QField = "index"
	Model             QField = "model"
	GpuType           QField = "gpuType"
	Arg               QField = "arg"
	ARCH              QField = "arch"
	DEVICE            QField = "device"
)

var (
	RequiredFields = []requiredField{
		{qField: MemoryFree, label: "memory.free"},
		{qField: MemoryUsed, label: "memory.used"},
		{qField: MemoryTotal, label: "memory.total"},
		{qField: UtilizationGpu, label: "utilization.memory"},
		{qField: UtilizationMemory, label: "utilization.memory"},
	}
	Fields = []requiredField{
		{qField: Index, label: "index"},
		{qField: Arg, label: "arg"},
		{qField: Model, label: "model"},
		{qField: GpuType, label: "gpuType"},
		{qField: ARCH, label: "arch"},
	}
	DiskFields = []requiredField{
		{qField: DEVICE, label: "device"},
	}
)

type QField string
type requiredField struct {
	qField QField
	label  string
}

func Register(name string, collector prometheus.Collector) {
	Collectors[name] = collector
}

func GetLabels(reqFields []requiredField) []string {
	r := make([]string, len(reqFields))
	for i, reqField := range reqFields {
		r[i] = reqField.label
	}
	return r
}
