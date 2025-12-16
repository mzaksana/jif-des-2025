#!/bin/bash

# Docker Scaling Demo - Startup Script
# JIF USK x Twibbonize Workshop

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

SCALE=${1:-1}

echo -e "${BLUE}========================================"
echo "  Docker Scaling Demo"
echo "  JIF USK x Twibbonize Workshop"
echo -e "========================================${NC}"
echo ""

# Check Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker not found. Please install Docker.${NC}"
    exit 1
fi

# Build and start
echo -e "${GREEN}Building and starting with $SCALE instance(s)...${NC}"
docker compose down --remove-orphans 2>/dev/null || true
docker compose up -d --build --scale app=$SCALE

echo ""
echo -e "${YELLOW}Waiting for services to be ready...${NC}"
sleep 5

# Check if running
if docker compose ps | grep -q "nginx.*running"; then
    echo ""
    echo -e "${GREEN}Service is running!${NC}"
    echo ""

    # Show running containers
    echo -e "${YELLOW}Running containers:${NC}"
    docker compose ps
    echo ""

    echo -e "${BLUE}========================================"
    echo "  Load balancer ready on port 8080"
    echo "  App instances: $SCALE"
    echo -e "========================================${NC}"
    echo ""
    echo "Test commands:"
    echo ""
    echo -e "${YELLOW}# Health check (shows instance ID)${NC}"
    echo "curl http://localhost:8080/health"
    echo ""
    echo -e "${YELLOW}# CPU work endpoint${NC}"
    echo "curl http://localhost:8080/api/work"
    echo ""
    echo -e "${YELLOW}# Scale to 5 instances${NC}"
    echo "docker compose up -d --scale app=5"
    echo ""
    echo -e "${YELLOW}# Stress test (requires 'hey')${NC}"
    echo "hey -n 500 -c 50 http://localhost:8080/api/work"
    echo ""
    echo -e "${YELLOW}# View live stats${NC}"
    echo "docker stats"
    echo ""
    echo -e "${YELLOW}# View logs${NC}"
    echo "docker compose logs -f app"
    echo ""
    echo -e "${YELLOW}# Stop service${NC}"
    echo "docker compose down"
    echo ""
else
    echo -e "${RED}Failed to start service. Check logs:${NC}"
    docker compose logs
    exit 1
fi
