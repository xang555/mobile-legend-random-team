BINARY := bin/random-ml-team
CONFIG ?= configs/config.yaml
IMAGE ?= random-ml-team:latest

.PHONY: build run test clean docker-build docker-run

build:
	go build -o $(BINARY) ./cmd/server

run:
	go run ./cmd/server -config $(CONFIG)

test:
	go test ./...

clean:
	rm -rf $(BINARY)

docker-build:
	docker build -f deployments/docker/Dockerfile -t $(IMAGE) .

docker-run:
	docker run --rm -p 8080:8080 -v $(PWD)/configs:/app/configs $(IMAGE) -config configs/config.yaml
