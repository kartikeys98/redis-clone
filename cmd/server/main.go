package main

import (
	"fmt"
	"log"

	"github.com/kartikey-singh/redis/internal/cache"
	"github.com/kartikey-singh/redis/internal/server"
)

func main() {
	fmt.Println("🚀 Starting Redis Clone...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Create cache
	c := cache.New()

	// Create server
	addr := ":6378"
	srv := server.New(addr, c)

	fmt.Printf("📡 Server address: %s\n", addr)
	fmt.Println("📝 Supported commands:")
	fmt.Println("   - SET key value  : Store a key-value pair")
	fmt.Println("   - GET key        : Retrieve a value")
	fmt.Println("   - DEL key        : Delete a key")
	fmt.Println("   - KEYS           : List all keys")
	fmt.Println("   - SIZE           : Get cache size")
	fmt.Println("   - FLUSH          : Clear all data")
	fmt.Println("   - PING           : Test connection")
	fmt.Println("\n🔗 Connect with: nc localhost 6378")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatal("Server error:", err)
	}
}
