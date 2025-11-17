# 🚀 START HERE - Your Redis Learning Journey

## What You're Building
A distributed Redis implementation from scratch to master:
- Distributed systems concepts
- System design patterns  
- Production-grade code skills
- Interview preparation

## 📁 Documentation Structure (Simple!)

**You need 3 files:**
1. **This file (START-HERE.md)** - Quick start guide ← You're reading it!
2. **README.md** - Project overview (for GitHub)
3. **vibing-history/week1-day-by-day.md** - Detailed daily instructions

**Reference docs** (read when needed):
- `vibing-history/collaboration-strategy.md` - Cursor vs Browser Claude tips
- `vibing-history/context-history.md` - Full 8-week concepts
- `vibing-history/action-history.md` - Your progress log

## Your First Task (Today - 2-4 hours)

### Day 1: Build a Thread-Safe Cache

#### Step 1: Read the detailed guide (10 min)
```bash
cat vibing-history/week1-day-by-day.md
# Focus on "Day 1: Core Cache Data Structure"
```

#### Step 2: Think Before Coding (10 min)
Answer these questions:
1. What Go data structure stores key-value pairs? → **map[string]string**
2. Why isn't it thread-safe? → **Concurrent access causes panics**
3. What makes it thread-safe? → **sync.RWMutex**
4. What should the Cache struct contain? → **map + RWMutex**

#### Step 3: Create Files & Start Coding
```bash
# Create your files
touch internal/cache/cache.go
touch internal/cache/cache_test.go
```

Implement in `cache.go`:
```go
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

func New() *Cache { ... }
func (c *Cache) Get(key string) (string, bool) { ... }
func (c *Cache) Set(key string, value string) { ... }
func (c *Cache) Delete(key string) bool { ... }
```

#### Step 4: Test
```bash
go test ./internal/cache -v
go test -race ./internal/cache -v  # Check for race conditions
```

#### When to Ask Me:
- "Should I use RWMutex or Mutex? What's the difference?"
- "Here's my code [paste]. Can you review it?"
- "Getting this error: [paste error]"
- "Why do I need defer for unlocking?"

## How We'll Work Together

### My Role (Cursor Claude):
- ✅ Answer your design questions
- ✅ Review your code  
- ✅ Help debug issues
- ✅ Guide you to solutions
- ❌ Won't write code for you (that's your job!)

### Your Role:
- 🛠️ Write the code yourself
- 🤔 Think through design decisions
- 🐛 Debug and iterate
- 💬 Ask questions when stuck
- 📊 Share work for review

## The Path Forward

```
Week 1: Single-node cache (Days 1-7)
  ├─ Day 1: Thread-safe cache ← YOU ARE HERE
  ├─ Day 2: TCP server
  ├─ Day 3: Testing & benchmarking
  ├─ Day 4-5: Load testing tool
  └─ Day 6-7: Polish & validation

Week 2: LRU eviction & TTL
Week 3-4: Master-Slave replication  
Week 5-6: Consistent hashing & sharding
Week 7-8: Raft consensus & leader election
```

## Context Window Strategy

### For Cursor (Me):
- We have ~100K tokens (plenty for weeks of work!)
- I'll warn you if we approach the limit
- Can reference your vibing-history docs for continuity

### For Browser Claude:
- Use for deep conceptual discussions
- Better for "why" questions vs "how" implementation
- Can use Projects feature for persistent memory
- See `collaboration-strategy.md` for details

## Quick Commands Reference

```bash
# Run tests
go test ./...

# Check for race conditions  
go test -race ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# Format code
go fmt ./...

# Check for issues
go vet ./...
```

## Success Metrics

### Day 1 Complete When:
- [ ] cache.go implemented with Get/Set/Delete
- [ ] Tests pass: `go test ./internal/cache`
- [ ] No race conditions: `go test -race ./internal/cache`
- [ ] Ready for code review

### Week 1 Complete When:
- [ ] TCP server running on port 6379
- [ ] Multiple clients can connect simultaneously
- [ ] Load tester shows >50K ops/sec
- [ ] All tests pass, no race conditions

## Need Help?

### Good Questions:
> "I'm implementing X. Should I use approach A or B? Here's what I've tried..."

> "Here's my code [paste]. Can you review for correctness/performance?"

> "Getting this error [paste]. I think it's related to Y. What am I missing?"

### Tips:
- Show your code when asking for reviews
- Explain what you've tried when stuck
- Ask "why" to understand concepts deeper
- It's okay to struggle - that's where learning happens!

---

## Ready to Start? 🎯

**Your roadmap:**
1. ✅ You've read this file
2. 📖 Read `vibing-history/week1-day-by-day.md` Day 1 section (detailed instructions)
3. 💻 Create `internal/cache/cache.go` and start coding
4. 🧪 Write tests and iterate
5. 💬 Ask me for code review when done

**Remember:** The goal is understanding, not just completing tasks. Take your time, think through decisions, and ask questions!

Let's build something awesome! 🚀


