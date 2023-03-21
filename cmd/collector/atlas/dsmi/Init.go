package dsmi

import "fmt"

import "github.com/fscomfs/gpu-exporter/cmd/utils/dl"

var dcmi *dl.DynamicLibrary

const (
	dcmiLibraryName      = "libdrvdsmi_host.so"
	dcmiLibraryName2     = "libdrvdsmi.so"
	dcmiLibraryLoadFlags = dl.RTLD_LAZY | dl.RTLD_GLOBAL
)

type Return int32

func Init() Return {
	lib := dl.New(dcmiLibraryName, dcmiLibraryLoadFlags)
	if lib == nil {
		lib = dl.New(dcmiLibraryName2, dcmiLibraryLoadFlags)
	}
	if lib == nil {
		panic(fmt.Sprintf("error instantiating DynamicLibrary for %s", dcmiLibraryName))
	}
	err := lib.Open()
	if err != nil {
		panic(fmt.Sprintf("error opening %s: %v", dcmiLibraryName, err))
	}
	dcmi = lib
	return 0
}

func Shutdown() Return {
	err := dcmi.Close()
	if err != nil {
		panic(fmt.Sprintf("error closing %s: %v", dcmiLibraryName, err))
	}
	return 0
}
