
.PHONY: build help compose-dev compose-dev-build compose-dev-down compose-prod compose-prod-build compose-prod-down compose-scale compose-logs compose-ps clean test-infra test test-coverage test-unit test-integration fmt fmt-check vet lint build-binary build-image

# Default HOST_IP for development
HOST_IP ?= 127.0.0.1

# Legacy build command (for backward compatibility)
build: # make build TAG=1.23;make build TAG=latest,自定义版本号,构建自己机器可以运行的镜像(因为M1机器拉取提供的镜像架构不同,只能自己构建)
	cd ./docker && docker build -t lockgit/gochat:${TAG} .

help:
	@echo "GoChat Multi-Container Deployment"
	@echo ""
	@echo "Available targets:"
	@echo "  build                - Legacy: Build single-container image (use TAG=version)"
	@echo "  compose-dev          - Start development environment"
	@echo "  compose-dev-build    - Build and start development environment"
	@echo "  compose-dev-down     - Stop development environment"
	@echo "  compose-prod         - Start production environment (detached)"
	@echo "  compose-prod-build   - Build and start production environment"
	@echo "  compose-prod-down    - Stop production environment"
	@echo "  compose-scale        - Start with scaled services (logic=3, connect-ws=2, task=2)"
	@echo "  compose-logs         - View logs from all services"
	@echo "  compose-ps           - Show running services"
	@echo "  test-infra           - Test infrastructure only (etcd + redis)"
	@echo "  clean                - Remove all containers, volumes, and images"
	@echo ""
	@echo "Testing & Quality:"
	@echo "  test                 - Run all tests"
	@echo "  test-coverage        - Run tests with coverage report"
	@echo "  test-unit            - Run unit tests only"
	@echo "  test-integration     - Run integration tests with Docker"
	@echo "  fmt                  - Format code with go fmt"
	@echo "  fmt-check            - Check if code is formatted"
	@echo "  vet                  - Run go vet"
	@echo "  lint                 - Run golangci-lint"
	@echo ""
	@echo "Building:"
	@echo "  build-binary         - Build gochat binary"
	@echo "  build-image          - Build Docker image"
	@echo ""
	@echo "Environment variables:"
	@echo "  HOST_IP=<ip>         - Set host IP address (default: 127.0.0.1)"

# Development commands
compose-dev:
	@echo "Starting GoChat in development mode..."
	HOST_IP=$(HOST_IP) docker compose -f docker-compose.yml -f deployments/docker-compose.dev.yml up

compose-dev-build:
	@echo "Building and starting GoChat in development mode..."
	HOST_IP=$(HOST_IP) docker compose -f docker-compose.yml -f deployments/docker-compose.dev.yml up --build

compose-dev-down:
	@echo "Stopping development environment..."
	docker compose -f docker-compose.yml -f deployments/docker-compose.dev.yml down

# Production commands
compose-prod:
	@echo "Starting GoChat in production mode..."
	@if [ -z "$(HOST_IP)" ] || [ "$(HOST_IP)" = "127.0.0.1" ]; then \
		echo "WARNING: HOST_IP is set to localhost. For production, use: make compose-prod HOST_IP=<your-server-ip>"; \
		echo "Continuing in 3 seconds..."; \
		sleep 3; \
	fi
	HOST_IP=$(HOST_IP) docker compose -f docker-compose.yml -f deployments/docker-compose.prod.yml up -d

compose-prod-build:
	@echo "Building and starting GoChat in production mode..."
	@if [ -z "$(HOST_IP)" ] || [ "$(HOST_IP)" = "127.0.0.1" ]; then \
		echo "WARNING: HOST_IP is set to localhost. For production, use: make compose-prod-build HOST_IP=<your-server-ip>"; \
		echo "Continuing in 3 seconds..."; \
		sleep 3; \
	fi
	HOST_IP=$(HOST_IP) docker compose -f docker-compose.yml -f deployments/docker-compose.prod.yml up -d --build

compose-prod-down:
	@echo "Stopping production environment..."
	docker compose -f docker-compose.yml -f deployments/docker-compose.prod.yml down

# Scaling example
compose-scale:
	@echo "Starting GoChat with scaled services..."
	HOST_IP=$(HOST_IP) docker compose -f docker-compose.yml up \
		--scale logic=3 \
		--scale connect-ws=2 \
		--scale connect-tcp=2 \
		--scale task=2 \
		--scale api=2

# Utility commands
compose-logs:
	docker compose -f docker-compose.yml logs -f

compose-ps:
	docker compose -f docker-compose.yml ps

# Test infrastructure only (etcd + redis)
test-infra:
	@echo "Starting infrastructure services only..."
	docker compose -f docker-compose.yml up etcd redis

# Clean up everything
clean:
	@echo "WARNING: This will remove all GoChat containers, volumes, and images."
	@echo "Press Ctrl+C to cancel, or wait 5 seconds to continue..."
	@sleep 5
	docker compose -f docker-compose.yml -f deployments/docker-compose.dev.yml down -v --rmi all || true
	docker compose -f docker-compose.yml -f deployments/docker-compose.prod.yml down -v --rmi all || true
	@echo "Clean complete!"

# Testing targets
test:
	@echo "Running all tests..."
	go test -v -race ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-unit:
	@echo "Running unit tests..."
	go test -v -short ./...

test-integration:
	@echo "Running integration tests with Docker..."
	@echo "Starting services..."
	docker compose -f docker-compose.yml -f deployments/docker-compose.test.yml up -d --build
	@echo "Waiting for services to be healthy..."
	@sleep 30
	@echo "Running integration tests from host..."
	go test -v ./... || (docker compose -f docker-compose.yml -f deployments/docker-compose.test.yml down && exit 1)
	@echo "Stopping services..."
	docker compose -f docker-compose.yml -f deployments/docker-compose.test.yml down

# Code quality targets
fmt:
	@echo "Formatting code..."
	go fmt ./...

fmt-check:
	@echo "Checking code formatting..."
	@test -z "$$(go fmt ./...)" || (echo "Code not formatted, run 'make fmt'" && exit 1)

vet:
	@echo "Running go vet..."
	go vet ./...

lint:
	@echo "Running golangci-lint..."
	golangci-lint run

# Build targets
build-binary:
	@echo "Building gochat binary..."
	CGO_ENABLED=1 GOOS=linux go build -tags=etcd -ldflags="-w -s" -o bin/gochat main.go
	@echo "Binary built: bin/gochat"

build-image:
	@echo "Building Docker image..."
	docker build -t gochat:latest -f docker/Dockerfile .
	@echo "Image built: gochat:latest"