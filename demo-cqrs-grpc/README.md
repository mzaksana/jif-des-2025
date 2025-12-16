# CQRS + gRPC Blog Service Demo

A demonstration of Command Query Responsibility Segregation (CQRS) pattern with gRPC for the **JIF USK x Twibbonize Modern Software Architecture** workshop.

## Architecture

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

## Concepts

- **Command Service**: Handles CREATE, UPDATE, DELETE (writes to SQLite)
- **Query Service**: Handles GET, LIST, SEARCH (reads from in-memory store)
- **Event Bus**: Syncs data from write store to read store
- **gRPC**: Fast binary protocol for service communication

## Prerequisites

```bash
# Go 1.21+
go version

# Protocol Buffer compiler
protoc --version

# Go gRPC plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# grpcurl for testing
brew install grpcurl  # macOS
# or: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

## Quick Start

```bash
# 1. Generate gRPC code
make proto

# 2. Run the server
make run

# 3. In another terminal, run the demo
make demo
```

## Manual Testing with grpcurl

### Create a Post (Command)
```bash
grpcurl -plaintext -d '{
    "title": "My First Post",
    "content": "Hello CQRS World!",
    "author": "Workshop",
    "tags": ["demo", "cqrs"]
}' localhost:50051 blog.CommandService/CreatePost
```

### List Posts (Query)
```bash
grpcurl -plaintext -d '{"limit": 10}' localhost:50051 blog.QueryService/ListPosts
```

### Search Posts (Query)
```bash
# By tag
grpcurl -plaintext -d '{"tags": ["demo"]}' localhost:50051 blog.QueryService/SearchPosts

# By query
grpcurl -plaintext -d '{"query": "hello"}' localhost:50051 blog.QueryService/SearchPosts
```

### Update a Post (Command)
```bash
grpcurl -plaintext -d '{
    "id": "<post-id>",
    "title": "Updated Title",
    "content": "Updated content",
    "tags": ["updated"]
}' localhost:50051 blog.CommandService/UpdatePost
```

### Delete a Post (Command)
```bash
grpcurl -plaintext -d '{"id": "<post-id>"}' localhost:50051 blog.CommandService/DeletePost
```

## Project Structure

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
├── Makefile
├── demo.sh                 # Demo script
└── README.md
```

## Key Demo Points

1. **Writes go to SQLite** (normalized, ACID compliant)
2. **Reads come from in-memory** (fast, denormalized)
3. **Event sync** keeps both stores consistent
4. **Separation of concerns** - commands and queries use different data stores optimized for their purpose

## Why CQRS?

| Aspect | Traditional | CQRS |
|--------|-------------|------|
| Read/Write | Same model | Separate models |
| Optimization | Compromise | Optimized per use case |
| Scaling | Together | Independent |
| Complexity | Lower | Higher (worth it at scale) |

## Production Considerations

In a real-world scenario:
- Replace in-memory read store with **Redis** or **Elasticsearch**
- Use **Kafka** or **RabbitMQ** for event bus
- Add **event sourcing** for full audit trail
- Implement **eventual consistency** handling
