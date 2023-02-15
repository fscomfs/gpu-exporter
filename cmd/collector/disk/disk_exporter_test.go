package disk

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"os"
	"testing"
)

func Test_dis(test *testing.T) {
	os.Setenv("HOST_ROOT", "/")
	stat, _ := disk.UsageWithContext(context.Background(), "/")
	partitionStat, _ := disk.Partitions(false)
	fmt.Println(stat.Total / 1024 / 1024)
	for _, p := range partitionStat {
		fmt.Println(p.Fstype)
	}
}
