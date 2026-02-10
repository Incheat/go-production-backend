FROM alpine:3.20

ARG GRPCURL_VERSION=1.9.1
ARG GRPC_HEALTH_PROBE_VERSION=0.4.26

RUN set -eux; \
    apk add --no-cache curl ca-certificates tar gzip; \
    arch="$(apk --print-arch)"; \
    case "$arch" in \
      x86_64)  grpc_arch="x86_64"; probe_arch="amd64" ;; \
      aarch64) grpc_arch="arm64";  probe_arch="arm64" ;; \
      *) echo "unsupported arch: $arch" >&2; exit 1 ;; \
    esac; \
    \
    # ---- install grpcurl ----
    grpcurl_url="https://github.com/fullstorydev/grpcurl/releases/download/v${GRPCURL_VERSION}/grpcurl_${GRPCURL_VERSION}_linux_${grpc_arch}.tar.gz"; \
    curl -fsSL "$grpcurl_url" -o /tmp/grpcurl.tgz; \
    tar -xzf /tmp/grpcurl.tgz -C /usr/local/bin grpcurl; \
    chmod +x /usr/local/bin/grpcurl; \
    grpcurl -version; \
    \
    # ---- install grpc_health_probe ----
    probe_url="https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-${probe_arch}"; \
    curl -fsSL "$probe_url" -o /usr/local/bin/grpc_health_probe; \
    chmod +x /usr/local/bin/grpc_health_probe; \
    grpc_health_probe -version

CMD ["sh", "-c", "sleep infinity"]
