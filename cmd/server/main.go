package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kartikey-singh/redis/internal/cache"
	"github.com/kartikey-singh/redis/internal/server"
)

func main() {
	fmt.Println("🚀 Starting Application Server ...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Create cache with 10,000 item limit
	c := cache.New(10000)

	// Create server
	port := flag.Int("port", 6379, "Port to listen on")
	role := flag.String("role", "standalone", "Role: master, slave or standalone")
	masterAddr := flag.String("master", "localhost:6380", "Master address")
	replicationPort := flag.Int("replication-port", 6380, "Replication port")
	flag.Parse()
	addr := fmt.Sprintf(":%d", *port)
	srv := server.New(addr, c, *role, *masterAddr, *replicationPort)

	fmt.Printf("📡 Server address: %s\n", addr)
	fmt.Println("📝 Supported commands:")
	fmt.Println("   - SET key value  : Store a key-value pair")
	fmt.Println("   - GET key        : Retrieve a value")
	fmt.Println("   - DEL key        : Delete a key")
	fmt.Println("   - KEYS           : List all keys")
	fmt.Println("   - SIZE           : Get cache size")
	fmt.Println("   - FLUSH          : Clear all data")
	fmt.Println("   - PING           : Test connection")
	fmt.Printf("\n🔗 Connect with: nc localhost:%d\n", *port)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("Role: %s, Master address: %s, Replication port: %d", *role, *masterAddr, *replicationPort)

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatal("Server error:", err)
	}
}
