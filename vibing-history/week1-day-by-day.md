# Week 1: Day-by-Day Build Guide (No Direct Code)

## Day 1: Core Cache Data Structure

### Your Mission:
Build a thread-safe in-memory key-value store that can handle concurrent access.

### Requirements:

Must support 3 operations:
- Store a key-value pair
- Retrieve a value by key
- Delete a key

Must be thread-safe:
- Multiple goroutines should be able to read/write simultaneously
- No race conditions
- No data corruption

Data types (keep it simple for now):
- Keys: strings
- Values: strings

### Questions to Think About:
Before coding, answer these:

1. **What Go data structure naturally stores key-value pairs?**
   - Hint: It's built into the language

2. **Why isn't that data structure thread-safe by default?**
   - Think: What happens if two goroutines modify it at the same time?

3. **What Go package provides synchronization primitives?**
   - Hint: Look up `sync` package in Go docs
   - Which type is best for multiple readers, single writer?

4. **Where should the lock be?**
   - Inside each operation? Outside? Why?

### Design Decisions You Need to Make:

**Error handling:**
- What should happen if someone tries to get a non-existent key?
- Return error? Return empty string? Return bool indicating existence?

**Package structure:**
- Should this be in its own package?
- What should you export (public) vs keep private?

### Checkpoint: How to Know You're Done

Write a test that:
1. Creates your cache
2. Sets 10 keys
3. Gets all 10 keys back
4. Deletes 5 keys
5. Verifies those 5 are gone, other 5 still exist

**Command to run:** `go test -v`  
**Expected output:** All tests pass

### Day 1 Hints (Only look if stuck):

<details>
<summary>Hint 1: What data structure?</summary>
You need a `map[string]string` - Go's built-in hash table.
</details>

<details>
<summary>Hint 2: Why isn't map thread-safe?</summary>
Go maps panic if accessed concurrently. You'll see: "fatal error: concurrent map writes"
</details>

<details>
<summary>Hint 3: What sync primitive?</summary>
`sync.RWMutex` - allows multiple readers OR one writer (not both)

- `Lock()` for writes
- `RLock()` for reads
</details>

<details>
<summary>Hint 4: Structure?</summary>
Create a struct that wraps the map + mutex together. Methods on that struct handle locking.
</details>

---

## Day 2: TCP Server

### Your Mission:
Create a server that listens on a TCP port and accepts client connections.

### Requirements:

Server should:
- Listen on port 6379 (Redis default)
- Accept multiple client connections simultaneously
- Read text commands from clients
- Send text responses back

Simple protocol (for now):
```
Client sends: "SET key value"
Server responds: "OK"

Client sends: "GET key"
Server responds: "value" or "(nil)"

Client sends: "DEL key"
Server responds: "OK"
```

Handle multiple clients:
- Each client should get their own connection handler
- Clients shouldn't block each other

### Questions to Think About:

1. **What Go package handles TCP networking?**
   - Hint: Look up `net` package
   - What function listens on a TCP port?

2. **How do you handle multiple connections?**
   - Should you handle them sequentially? (slow!)
   - Or concurrently? How?

3. **How do you parse commands?**
   - Commands are text: "SET key value"
   - How do you split them into parts?
   - What if someone sends malformed input?

4. **How do you read from a connection?**
   - Byte by byte? Line by line?
   - What if the client sends a partial command?

5. **Integration with Day 1:**
   - How does the server use your cache?
   - Should server create the cache? Or receive it?

### Design Decisions:

**Error handling:**
- What if a client disconnects mid-command?
- What if they send garbage input?

**Command parsing:**
- Should you support commands with spaces in values?
- Case sensitive? ("SET" vs "set")

**Connection lifecycle:**
- Keep connection open for multiple commands? (yes, like Redis)
- Or close after each command? (simpler but inefficient)

### Checkpoint: How to Know You're Done

**Test manually:**
```bash
# Terminal 1: Start your server
go run main.go

# Terminal 2: Connect with netcat
nc localhost 6379

# Type these commands:
SET mykey myvalue
# Expect: OK

GET mykey
# Expect: myvalue

GET notexist
# Expect: (nil)

DEL mykey
# Expect: OK

GET mykey
# Expect: (nil)
```

**Advanced test:**
```bash
# Terminal 2: Open connection 1
nc localhost 6379
SET key1 value1

# Terminal 3: Open connection 2 (simultaneously!)
nc localhost 6379
SET key2 value2
GET key1
# Expect: value1 (proving both clients work)
```

### Day 2 Hints:

<details>
<summary>Hint 1: Listening</summary>
`net.Listen("tcp", ":6379")` returns a listener.
Then loop calling `listener.Accept()` to get connections.
</details>

