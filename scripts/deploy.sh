#!/bin/bash

# GoChat Deployment Script
# Usage: ./scripts/deploy.sh <environment> <image-tag>
# Example: ./scripts/deploy.sh dev latest

set -e

# Configuration
ENVIRONMENT=$1
IMAGE_TAG=$2
DOCKER_IMAGE="${DOCKERHUB_USERNAME:-yourname}/gochat"

# Validate arguments
if [ -z "$ENVIRONMENT" ] || [ -z "$IMAGE_TAG" ]; then
    echo "Usage: $0 <environment> <image-tag>"
    echo "Example: $0 dev latest"
    echo ""
    echo "Environments: dev, staging, prod"
    echo "Image tags: latest, dev, staging, <git-sha>"
    exit 1
fi

# Validate environment
if [[ ! "$ENVIRONMENT" =~ ^(dev|staging|prod)$ ]]; then
    echo "Error: Invalid environment '$ENVIRONMENT'"
    echo "Valid environments: dev, staging, prod"
    exit 1
fi

echo "================================================"
echo "GoChat Deployment"
echo "================================================"
echo "Environment: $ENVIRONMENT"
echo "Image Tag: $IMAGE_TAG"
echo "Docker Image: $DOCKER_IMAGE:$IMAGE_TAG"
echo "================================================"

# Pull latest image
echo "Pulling Docker image..."
docker pull $DOCKER_IMAGE:$IMAGE_TAG

# Stop existing services
echo "Stopping existing services..."
docker compose -f docker-compose.yml -f deployments/docker-compose.$ENVIRONMENT.yml down

# Start new services
echo "Starting services..."
docker compose -f docker-compose.yml -f deployments/docker-compose.$ENVIRONMENT.yml up -d

# Wait for services to start
echo "Waiting for services to start..."
sleep 10

# Health check
echo "Checking service status..."
docker compose ps

echo "================================================"
echo "Deployment complete!"
echo "================================================"

# Show service URLs
if [ "$ENVIRONMENT" = "dev" ]; then
    echo "WebSocket: http://localhost:7000"
    echo "API: http://localhost:7070"
    echo "Site: http://localhost:8080"
elif [ "$ENVIRONMENT" = "staging" ]; then
    echo "Check your staging server for service URLs"
elif [ "$ENVIRONMENT" = "prod" ]; then
    echo "Check your production server for service URLs"
fi
