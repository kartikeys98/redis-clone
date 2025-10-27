package cache

import (
	"fmt"
	"sync"
	"testing"
)

func TestBasicOperations(t *testing.T) {
	c := New()

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
	c := New()

	// Get a key that was never set
	_, found := c.Get("nonexistent")
	if found {
		t.Error("should not find nonexistent key")
	}
}

func TestConcurrentOperations(t *testing.T) {
	c := New()
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
			t.Errorf("expected to find key%d", i) // ✅ Changed to Errorf
		}
		if got != fmt.Sprintf("value%d", i) {
			t.Errorf("expected 'value%d', got '%s'", i, got)
		}
	}
}
