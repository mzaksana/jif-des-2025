# Docker Scaling + Stress Test Demo

A demonstration of horizontal scaling with Docker Compose and nginx load balancing for the **JIF USK x Twibbonize Modern Software Architecture** workshop.

## Architecture

```
                              ┌─────────────────┐
                              │   Load Balancer │
                              │     (nginx)     │
                              │    :8080        │
                              └────────┬────────┘
                                       │
              ┌────────────────────────┼────────────────────────┐
              │                        │                        │
              ▼                        ▼                        ▼
       ┌─────────────┐          ┌─────────────┐          ┌─────────────┐
       │   App:1     │          │   App:2     │          │   App:3     │
       │  (Go API)   │          │  (Go API)   │          │  (Go API)   │
       └─────────────┘          └─────────────┘          └─────────────┘

Stress Test Tool: hey / wrk / ab
```

## Concepts

- **Horizontal Scaling**: Adding more instances to handle load
- **Load Balancing**: nginx distributes requests across instances
- **Stress Testing**: Measure performance under load
- **Container Orchestration**: Docker Compose manages multiple containers

## Prerequisites

```bash
# Docker & Docker Compose
docker --version
docker compose version

# Load testing tool (choose one)
go install github.com/rakyll/hey@latest  # Recommended
# or: brew install wrk
# or: apt install apache2-utils (for ab)
```

## Quick Start

```bash
# Run the full demo (recommended)
./demo.sh

# Or manually:
# Start with 1 instance
docker compose up -d --scale app=1

# Scale to 5 instances
docker compose up -d --scale app=5

# Run stress test
./stress-test.sh 1000 50
```

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /` | Service info |
| `GET /health` | Health check with instance ID |
| `GET /api/work` | CPU-intensive task (fibonacci) |
| `GET /api/metrics` | Request count per instance |

## Demo Script Walkthrough

The `demo.sh` script performs:

1. **Start with 1 instance** - Single container handling all requests
2. **Stress test** - Send 500 requests with 50 concurrent connections
3. **Scale to 5 instances** - Add 4 more containers
4. **Stress test again** - Compare performance
5. **Show load distribution** - Different instances handling requests

## Manual Commands

### Start Services
```bash
# Build and start
docker compose up -d --build

# Scale to N instances
docker compose up -d --scale app=3
```

### Monitor
```bash
# View logs
docker compose logs -f app

# Live stats
docker stats

# Check running containers
docker compose ps
```

### Stress Test
```bash
# Using hey
hey -n 1000 -c 50 http://localhost:8080/api/work

# Using wrk
wrk -t4 -c50 -d30s http://localhost:8080/api/work

# Using Apache Bench
ab -n 1000 -c 50 http://localhost:8080/api/work
```

### Verify Load Balancing
```bash
# Multiple requests show different instance IDs
for i in {1..10}; do curl -s http://localhost:8080/health | jq .instance; done
```

### Cleanup
```bash
docker compose down
```

## Project Structure

```
demo-scaling/
├── app/
│   ├── main.go             # Go API with Fiber
│   ├── Dockerfile          # Multi-stage build
│   ├── go.mod
│   └── go.sum
├── nginx/
│   └── nginx.conf          # Load balancer config
├── docker-compose.yml      # Service definitions
├── stress-test.sh          # Load test script
├── demo.sh                 # Full demo script
└── README.md
```

## Expected Results

### With 1 Instance
- Higher latency (100-500ms average)
- Lower requests/sec
- Single point of failure
- CPU maxed on one container

### With 5 Instances
- Lower latency (20-100ms average)
- ~5x requests/sec improvement
- Load distributed across containers
- CPU spread evenly

## Key Takeaways

| Aspect | 1 Instance | 5 Instances |
|--------|------------|-------------|
| Throughput | ~100 req/s | ~500 req/s |
| Latency (avg) | ~200ms | ~50ms |
| CPU per container | 100% | ~20% each |
| Fault tolerance | None | 4 can fail |

## Production Considerations

In a real-world scenario:
- Use **Kubernetes** for orchestration
- Implement **auto-scaling** based on metrics
- Add **health checks** and **readiness probes**
- Use **sticky sessions** if needed
- Monitor with **Prometheus/Grafana**
