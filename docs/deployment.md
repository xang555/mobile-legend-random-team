# Deployment Guide

This service ships with a production-grade Docker image built from `deployments/docker/Dockerfile`.

## Build
```bash
docker build -f deployments/docker/Dockerfile -t ghcr.io/<org>/random-ml-team:$(git rev-parse --short HEAD) .
```

### Build Arguments
- Uses multi-stage build with Go 1.21.
- Produces a statically linked binary compressed with `upx` to reduce size.

## Run
```bash
docker run --rm \
  -p 8080:8080 \
  -e RMT_SERVER_PORT=8080 \
  ghcr.io/<org>/random-ml-team:latest
```

Mount a custom config file:
```bash
docker run --rm \
  -p 8080:8080 \
  -v $(pwd)/configs:/app/configs \
  ghcr.io/<org>/random-ml-team:latest \
  -config configs/config.yaml
```

## Kubernetes
For Kubernetes deployments, expose `/healthz` for readiness/liveness probes and configure resource limits based on expected RPS. Sample manifests can be derived from the Docker image; ensure configuration is injected via ConfigMap or environment variables.

## Observability
- Logs: structured JSON to stdout.
- Health probe: `/healthz` returns 200 on success.
- Add metrics/tracing by extending `internal/http/router` middleware.
