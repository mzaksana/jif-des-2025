#!/bin/bash

# CQRS + gRPC Blog Service Demo Script
# JIF USK x Twibbonize Workshop

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  CQRS + gRPC Demo Script${NC}"
echo -e "${BLUE}  JIF USK x Twibbonize Workshop${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check if server is running
if ! grpcurl -plaintext localhost:50051 list &>/dev/null; then
    echo -e "${YELLOW}Server not running. Please start it first with: make run${NC}"
    exit 1
fi

echo -e "${GREEN}Step 1: Create 3 posts via CommandService${NC}"
echo "----------------------------------------"
echo ""

echo -e "${YELLOW}Creating post 1: Introduction to CQRS${NC}"
grpcurl -plaintext -d '{
    "title": "Introduction to CQRS",
    "content": "CQRS stands for Command Query Responsibility Segregation. It separates read and write operations.",
    "author": "Ahmad",
    "tags": ["architecture", "cqrs", "tutorial"]
}' localhost:50051 blog.CommandService/CreatePost
echo ""

sleep 1

echo -e "${YELLOW}Creating post 2: gRPC Best Practices${NC}"
grpcurl -plaintext -d '{
    "title": "gRPC Best Practices",
    "content": "gRPC uses Protocol Buffers for fast, efficient communication between services.",
    "author": "Budi",
    "tags": ["grpc", "microservices", "tutorial"]
}' localhost:50051 blog.CommandService/CreatePost
echo ""

sleep 1

echo -e "${YELLOW}Creating post 3: Event-Driven Architecture${NC}"
grpcurl -plaintext -d '{
    "title": "Event-Driven Architecture",
    "content": "Events enable loose coupling between services and async processing.",
    "author": "Citra",
    "tags": ["architecture", "events", "async"]
}' localhost:50051 blog.CommandService/CreatePost
echo ""

sleep 1

echo -e "${GREEN}Step 2: Query all posts via QueryService${NC}"
echo "----------------------------------------"
echo ""

echo -e "${YELLOW}Listing all posts (from in-memory read store):${NC}"
grpcurl -plaintext -d '{"limit": 10, "offset": 0}' localhost:50051 blog.QueryService/ListPosts
echo ""

sleep 1

echo -e "${GREEN}Step 3: Search posts by tag${NC}"
echo "----------------------------------------"
echo ""

echo -e "${YELLOW}Searching for posts with tag 'architecture':${NC}"
grpcurl -plaintext -d '{"tags": ["architecture"]}' localhost:50051 blog.QueryService/SearchPosts
echo ""

sleep 1

echo -e "${YELLOW}Searching for posts with tag 'tutorial':${NC}"
grpcurl -plaintext -d '{"tags": ["tutorial"]}' localhost:50051 blog.QueryService/SearchPosts
echo ""

echo -e "${GREEN}Step 4: Search by query string${NC}"
echo "----------------------------------------"
echo ""

echo -e "${YELLOW}Searching for 'gRPC' in title/content:${NC}"
grpcurl -plaintext -d '{"query": "grpc"}' localhost:50051 blog.QueryService/SearchPosts
echo ""

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Demo Complete!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Key Observations:"
echo "  1. Writes went to SQLite (persistent, normalized)"
echo "  2. Reads came from in-memory store (fast, denormalized)"
echo "  3. Event sync kept both stores in sync"
echo "  4. Check server logs to see the separation!"
echo ""
