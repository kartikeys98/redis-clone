package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/kartikey-singh/redis/internal/cache"
)

// Helper function to start test server
// Each test gets its own cache, so tests don't interfere with each other
var testPortCounter = 17000

func startTestServer(t *testing.T) (*Server, string, func()) {
	c := cache.New(1000) // Test cache with 1000 item limit

	// Use a unique port for each test
	testPortCounter++
	addr := fmt.Sprintf("localhost:%d", testPortCounter)
	srv := New(addr, c, "standalone", "", 0)

	// Start server in goroutine
	done := make(chan error, 1)
	go func() {
		done <- srv.Start()
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Cleanup function
	cleanup := func() {
		srv.cache.Flush()
		close(done)
	}

	return srv, addr, cleanup
}

// Helper function to connect and send command
func sendCommand(t *testing.T, addr, command string) string {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Send command
	_, err = conn.Write([]byte(command + "\n"))
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}

	// Read response
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		return scanner.Text()
	}

	return ""
}

func TestServerBasicCommands(t *testing.T) {
	srv, addr, cleanup := startTestServer(t)
	defer cleanup()
	_ = srv // Keep reference to prevent garbage collection

	// Test PING
	response := sendCommand(t, addr, "PING")
	if response != "+PONG" {
		t.Errorf("PING failed: expected '+PONG', got '%s'", response)
	}

	// Test SET
	response = sendCommand(t, addr, "SET testkey testvalue")
	if !strings.Contains(response, "OK") {
		t.Errorf("SET failed: expected 'OK', got '%s'", response)
	}

	// Test GET
	response = sendCommand(t, addr, "GET testkey")
	if response != "testvalue" {
		t.Errorf("GET failed: expected 'testvalue', got '%s'", response)
	}

	// Test GET non-existent
	response = sendCommand(t, addr, "GET nonexistent")
	if response != "(nil)" {
		t.Errorf("GET non-existent failed: expected '(nil)', got '%s'", response)
	}

	// Test DEL
	response = sendCommand(t, addr, "DEL testkey")
	if !strings.Contains(response, "OK") {
		t.Errorf("DEL failed: expected 'OK', got '%s'", response)
	}

	// Verify key is deleted
	response = sendCommand(t, addr, "GET testkey")
	if response != "(nil)" {
		t.Errorf("GET after DELETE failed: expected '(nil)', got '%s'", response)
	}
}

func TestServerKEYSCommand(t *testing.T) {
	srv, addr, cleanup := startTestServer(t)
	defer cleanup()
	_ = srv

	// Set multiple keys
	sendCommand(t, addr, "SET key1 value1")
	sendCommand(t, addr, "SET key2 value2")
	sendCommand(t, addr, "SET key3 value3")

	// Test KEYS
	response := sendCommand(t, addr, "KEYS")
	if !strings.Contains(response, "key1") ||
		!strings.Contains(response, "key2") ||
		!strings.Contains(response, "key3") {
		t.Errorf("KEYS failed: expected all keys, got '%s'", response)
	}
}

func TestServerSIZECommand(t *testing.T) {
	srv, addr, cleanup := startTestServer(t)
	defer cleanup()
	_ = srv

	// Test SIZE on empty cache
	response := sendCommand(t, addr, "SIZE")
	if response != "0" {
		t.Errorf("SIZE on empty cache: expected '0', got '%s'", response)
	}

	// Add some keys
	sendCommand(t, addr, "SET key1 value1")
	sendCommand(t, addr, "SET key2 value2")

	// Test SIZE
	response = sendCommand(t, addr, "SIZE")
	if response != "2" {
		t.Errorf("SIZE after 2 sets: expected '2', got '%s'", response)
	}
}

