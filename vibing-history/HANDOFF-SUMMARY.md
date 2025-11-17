# Handoff Summary - Redis Clone Project

**Copy this entire file when starting a new conversation!**

---

## 📍 Current Position

**Week:** 2 - LRU Eviction + TTL  
**Day:** 1 (LRU Data Structure) ✅ COMPLETED  
**Next:** Day 2 - Integrate LRU with Cache  
**Project:** Building Redis from scratch in Go for learning system design

---

## ✅ What's Been Built (Summary)

### Week 1: Core Cache + Server + Load Testing (ALL DONE ✅)

**1. Thread-Safe Cache** (`internal/cache/cache.go`)
- In-memory key-value store with `sync.RWMutex`
- Methods: `Get()`, `Set()`, `Delete()`, `Keys()`, `Flush()`, `Size()`
- Performance: 10M ops/sec, 0 allocations
- All tests pass, no race conditions

**2. TCP Server** (`internal/server/server.go`, `cmd/server/main.go`)
- Listens on port 6378
- 7 commands: SET, GET, DEL, KEYS, SIZE, FLUSH, PING
- Goroutine-per-connection model
- Simple text protocol (not RESP)

**3. Load Testing Tool** (`cmd/loadtest/main.go`)
- Configurable connections, duration, read/write ratio
- Measures throughput + latency percentiles (p50, p95, p99)
- Performance: 75K ops/sec, <1ms latency
- Key finding: Network-bound, not CPU-bound

### Week 2 Progress: LRU Data Structure (JUST COMPLETED ✅)

**4. LRU Doubly-Linked List** (`internal/cache/lru.go`)
```go
type Node struct {
    Key  string  // ⚠️ ONLY stores key, NOT value!
    Prev *Node
    Next *Node
}

type LRUList struct {
    Head *Node  // Most recently used
    Tail *Node  // Least recently used  
    Size int
}
```

**Methods (all O(1)):**
- `AddToFront(key string) *Node` - Add new entry
- `MoveToFront(node *Node)` - Mark as recently used
- `RemoveLRU() string` - Evict least recently used
- `Remove(node *Node)` - Remove specific node

**Tests:** 15/15 passing (`internal/cache/lru_test.go`)
- Empty list, single node, multiple nodes
- All edge cases covered

---

## 🎯 What's Next (Your Current Task)

### Task: Integrate LRU with Cache

**Goal:** Add memory limits and automatic eviction to the cache

**What to Implement:**

1. **Create `CacheEntry` struct:**
```go
type CacheEntry struct {
    Value   string
    LRUNode *Node  // ← Pointer to LRU list position!
}
```

2. **Update `Cache` struct:**
```go
type Cache struct {
    mu      sync.RWMutex
    data    map[string]*CacheEntry  // ← Changed from map[string]string
    lruList *LRUList                // ← New: LRU tracking
    maxSize int                     // ← New: memory limit
}
```

3. **Update Methods:**
- `New()` - Add `maxSize` parameter, initialize `lruList`
- `Set()` - Add to LRU front, evict if `len(data) >= maxSize`
- `Get()` - Move node to LRU front (mark as recently used)
- `Delete()` - Remove from both `data` map AND `lruList`

4. **Add Tests:**
- Test eviction (set maxSize=3, add 4 items, verify oldest evicted)
- Test LRU ordering (access pattern affects eviction order)
- Test GET updates LRU order

---

## 🔑 Key Design Insight

**Why Node doesn't store value:**

```
HashMap (data storage):              LRU List (access order):

┌─────────────────────┐             HEAD (most recent)
│                     │               ↓
│ "key1" → Entry ─────┼─────────→  [key1] ←─┐
│          - Value: "A"│               ↕      │
│          - LRUNode ──┼───────────────────┘  │
│                     │               ↕       │
│ "key2" → Entry ─────┼─────────→  [key2] ←──┤
│          - Value: "B"│               ↕      │
│          - LRUNode ──┼───────────────────┘  │
└─────────────────────┘             TAIL (least recent)

- Value stored ONCE (in HashMap)
- LRU list stores ONLY keys (lightweight)
- CacheEntry has pointer to its LRU node
- O(1) operations everywhere!
```

