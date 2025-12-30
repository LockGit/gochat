# GoChat Multi-Container Deployment - Change Log

## 2025-12-30

### Session 3: Added CI/CD Pipeline

**Makefile**
- Added `test` - Run all tests with race detector
- Added `test-coverage` - Generate HTML coverage report
- Added `test-unit` - Run unit tests only
- Added `test-integration` - Run integration tests with Docker
- Added `fmt` - Format code with go fmt
- Added `fmt-check` - Check if code is properly formatted
- Added `vet` - Run go vet static analysis
- Added `lint` - Run golangci-lint
- Added `build-binary` - Build gochat binary
- Added `build-image` - Build Docker image

**.github/workflows/ci-cd.yml**
- Created GitHub Actions CI/CD workflow
- Test job: go fmt, go vet, tests with Redis service container
- Build job: Docker image with caching
- Push job: Push to Docker Hub with multiple tags (latest, branch, git-sha)
- Deploy jobs: Optional deployment to dev/staging/prod (commented out, requires server setup)

**docker-compose.test.yml**
- Created test environment with Redis and etcd service containers
- Test runner service for integration tests

**docker/Dockerfile.test**
- Created test Dockerfile for running tests in Docker

**scripts/deploy.sh**
- Created deployment script for manual deployments
- Supports dev, staging, prod environments
- Pulls image, stops old services, starts new services

**config/staging/**
- Created staging configuration (copied from dev)

**docker-compose.staging.yml**
- Created staging environment override
- Uses staging config with info log level

**README.md**
- Added CI/CD Pipeline section with workflow description
- Added GitHub Secrets setup instructions
- Added branch strategy documentation
- Added manual deployment commands
- Added testing and building commands to Commands section

---

### Session 2: Removed entrypoint.sh

**tools/network.go**
- Added `GetContainerIP()` - Returns container's actual IPv4 address
- Added `GetServiceAddress()` - Replaces 0.0.0.0 with container IP for etcd registration

**logic/publish.go**
- Changed `addRegistryPlugin()` - ServiceAddress now uses `tools.GetServiceAddress(network, addr)`

**connect/rpc.go**
- Changed `addRegistryPlugin()` - ServiceAddress now uses `tools.GetServiceAddress(network, addr)`

**docker/Dockerfile**
- Removed entrypoint.sh copy and chmod commands
- Changed from `ENTRYPOINT ["/app/entrypoint.sh"]` to `CMD ["/app/gochat", "-module", "api"]`

**docker-compose.yml**
- All services: Removed `ETCD_HOST` and `REDIS_HOST` environment variables
- All services: Changed command from `["-module", "xxx"]` to `["/app/gochat", "-module", "xxx"]`

**docker/entrypoint.sh**
- Deleted (no longer needed)

---

### Session 1: Cleaned Up Legacy Files

**Removed directories:**
- `docker/dev/` - Supervisord configs for development (10 files)
- `docker/prod/` - Supervisord configs for production (10 files)

**Removed files:**
- `run.sh` - Legacy single-container deployment script
- `reload.sh` - Legacy reload script

**README.md**
- Removed "Legacy Deployment" section mentioning `run.sh`

---

## Future Updates

Add entries above in reverse chronological order (newest first).
