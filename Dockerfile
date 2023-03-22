FROM ubuntu:18.04
ARG TARGETARCH
RUN mkdir -p /cvmart-exporter
COPY ./build/bin/gpu-exporter-${TARGETARCH} /cvmart-exporter/gpu-exporter
RUN chmod +x /cvmart-exporter/gpu-exporter
ENV LD_LIBRARY_PATH=/host/usr/local/lib:/usr/local/Ascend/driver/lib64:/host/usr/local/Ascend/driver/lib64:/usr/local/dcmi:/host/usr/lib/:/usr/lib64:$LD_LIBRARY_PATH
ENTRYPOINT ["/cvmart-exporter/gpu-exporter"]