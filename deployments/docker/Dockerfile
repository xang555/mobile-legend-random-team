# syntax=docker/dockerfile:1.4

# --- Build stage ---
FROM golang:1.21-alpine AS build

WORKDIR /src

# Install build tools
RUN apk add --no-cache ca-certificates tzdata upx

# Cache go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy source
COPY . .

# Build static binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o /bin/random-ml-team ./cmd/server && \
    upx /bin/random-ml-team

# --- Runtime stage ---
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=build /bin/random-ml-team /app/random-ml-team
COPY configs /app/configs

EXPOSE 8080

ENTRYPOINT ["/app/random-ml-team", "-config", "configs/config.yaml"]
