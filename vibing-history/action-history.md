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

---

### Day 2 (Week 2): Integrate LRU with Cache (✅ COMPLETED)

11 | internal/cache/cache.go | Integrated LRU list with Cache for memory limits and automatic eviction | Enable O(1) cache operations with LRU eviction policy when maxSize is reached

12 | internal/cache/cache_test.go | Added 4 comprehensive LRU eviction tests | Ensure eviction correctness, test LRU ordering, verify Get() affects eviction order

13 | cmd/server/main.go, internal/server/server_test.go | Updated cache.New() calls to pass maxSize parameter | Fix compilation after cache constructor signature change

**Features Implemented:**
- `CacheEntry` struct linking values to LRU nodes
- Updated `Cache` struct with `lruList *LRUList` and `maxSize int` fields
- Constructor `New(maxSize int)` with validation (panic on negative, 0 = unlimited)
- `Set()` with eviction logic: checks for existing key, evicts if full, then adds
- `Get()` uses write lock and moves node to front (marks as recently used)
- `Delete()` removes from both HashMap and LRU list
- Defensive panic checks for data structure corruption

**Test Coverage:**
- TestCacheWithLRU_BasicEviction: Verifies oldest item evicted when full
- TestCacheWithLRU_GetAffectsEviction: Confirms Get() updates LRU order
- TestCacheWithLRU_UpdateDoesntEvict: Ensures updating existing key doesn't evict
- TestCacheWithLRU_UpdateExistingKey: Validates Set() on existing key behavior
- All 22 tests passing (cache + server + LRU)
- 0 race conditions detected

**Performance Results:**
- Get: 32ns/op (down from 120ns despite write lock!)
- Set: 37ns/op
- Mixed (80/20 read/write): 89ns/op
- Still 0 allocations per operation

**Key Design Decisions:**
- Single shared lock (prevents HashMap/LRU drift)
- Write lock in Get() for accurate LRU tracking (correctness over speed)
- HashMap → LRU pointer direction (enables O(1) MoveToFront)
- Evict-then-add pattern (prevents temporary size overflow)
- maxSize = 0 means unlimited (explicit design choice)

**Key Learnings:**
- Why pointer direction matters in data structure integration
- Thread safety trade-offs: one lock vs two locks
- Eviction timing: check existing key before evicting
- Write lock in Get() is necessary for accurate LRU
- Defensive programming with panic checks catches bugs early
- Performance can improve with better memory locality

---

### Day 3 (Week 2): TTL (Time-To-Live) Implementation (✅ COMPLETED)

14 | internal/cache/cache.go | Added TTL expiration support with passive and active cleanup | Enable keys to expire automatically, freeing memory for expired entries

15 | internal/cache/cache_test.go | Added 7 comprehensive TTL expiration tests | Ensure TTL correctness, test passive expiration, verify expired keys evicted before LRU

16 | internal/server/server.go | Updated SET command to support EX seconds parameter | Enable clients to set TTL via network protocol: SET key value EX 60

**Features Implemented:**
- `ExpiryTime time.Time` field in `CacheEntry` (zero value = never expires)
- `SetWithTTL(key, value, ttl)` method for setting keys with expiration
- `Set(key, value)` calls `SetWithTTL` with ttl=0 (backward compatible)
- Passive expiration: `Get()` checks and deletes expired keys on access
- Passive expiration: `Keys()` filters out expired keys
- Active expiration: Background goroutine runs every 10 seconds
- `cleanupExpiredKeys()` scans all keys and removes expired ones
- `Close()` method to gracefully stop background cleanup goroutine
- `deleteWithoutLocking()` helper to prevent deadlocks
- Server protocol: `SET key value EX seconds` parses and sets TTL
- Expired keys evicted before LRU eviction (optimization)

**Test Coverage:**
- TestTTL_BasicExpiration: Keys expire after TTL duration
- TestTTL_SetWithoutTTL: Keys without TTL don't expire
- TestTTL_SetWithTTLAndUpdateWithoutTTLClearsTTL: Update clears TTL
- TestTTL_UpdateWithNewTTL: Update with new TTL replaces old TTL
- TestTTL_KeysFiltersExpired: Keys() doesn't return expired keys
- TestTTL_ExpiredEvictedBeforeLRU: Expired keys evicted before LRU
- TestTTL_MultipleExpiredEvicted: Multiple expired keys handled correctly
- TestServerTTL: Server protocol SET key value EX seconds works
- All 30 tests passing (cache + server + LRU + TTL)
- 0 race conditions detected

