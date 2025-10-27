# Context History - Build Your Own Redis

## Project Overview
Building a distributed Redis implementation from scratch to gain deep knowledge of:
- System design principles
- Distributed systems concepts
- Low-level implementation details
- Interview preparation for system design questions

## Goals
- **Learning**: Understand how distributed databases and caches work at a fundamental level
- **Practical Skills**: Build production-grade distributed systems knowledge
- **Interview Prep**: Master common system design concepts (LRU, consistent hashing, consensus algorithms)

## Project Scope

### Week 1-2: Single Node Cache + LRU/TTL
**Core Concepts to Learn:**
- Concurrency control (RWMutex, race conditions)
- TCP protocol design (RESP protocol)
- Connection handling (goroutine per connection)
- LRU cache algorithm (doubly-linked list + HashMap)
- TTL implementation (passive + active expiration)
- Memory limits & eviction policies

**Build Yourself:** ✅ All components
- Core cache data structure
- LRU algorithm
- TCP server
- Concurrency patterns

### Week 3-4: Replication & Fault Tolerance
**Core Concepts to Learn:**
- Master-Slave replication model
- Asynchronous replication
- Statement-based replication protocol
- Fault detection (heartbeat)
- Manual failover
- Replication lag handling

**Build Yourself:** ✅ Replication protocol
- Master-slave setup
- Replication log
- Failure detection
- Manual promotion

### Week 5-6: Sharding & Consistent Hashing
**Core Concepts to Learn:**
- Horizontal scaling through data partitioning
- Hash-based sharding problems
- Consistent hashing algorithm
- Virtual nodes for even distribution
- Minimal rebalancing on node changes

**Build Yourself:** ✅ Consistent hashing implementation
- Hash ring data structure
- Virtual nodes
- Key placement algorithm

### Week 7-8: Raft Consensus for Leader Election
**Core Concepts to Learn:**
- Consensus algorithms
- Leader election
- Log replication
- Strong consistency guarantees
- Split-brain prevention

**Build Yourself:** ❌ Use library (hashicorp/raft)
- Integration of Raft into the system
- Understanding when to use consensus vs simpler replication

## Technical Stack
- **Language**: Go (for concurrency primitives and performance)
- **Protocol**: Start with simple text protocol, optionally implement RESP later
- **Libraries**: Minimal - build most things from scratch for learning

## Design Decisions

### Build vs Use Library
- **Build**: Cache structures, LRU, TCP server, replication protocol, consistent hashing
- **Use Library**: Raft consensus, serialization (Protocol Buffers)
- **Rationale**: Build what teaches core concepts; use battle-tested implementations for complex, error-prone components

### Architecture Progression
1. Single-node cache (Week 1-2)
2. Add replication (Week 3-4)
3. Add sharding (Week 5-6)
4. Add consensus (Week 7-8)

Each phase builds on the previous, allowing incremental learning and testing.

## Real-World Systems to Understand
By completing this project, will gain deep understanding of:
- Redis (obviously)
- Memcached
- DynamoDB (consistent hashing + leaderless)
- Cassandra (consistent hashing + gossip)
- etcd (Raft-based KV store)
- MongoDB (replica sets with auto-failover)
- CDNs (consistent hashing for routing)





