1. sync.RwMutex - RLock() vs Lock()
2. TCP buffer - I was not reading from it which could accumulate
3. atomic in go - goroutine safe operation, faster than mutex
    // Assembly level:
    // 1. READ:  Load stats.Errors into CPU register
    // 2. ADD:   Increment register
    // 3. WRITE: Store register back to memory
4. Sequential:
    - 100 ns/op
    - 1 operation at a time
    - Throughput: 10 million ops/second

    Parallel (8 cores):
    - 200 ns/op (2x slower per op)
    - 8 operations simultaneously
    - Throughput: 40 million ops/second (4x faster total!)
    So parallel is:
    ❌ Slower per operation (individual latency)
    ✅ Faster overall (total throughput)
5. Benchmarking tests.
6. Race condition test with --race flag.
7. ExpiresAt - time.Time vs ExpiresAt - *time.Time. 
    pointer is less efficient than value as it requires extra 8 bytes and pointers are stored in heap which is slower than stack.
8. Go mutexes are NOT reentrant! You can't lock the same mutex twice from the same goroutine. What i was doing - locking in GET() and DELETE() both and GET() calls DELETE() resulting in deadlock. 
9. Snapshot under lock pattern: 
    Acquire lock
    Make a shallow copy of the data structure
    Release lock immediately
    Work with the copy
    eg: 
    m.mu.RLock()
	slaves := make([]*SlaveConnection, len(m.slaves))
	copy(slaves, m.slaves)
	m.mu.RUnlock()
10. ticker := time.NewTicker(10 * time.Second) & <-ticker.C