# Workshop Demo Projects - Code Prompts

---

## Demo List

| # | Demo | Purpose | Duration |
|---|------|---------|----------|
| 1 | CQRS + gRPC Blog Service | Show read/write separation with gRPC | 10-15 min |
| 2 | Docker Scaling + Stress Test | Show horizontal scaling under load | 10-15 min |

---

## Demo 1: CQRS + gRPC Blog Service

### Concept
- **Command Service**: Handles CREATE, UPDATE, DELETE (writes to primary DB)
- **Query Service**: Handles GET, LIST, SEARCH (reads from optimized read DB)
- **gRPC**: Fast binary protocol for service communication

### Architecture
```
┌─────────────┐     gRPC      ┌─────────────────────────────────┐
│   Client    │──────────────▶│          Blog Service           │
│  (grpcurl)  │               ├─────────────────────────────────┤
└─────────────┘               │  CommandService  │ QueryService │
                              │    (writes)      │   (reads)    │
                              └────────┬─────────┴───────┬──────┘
                                       │                 │
                                       ▼                 ▼
                              ┌─────────────┐   ┌─────────────┐
                              │  Write DB   │──▶│  Read DB    │
                              │  (SQLite)   │   │  (In-Memory)│
                              └─────────────┘   └─────────────┘
                                       │
                                  Event Sync
```

### Code Prompt

```
Create a CQRS + gRPC demo for a workshop presentation. Requirements:

**Tech Stack:**
- Go with gRPC
- SQLite for write database
- In-memory map for read database (simulating Redis/Elasticsearch)

**Structure:**
```
demo-cqrs-grpc/
├── proto/
│   └── blog.proto          # gRPC service definitions
├── cmd/
│   └── server/
│       └── main.go         # Entry point
├── internal/
│   ├── command/            # Command handlers (Create, Update, Delete)
│   ├── query/              # Query handlers (Get, List, Search)
│   ├── event/              # Event bus to sync write→read DB
│   └── store/
│       ├── write_store.go  # SQLite implementation
│       └── read_store.go   # In-memory implementation
├── go.mod
├── go.sum
├── Makefile                # proto gen, build, run commands
└── README.md               # How to run + demo script
```

**Proto Services:**
1. CommandService
   - CreatePost(title, content, author, tags) → id
   - UpdatePost(id, title, content, tags) → success
   - DeletePost(id) → success

2. QueryService
   - GetPost(id) → Post
   - ListPosts(limit, offset) → []Post
   - SearchPosts(query, tags) → []Post

**Key Demo Points:**
1. Show that writes go to SQLite (normalized, ACID)
2. Show that reads come from in-memory store (fast, denormalized)
3. Show event sync: after CreatePost, data appears in read store
4. Use grpcurl or evans CLI to demonstrate the calls

**Include demo script in README:**
- Create 3 posts via CommandService
- Query via QueryService (show fast reads)
- Search by tag
- Show the separation visually in logs

Keep it simple and focused for a 10-15 minute demo.
```

---

## Demo 2: Docker Scaling + Stress Test

### Concept
- Start with 1 container
- Stress test → show it struggling
- Scale to 3-5 containers
- Stress test again → show improved performance
- Visualize with live metrics

### Architecture
```
                              ┌─────────────────┐
                              │   Load Balancer │
                              │     (nginx)     │
                              └────────┬────────┘
                                       │
              ┌────────────────────────┼────────────────────────┐
              │                        │                        │
              ▼                        ▼                        ▼
       ┌─────────────┐          ┌─────────────┐          ┌─────────────┐
       │   App:1     │          │   App:2     │          │   App:3     │
       │  (Go API)   │          │  (Go API)   │          │  (Go API)   │
       └─────────────┘          └─────────────┘          └─────────────┘

Stress Test Tool: wrk / hey / k6
```

### Code Prompt

```
Create a Docker scaling demo for a workshop. Show horizontal scaling under load.

**Tech Stack:**
- Go + Fiber (simple REST API)
- Docker Compose with nginx load balancer
- `hey` or `wrk` for stress testing

**Structure:**
```
demo-scaling/
├── app/
│   ├── main.go             # Simple API with /health, /api/work (CPU task)
│   ├── Dockerfile
│   └── go.mod
├── nginx/
│   └── nginx.conf          # Load balancer config
├── docker-compose.yml      # Scalable service definition
├── stress-test.sh          # Script to run load tests
├── demo.sh                 # Full demo script with commentary
└── README.md               # Instructions
```

**API Endpoints:**
1. GET /health → { "status": "ok", "instance": "<container_id>" }
2. GET /api/work → Simulates CPU work (fibonacci or prime calculation)
3. GET /api/metrics → Current request count for this instance

**Docker Compose:**
- Service `app` with deploy.replicas configurable
- Service `nginx` as load balancer on port 8080
- Show container ID in response so we can see load distribution

**Demo Script (demo.sh):**
```bash
# Step 1: Start with 1 instance
docker compose up -d --scale app=1

# Step 2: Stress test (show it struggling)
echo "Testing with 1 instance..."
hey -n 1000 -c 50 http://localhost:8080/api/work

# Step 3: Scale to 5 instances
docker compose up -d --scale app=5

# Step 4: Stress test again (show improvement)
echo "Testing with 5 instances..."
hey -n 1000 -c 50 http://localhost:8080/api/work

# Step 5: Show load distribution
curl http://localhost:8080/health  # Run multiple times, see different instances
```

**Key Demo Points:**
1. Show response time degradation under load (1 instance)
2. Scale up with single command
3. Show response time improvement (5 instances)
4. Show requests distributed across instances
5. Bonus: docker stats to show CPU distribution

**Output should show:**
- Requests/sec comparison
- Average latency comparison
- Instance IDs in responses proving load balancing works

Keep simple, visual, impactful for 10-15 minute demo.
```

---

## Quick Reference for Workshop

### Running Demo 1 (CQRS + gRPC)
```bash
cd demo-cqrs-grpc
make proto        # Generate gRPC code
make run          # Start server

# In another terminal
grpcurl -plaintext localhost:50051 blog.CommandService/CreatePost
grpcurl -plaintext localhost:50051 blog.QueryService/ListPosts
```

### Running Demo 2 (Scaling)
```bash
cd demo-scaling
./demo.sh         # Runs full demo with commentary
```

---

## Pre-Workshop Checklist

- [ ] Go 1.21+ installed
- [ ] Docker + Docker Compose installed
- [ ] protoc + protoc-gen-go installed
- [ ] grpcurl or evans installed
- [ ] hey or wrk installed for stress testing
- [ ] Test both demos before presentation
