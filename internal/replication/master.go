package replication

import (
	"bufio"
	"log"
	"net"
	"sync"
	"time"

	"github.com/kartikey-singh/redis/internal/cache"
)

type Master struct {
	cache  *cache.Cache
	slaves []*SlaveConnection
	mu     sync.RWMutex
}

type SlaveConnection struct {
	conn          net.Conn
	writer        *bufio.Writer
	mu            sync.RWMutex
	health        *HealthMonitor
	pongReceived  chan int64
	stopHeartbeat chan struct{}
	closeOnce     sync.Once
}

func NewMaster(c *cache.Cache) *Master {
	return &Master{
		cache:  c,
		slaves: make([]*SlaveConnection, 0),
	}
}

// Cache functions
// Set wraps cache.SetWithTTL and broadcasts to slaves
func (m *Master) Set(key, value string, ttl time.Duration) error {
	m.cache.SetWithTTL(key, value, ttl)
	m.broadcast(&Operation{
		Type:      OpSet,
		Key:       key,
		Value:     value,
		TTL:       ttl,
		Timestamp: time.Now().Unix(),
	})
	return nil
}

// Delete wraps cache.Delete and broadcasts to slaves
func (m *Master) Delete(key string) error {
	m.cache.Delete(key)
	m.broadcast(&Operation{
		Type:      OpDelete,
		Key:       key,
		Timestamp: time.Now().Unix(),
	})
	return nil
}

// Flush wraps cache.Flush and broadcasts to slaves
func (m *Master) Flush() error {
	m.cache.Flush()
	m.broadcast(&Operation{
		Type:      OpFlush,
		Timestamp: time.Now().Unix(),
	})
	return nil
}

// Get reads from cache (no replication needed for reads)
func (m *Master) Get(key string) (string, bool) {
	return m.cache.Get(key)
}

// broadcast sends operation to all connected slaves
func (m *Master) broadcast(op *Operation) {
	m.mu.RLock()
	slaves := make([]*SlaveConnection, len(m.slaves))
	copy(slaves, m.slaves)
	m.mu.RUnlock()

	for _, slave := range slaves {
		go func(s *SlaveConnection) {
			if err := s.Send(op); err != nil {
				//m.removeSlave(s)
				log.Printf("Failed to send operation to slave: %s", s.conn.RemoteAddr())
			}
		}(slave)
	}
}

// Goroutine 1: Listen for PONGs
func (s *SlaveConnection) ListenForPongs() {
	scanner := bufio.NewScanner(s.conn)
	for scanner.Scan() {
		op, err := ParseOperation(scanner.Text())
		if err == nil && op.Type == OpPong {
			select {
			case s.pongReceived <- op.Timestamp:
				// Successfully sent
			default:
				// Receiver not ready, drop this late pong
			}
		}
	}
	close(s.pongReceived)
}

// Heartbeat functions
func (m *Master) StartHeartbeatForSlave(s *SlaveConnection, pingInterval time.Duration, maxMissedHeartbeats int) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopHeartbeat:
			return
		case <-ticker.C:
			timestamp := time.Now().Unix()
			op := &Operation{Type: OpPing, Timestamp: timestamp}
			if err := s.Send(op); err != nil {
				log.Printf("Heartbeat failed for slave: %s", s.conn.RemoteAddr())
				s.health.RecordFailure()
				if !s.health.IsHealthy() {
					m.removeSlave(s)
					return
				}
				// If threshold not breached, continue for next tick
				continue
			}
			select {
			case <-s.stopHeartbeat:
				return
			case pongTimestamp := <-s.pongReceived:
				if pongTimestamp != timestamp {
					log.Printf("Pong timestamp mismatch for slave: %s", s.conn.RemoteAddr())
					s.health.RecordFailure()
					if !s.health.IsHealthy() {
						m.removeSlave(s)
						return
					}
				} else {
					log.Printf("Slave is healthy: %s", s.conn.RemoteAddr())
					s.health.RecordSuccess()
				}
			case <-time.After(pingInterval):
				log.Printf("Heartbeat failed for slave: %s", s.conn.RemoteAddr())
				s.health.RecordFailure()
				if !s.health.IsHealthy() {
					m.removeSlave(s)
					return
				}
			}
		}
	}
}

// ListenForSlaves accepts slave connections on the given port
func (m *Master) ListenForSlaves(port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go m.addSlave(conn)
		log.Printf("New slave connected: %s", conn.RemoteAddr())
	}
}

// addSlave adds a new slave connection
func (m *Master) addSlave(conn net.Conn) {
	slave := &SlaveConnection{
		conn:          conn,
		writer:        bufio.NewWriter(conn),
		health:        NewHealthMonitor(5*time.Second, 3),
		pongReceived:  make(chan int64),
		stopHeartbeat: make(chan struct{}),
	}

	// Send all existing data first
	m.mu.RLock()
	keys := m.cache.Keys()
	m.mu.RUnlock()

	for _, key := range keys {
		value, ttl, found := m.cache.GetWithTTL(key)
		if found {
			op := &Operation{
				Type: OpSet, Key: key, Value: value, TTL: ttl,
				Timestamp: time.Now().Unix(),
			}
			if err := slave.Send(op); err != nil {
				log.Printf("Failed to send initial state: %v", err)
				conn.Close()
				return
			}
		}
	}

	// Now add to slave list for ongoing replication
	m.mu.Lock()
	m.slaves = append(m.slaves, slave)
	m.mu.Unlock()

	// Add health monitoring
	go m.StartHeartbeatForSlave(slave, 5*time.Second, 3)
	go slave.ListenForPongs()
}

// removeSlave removes a disconnected slave
func (m *Master) removeSlave(slave *SlaveConnection) {
	log.Printf("Slave is unhealthy, removing: %s", slave.conn.RemoteAddr())
	slave.closeOnce.Do(func() {
		m.mu.Lock()
		defer m.mu.Unlock()

		for i, s := range m.slaves {
			if s == slave {
				m.slaves = append(m.slaves[:i], m.slaves[i+1:]...)
				break
			}
		}
		slave.conn.Close()
		log.Printf("Slave disconnected: %s (total: %d)", slave.conn.RemoteAddr(), len(m.slaves))
		close(slave.stopHeartbeat)
	})
}

// Send sends an operation to this slave
func (s *SlaveConnection) Send(op *Operation) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.writer.WriteString(op.String())
	if err != nil {
		return err
	}
	err = s.writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
