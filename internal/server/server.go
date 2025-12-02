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
)

type Server struct {
	addr  string
	cache *cache.Cache
}

func New(addr string, cache *cache.Cache) *Server {
	return &Server{
		addr:  addr,
		cache: cache,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Server listening on %s", s.addr)

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
			s.cache.SetWithTTL(key, value, ttl)
			conn.Write([]byte("+OK\n"))

		case "GET":
			if len(parts) < 2 {
				conn.Write([]byte("ERR wrong number of arguments for 'get' command\n"))
				continue
			}
			value, found := s.cache.Get(parts[1])
			if !found {
				conn.Write([]byte("(nil)\n"))
			} else {
				conn.Write([]byte(value + "\n"))
			}

		case "DEL":
			if len(parts) < 2 {
				conn.Write([]byte("ERR wrong number of arguments for 'del' command\n"))
				continue
			}
			deleted := s.cache.Delete(parts[1])
			if deleted {
				conn.Write([]byte("+OK\n"))
			} else {
				conn.Write([]byte("(nil)\n"))
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
			s.cache.Flush()
			conn.Write([]byte("+OK\n"))

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