---

## 📊 Performance Achieved So Far

- **Cache (in-memory):** 10M ops/sec, 0 allocations
- **TCP Server:** 75K ops/sec, 130µs p50 latency
- **Network-bound:** Network I/O is bottleneck, not cache or locks

---

## 🎓 Key Concepts Mastered

**Week 1:**
- ✅ RWMutex vs Mutex (concurrent reads, exclusive writes)
- ✅ Race detector (`go test -race`)
- ✅ Benchmarking (`go test -bench=. -benchmem`)
- ✅ TCP networking in Go (`net` package)
- ✅ Goroutine-per-connection pattern
- ✅ Atomic operations (`sync/atomic`)
- ✅ Load testing methodology
- ✅ Performance bottleneck analysis

**Week 2 (In Progress):**
- ✅ Doubly-linked list pointer manipulation
- ✅ LRU cache algorithm design
- ✅ HashMap + LRU integration pattern
- 🎯 Next: Cache eviction policies

---

## 📁 File Structure

```
redis/
├── cmd/
│   ├── server/main.go          # ✅ Server entry point
│   └── loadtest/main.go        # ✅ Load testing tool
├── internal/
│   ├── cache/
│   │   ├── cache.go            # ✅ Core cache (needs LRU integration)
│   │   ├── cache_test.go       # ✅ Cache tests (needs update)
│   │   ├── lru.go              # ✅ LRU list (just completed!)
│   │   └── lru_test.go         # ✅ LRU tests (15/15 passing)
│   └── server/
│       ├── server.go           # ✅ TCP server
│       └── server_test.go      # ✅ Server tests
├── vibing-history/             # 📚 Documentation
│   ├── action-history.md       # Detailed log of all changes
│   ├── CURRENT-STATUS.md       # Quick status snapshot
│   ├── context-history.md      # Overall project roadmap
│   └── week1-day-by-day.md     # Week 1 guide
├── .cursorrules                # Auto-loaded context
├── go.mod                      # Module: github.com/kartikey-singh/redis
└── README.md                   # Project overview
```

---

## 🧑‍🏫 Teaching Style Preferences

**Student's Learning Style:**
- Writes code independently (don't write code for them!)
- Asks thoughtful design questions
- Appreciates WHY explanations, not just HOW
- Likes thorough code reviews
- Learns from mistakes

**Your Role:**
- ✅ Ask guiding questions
- ✅ Explain concepts and tradeoffs
- ✅ Point out edge cases
- ✅ Give detailed reviews
- ❌ Don't write code unless truly stuck

---

## 🚀 How to Resume in New Chat

**Start new conversation with:**

> "I'm continuing building Redis from scratch in Go. I just completed implementing the LRU data structure (doubly-linked list) with 15 passing tests. I'm ready to integrate LRU with my existing Cache to add memory limits and eviction.
> 
> Current status:
> - Week 2, Day 1 (LRU) completed
> - Files: `internal/cache/lru.go` (done), `internal/cache/cache.go` (needs update)
> - Next task: Add `CacheEntry` struct, update `Cache` to use LRU, implement eviction
> 
> See `vibing-history/CURRENT-STATUS.md` and `vibing-history/action-history.md` for full context.
> 
> Ready to start! What should I think about before modifying `cache.go`?"

---

## 📝 Important Files to Reference

When the new AI responds, it will automatically read:
- `.cursorrules` (project overview, role, context pointers)

You should reference:
- `vibing-history/CURRENT-STATUS.md` (quick status)
- `vibing-history/action-history.md` (detailed log)
- `vibing-history/context-history.md` (8-week roadmap)

---

## 🎯 Immediate Next Questions to Ask AI

1. "Before I start integrating LRU with cache.go, what are the key things I should consider?"

2. "Should I create a new constructor like `NewWithLRU(maxSize int)` or modify the existing `New()` function?"

3. "How do I handle the thread-safety when updating both the map and LRU list? Should they share the same lock?"

4. "What tests should I write first to validate the integration?"

---

**Good luck! You're doing great - your LRU implementation was production-quality! 🚀**

