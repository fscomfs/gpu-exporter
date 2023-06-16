package atlas

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"testing"
)

func TestSum(t *testing.T) {
	//a := dsmi.AtlasInfo{
	//	Total: uint32(9000),
	//	Used:  222,
	//}
	//var w C.int32
	//w = 999
	//b := w * 1024 * 1024
	//fmt.Printf("total:%+v", b)

	fmt.Println("test1")
	export := new(true)
	metrics := make(chan prometheus.Metric)
	export.Collect(metrics)

}

func TestSum2(t *testing.T) {
	//a := dsmi.AtlasInfo{
	//	Total: uint32(9000),
	//	Used:  222,
	//}
	//var w C.int32
	//w = 999
	//b := w * 1024 * 1024
	//fmt.Printf("total:%+v", b)
	fmt.Println("test2")

}
