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
    Q) Why need lock we can send broadcast event to all slaves from for loop even if another gorutine removes it?
    A) If another gorutine removes it then during for loop it might access some undefined memory location + the go race flag will not allow it to pass.
10. ticker := time.NewTicker(10 * time.Second) & <-ticker.C
11. 
    a. Array of pointers: 
    slaves := []*SlaveConnection{&slave1, &slave2}
    slaves[0].conn = newConn  
    // Modifies the actual SlaveConnection
    b. Pointer to a slice
    slaves := &[]SlaveConnection{slave1, slave2}
    (*slaves)[0].conn = newConn  
    // Modifies a copy, not the original!
12. Understanding of buffer:
    writer := bufio.NewWriter(conn)
    // io.Writer -> interface, net.Conn -> interface
    // anything that is net.Conn is also an io.writer as it implements Write() method
    writer.WriteString(str)
    writer.Flush()
    Q) Why does Go allow this without any conversion?
        Because Go uses implicit interface satisfaction:
        You don’t write: type TCPConn implements io.Writer
        Go checks automatically: “does this value have the methods the interface needs?”
        So when you write:
        w := bufio.NewWriter(conn)
        Go is effectively saying:
        “NewWriter needs something that can Write([]byte)”
        “conn can Write([]byte)”
        “ok, pass it”
13. Note: when writing to a connection; always lock it so it doesn't result in race condition - Hence it's better to use a connection buffer so it's easy to take lock from different parts of code.
14. sync.Once for code that should only execute once. Internally it uses a done flag and a mutex.
15. Channel: 
    a. Closing a channel i.e. close(chan) -> sends events to all gorutines listening on that channel
    b. Listening on a channel -> only sends event to a single goroutine even if multiple are listening.
    c. Non blocking channel pattern using default case in select i.e.
    select {
    case s.pongReceived <- op.Timestamp:
        // Successfully sent
    default:
        // Receiver not ready, drop this late pong
    }