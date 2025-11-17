# Week-by-Week Implementation Guide

## Week 1-2: Single Node Cache + LRU/TTL

### Learning Objectives
- Master Go's concurrency primitives (RWMutex, goroutines)
- Understand TCP server implementation
- Implement LRU cache algorithm from scratch
- Build TTL expiration system
- Handle memory eviction policies

### Implementation Checklist

#### Week 1: Core Cache + TCP Server
- [ ] Project setup (Go module, directory structure)
- [ ] Core cache data structure with RWMutex
- [ ] Basic operations: GET, SET, DELETE
- [ ] TCP server setup
- [ ] Simple text-based protocol parser
- [ ] Goroutine per connection model
- [ ] Basic command handling

#### Week 2: LRU + TTL + Eviction
- [ ] LRU implementation (doubly-linked list + HashMap)
- [ ] TTL support (passive expiration)
- [ ] Background TTL cleanup (active expiration)
- [ ] Memory tracking
- [ ] Eviction policies:
  - [ ] allkeys-lru
  - [ ] volatile-lru
  - [ ] volatile-ttl
- [ ] Testing & benchmarking

### Key Concepts to Master

#### 1. Concurrency Control
```go
// Why RWMutex?
// - Multiple readers can read simultaneously
// - Only one writer at a time
// - Writers block readers and other writers

type Cache struct {
    mu    sync.RWMutex
    data  map[string]*CacheEntry
}

// Read operation (allows concurrent reads)
func (c *Cache) Get(key string) {
    c.mu.RLock()         // Read lock
    defer c.mu.RUnlock()
    // ... read data
}

// Write operation (exclusive)
func (c *Cache) Set(key string, value string) {
    c.mu.Lock()          // Write lock
    defer c.mu.Unlock()
    // ... write data
}
```

#### 2. LRU Algorithm
**Why doubly-linked list + HashMap?**
- HashMap: O(1) key lookup
- Doubly-linked list: O(1) move to front, O(1) remove from tail
- Combined: O(1) all operations!

**Structure:**
```
HashMap: key -> *Node
        
Linked List (access order):
[Most Recent] <-> node <-> node <-> node <-> [Least Recent]
    HEAD                                           TAIL

On access: Move node to HEAD
On eviction: Remove from TAIL
```

#### 3. TTL Strategies
**Passive Expiration:**
- Check expiry when key is accessed
- Pros: No CPU overhead
- Cons: Memory held until access

**Active Expiration:**
- Background goroutine scans periodically
- Pros: Timely memory release
- Cons: CPU overhead

**Hybrid (Redis approach):**
- Passive check on every access
- Active scan: Sample random keys, clean expired ones

#### 4. TCP Protocol Design
**Simple text protocol:**
```
Client -> Server:  GET key\r\n
Server -> Client:  +OK value\r\n
                   OR
                   -ERR not found\r\n

Client -> Server:  SET key value\r\n
Server -> Client:  +OK\r\n

Client -> Server:  SET key value EX 60\r\n  (60 sec TTL)
Server -> Client:  +OK\r\n
```

### Testing Strategy

#### Unit Tests
1. **Cache operations**: Test GET, SET, DELETE with concurrent access
2. **LRU behavior**: Verify correct eviction order
3. **TTL expiration**: Test passive and active expiration
4. **Race conditions**: Run with `go test -race`

#### Integration Tests
1. **TCP client-server**: Connect multiple clients simultaneously
2. **Load testing**: 10,000 concurrent operations
3. **Memory limits**: Verify eviction when limit reached

#### Benchmarks
```bash
go test -bench=. -benchmem
```
Target performance:
- SET: < 1ms per operation
- GET: < 0.5ms per operation
- 100,000 ops/sec on single core

### Directory Structure
```
redis/
├── vibing-history/           # Documentation
├── cmd/
│   └── server/
│       └── main.go          # Server entry point
├── internal/
│   ├── cache/
│   │   ├── cache.go         # Core cache with RWMutex
│   │   ├── lru.go           # LRU implementation
│   │   ├── ttl.go           # TTL manager
│   │   └── eviction.go      # Eviction policies
│   ├── protocol/
│   │   ├── parser.go        # Command parser
│   │   └── responder.go     # Response formatter
│   └── server/
│       └── tcp.go           # TCP server
├── pkg/
│   └── client/
│       └── client.go        # Client library for testing
└── tests/
    ├── cache_test.go
    ├── lru_test.go
    └── integration_test.go
```

### Common Pitfalls to Avoid
1. **Deadlocks**: Always use defer for unlocks
2. **Goroutine leaks**: Ensure all goroutines can exit
3. **Memory leaks**: Properly clean up expired keys
4. **Race conditions**: Use `go test -race` frequently

### Interview Questions You'll Be Able to Answer
- How does LRU cache work? Implement it.
- How would you implement a cache with TTL?
- Explain the difference between RWMutex and Mutex
- How do you handle concurrent connections in a TCP server?
- What are different cache eviction policies?
- How would you test a concurrent system?

---

## Week 3-4: Replication & Fault Tolerance
[To be implemented after Week 1-2 completion]

## Week 5-6: Sharding & Consistent Hashing
[To be implemented after Week 3-4 completion]

## Week 7-8: Raft Consensus
[To be implemented after Week 5-6 completion]





