#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Configuration ---
# These can still be defined here or passed as environment variables
DOCKER_COMPOSE_FILE="docker-compose.yml" # Your docker-compose file name

# --- Deployment Steps ---
echo "Starting Docker-Compose deployment on VPS..."

echo "Building Docker images..."
docker-compose -f "${DOCKER_COMPOSE_FILE}" build --no-cache

echo "Stopping current Docker Compose services..."
docker-compose -f "${DOCKER_COMPOSE_FILE}" down

echo "Starting Docker Compose services..."
docker-compose -f "${DOCKER_COMPOSE_FILE}" up -d --remove-orphans

echo "Deployment complete on VPS!"
