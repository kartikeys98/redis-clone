package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBasicOperations(t *testing.T) {
	c := New(1000)

	// Test SET and GET
	c.Set("key1", "value1")

	got, found := c.Get("key1")
	if !found {
		t.Error("expected to find key1")
	}
	if got != "value1" {
		t.Errorf("expected 'value1', got '%s'", got)
	}

	// Test DELETE
	deleted := c.Delete("key1")
	if !deleted {
		t.Error("expected delete to return true")
	}

	_, found = c.Get("key1")
	if found {
		t.Error("key1 should not exist after delete")
	}
}

func TestGetNonExistent(t *testing.T) {
	c := New(1000)

	// Get a key that was never set
	_, found := c.Get("nonexistent")
	if found {
		t.Error("should not find nonexistent key")
	}
}

func TestConcurrentOperations(t *testing.T) {
	c := New(1000)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			c.Set(fmt.Sprintf("key%d", id), fmt.Sprintf("value%d", id))
		}(i)
	}
	wg.Wait()

	for i := 0; i < 100; i++ {
		got, found := c.Get(fmt.Sprintf("key%d", i))
		if !found {
			t.Errorf("expected to find key%d", i)
		}
		if got != fmt.Sprintf("value%d", i) {
			t.Errorf("expected 'value%d', got '%s'", i, got)
		}
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	c := New(1000)
	var wg sync.WaitGroup

	// Pre-populate
	for i := 0; i < 50; i++ {
		c.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	// 50 goroutines writing
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Set(fmt.Sprintf("key%d", id), fmt.Sprintf("value%d-%d", id, j))
			}
		}(i)
	}

	// 50 goroutines reading
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Get(fmt.Sprintf("key%d", id))
			}
		}(i)
	}

	wg.Wait()
	// If we got here without panic/race, it's thread-safe!
}

// Test 1: Basic Eviction
// Verifies that when cache is full, adding a new key evicts the oldest (LRU) key
func TestCacheWithLRU_BasicEviction(t *testing.T) {
	c := New(3)
	c.Set("A", "A")
	c.Set("B", "B")
	c.Set("C", "C")
	// Cache is now full (3 items)
	c.Set("D", "D") // Should evict "A" (oldest)

	// Verify A was evicted
	_, found := c.Get("A")
	if found {
		t.Error("expected A to be evicted")
	}

	// Verify B, C, D still exist
	for _, key := range []string{"B", "C", "D"} {
		val, found := c.Get(key)
		if !found {
			t.Errorf("expected %s to exist after eviction", key)
		}
		if val != key {
			t.Errorf("expected %s to have value %s, got %s", key, key, val)
		}
	}
}

// Test 2: LRU Order (Get affects eviction)
// Verifies that Get operations update the LRU order, affecting which key gets evicted
func TestCacheWithLRU_GetAffectsEviction(t *testing.T) {
	c := New(3)
	c.Set("A", "A")
	c.Set("B", "B")
	c.Set("C", "C")
	// Cache is full: A (oldest) -> B -> C (newest)

	// Get("A") makes A the most recently used
	c.Get("A")
	// Order is now: B (oldest) -> C -> A (newest)

	// Adding D should evict B (now oldest), not A
	c.Set("D", "D")

	// Verify B was evicted
	_, found := c.Get("B")
	if found {
		t.Error("expected B to be evicted (was oldest after Get(A))")
	}

	// Verify A, C, D still exist
	for _, key := range []string{"A", "C", "D"} {
		val, found := c.Get(key)
		if !found {
			t.Errorf("expected %s to exist after eviction", key)
		}
		if val != key {
			t.Errorf("expected %s to have value %s, got %s", key, key, val)
		}
	}
}

// Test 3: Update doesn't evict
// Verifies that updating an existing key (Set on existing key) doesn't cause eviction
// and that Get operations correctly affect which key gets evicted
func TestCacheWithLRU_UpdateDoesntEvict(t *testing.T) {
	c := New(3)
	c.Set("A", "A")
	c.Set("B", "B")
	c.Set("C", "C")
	// Cache is full: A (oldest) -> B -> C (newest)

	// Get("A") makes A most recent
	c.Get("A")
	// Order is now: B (oldest) -> C -> A (newest)

	// Set("D") should evict B (now oldest), not A
	c.Set("D", "D")

	// Verify B was evicted
	_, found := c.Get("B")
	if found {
		t.Error("expected B to be evicted (was oldest after Get(A))")
	}

	// Verify A, C, D exist
	for _, key := range []string{"A", "C", "D"} {
		val, found := c.Get(key)
		if !found {
			t.Errorf("expected %s to exist after eviction", key)
		}
		if val != key {
			t.Errorf("expected %s to have value %s, got %s", key, key, val)
		}
	}
}

// Test 4: Multiple evictions
// Verifies that updating an existing key doesn't cause eviction
func TestCacheWithLRU_UpdateExistingKey(t *testing.T) {
	c := New(3)
	c.Set("A", "A")
	c.Set("B", "B")
	c.Set("C", "C")
	// Cache is full: A (oldest) -> B -> C (newest)

	// Update existing key B - should not evict anything
	c.Set("B", "updated")

	// Verify all three keys still exist
	keys := []string{"A", "B", "C"}
	for _, key := range keys {
		val, found := c.Get(key)
		if !found {
			t.Errorf("expected %s to exist after update", key)
		}
		if key == "B" && val != "updated" {
			t.Errorf("expected B to have value 'updated', got %s", val)
		} else if key != "B" && val != key {
			t.Errorf("expected %s to have value %s, got %s", key, key, val)
		}
	}

	// Verify cache size is still 3
	if c.Size() != 3 {
		t.Errorf("expected cache size to be 3, got %d", c.Size())
	}
}

