# Modern Software Architecture - Project Examples

A practical project list to demonstrate each concept from the JIF USK x Twibbonize presentation.

---

## 1. Foundational Attributes (Slide 2)

### Project: **Simple URL Shortener Service**
Demonstrates all 6 attributes in one project:

| Attribute | Implementation |
|-----------|----------------|
| **Scalability** | Redis caching for hot URLs, horizontal scaling ready |
| **Reliability** | Health checks, graceful shutdown, retry logic |
| **Maintainability** | Clean folder structure, typed interfaces, documentation |
| **Security** | Rate limiting, input validation, SQL injection prevention |
| **Performance** | Response time < 50ms, connection pooling |
| **Observability** | Prometheus metrics, structured logging, tracing IDs |

**Tech Stack:** Go + Fiber, Redis, MySQL, Prometheus + Grafana

---

## 2. Monolithic vs Microservices (Slide 3)

### Project A: **E-Commerce Monolith**
Single codebase handling:
- User authentication
- Product catalog
- Shopping cart
- Order processing
- Payment integration

**Tech:** Next.js full-stack, PostgreSQL, single Docker container

### Project B: **E-Commerce Microservices**
Same features, split into independent services:

```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   Auth      │  │   Catalog   │  │   Cart      │
│   Service   │  │   Service   │  │   Service   │
└─────────────┘  └─────────────┘  └─────────────┘
       │                │                │
       └────────────────┼────────────────┘
                        │
              ┌─────────────────┐
              │   API Gateway   │
              │   (KrakenD)     │
              └─────────────────┘
```

**Tech:** Go services, KrakenD gateway, separate databases per service

---

## 3. Database Sharding (Slide 4)

### Project: **User Analytics Platform**

Shard 10M+ user events by user_id:

```go
func getShardID(userID string) int {
    hash := fnv.New32a()
    hash.Write([]byte(userID))
    return int(hash.Sum32() % 3) // 3 shards
}
```

**Implementation:**
- Shard 1: Users A-I (user_id hash % 3 == 0)
- Shard 2: Users J-R (user_id hash % 3 == 1)  
- Shard 3: Users S-Z (user_id hash % 3 == 2)

**Challenges to demonstrate:**
- Cross-shard queries for reporting
- Hot partition detection
- Shard rebalancing strategy

**Tech:** MySQL cluster, custom shard router in Go

---

## 4. CAP Theorem & Distributed Consistency (Slide 5)

### Project: **Distributed Counter Service**

Three implementations showing CAP trade-offs:

| Mode | Guarantees | Use Case |
|------|-----------|----------|
| **CP Mode** | Strong consistency, may reject writes during partition | Bank balance |
| **AP Mode** | Always available, eventual consistency | Like counts |
| **CA Mode** | Single node, no partition tolerance | Local cache |

**Demo scenarios:**
1. Kill one node, observe behavior
2. Network partition simulation
3. Conflict resolution with CRDTs

**Tech:** Go + Redis Cluster or custom implementation with Raft consensus

---

## 5. Communication Patterns (Slide 6)

### Project: **Order Processing System**

#### Synchronous (REST/gRPC):
```
Client → Order API → Inventory API → Payment API → Response
         (waits)      (waits)         (waits)
```

#### Asynchronous (Event-Driven):
```
Client → Order API → [Kafka: order.created]
                            ↓
              ┌─────────────┼─────────────┐
              ↓             ↓             ↓
         Inventory      Payment      Notification
         Consumer       Consumer       Consumer
```

**Compare:**
- Response time
- Failure handling
- System coupling

**Tech:** Go + Fiber (REST), gRPC, Kafka/Redis Streams

---

## 6. Stability & Performance Patterns (Slide 7)

### Project: **Weather Aggregator API**

Fetches from multiple external APIs with resilience patterns:

```go
// Caching Layer
cache := redis.Get("weather:jakarta")
if cache != nil {
    return cache // Cache hit
}

// Circuit Breaker
result := circuitBreaker.Execute(func() {
    return fetchFromWeatherAPI()
})

// Retry with Exponential Backoff
retry.Do(
    fetchFromBackupAPI,
    retry.Attempts(3),
    retry.Delay(100*time.Millisecond),
    retry.DelayType(retry.BackOffDelay),
)
```

**Demo scenarios:**
- Cache hit/miss performance comparison
- Circuit breaker state transitions (Closed → Open → Half-Open)
- Retry behavior visualization

**Tech:** Go, Redis, sony/gobreaker, external weather APIs

---

## 7. CQRS & Saga Pattern (Slide 8)

### Project A: **Blog Platform with CQRS**