func TestServerFLUSHCommand(t *testing.T) {
	srv, addr, cleanup := startTestServer(t)
	defer cleanup()
	_ = srv

	// Add some keys
	sendCommand(t, addr, "SET key1 value1")
	sendCommand(t, addr, "SET key2 value2")

	// Verify they exist
	response := sendCommand(t, addr, "SIZE")
	if response != "2" {
		t.Errorf("SIZE before FLUSH: expected '2', got '%s'", response)
	}

	// FLUSH
	response = sendCommand(t, addr, "FLUSH")
	if !strings.Contains(response, "OK") {
		t.Errorf("FLUSH failed: expected 'OK', got '%s'", response)
	}

	// Verify cache is empty
	response = sendCommand(t, addr, "SIZE")
	if response != "0" {
		t.Errorf("SIZE after FLUSH: expected '0', got '%s'", response)
	}
}

func TestServerValuesWithSpaces(t *testing.T) {
	srv, addr, cleanup := startTestServer(t)
	defer cleanup()
	_ = srv

	// Set value with spaces
	sendCommand(t, addr, "SET greeting hello world from redis")

	// Get it back
	response := sendCommand(t, addr, "GET greeting")
	if response != "hello world from redis" {
		t.Errorf("GET with spaces: expected 'hello world from redis', got '%s'", response)
	}
}

func TestServerConcurrentConnections(t *testing.T) {
	srv, addr, cleanup := startTestServer(t)
	defer cleanup()
	_ = srv

	// Spawn 10 concurrent clients
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			key := fmt.Sprintf("key%d", id)
			value := fmt.Sprintf("value%d", id)

			// SET
			response := sendCommand(t, addr, fmt.Sprintf("SET %s %s", key, value))
			if !strings.Contains(response, "OK") {
				t.Errorf("Concurrent SET failed for %s", key)
			}

			// GET
			response = sendCommand(t, addr, fmt.Sprintf("GET %s", key))
			if response != value {
				t.Errorf("Concurrent GET failed: expected '%s', got '%s'", value, response)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all keys exist
	response := sendCommand(t, addr, "SIZE")
	if response != "10" {
		t.Errorf("SIZE after concurrent writes: expected '10', got '%s'", response)
	}
}

func TestServerErrorHandling(t *testing.T) {
	srv, addr, cleanup := startTestServer(t)
	defer cleanup()
	_ = srv

	// Test invalid command
	response := sendCommand(t, addr, "INVALID")
	if !strings.Contains(response, "ERR") {
		t.Errorf("Invalid command: expected error, got '%s'", response)
	}

	// Test SET with missing argument
	response = sendCommand(t, addr, "SET key")
	if !strings.Contains(response, "ERR") {
		t.Errorf("SET with missing arg: expected error, got '%s'", response)
	}

	// Test GET with missing argument
	response = sendCommand(t, addr, "GET")
	if !strings.Contains(response, "ERR") {
		t.Errorf("GET with missing arg: expected error, got '%s'", response)
	}
}

func TestServerCaseInsensitivity(t *testing.T) {
	srv, addr, cleanup := startTestServer(t)
	defer cleanup()
	_ = srv

	// Test lowercase commands
	response := sendCommand(t, addr, "set mykey myvalue")
	if !strings.Contains(response, "OK") {
		t.Errorf("Lowercase SET failed: got '%s'", response)
	}

	response = sendCommand(t, addr, "get mykey")
	if response != "myvalue" {
		t.Errorf("Lowercase GET failed: expected 'myvalue', got '%s'", response)
	}

	// Test mixed case
	response = sendCommand(t, addr, "DeL mykey")
	if !strings.Contains(response, "OK") {
		t.Errorf("Mixed case DEL failed: got '%s'", response)
	}
}

func TestServerTTL(t *testing.T) {
	srv, addr, cleanup := startTestServer(t)
	defer cleanup()
	_ = srv

	// Test SET with TTL
	response := sendCommand(t, addr, "SET key value EX 10")
	if !strings.Contains(response, "OK") {
		t.Errorf("SET with TTL failed: got '%s'", response)
	}
}