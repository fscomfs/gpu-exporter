FROM ubuntu:18.04
ARG TARGETARCH
RUN mkdir -p /cvmart-exporter
COPY build/bin/gpu-exporter-${TARGETARCH} /cvmart-exporter/gpu-exporter
ENV LD_LIBRARY_PATH=/host/usr/local/lib:/usr/local/Ascend/driver/lib64:/host/usr/local/Ascend/driver/lib64:/usr/local/dcmi:/host/usr/lib/:$LD_LIBRARY_PATH
ENTRYPOINT ["/cvmart-expoter/gpu-exporter","-listen-address",":9101"]