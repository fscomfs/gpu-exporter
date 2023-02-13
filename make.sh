OOS=linux GOARCH=amd64 go build -o  ./build/bin/gpu-exporter-amd64
GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -o  ./build/bin/gpu-exporter-arm64