<details>
<summary>Hint 2: Concurrency</summary>
For each accepted connection, launch a goroutine: `go handleConnection(conn)`
</details>

<details>
<summary>Hint 3: Reading data</summary>
Use `bufio.Scanner` to read line-by-line:
```go
scanner := bufio.NewScanner(conn)
for scanner.Scan() {
    line := scanner.Text() // One line of input
}
```
</details>

<details>
<summary>Hint 4: Parsing</summary>
`strings.Fields(line)` splits on whitespace into `[]string`
</details>

---

## Day 3: Testing & Benchmarking

### Your Mission:
Verify your cache works correctly under load and measure performance.

### Part A: Correctness Tests

Write tests for:

1. **Race Condition Test:**
   - Spawn 100 goroutines
   - Each writes 100 keys
   - All simultaneously
   - Verify: No crashes, all keys exist

2. **Concurrent Read/Write Test:**
   - 50 goroutines writing
   - 50 goroutines reading
   - Mix of operations
   - Verify: No corruption

**Run with race detector:**
```bash
go test -race ./cache
```
**What to look for:** Any output about data races? If yes, fix your locking!

### Part B: Performance Benchmarks

Write benchmarks for:

1. **Read performance:**
   - How many GET operations per second?
   - Run with multiple goroutines (parallel benchmark)

2. **Write performance:**
   - How many SET operations per second?

3. **Mixed workload:**
   - 80% reads, 20% writes (typical cache pattern)

**Run benchmarks:**
```bash
go test -bench=. -benchmem ./cache
```

**What to look for:**
- Operations per second (should be >100,000)
- Memory allocations per operation (fewer is better)

### Questions to Think About:

1. **What's the difference between a test and a benchmark?**
   - Test: Verifies correctness
   - Benchmark: Measures performance

2. **Why use -race flag?**
   - Detects race conditions (subtle bugs)

3. **What's b.RunParallel?**
   - Runs benchmark with multiple goroutines
   - Simulates realistic concurrent load

### Checkpoint: How to Know You're Done

**Checklist:**
- ✅ All tests pass
- ✅ No race conditions detected
- ✅ Benchmarks show >100K ops/sec for GET
- ✅ Benchmarks show >50K ops/sec for SET

**If your numbers are lower:**
- That's okay! Understanding why is the learning goal
- Common issues: Lock contention, memory allocations

### Day 3 Hints:

<details>
<summary>Hint 1: Race test structure</summary>
```go
func TestConcurrent(t *testing.T) {
    cache := New()
    var wg sync.WaitGroup
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            // ... write and read ...
        }(i)
    }
    
    wg.Wait() // Wait for all goroutines
}
```
</details>

<details>
<summary>Hint 2: Benchmark structure</summary>
```go
func BenchmarkGet(b *testing.B) {
    cache := New()
    cache.Set("key", "value")
    
    b.ResetTimer() // Don't count setup time
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            cache.Get("key")
        }
    })
}
```
</details>

---

## Day 4-5: Load Testing Tool

### Your Mission:
Build a tool to stress-test your cache server and measure throughput.

### Requirements:

Tool should:
- Connect to your server (localhost:6379)
- Send commands as fast as possible
- Use multiple concurrent connections (simulate many clients)
- Run for a fixed duration (e.g., 10 seconds)
- Report results

Metrics to track:
- Total operations completed
- Operations per second (throughput)
- Errors encountered
- Average latency

Configuration:
- Number of concurrent connections (e.g., 50)
- Test duration (e.g., 30 seconds)
- Mix of operations (e.g., 80% GET, 20% SET)

### Questions to Think About:

1. **How do you measure time accurately?**
   - Start time vs end time
   - Per-operation timing vs total timing

2. **How do you coordinate multiple connections?**
   - Each connection in a goroutine?
   - How do you wait for all to finish?
   - How do you aggregate results from all goroutines?

3. **What should you test?**
   - Only GET? Only SET? Mixed?
   - Real-world caches see mostly reads (80-90%)

4. **Atomic counters:**
   - Multiple goroutines incrementing same counter = race condition
   - How do you safely count operations across goroutines?
   - Hint: `sync/atomic` package

### Design Decisions:

**Connection reuse:**
- Should each goroutine open new connection per request? (slow!)
- Or keep connection open? (like persistent HTTP connections)

**Error handling:**
- If a connection fails, retry? Or just count as error?

**Realistic data:**
- Use same key for all operations? (unrealistic)
- Or generate random keys? (more realistic)

### Checkpoint: How to Know You're Done

**Run your load tester:**
```bash
# Terminal 1: Start server
go run main.go

# Terminal 2: Run load test
go run cmd/loadtest/main.go

# Expected output:
===  Load Test Results ===
Duration: 10s
Total Operations: 650,000
Errors: 0
Throughput: 65,000 ops/sec
Avg Latency: 0.015 ms
```

