package main

import (
	"flag"
	"fmt"
	"github.com/fscomfs/gpu-exporter/cmd/collector"
	_ "github.com/fscomfs/gpu-exporter/cmd/collector/atlas"
	_ "github.com/fscomfs/gpu-exporter/cmd/collector/disk"
	_ "github.com/fscomfs/gpu-exporter/cmd/collector/mlu"
	_ "github.com/fscomfs/gpu-exporter/cmd/collector/nvidia"
	_ "github.com/fscomfs/gpu-exporter/cmd/collector/sophgo"
	_ "github.com/fscomfs/gpu-exporter/cmd/collector/test"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	stdLog "log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

const (
	redirectPageTemplate = `<html lang="en">
<head><title>Nvidia GPU Exporter</title></head>
<body>
<h1>CVMART GPU Exporter</h1>
<p><a href="%s">Metrics</a></p>
</body>
</html>
`
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func main() {

	flag.Parse()

	// Create non-global registry.
	reg := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	for n, p := range collector.Collectors {
		stdLog.Printf("register name:%+v,process:%+v", n, p)
		reg.MustRegister(p)
	}

	evnListen := os.Getenv("EXPORTER_LISTEN")
	listenAddr := *addr
	if evnListen != "" {
		listenAddr = evnListen
	}
	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)
	rootHandler := NewRootHandler(logger, "/metrics")
	http.Handle("/", rootHandler)
	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	stdLog.Printf("http.ListenAndServe %+v", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

type RootHandler struct {
	response []byte
	logger   log.Logger
}

func NewRootHandler(logger log.Logger, metricsPath string) *RootHandler {
	return &RootHandler{
		response: []byte(fmt.Sprintf(redirectPageTemplate, metricsPath)),
		logger:   logger,
	}
}

func (r *RootHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write(r.response); err != nil {
		_ = level.Error(r.logger).Log("msg", "Error writing redirect", "err", err)
	}
}
