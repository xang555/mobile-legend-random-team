# Random Mobile Legends Team API

Production-ready Go service that returns random team compositions for Mobile Legends. The project follows the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) guidelines, includes structured logging, configuration management, graceful shutdown, and containerized deployment.

## Features
- Generate random Mobile Legends teams based on configurable roles and hero pools
- REST API with `/api/v1/team/random` and `/healthz`
- Structured logging with Zap
- Configuration via YAML and environment overrides (`RMT_` prefix)
- Graceful shutdown and server timeouts
- Production-focused Docker image and Make targets

## Quick Start
```bash
# Install dependencies and generate go.sum
go mod tidy

# Build binary
make build

# Run locally
make run

# Run tests
make test
```

## Configuration
Default configuration lives in `configs/config.yaml`. Override any value with environment variables prefixed with `RMT_`, for example:
```bash
export RMT_SERVER_PORT=9090
export RMT_LOGGING_LEVEL=debug
```

See `docs/configuration.md` for the full reference.

## API
- `GET /healthz` — readiness probe
- `GET /api/v1/team/random` — returns a random team composition

Example response:
```json
{
  "members": [
    {"role": "Tank", "hero": "Atlas"},
    {"role": "Mage", "hero": "Lunox"},
    {"role": "Marksman", "hero": "Beatrix"},
    {"role": "Fighter", "hero": "Paquito"},
    {"role": "Assassin", "hero": "Gusion"}
  ]
}
```

## Docker Deployment
A multi-stage Dockerfile lives in `deployments/docker/Dockerfile`. Build and run with:
```bash
docker build -f deployments/docker/Dockerfile -t random-ml-team:latest .
docker run --rm -p 8080:8080 random-ml-team:latest
```

The container copies the default config; mount your own configs if required:
```bash
docker run --rm -p 8080:8080 -v $(pwd)/configs:/app/configs random-ml-team:latest
```

## Operations
- Production probes should hit `/healthz` for liveness/readiness
- Logs emit structured JSON by default; configure `logging.encoding` to `console` for local dev
- All timeouts are configurable in `configs/config.yaml`

Additional documentation is available under `docs/`.
