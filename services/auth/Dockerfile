# syntax=docker/dockerfile:1.19.0

ARG GO_VERSION=1.25.5

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS builder
WORKDIR /src

# Build args provided by BuildKit/Buildx
ARG SERVICE
ARG TARGETOS
ARG TARGETARCH

# Cache deps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    go mod download

# Copy source
COPY services/ ./services/
COPY api/ ./api/

# Build (repro / CI-friendly)
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -buildvcs=false -ldflags="-s -w" \
    -o /out/${SERVICE} ./services/${SERVICE}/cmd

# Distroless (supported / updated tag)
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
ARG SERVICE
COPY --from=builder /out/${SERVICE} /app

EXPOSE 8080 9090
USER nonroot:nonroot
ENTRYPOINT ["/app"]
