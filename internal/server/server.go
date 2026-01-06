package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/kartikey-singh/redis/internal/cache"
	"github.com/kartikey-singh/redis/internal/replication"
)

type Server struct {
	addr            string
	cache           *cache.Cache
	role            string
	masterAddr      string
	replicationPort int
	master          *replication.Master
	slave           *replication.Slave
}

func New(addr string, cache *cache.Cache, role string, masterAddr string, replicationPort int) *Server {
	s := &Server{
		addr:            addr,
		cache:           cache,
		role:            role,
		masterAddr:      masterAddr,
		replicationPort: replicationPort,
	}
	if role == "master" {
		s.master = replication.NewMaster(cache)
	} else if role == "slave" {
		s.slave = replication.NewSlave(cache, masterAddr)
	}
	return s
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Server listening on %s", s.addr)

	switch s.role {
	case "master":
		go s.master.ListenForSlaves(fmt.Sprintf(":%d", s.replicationPort))
	case "slave":
		if err := s.slave.ConnectToMaster(); err != nil {
			log.Printf("Error connecting to master: %v", err)
			return err
		}
		go s.slave.StartReplication()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue // Don't exit on accept error
		}
		log.Printf("New connection from %s", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		log.Printf("Connection closed from %s", conn.RemoteAddr())
		conn.Close()
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])
		log.Printf("[%s] Command: %s", conn.RemoteAddr(), line)

		switch command {
		case "SET":
			if len(parts) < 3 {
				conn.Write([]byte("ERR wrong number of arguments for 'set' command\n"))
				continue
			}
			key := parts[1]
			var value string
			var ttl time.Duration

			// Check for TTL
			if len(parts) >= 5 && strings.ToUpper(parts[len(parts)-2]) == "EX" {
				ttlValue := parts[len(parts)-1]
				t, err := strconv.Atoi(ttlValue)
				if err != nil || t <= 0 {
					conn.Write([]byte("ERR invalid TTL value\n"))
					continue
				}
				value = strings.Join(parts[2:len(parts)-2], " ")
				ttl = time.Duration(t) * time.Second
			} else {
				value = strings.Join(parts[2:], " ")
				ttl = 0
			}
			switch s.role {
			case "master":
				err := s.master.Set(key, value, ttl)
				if err != nil {
					conn.Write([]byte("ERR " + err.Error() + "\n"))
					continue
				}
				conn.Write([]byte("+OK\n"))
			case "slave":
				conn.Write([]byte("+ERR: Slave is not allowed to set keys\n"))
			case "standalone":
				s.cache.SetWithTTL(key, value, ttl)
				conn.Write([]byte("+OK\n"))
			}

		case "GET":
			if len(parts) < 2 {
				conn.Write([]byte("ERR wrong number of arguments for 'get' command\n"))
				continue
			}
			switch s.role {
			case "master":
				value, found := s.master.Get(parts[1])
				if !found {
					conn.Write([]byte("(nil)\n"))
				} else {
					conn.Write([]byte(value + "\n"))
				}
			case "slave":
				value, found := s.slave.Get(parts[1])
				if !found {
					conn.Write([]byte("(nil)\n"))
				} else {
					conn.Write([]byte(value + "\n"))
				}
			case "standalone":
				value, found := s.cache.Get(parts[1])
				if !found {
					conn.Write([]byte("(nil)\n"))
				} else {
					conn.Write([]byte(value + "\n"))
				}
			}

		case "DEL":
			if len(parts) < 2 {
				conn.Write([]byte("ERR wrong number of arguments for 'del' command\n"))
				continue
			}
			switch s.role {
			case "master":
				deleted := s.master.Delete(parts[1])
				if deleted == nil {
					conn.Write([]byte("+OK\n"))
				} else {
					conn.Write([]byte("+ERR: " + deleted.Error() + "\n"))
				}
			case "slave":
				conn.Write([]byte("+ERR: Slave is not allowed to delete keys\n"))

			case "standalone":
				deleted := s.cache.Delete(parts[1])
				if deleted {
					conn.Write([]byte("+OK\n"))
				} else {
					conn.Write([]byte("+ERR: Key not found\n"))
				}
			}

		case "PING":
			conn.Write([]byte("+PONG\n"))

		case "KEYS":
			keys := s.cache.Keys()
			if len(keys) == 0 {
				conn.Write([]byte("(empty)\n"))
			} else {
				response := strings.Join(keys, ", ") + "\n"
				conn.Write([]byte(response))
			}

		case "FLUSH":
			switch s.role {
			case "master":
				err := s.master.Flush()
				if err != nil {
					conn.Write([]byte("+ERR: " + err.Error() + "\n"))
					continue
				}
				conn.Write([]byte("+OK\n"))
			case "slave":
				conn.Write([]byte("+ERR: Slave is not allowed to flush the cache\n"))
			case "standalone":
				s.cache.Flush()
				conn.Write([]byte("+OK\n"))
			}
		case "SIZE":
			size := s.cache.Size()
			conn.Write([]byte(fmt.Sprintf("%d\n", size)))
		default:
			conn.Write([]byte("ERR unknown command '" + command + "'\n"))
		}
	}
	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		log.Printf("[%s] Scanner error: %v", conn.RemoteAddr(), err)
		return
	}
}