```
WRITE (Commands)              READ (Queries)
     │                              │
     ▼                              ▼
┌─────────┐                  ┌─────────────┐
│ MongoDB │ ──sync event──→  │ Elasticsearch│
│ (posts) │                  │ (search)    │
└─────────┘                  └─────────────┘
```

- Commands: CreatePost, UpdatePost, DeletePost
- Queries: SearchPosts, GetPostsByTag, GetTrending

### Project B: **Travel Booking Saga**

```
Book Flight → Book Hotel → Book Car → Confirm
     │            │           │
     ▼            ▼           ▼
  (fail?)      (fail?)     (fail?)
     │            │           │
     ▼            ▼           ▼
Cancel All ← Cancel Hotel ← Cancel Car
```

**Saga Orchestrator** manages compensation on failures.

**Tech:** Go, MongoDB, Elasticsearch, Redis for event bus

---

## 8. Containers & Orchestration (Slide 9)

### Project: **Kubernetes Microservices Deployment**

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: getter-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: getter
  template:
    spec:
      containers:
      - name: getter
        image: twibbonize/getter:v1.0
        resources:
          limits:
            memory: "128Mi"
            cpu: "500m"
```

**Demonstrate:**
- Multi-container pods
- Service discovery
- Horizontal Pod Autoscaler
- Rolling updates
- Self-healing (pod restart on crash)

**Tech:** Docker, Kubernetes (minikube or GKE), Helm charts

---

## 9. CI/CD & Automation (Slide 10)

### Project: **Full CI/CD Pipeline**

```yaml
# .github/workflows/deploy.yml
name: Deploy Pipeline

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: go test ./...
      
  build:
    needs: test
    steps:
      - run: docker build -t app:${{ github.sha }} .
      - run: docker push registry/app:${{ github.sha }}
      
  deploy:
    needs: build
    steps:
      - run: kubectl set image deployment/app app=registry/app:${{ github.sha }}
```

**Infrastructure as Code with Pulumi:**

```typescript
const cluster = new aws.eks.Cluster("twibbonize-cluster", {
    vpcConfig: { subnetIds: vpc.subnetIds }
});

const deployment = new k8s.apps.v1.Deployment("api", {
    spec: { replicas: 3 }
});
```

**Tech:** GitHub Actions, Docker, Kubernetes, Pulumi/Terraform

---

## 10. Frontend Rendering (Slide 11)

### Project: **Same App, Three Rendering Strategies**

Build identical "Product Listing" page with:

| Strategy | Framework | Initial Load | Interactivity | SEO |
|----------|-----------|--------------|---------------|-----|
| **CSR** | React SPA | Slow (JS bundle) | Instant | Poor |
| **SSR** | Next.js | Fast (HTML) | After hydration | Good |
| **SSG** | Next.js Static | Instant (CDN) | After hydration | Good |

**Metrics to compare:**
- First Contentful Paint (FCP)
- Time to Interactive (TTI)
- Lighthouse score

**Tech:** React, Next.js, Vercel/Cloudflare Pages

---

## 11. Micro-Frontends (Slide 12)

### Project: **Dashboard with Independent Frontends**

```
┌──────────────────────────────────────────────┐
│                Shell App (React)              │
├──────────────────────────────────────────────┤
│  Nav (Team A)  │        │                    │
│    React       │ Sidebar│   Main Content     │
│                │ (Vue)  │   (Svelte)         │
│                │ Team B │   Team C           │
└────────────────┴────────┴────────────────────┘
```

**Implementation approaches:**
1. **Webpack Module Federation** - Runtime integration
2. **iframe** - Simple isolation (old school)
3. **Web Components** - Framework agnostic

**Communication:**
- Custom events for cross-app messaging
- Shared state via localStorage/Redux

**Tech:** React (shell), Vue (sidebar), Svelte (content), Module Federation

---

## Suggested Learning Order

| Week | Project | Key Concepts |
|------|---------|--------------|
| 1 | URL Shortener | Foundational attributes, caching |
| 2 | E-Commerce Monolith → Microservices | Architecture evolution |
| 3 | Weather Aggregator | Resilience patterns |
| 4 | Order Processing (sync vs async) | Communication patterns |
| 5 | Blog with CQRS | Read/write separation |
| 6 | K8s Deployment | Containers & orchestration |
| 7 | Frontend Rendering Comparison | CSR/SSR/SSG |
| 8 | Micro-Frontends Dashboard | Component composition |

---

## Resources

- **Go:** https://go.dev/doc/
- **Kubernetes:** https://kubernetes.io/docs/tutorials/
- **Next.js:** https://nextjs.org/learn
- **System Design:** https://github.com/donnemartin/system-design-primer
