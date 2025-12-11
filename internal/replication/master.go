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
	conn   net.Conn
	writer *bufio.Writer
	mu     sync.Mutex // Protects writer
}

func NewMaster(c *cache.Cache) *Master {
	return &Master{
		cache:  c,
		slaves: make([]*SlaveConnection, 0),
	}
}

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
				m.removeSlave(s)
			}
		}(slave)
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
	m.mu.Lock()
	defer m.mu.Unlock()

	m.slaves = append(m.slaves, &SlaveConnection{
		conn:   conn,
		writer: bufio.NewWriter(conn),
	})
	log.Printf("Slave connected: %s (total: %d)", conn.RemoteAddr(), len(m.slaves))
}

// removeSlave removes a disconnected slave
func (m *Master) removeSlave(slave *SlaveConnection) {
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