// Test 5: TTL Basic Expiration
// Verifies that a key expires after the TTL
func TestTTL_BasicExpiration(t *testing.T) {
    c := New(10)
    c.SetWithTTL("key", "value", 100*time.Millisecond)
    
    // Should exist immediately
    val, found := c.Get("key")
    if !found || val != "value" {
        t.Error("key should exist before expiration")
    }
    
    // Wait for expiration
    time.Sleep(150 * time.Millisecond)
    
    // Should be gone
    _, found = c.Get("key")
    if found {
        t.Error("key should be expired")
    }
}

// Test 6: Set without TTL
// Verifies that a key without a TTL does not expire
func TestTTL_SetWithoutTTL(t *testing.T) {
    c := New(10)
    c.Set("key", "value")  // No TTL
    
    time.Sleep(100 * time.Millisecond)
    
    // Should still exist (no expiration)
    val, found := c.Get("key")
    if !found || val != "value" {
        t.Error("key without TTL should not expire")
    }
}

// Test 7: Set with TTL and update without TTL clears TTL
// Verifies that a key with a TTL set and updated without a TTL clears the TTL
func TestTTL_SetWithTTLAndUpdateWithoutTTLClearsTTL(t *testing.T) {
    c := New(10)
    c.SetWithTTL("key", "value1", 100*time.Millisecond)
    c.Set("key", "value2")  // Update without TTL
    
    time.Sleep(150 * time.Millisecond)
    
    // Should still exist (TTL was cleared)
    val, found := c.Get("key")
    if !found || val != "value2" {
        t.Error("updated key should not expire when TTL cleared")
    }
}

// Test 8: Set with TTL and update with new TTL
// Verifies that a key with a TTL set and updated with a new TTL extends the TTL
func TestTTL_UpdateWithNewTTL(t *testing.T) {
    c := New(10)
    c.SetWithTTL("key", "value1", 500*time.Millisecond)
    c.SetWithTTL("key", "value2", 100*time.Millisecond)  // Shorter TTL
    
    time.Sleep(150 * time.Millisecond)
    
    // Should be expired (new TTL applied)
    _, found := c.Get("key")
    if found {
        t.Error("key should expire with new shorter TTL")
    }
}

// Test 9: Keys filters expired keys
func TestTTL_KeysFiltersExpired(t *testing.T) {
    c := New(10)
    c.Set("permanent", "value")
    c.SetWithTTL("temporary", "value", 50*time.Millisecond)
    
    time.Sleep(100 * time.Millisecond)
    
    keys := c.Keys()
    if len(keys) != 1 || keys[0] != "permanent" {
        t.Errorf("Keys() should only return non-expired keys, got %v", keys)
    }
}

// Test 10: Expired evicted before LRU
// Verifies that expired keys are evicted before the least recently used key
func TestTTL_ExpiredEvictedBeforeLRU(t *testing.T) {
    c := New(3)
    c.SetWithTTL("A", "1", 50*time.Millisecond)  // Will expire
    c.Set("B", "2")  // Permanent
    c.Set("C", "3")  // Permanent
    
    time.Sleep(100 * time.Millisecond)  // A expires
    
    c.Set("D", "4")  // Should evict expired A, not B
    
    // B and C should still exist (A was evicted, not them)
    if _, found := c.Get("B"); !found {
        t.Error("B should still exist")
    }
    if _, found := c.Get("C"); !found {
        t.Error("C should still exist")
    }
    if _, found := c.Get("D"); !found {
        t.Error("D should exist")
    }
}

// Test 11: Multiple expired evicted
// Verifies that expired keys are evicted before the least recently used key
func TestTTL_MultipleExpiredEvicted(t *testing.T) {
    c := New(3)
    c.SetWithTTL("A", "1", 50*time.Millisecond)  // Will expire
    c.SetWithTTL("B", "2", 50*time.Millisecond)  // Will expire
    c.Set("C", "3")  // Permanent
    
    time.Sleep(100 * time.Millisecond)  // A and B expire
    
    c.Set("D", "4")  // Should evict expired A and B, not C
    
    // C should still exist (A and B were evicted, not them)
    if _, found := c.Get("C"); !found {
        t.Error("C should still exist")
    }
    if _, found := c.Get("D"); !found {
        t.Error("D should exist")
    }
}


///////////////////////////////
// Benchmarks
///////////////////////////////

// go test -bench=. -benchmem
func BenchmarkCacheGet(b *testing.B) {
	c := New(1000)
	c.Set("key", "value")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get("key")
		}
	})
}

func BenchmarkCacheSet(b *testing.B) {
	c := New(1000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set("key", "value")
		}
	})
}

func BenchmarkCacheMixed(b *testing.B) {
	c := New(1000)
	// Pre-populate
	for i := 0; i < 100; i++ {
		c.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// 80% reads, 20% writes (realistic workload)
			if i%5 == 0 {
				c.Set(fmt.Sprintf("key%d", i%100), "value")
			} else {
				c.Get(fmt.Sprintf("key%d", i%100))
			}
			i++
		}
	})
}
