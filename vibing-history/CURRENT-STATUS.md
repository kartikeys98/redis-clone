# Current Status - Quick Context for New Conversations

**Last Updated:** Week 2, Day 2 (LRU Integration) Completed

---

## 🎯 Where We Are

**Current Week:** Week 2 - LRU Eviction + TTL  
**Current Day:** Day 2 ✅ LRU + Cache Integration COMPLETED - Ready for TTL!  
**Role:** Instructor mode - Student codes, I guide and review

---

## ✅ What's Been Completed

### Week 1 Summary (ALL COMPLETED ✅)

**Day 1: Thread-Safe Cache**
- `internal/cache/cache.go` - Core cache with RWMutex
- `internal/cache/cache_test.go` - Unit tests + benchmarks
- Performance: 10M ops/sec mixed workload, 0 allocations

**Day 2: TCP Server**
- `internal/server/server.go` - TCP server with 7 commands
- Commands: SET, GET, DEL, KEYS, SIZE, FLUSH, PING
- `cmd/server/main.go` - Server entry point
- Runs on port 6378

**Day 3: Load Testing**
- `cmd/loadtest/main.go` - Professional load testing tool
- Performance: 75K ops/sec, sub-millisecond latency
- Key finding: Network-bound, not CPU-bound

### Week 2 Progress (IN PROGRESS)

**Day 1: LRU Data Structure ✅ COMPLETED**

**Files:**
- `internal/cache/lru.go` - Doubly-linked list implementation
- `internal/cache/lru_test.go` - 15 comprehensive tests

**Implementation:**
```go
type Node struct {
    Key  string  // Only stores key, not value!
    Prev *Node
    Next *Node
}

type LRUList struct {
    Head *Node  // Most recently used
    Tail *Node  // Least recently used
    Size int
}

// Methods: AddToFront(), MoveToFront(), RemoveLRU(), Remove()
```

**Design Decision (CRITICAL):**
- Node stores ONLY the key (not value)
- Value stored in HashMap, LRU list tracks access order
- CacheEntry struct will hold: Value + pointer to LRU Node
- This avoids memory duplication and enables O(1) operations

**Tests:**
- ✅ 15/15 tests pass
- ✅ AddToFront (empty, multiple nodes)
- ✅ RemoveLRU (empty, single, multiple)
- ✅ MoveToFront (head, tail, middle, two nodes)
- ✅ Remove (single, head, tail, middle)
- ✅ Complex scenario (mix of operations)

**Key Concepts Learned:**
- Doubly-linked list pointer manipulation
- How LRU + HashMap integrate via pointers
- Importance of edge case testing
- O(1) cache operations with proper design

**Day 2: LRU + Cache Integration ✅ COMPLETED**

**Files:**
- `internal/cache/cache.go` - Updated with LRU integration
- `internal/cache/cache_test.go` - Added 4 LRU eviction tests
- `cmd/server/main.go` - Updated cache.New(10000)
- `internal/server/server_test.go` - Updated cache.New(1000)

**Implementation:**
```go
type CacheEntry struct {
    Value   string
    lruNode *Node  // Pointer to LRU list position
}

type Cache struct {
    data    map[string]*CacheEntry
    lock    sync.RWMutex
    maxSize int       // 0 = unlimited
    lruList *LRUList
}

func New(maxSize int) *Cache {
    if maxSize < 0 {
        panic("cache: maxSize cannot be negative")
    }
    // maxSize = 0 means unlimited
    return &Cache{...}
}
```

**Key Operations:**
- `Set()`: Checks for existing key, evicts if full, then adds to front
- `Get()`: Uses write lock, moves node to front (marks as recently used)
- `Delete()`: Removes from both HashMap and LRU list

