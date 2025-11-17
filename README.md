# Build Your Own Redis

A learning project to deeply understand distributed systems, caching, and system design by building Redis from scratch in Go.

## 🎯 Learning Goals

- Master distributed systems concepts (replication, sharding, consensus)
- Understand cache internals (LRU, TTL, eviction policies)
- Learn concurrency patterns in Go
- Prepare for system design interviews
- Build production-grade distributed system knowledge

## 📚 Documentation

**Start here:**
- **`START-HERE.md`** - Begin your journey (read this first!)
- **`vibing-history/week1-day-by-day.md`** - Daily implementation guide

**Reference:**
- `vibing-history/context-history.md` - Full 8-week roadmap & concepts
- `vibing-history/collaboration-strategy.md` - How to use Cursor vs Browser Claude
- `vibing-history/action-history.md` - Track your progress

## 🗓️ Roadmap

### Week 1-2: Single Node Cache + LRU/TTL
- Core cache with RWMutex concurrency
- LRU algorithm implementation
- TTL with passive + active expiration
- TCP server with simple protocol
- Memory limits and eviction policies

### Week 3-4: Replication & Fault Tolerance
- Master-slave replication
- Asynchronous replication protocol
- Failure detection and manual failover
- Replication lag handling

### Week 5-6: Sharding & Consistent Hashing
- Consistent hashing algorithm
- Virtual nodes
- Data partitioning
- Rebalancing strategies

### Week 7-8: Raft Consensus
- Leader election
- Automatic failover
- Strong consistency guarantees
- Split-brain prevention

## 🏗️ Project Structure

```
redis/
├── vibing-history/          # Documentation and guides
├── cmd/
│   └── server/             # Server entry point
├── internal/
│   ├── cache/              # Core cache logic
│   ├── protocol/           # Protocol parser/responder
│   └── server/             # TCP server
├── pkg/
│   └── client/             # Client library
└── tests/                  # Integration tests
```

## 🚀 Getting Started

```bash
# 1. Read the starting guide
cat START-HERE.md

# 2. Read Day 1 detailed instructions
cat vibing-history/week1-day-by-day.md

# 3. Create your first file
touch internal/cache/cache.go

# 4. Start coding!
```

## 🛠️ Development

```bash
# Run tests
go test ./...

# Run with race detector
go test -race ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# Run server
go run cmd/server/main.go
```

## 📖 What You'll Learn

By building this, you'll understand how these real-world systems work:
- Redis, Memcached (caching)
- DynamoDB, Cassandra (consistent hashing, eventual consistency)
- etcd, Consul (Raft consensus)
- MongoDB (replica sets with auto-failover)
- CDNs (consistent hashing for routing)

## 🎓 Interview Prep

This project covers common interview topics:
- LRU cache implementation
- Concurrency and thread safety
- TCP server design
- Distributed systems concepts
- CAP theorem (practical application)
- Consensus algorithms
- Sharding strategies


