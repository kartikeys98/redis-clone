package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	ServerAddress  string
	NumConnections int
	Duration       time.Duration
	ReadRatio      float64
}

type Stats struct {
	mutex           sync.Mutex
	TotalOperations uint64
	Errors          uint64
	Latencies       []time.Duration
}

func worker(id int, config Config, stats *Stats, done chan bool) {
	conn, err := net.Dial("tcp", config.ServerAddress)
	if err != nil {
		log.Printf("Worker %d: Failed to connect: %v", id, err)
		atomic.AddUint64(&stats.Errors, 1)
		done <- true
		return
	}
	reader := bufio.NewReader(conn)
	defer conn.Close()

	start := time.Now()
	cmd := ""
	for time.Since(start) < config.Duration {
		startTime := time.Now()
		key := rand.IntN(100) // Use same random key for both
		if rand.Float64() < config.ReadRatio {
			cmd = fmt.Sprintf("GET key%d", key)
		} else {
			cmd = fmt.Sprintf("SET key%d value%d", key, key)
		}
		_, err := conn.Write([]byte(cmd + "\n"))
		if err != nil {
			log.Printf("Worker %d: Failed to write: %v", id, err)
			atomic.AddUint64(&stats.Errors, 1)
			continue
		}
		_, _ = reader.ReadString('\n')
		atomic.AddUint64(&stats.TotalOperations, 1)
		stats.mutex.Lock()
		stats.Latencies = append(stats.Latencies, time.Since(startTime))
		stats.mutex.Unlock()
	}
	done <- true
}

func main() {
	// 1. Parse command-line flags
	// 2. Print configuration
	// 3. Spawn worker goroutines
	// 4. Wait for completion
	// 5. Calculate and print results
	addr := flag.String("addr", "localhost:6378", "Server address")
	conn := flag.Int("conn", 100, "Connections")
	duration := flag.Duration("duration", 10*time.Second, "Duration")
	ratio := flag.Float64("ratio", 0.8, "Read ratio")
	flag.Parse()

	config := Config{
		ServerAddress:  *addr,
		NumConnections: *conn,
		Duration:       *duration,
		ReadRatio:      *ratio,
	}
	log.Printf("Starting load test with configuration: %+v", config)

	stats := Stats{}
	done := make(chan bool, config.NumConnections)
	for i := 0; i < config.NumConnections; i++ {
		go worker(i, config, &stats, done)
	}
	for i := 0; i < config.NumConnections; i++ {
		<-done
	}
	log.Printf("Load test completed with results: ")
	log.Printf("Total operations: %d", stats.TotalOperations)
	log.Printf("Errors: %d", stats.Errors)
	log.Printf("Throughput: %f ops/s", float64(stats.TotalOperations)/config.Duration.Seconds())
	if len(stats.Latencies) > 0 {
		sort.Slice(stats.Latencies, func(i, j int) bool {
			return stats.Latencies[i] < stats.Latencies[j]
		})
		p50 := stats.Latencies[len(stats.Latencies)*50/100]
		p95 := stats.Latencies[len(stats.Latencies)*95/100]
		p99 := stats.Latencies[len(stats.Latencies)*99/100]
		log.Printf("Latency p50: %v, p95: %v, p99: %v", p50, p95, p99)
	}
}
