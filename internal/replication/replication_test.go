package replication

import (
	"testing"
	"time"

	"github.com/kartikey-singh/redis/internal/cache"
)

func TestMasterSlaveReplication(t *testing.T) {
	// Create master cache and master
	masterCache := cache.New(100)
	defer masterCache.Close()
	master := NewMaster(masterCache)

	// Start master listening on port 0 (random port)
	masterPort := ":19000"
	go master.ListenForSlaves(masterPort)

	// Give master time to start
	time.Sleep(100 * time.Millisecond)

	// Create slave cache and slave
	slaveCache := cache.New(100)
	defer slaveCache.Close()
	slave := NewSlave(slaveCache, "localhost"+masterPort)

	// Connect slave to master
	err := slave.ConnectToMaster()
	if err != nil {
		t.Fatalf("Failed to connect to master: %v", err)
	}
	defer slave.Close()

	// Start slave replication in background
	go slave.StartReplication()

	// Give connection time to establish
	time.Sleep(100 * time.Millisecond)

	// Test 1: SET operation replicates
	err = master.Set("key1", "value1", 0)
	if err != nil {
		t.Fatalf("Master Set failed: %v", err)
	}

	// Give time for replication
	time.Sleep(100 * time.Millisecond)

	// Verify on slave
	val, found := slave.Get("key1")
	if !found {
		t.Error("key1 should exist on slave after replication")
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %s", val)
	}

	// Test 2: DELETE operation replicates
	err = master.Delete("key1")
	if err != nil {
		t.Fatalf("Master Delete failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	_, found = slave.Get("key1")
	if found {
		t.Error("key1 should not exist on slave after DELETE")
	}

	// Test 3: SET with TTL replicates
	err = master.Set("key2", "value2", 1*time.Second)
	if err != nil {
		t.Fatalf("Master Set with TTL failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	val, found = slave.Get("key2")
	if !found {
		t.Error("key2 should exist on slave before expiration")
	}
	if val != "value2" {
		t.Errorf("Expected value2, got %s", val)
	}

	// Wait for expiration (1s original TTL + buffer)
	time.Sleep(1100 * time.Millisecond)

	_, found = slave.Get("key2")
	if found {
		t.Error("key2 should be expired on slave")
	}
}

func TestMultipleSlaves(t *testing.T) {
	// Create master
	masterCache := cache.New(100)
	defer masterCache.Close()
	master := NewMaster(masterCache)

	masterPort := ":19001"
	go master.ListenForSlaves(masterPort)
	time.Sleep(100 * time.Millisecond)

	// Create 3 slaves
	slaves := make([]*Slave, 3)
	for i := 0; i < 3; i++ {
		slaveCache := cache.New(100)
		defer slaveCache.Close()

		slave := NewSlave(slaveCache, "localhost"+masterPort)
		err := slave.ConnectToMaster()
		if err != nil {
			t.Fatalf("Slave %d failed to connect: %v", i, err)
		}
		defer slave.Close()

		go slave.StartReplication()
		slaves[i] = slave
	}

	time.Sleep(100 * time.Millisecond)

	// Set on master
	master.Set("testkey", "testvalue", 0)

	time.Sleep(200 * time.Millisecond)

	// Verify all slaves have it
	for i, slave := range slaves {
		val, found := slave.Get("testkey")
		if !found {
			t.Errorf("Slave %d should have testkey", i)
		}
		if val != "testvalue" {
			t.Errorf("Slave %d: expected testvalue, got %s", i, val)
		}
	}
}
