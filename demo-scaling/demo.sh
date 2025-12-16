#!/bin/bash

# Docker Scaling Demo Script
# JIF USK x Twibbonize Workshop

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

REQUESTS=500
CONCURRENCY=50

echo -e "${BLUE}========================================"
echo "  Docker Scaling Demo"
echo "  JIF USK x Twibbonize Workshop"
echo -e "========================================${NC}"
echo ""

# Check prerequisites
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker not found. Please install Docker.${NC}"
    exit 1
fi

if ! command -v hey &> /dev/null; then
    echo -e "${YELLOW}Warning: 'hey' not found. Install with: go install github.com/rakyll/hey@latest${NC}"
    echo "Falling back to curl-based test..."
    USE_CURL=true
fi

# Step 1: Build and start with 1 instance
echo -e "${GREEN}Step 1: Starting with 1 instance${NC}"
echo "----------------------------------------"
docker compose down --remove-orphans 2>/dev/null || true
docker compose up -d --build --scale app=1
echo ""
echo "Waiting for services to be ready..."
sleep 5
echo ""

# Verify it's running
echo -e "${YELLOW}Testing single instance...${NC}"
curl -s http://localhost:8080/health | python3 -m json.tool 2>/dev/null || curl -s http://localhost:8080/health
echo ""
echo ""

# Step 2: Stress test with 1 instance
echo -e "${GREEN}Step 2: Stress test with 1 instance${NC}"
echo "----------------------------------------"
if [ "$USE_CURL" = true ]; then
    echo "Running simplified test with curl..."
    for i in {1..10}; do
        curl -s http://localhost:8080/api/work > /dev/null &
    done
    wait
    echo "Simple test completed"
else
    hey -n $REQUESTS -c $CONCURRENCY http://localhost:8080/api/work
fi
echo ""

# Show current stats
echo -e "${YELLOW}Current container stats:${NC}"
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
echo ""

# Step 3: Scale to 5 instances
echo -e "${GREEN}Step 3: Scaling to 5 instances${NC}"
echo "----------------------------------------"
docker compose up -d --scale app=5
echo ""
echo "Waiting for new instances..."
sleep 5
echo ""

# Show running containers
echo -e "${YELLOW}Running containers:${NC}"
docker compose ps
echo ""

# Step 4: Stress test with 5 instances
echo -e "${GREEN}Step 4: Stress test with 5 instances${NC}"
echo "----------------------------------------"
if [ "$USE_CURL" = true ]; then
    echo "Running simplified test with curl..."
    for i in {1..10}; do
        curl -s http://localhost:8080/api/work > /dev/null &
    done
    wait
    echo "Simple test completed"
else
    hey -n $REQUESTS -c $CONCURRENCY http://localhost:8080/api/work
fi
echo ""

# Show stats after scaling
echo -e "${YELLOW}Container stats after scaling:${NC}"
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
echo ""

# Step 5: Demonstrate load balancing
echo -e "${GREEN}Step 5: Demonstrating load distribution${NC}"
echo "----------------------------------------"
echo "Calling /health 10 times to show different instances:"
echo ""
for i in {1..10}; do
    RESPONSE=$(curl -s http://localhost:8080/health)
    INSTANCE=$(echo $RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin)['instance'][:12])" 2>/dev/null || echo "unknown")
    echo "  Request $i -> Instance: $INSTANCE"
done
echo ""

# Summary
echo -e "${BLUE}========================================"
echo "  Demo Complete!"
echo -e "========================================${NC}"
echo ""
echo "Key Observations:"
echo "  1. Single instance: Higher latency under load"
echo "  2. After scaling: Lower latency, higher throughput"
echo "  3. Load balancing: Requests distributed across instances"
echo "  4. CPU distribution: Load shared among containers"
echo ""
echo "Commands to explore:"
echo "  - docker compose logs -f app     # View app logs"
echo "  - docker stats                   # Live resource usage"
echo "  - docker compose ps              # Running services"
echo ""
echo -e "${YELLOW}To clean up: docker compose down${NC}"
echo ""
