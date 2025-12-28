package replication

import (
	"bufio"
	"log"
	"net"
	"sync"
	"time"

	"github.com/kartikey-singh/redis/internal/cache"
)

type Slave struct {
	cache      *cache.Cache
	masterAddr string
	conn       net.Conn
	mu         sync.RWMutex
}

func NewSlave(c *cache.Cache, masterAddr string) *Slave {
	return &Slave{
		cache:      c,
		masterAddr: masterAddr,
	}
}

// ConnectToMaster establishes connection to master
func (s *Slave) ConnectToMaster() error {
	conn, err := net.Dial("tcp", s.masterAddr)
	if err != nil {
		return err
	}
	s.conn = conn
	log.Printf("Connected to master: %s", s.masterAddr)
	return nil
}

// StartReplication receives and applies operations from master
func (s *Slave) StartReplication() error {
	scanner := bufio.NewScanner(s.conn)
	for scanner.Scan() {
		line := scanner.Text()
		op, err := ParseOperation(line)
		if err != nil {
			log.Printf("Error parsing operation: %v", err)
			continue
		}
		s.apply(op) // Apply synchronously to maintain order
	}
	return scanner.Err()
}

// apply executes an operation on the local cache
func (s *Slave) apply(op *Operation) {
	switch op.Type {
	case OpSet:
		if op.TTL > 0 {
			// Calculate remaining TTL to account for replication lag
			elapsed := time.Since(time.Unix(op.Timestamp, 0))
			remaining := op.TTL - elapsed

			log.Printf("SET %s: original TTL=%v, elapsed=%v, remaining=%v", op.Key, op.TTL, elapsed, remaining)

			if remaining <= 0 {
				// Already expired during replication, skip
				log.Printf("Skipping expired key: %s", op.Key)
				return
			}

			s.cache.SetWithTTL(op.Key, op.Value, remaining)
			log.Printf("Applied SET with TTL: %s (remaining=%v)", op.Key, remaining)
		} else {
			// No TTL, set without expiration
			s.cache.Set(op.Key, op.Value)
			log.Printf("Applied SET without TTL: %s", op.Key)
		}
	case OpDelete:
		s.cache.Delete(op.Key)
		log.Printf("Applied DELETE: %s", op.Key)
	case OpFlush:
		s.cache.Flush()
		log.Printf("Applied FLUSH")
	case OpPing:
		log.Printf("Received PING from master")
	default:
		log.Printf("Unknown operation: %s", op.Type)
	}
}

// Get reads from local cache (slaves can serve reads!)
func (s *Slave) Get(key string) (string, bool) {
	return s.cache.Get(key)
}

// Close closes connection to master
func (s *Slave) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