**Tests:**
- ✅ 4 LRU eviction tests (basic, Get affects order, update doesn't evict)
- ✅ 22/22 tests passing across entire project
- ✅ 0 race conditions
- ✅ Server integration tests passing

**Performance:**
- Get: 32ns/op (0 allocs)
- Set: 37ns/op (0 allocs)
- Mixed: 89ns/op

**Key Concepts Learned:**
- HashMap → LRU pointer direction for O(1) operations
- Single shared lock prevents data structure drift
- Write lock in Get() maintains accurate LRU order
- Eviction timing: check existing before evicting
- Defensive programming with panic checks

---

## 🚀 What's Next

### Day 3 (Week 2): TTL (Time-To-Live) Implementation
**Goal:** Add expiration timestamps and automatic cleanup for keys

**What to implement:**
1. Add expiration to `CacheEntry`:
```go
type CacheEntry struct {
    Value      string
    lruNode    *Node
    ExpiresAt  time.Time  // New: expiration timestamp
}
```

2. Update `Set()` to support TTL:
```go
// New method signature
func (c *Cache) SetWithTTL(key, value string, ttl time.Duration) {
    // Set entry with ExpiresAt = time.Now().Add(ttl)
}
```

3. Implement passive expiration:
- `Get()`: Check if key expired, return nil and delete if so
- Lazy cleanup on access

4. Implement active expiration:
- Background goroutine that periodically scans
- Removes expired keys proactively
- Use ticker for periodic cleanup

5. Update server protocol:
- `SET key value EX seconds` - set with TTL
- Backward compatible with existing `SET key value`

6. Add tests:
- Keys expire after TTL
- Passive expiration (Get returns nil for expired)
- Active cleanup runs in background
- TTL works with LRU eviction

**This combines two expiration strategies used by Redis!**

---

## 📊 Student's Learning Style

**Strengths:**
- Writes code independently
- Asks for reviews
- Implements tests proactively
- Good with concurrency concepts (WaitGroup usage was correct!)

**Areas of Growth:**
- Small syntax errors (t.Error vs t.Errorf) - normal for learning!
- Could explore more edge cases

**Teaching Approach:**
- Guide with questions, not direct answers
- Explain WHY behind decisions
- Point out edge cases to consider
- Give detailed code reviews
- Encourage experimentation

---

## 🔧 Technical Context

**Project:** Redis clone in Go  
**Module:** `github.com/kartikey-singh/redis`  
**Go Version:** 1.24.3

**Directory Structure:**
```
redis/
├── cmd/server/          # TODO: Day 2
├── internal/
│   ├── cache/          # ✅ Day 1 - Complete
│   ├── protocol/       # TODO: Day 2
│   └── server/         # TODO: Day 2
├── pkg/client/         # TODO: Later weeks
├── tests/              # TODO: Integration tests
└── vibing-history/     # Documentation
```

---

## 💬 Recent Discussions

**Topics covered:**
- Why use organized project structure vs single main.go
- Why write unit tests vs only integration tests
- Difference between RWMutex and Mutex
- How go test works
- WaitGroup usage for concurrent testing

**Student questions to remember:**
- Thinking about context management (why we created this file!)
- Wants to understand tradeoffs, not just implementations

---

## 🎓 Interview Prep Goals

**Already Mastered:**
- ✅ Thread-safe data structures
- ✅ RWMutex usage
- ✅ Concurrent testing
- ✅ Race condition detection

**Coming Up:**
- TCP server architecture (Day 2)
- Protocol design (Day 2)
- Load testing (Day 4-5)
- LRU algorithm (Week 2)
- Distributed systems (Weeks 3-8)

---

## 📝 Notes for Next Conversation

**If context window resets, start with:**
1. Read `.cursorrules` (automatic)
2. Read this file (CURRENT-STATUS.md)
3. Read `action-history.md` for detailed log
4. Ask student: "Ready to continue with Day 2?"

**Don't re-explain:**
- Project structure (already understood)
- Unit testing value (already convinced)
- Basic Go concepts (student is competent)

**Do focus on:**
- New concepts for current task
- Code review and feedback
- Edge cases and best practices
- System design thinking

