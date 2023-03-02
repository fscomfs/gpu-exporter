docker buildx build -t 192.168.1.76:8099/evtrain/cvmart-gpu-exporter-arm64:v1 --platform linux/arm64 -f Dockerfile . --load
docker build --build-arg TARGETARCH="amd64" -t 192.168.1.76:8099/evtrain/cvmart-gpu-exporter-amd64:v1 -f Dockerfile .
docker push 192.168.1.76:8099/evtrain/cvmart-gpu-exporter-arm64:v1
docker push 192.168.1.76:8099/evtrain/cvmart-gpu-exporter-amd64:v1
docker manifest rm 192.168.1.76:8099/evtrain/cvmart-gpu-exporter:v1
docker manifest create  192.168.1.76:8099/evtrain/cvmart-gpu-exporter:v1 \
192.168.1.76:8099/evtrain/cvmart-gpu-exporter-arm64:v1 \
192.168.1.76:8099/evtrain/cvmart-gpu-exporter-amd64:v1 -a --insecure

docker manifest annotate 192.168.1.76:8099/evtrain/cvmart-gpu-exporter:v1 \
192.168.1.76:8099/evtrain/cvmart-gpu-exporter-amd64:v1 --os linux --arch amd64

docker manifest annotate 192.168.1.76:8099/evtrain/cvmart-gpu-exporter:v1 \
192.168.1.76:8099/evtrain/cvmart-gpu-exporter-arm64:v1 --os linux --arch arm64

docker manifest push 192.168.1.76:8099/evtrain/cvmart-gpu-exporter:v1 --insecure