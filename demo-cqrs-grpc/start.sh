#!/bin/bash

# CQRS + gRPC Blog Service - Docker Startup Script
# JIF USK x Twibbonize Workshop

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo -e "${BLUE}========================================"
echo "  CQRS + gRPC Blog Service"
echo "  JIF USK x Twibbonize Workshop"
echo -e "========================================${NC}"
echo ""

# Check Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker not found. Please install Docker.${NC}"
    exit 1
fi

# Create data directory
mkdir -p data

# Build and start
echo -e "${GREEN}Building and starting CQRS service...${NC}"
docker compose down --remove-orphans 2>/dev/null || true
docker compose up -d --build

echo ""
echo -e "${YELLOW}Waiting for service to be ready...${NC}"
sleep 5

# Check if running
if docker compose ps | grep -q "grpc-server.*running"; then
    echo ""
    echo -e "${GREEN}Service is running!${NC}"
    echo ""
    echo -e "${BLUE}========================================"
    echo "  Server ready on port 50051"
    echo -e "========================================${NC}"
    echo ""
    echo "Test commands (run in another terminal):"
    echo ""
    echo -e "${YELLOW}# Create a post${NC}"
    echo 'grpcurl -plaintext -d '\''{"title":"Hello CQRS","content":"This is a test","author":"Workshop","tags":["demo"]}'\'' localhost:50051 blog.CommandService/CreatePost'
    echo ""
    echo -e "${YELLOW}# List posts${NC}"
    echo 'grpcurl -plaintext -d '\''{"limit":10}'\'' localhost:50051 blog.QueryService/ListPosts'
    echo ""
    echo -e "${YELLOW}# Search posts${NC}"
    echo 'grpcurl -plaintext -d '\''{"tags":["demo"]}'\'' localhost:50051 blog.QueryService/SearchPosts'
    echo ""
    echo -e "${YELLOW}# View logs${NC}"
    echo "docker compose logs -f grpc-server"
    echo ""
    echo -e "${YELLOW}# Stop service${NC}"
    echo "docker compose down"
    echo ""
else
    echo -e "${RED}Failed to start service. Check logs:${NC}"
    docker compose logs
    exit 1
fi
