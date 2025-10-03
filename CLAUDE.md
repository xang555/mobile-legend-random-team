# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

This Go project uses standard Go tools and Make targets:

```bash
# Build the binary
make build

# Run the server locally (uses configs/config.yaml by default)
make run

# Run all tests
make test

# Clean build artifacts
make clean

# Docker operations
make docker-build
make docker-run
```

The server accepts a `-config` flag to specify configuration file path.

## Project Architecture

This project follows the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) structure:

- **`cmd/server/`** - Application entry point and main function
- **`internal/`** - Private application code, not importable by other projects
  - `app/` - HTTP server wrapper with graceful shutdown
  - `config/` - Configuration loading using Viper with YAML + environment variables
  - `http/` - HTTP layer (handlers, routing with chi)
  - `random/` - Core business logic for team generation
- **`pkg/`** - Library code that could be imported by external applications
  - `logger/` - Zap logger configuration wrapper
- **`configs/`** - Default configuration files
- **`docs/`** - Additional documentation

## Configuration System

Configuration uses Viper with:
- Base config from `configs/config.yaml`
- Environment variable overrides with `RMT_` prefix
- Dot notation converted to underscores (e.g., `RMT_SERVER_PORT` overrides `server.port`)

Key config sections: `server`, `logging`, `team`. See `docs/configuration.md` for full reference.

## Dependencies

- **Chi v5** - HTTP router
- **Viper** - Configuration management
- **Zap** - Structured logging
- Go 1.21+

## API Endpoints

- `GET /healthz` - Health check endpoint
- `GET /api/v1/team/random` - Random team generation

The application serves a REST API for Mobile Legends team composition generation with configurable roles and hero pools.