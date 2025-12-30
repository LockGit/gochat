# GoChat - Multi-Container Deployment

> This is a fork of [LockGit/gochat](https://github.com/LockGit/gochat) with Docker Compose multi-container deployment support.

## What's New

This fork adds **production-ready multi-container deployment** to the original gochat project:

- ✅ **8 Independent Containers**: Each service (etcd, redis, logic, connect-ws, connect-tcp, task, api, site) runs in its own container
- ✅ **Horizontal Scaling**: Scale Logic, Connect, Task, and API services independently
- ✅ **Docker Compose**: One-command deployment with `make compose-dev` or `make compose-prod`
- ✅ **Auto-Configuration**: Container IPs automatically registered to etcd for proper RPC communication
- ✅ **Health Checks**: Each service monitors its own health with auto-restart
- ✅ **Dev/Prod Configs**: Separate configurations for development and production environments

### Quick Start

```bash
# Start all services
make compose-dev

# Or manually
docker compose -f docker-compose.yml -f deployments/docker-compose.dev.yml up

# Visit the chat app
http://localhost:8080
```

### Scaling Services

```bash
# Scale specific services
docker compose up --scale logic=3 --scale connect-ws=2 --scale task=2
```

### Production Deployment

```bash
make compose-prod HOST_IP=<your-server-ip>
```

## Architecture

```
├── etcd          → Service discovery
├── redis         → Message queue & cache
├── logic × N     → Business logic RPC (scalable)
├── connect-ws × N→ WebSocket handler (scalable)
├── connect-tcp × N→ TCP handler (scalable)
├── task × N      → Message worker (scalable)
├── api × N       → REST API (scalable)
└── site          → Frontend
```

## Documentation

- **Multi-Container Guide**: See [README-compose.md](./README-compose.md)
- **Original Documentation**: See below or visit [github.com/LockGit/gochat](https://github.com/LockGit/gochat)

---

## Original README

### gochat is a lightweight IM server implemented using pure go

* gochat is an instant messaging system based on go. It supports private messages and room broadcast messages. It communicates between layers through RPC and supports horizontal expansion.
* Supports websocket and TCP access, and supports websocket and TCP message interworking.
* Based on etcd service discovery, making it convenient to scale and deploy.
* Using redis as the carrier of message storage and delivery is very lightweight.
* Because of Go's cross-compilation features, compiled binaries can run quickly on various platforms.
* Clear architecture and directory structure.

### Components

```
Language: golang
Database: sqlite3 (can be replaced with MySQL or other databases)
Database ORM: gorm
Service Discovery: etcd
RPC Communication: rpcx
Queue: redis (can be replaced with kafka or rabbitmq)
Cache: redis
Message ID: snowflake algorithm
```

### Architecture Design
![](./architecture/gochat.png)

### Service Discovery
![](./architecture/gochat_discovery.png)

### Message Delivery
![](./architecture/single_send.png)

### Features

- Private messaging
- Room broadcast messaging
- Websocket support
- TCP support
- Horizontal scaling
- Service discovery via etcd
- Message queuing via Redis
- Lightweight and fast

### Directory Structure

```
.
├── api          # API layer, provides REST API services
├── connect      # Connection layer, handles long connections
├── logic        # Logic layer, handles business logic
├── task         # Task layer, consumes queue messages
├── site         # Frontend static files
├── config       # Configuration files
├── db           # Database initialization
├── docker       # Docker related files
├── proto        # RPC proto files
└── tools        # Utility methods
```

### Test Users

Default test users (username/password):
- demo / 111111
- test / 111111
- admin / 111111

### Commands

```bash
# Development
make compose-dev              # Start dev environment
make compose-logs             # View logs
make compose-ps               # Check status

# Production
make compose-prod HOST_IP=x.x.x.x    # Deploy to server

# Testing & Quality
make test                     # Run all tests
make test-coverage            # Run tests with coverage
make fmt                      # Format code
make vet                      # Run go vet
make lint                     # Run linter

# Building
make build-binary             # Build gochat binary
make build-image              # Build Docker image

# Utilities
make clean                    # Clean up everything
```

### CI/CD Pipeline

#### GitHub Actions Workflow

Automated CI/CD pipeline for testing, building, and deploying GoChat:

- **Test**: Runs on every push and pull request
  - Go fmt check
  - Go vet static analysis
  - Unit and integration tests with coverage
  - Redis service container for tests

- **Build**: Builds Docker image after tests pass
  - Multi-stage Docker build
  - Image caching for faster builds

- **Push**: Pushes to Docker Hub (on push to branches)
  - Multiple tags: `latest`, `<branch>`, `<git-sha>`
  - Credentials via GitHub Secrets

- **Deploy**: Optional deployment to environments
  - Development: Auto-deploy on `dev` branch
  - Staging: Auto-deploy on `staging` branch
  - Production: Manual approval on `master` branch

#### Setup GitHub Secrets

Configure these in repository Settings → Secrets and variables → Actions:

**Required for Docker Hub:**
- `DOCKERHUB_USERNAME` - Your Docker Hub username
- `DOCKERHUB_TOKEN` - Docker Hub access token

**Optional for server deployment:**
- `DEV_SERVER_HOST`, `DEV_SERVER_USER`, `DEV_SERVER_KEY`
- `STAGING_SERVER_HOST`, `STAGING_SERVER_USER`, `STAGING_SERVER_KEY`
- `PROD_SERVER_HOST`, `PROD_SERVER_USER`, `PROD_SERVER_KEY`

#### Branch Strategy

- `dev` → Development environment (auto-deploy)
- `staging` → Staging environment (auto-deploy)
- `master` → Production environment (manual approval)

#### Manual Deployment

```bash
# Pull and deploy specific version
./scripts/deploy.sh dev latest
./scripts/deploy.sh staging abc123def
./scripts/deploy.sh prod latest
```

### License

MIT License

### Credits

**Original Project**: [LockGit/gochat](https://github.com/LockGit/gochat)

Special thanks to LockGit for creating this excellent IM system.
