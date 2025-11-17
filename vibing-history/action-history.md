# Action History - Build Your Own Redis

This file tracks all significant actions performed during the development of this Redis implementation.

Format: `<Number> | <Files Changed> | <Summary of Action> | <Purpose of Action>`

---

## Setup Phase

0 | vibing-history/*.md, README.md, START-HERE.md, go.mod | Created project structure and comprehensive documentation framework | Establish learning environment with day-by-day guides, collaboration strategies, and context tracking for 8-week Redis implementation journey

---

## Week 1: Single Node Cache + LRU/TTL

### Day 1: Thread-Safe Cache (✅ COMPLETED)

1 | internal/cache/cache.go | Implemented thread-safe cache with Get/Set/Delete operations | Create foundation for in-memory key-value store with concurrent access support using sync.RWMutex

2 | internal/cache/cache_test.go | Added comprehensive unit tests and benchmarks | Ensure cache correctness and thread-safety, validate RWMutex prevents race conditions, measure performance

**Tests Added:**
- TestBasicOperations (set, get, delete)
- TestGetNonExistent (edge cases)
- TestConcurrentOperations (100 goroutines writing)
- TestConcurrentReadWrite (50 readers + 50 writers)

**Benchmarks:**
- BenchmarkCacheGet: 122ns/op, 8.2M ops/sec
- BenchmarkCacheSet: 79ns/op, 12.6M ops/sec
- BenchmarkCacheMixed: 100ns/op, 10M ops/sec (80% read, 20% write)

**Performance:** 10M ops/sec mixed workload, 0 allocations, thread-safe

**Key Learnings:**
- RWMutex allows concurrent reads, exclusive writes
- WaitGroup coordinates concurrent goroutines
- defer ensures locks are released (prevents deadlocks)
- Race detector catches concurrency bugs
- Benchmarking reveals real performance characteristics
- Zero allocations = no GC pressure

---

### Day 2: TCP Server (✅ COMPLETED)

3 | internal/cache/cache.go | Added Keys(), Flush(), Size() methods | Support cache inspection and management operations

4 | internal/server/server.go | Built TCP server with 7 commands (SET, GET, DEL, KEYS, SIZE, FLUSH, PING) + logging | Enable network access to cache, support multiple concurrent clients, add observability

5 | internal/server/server_test.go | Created 8 comprehensive server tests | Ensure server correctness, test concurrent connections, validate error handling

6 | cmd/server/main.go | Improved startup with user-friendly messages | Better developer experience

7 | test-server.sh | Created automated manual testing script | Quick validation of all features

**Features Implemented:**
- TCP server on port 6378 with goroutine-per-connection
- 7 commands: SET, GET, DEL, KEYS, SIZE, FLUSH, PING
- Connection logging (new connection, commands, close)
- Error handling with proper error messages
- Support for values with spaces
- Case-insensitive commands
- Comprehensive test coverage

**Test Results:**
- All server tests pass (8 tests)
- No race conditions
- Handles concurrent connections correctly

**Key Learnings:**
- TCP networking in Go (net package)
- Goroutine-per-connection pattern
- Protocol design decisions
- Connection lifecycle management
- Integration testing strategies
- Server logging for observability

---

### Day 3: Load Testing & Performance Analysis (✅ COMPLETED)

8 | cmd/loadtest/main.go | Built professional load testing tool with configurable workloads | Measure throughput, latency percentiles, and identify performance bottlenecks

**Features Implemented:**
- Configurable parameters (connections, duration, read/write ratio)
- Atomic operations for thread-safe stats
- Latency percentile calculation (p50, p95, p99)
- Multiple concurrent connections
- Random key distribution

**Performance Results:**
- Single connection: 38K ops/sec, 23µs p50 latency
- 10 connections: 74K ops/sec, 128µs p50 latency
- Peak throughput: 78K ops/sec (50 connections)
- 100% reads vs 100% writes: Nearly identical performance!
- 0 errors across all tests

**Key Discoveries:**
- Redis is network-bound, not CPU-bound
- Bottleneck is network I/O (~100µs per op), not the cache (~5µs)
- RWMutex lock contention is minimal (writes aren't slower than reads)
- Performance plateaus at ~75K ops/sec due to network protocol
- Latency grows with connections (goroutine scheduling overhead)
- Performing at 75% of production Redis speed!

**Key Learnings:**
- How to measure throughput and latency
- Understanding performance bottlenecks
- Network I/O dominates in distributed systems
- Importance of percentiles (p50, p95, p99) over averages
- Lock-free isn't always better (network is the real bottleneck)

---

## Week 2: LRU Eviction + TTL

### Day 1 (Week 2): LRU Data Structure (✅ COMPLETED)

9 | internal/cache/lru.go | Implemented doubly-linked list LRU data structure with Node and LRUList | Build foundation for cache eviction policy with O(1) add, move, and remove operations

10 | internal/cache/lru_test.go | Wrote comprehensive test suite for LRU (15 tests covering all operations and edge cases) | Ensure correctness of doubly-linked list pointer manipulation, test all scenarios (empty, single node, multiple nodes)

**Features Implemented:**
- `Node` struct with Key, Prev, Next pointers (stores only key, not value)
- `LRUList` struct with Head, Tail, Size tracking
- `AddToFront(key)` - O(1) insertion at head
- `MoveToFront(node)` - O(1) promotion (for cache hits)
- `RemoveLRU()` - O(1) eviction from tail
- `Remove(node)` - O(1) arbitrary removal (for DELETE operations)

**Test Coverage:**
- AddToFront: Empty list, multiple nodes, pointer integrity
- RemoveLRU: Empty list, single node, multiple nodes
- MoveToFront: Already at head, from tail, from middle, two nodes
- Remove: Single node, head, tail, middle
- Complex scenario: Mix of all operations

**All Tests Pass:** 15/15 tests ✅

**Key Design Decisions:**
- Node stores ONLY the key (not value) to avoid memory duplication
- Value stored in HashMap, Node stored in LRU list
- Bidirectional pointers enable O(1) operations
- CacheEntry will hold both value and pointer to LRU node

**Key Learnings:**
- Doubly-linked list pointer manipulation
- Importance of edge case testing (empty, single node)
- How LRU and HashMap integrate via pointers
- O(1) cache operations with proper data structure design