**Key Design Decisions:**
- `time.Time` zero value for "never expires" (no allocations)
- Passive expiration in Get() and Keys() (lazy cleanup)
- Active expiration every 10 seconds (proactive cleanup)
- Background goroutine with ticker pattern (standard Go concurrency)
- Channel-based graceful shutdown (`stopCleanup chan struct{}`)
- Expired keys checked before LRU eviction (memory optimization)
- Update existing key clears TTL if ttl=0, sets new TTL if ttl>0
- Server protocol: EX parameter at end, value can contain spaces

**Key Learnings:**
- `time.Ticker` for periodic background tasks
- Goroutines and channels for concurrent operations
- Select statement for listening to multiple channels
- Graceful shutdown pattern with channels
- Deadlock prevention: internal methods without locking
- Zero value semantics in Go (time.Time{} = never expires)
- Protocol parsing: handling optional parameters
- Two-phase expiration: passive (on access) + active (background)
- How Redis implements expiration (similar pattern!)

---

## Week 3: Replication & Fault Tolerance

### Day 1 (Week 3): Master-Slave Replication (✅ COMPLETED)

17 | internal/replication/protocol.go | Designed and implemented replication protocol with Operation struct and serialization | Enable master-slave communication with text-based protocol supporting SET, DELETE, FLUSH operations with TTL and timestamps

18 | internal/replication/protocol_test.go | Added comprehensive protocol tests for serialization/deserialization | Ensure protocol correctness, test all operation types, validate error handling

19 | internal/replication/master.go | Implemented Master node with asynchronous broadcasting to slaves | Enable master to replicate all write operations (SET, DELETE, FLUSH) to connected slaves

20 | internal/replication/slave.go | Implemented Slave node with replication receiver and TTL compensation | Enable slave to receive operations from master, apply them in order, and compensate for replication lag in TTL

21 | internal/replication/replication_test.go | Added integration tests for master-slave replication | Verify end-to-end replication works, test SET, DELETE, FLUSH, TTL expiration with lag compensation

22 | cmd/server/main.go | Added command-line flags for replication mode (--role, --port, --replication-port, --master) | Enable running server as master, slave, or standalone mode

23 | internal/server/server.go | Integrated Master and Slave into server with role-based command handling | Enable server to operate in master/slave/standalone modes, enforce read-only on slaves, route commands appropriately

24 | internal/server/server_test.go | Fixed server.New() calls to match new signature | Update tests after server constructor changes

**Features Implemented:**
- Replication protocol: Text-based with Operation struct (Type, Key, Value, TTL, Timestamp)
- Master node: Wraps cache operations, broadcasts to all slaves asynchronously
- Slave node: Connects to master, receives operations, applies in order
- TTL compensation: Slave calculates remaining TTL after replication lag
- Server integration: Three modes (master, slave, standalone) with command-line flags
- Read-only slaves: Rejects SET, DELETE, FLUSH from clients (writes only via replication)
- Thread-safe slave management: RWMutex for slave list, Mutex for per-slave writer
- Graceful error handling: Auto-removes disconnected slaves, logs errors

**Test Coverage:**
- Protocol tests: SET (with/without TTL), DELETE, FLUSH, PING serialization/deserialization
- Protocol error tests: Invalid formats, missing fields, invalid TTL/timestamp
- Replication integration: Master-slave SET, DELETE, FLUSH, TTL expiration
- Multiple slaves test: Master broadcasts to multiple slaves correctly
- All 47 tests passing (cache + server + replication)
- 0 race conditions detected with race detector

**End-to-End Testing Results:**
- Master on port 6379, slave on port 6378
- SET operations replicated correctly ✅
- DELETE operations replicated correctly ✅
- FLUSH operations replicated correctly ✅
- TTL expiration with replication lag compensation ✅
- Read-only enforcement on slaves (SET/DEL/FLUSH rejected) ✅

**Key Design Decisions:**
- Millisecond precision for TTL (prevents truncation: int64(0.2 seconds) = 0 bug)
- Synchronous operation application on slave (no `go s.apply()` - maintains order)
- Asynchronous broadcasting on master (non-blocking, spawns goroutine per slave)
- Thread-safe writer per slave connection (bufio.Writer not thread-safe)
- Copy slaves slice under read lock before iterating (avoids holding lock too long)
- TTL compensation: `remaining = op.TTL - elapsed` to account for network delay
- Skip expired keys on slave (don't apply operations for keys already expired)

**Key Learnings:**
- Precision loss bugs: Always use smallest unit (milliseconds) for time serialization
- Thread safety: Most Go types (bufio.Writer, slices, maps) need explicit mutex protection
- Race detector catches concurrency bugs that are hard to find manually
- Async replication challenges: Must apply operations in order, compensate for lag
- Design patterns: Master wraps writes + broadcasts, reads go direct to cache
- Error handling: Check connection errors before starting replication, log but don't crash
- Scanner error checking: Must be outside the loop, not inside
- Read-only slaves: Standard pattern for master-slave replication (single source of truth)