**Goal:** Achieve >50,000 ops/sec

**Compare with real Redis:**
```bash
# If you have redis-server installed
redis-server &
redis-benchmark -t get,set -n 100000 -q

# Compare your numbers!
```

### Day 4-5 Hints:

<details>
<summary>Hint 1: Atomic counters</summary>
```go
import "sync/atomic"

var totalOps uint64
var errors uint64

// In goroutine:
atomic.AddUint64(&totalOps, 1)
```
</details>

<details>
<summary>Hint 2: Timing</summary>
```go
start := time.Now()
// ... run test ...
elapsed := time.Since(start)

opsPerSec := float64(totalOps) / elapsed.Seconds()
```
</details>

<details>
<summary>Hint 3: Multiple connections</summary>
```go
for i := 0; i < numConnections; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        conn, _ := net.Dial("tcp", "localhost:6379")
        defer conn.Close()
        
        // Send commands in a loop
        for time.Since(start) < duration {
            // ... send command ...
            atomic.AddUint64(&totalOps, 1)
        }
    }()
}
wg.Wait()
```
</details>

---

## Day 6-7: Documentation & Validation

### Your Mission:
Polish your work and document what you've learned.

### Tasks:

1. **Write README.md:**
   - What does your cache do?
   - How to run it?
   - How to test it?
   - What are the performance characteristics?

2. **Code cleanup:**
   - Add comments to non-obvious parts
   - Remove debug print statements
   - Organize into proper packages

3. **Create examples:**
   - Show how to use it as a library
   - Show how to run the server

4. **Git setup:**
   - Initialize repo
   - Add .gitignore (ignore binaries)
   - Commit your work

### Validation Checklist:

**Run these commands and verify all pass:**
```bash
# 1. Build succeeds
go build

# 2. All tests pass
go test ./...

# 3. No race conditions
go test -race ./...

# 4. Benchmarks run
go test -bench=. ./cache

# 5. Server starts
./your-cache-server &

# 6. Can connect
echo "SET test value" | nc localhost 6379

# 7. Load test achieves target
go run cmd/loadtest/main.go
# Should show >50K ops/sec

# 8. Code is formatted
go fmt ./...

# 9. No obvious issues
go vet ./...
```

### Write a Reflection:

Answer these questions (for yourself):

1. **What was hardest?**
   - Concurrency? Networking? Testing?

2. **What surprised you?**
   - Performance? Simplicity? Complexity?

3. **What would you do differently?**
   - Different data structures? Different protocol?

4. **What did you learn?**
   - About Go? About distributed systems? About yourself?

---

## End of Week 1: Share Your Work

### What to Share:

**Test results:**
```
Paste output of:
- go test -v ./...
- go test -bench=. ./cache
- Load test results
```

**Challenges faced:**
- What problems did you encounter?
- How did you solve them?

**Questions for Week 2:**
- What confused you?
- What do you want to understand better?

**Architecture decisions:**
- Why did you structure code a certain way?
- What tradeoffs did you make?

---

## How I'll Help You Build This

### My Role:

✅ **Answer specific questions:**
- "Should I use X or Y approach?"
- "Why am I getting this error?"
- "Is this the right way to handle Z?"

✅ **Review your code:**
- "Here's what I built, what do you think?"
- "Is my locking strategy correct?"
- "How can I improve performance?"

✅ **Explain concepts:**
- "Why do we need RWMutex vs Mutex?"
- "What's the tradeoff here?"

✅ **Debug together:**
- "My test fails with this error..."
- "Race detector shows this..."

❌ **Won't give you direct solutions unless you're truly stuck**

### How to Ask for Help:

**Good question:**
> "I'm implementing the TCP server. I can accept connections, but when I try to read data, I only get partial commands. Should I be buffering somewhere? Here's my current approach: [explain what you tried]"

**Even better:**
> "Here's my connection handler code [paste]. It works for single commands but breaks when client sends multiple commands quickly. I think it's a buffering issue - am I reading from the connection correctly?"

**I'll respond with:**
- Probing questions to make you think
- Hints about what to look up
- Explanations of concepts
- Only if totally stuck: skeleton/pseudocode

---

## Ready to Start?

1. **Begin with Day 1:**
   - Think about the questions I posed
   - Start building your cache data structure
   - When you have questions or want a review, come back!

2. **Cursor will help with:**
   - Syntax
   - Autocomplete
   - Quick fixes

3. **I'll help with:**
   - Design decisions
   - Architecture
   - Understanding concepts
   - Debugging tricky issues

**Start Day 1 and let me know when you have your first version or hit your first question!** 🚀





