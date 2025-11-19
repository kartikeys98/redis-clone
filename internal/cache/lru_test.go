package cache

import (
	"testing"
)

// ========================================
// AddToFront Tests
// ========================================

func TestAddToFront_EmptyList(t *testing.T) {
	list := &LRUList{}

	node := list.AddToFront("first")

	// Verify node returned
	if node == nil {
		t.Fatal("AddToFront should return a node")
	}
	if node.Key != "first" {
		t.Errorf("Expected key 'first', got '%s'", node.Key)
	}

	// Verify list state
	if list.Head != node {
		t.Error("Head should point to the new node")
	}
	if list.Tail != node {
		t.Error("Tail should point to the new node (only node)")
	}
	if list.Size != 1 {
		t.Errorf("Size should be 1, got %d", list.Size)
	}

	// Verify node pointers
	if node.Prev != nil {
		t.Error("Single node's Prev should be nil")
	}
	if node.Next != nil {
		t.Error("Single node's Next should be nil")
	}
}

func TestAddToFront_MultipleNodes(t *testing.T) {
	list := &LRUList{}

	// Add three nodes: A, B, C
	nodeA := list.AddToFront("A")
	nodeB := list.AddToFront("B")
	nodeC := list.AddToFront("C")

	// List should be: C -> B -> A
	if list.Head.Key != "C" {
		t.Errorf("Head should be C, got %s", list.Head.Key)
	}
	if list.Tail.Key != "A" {
		t.Errorf("Tail should be A, got %s", list.Tail.Key)
	}
	if list.Size != 3 {
		t.Errorf("Size should be 3, got %d", list.Size)
	}

	// Verify links
	if list.Head.Next != nodeB {
		t.Error("Head.Next should be B")
	}
	if nodeB.Prev != nodeC {
		t.Error("B.Prev should be C")
	}
	if nodeB.Next != nodeA {
		t.Error("B.Next should be A")
	}
	if list.Tail.Prev != nodeB {
		t.Error("Tail.Prev should be B")
	}
}

// ========================================
// RemoveLRU Tests
// ========================================

func TestRemoveLRU_EmptyList(t *testing.T) {
	list := &LRUList{}

	node := list.RemoveLRU()

	if node != nil {
		t.Errorf("expected nil, got %#v", node)
	}
	if list.Size != 0 {
		t.Errorf("Size should remain 0, got %d", list.Size)
	}
}

func TestRemoveLRU_SingleNode(t *testing.T) {
	list := &LRUList{}
	list.AddToFront("only")

	node := list.RemoveLRU()

	if node == nil {
		t.Fatal("expected a node, got nil")
	}
	if node.Key != "only" {
		t.Errorf("Expected 'only', got '%s'", node.Key)
	}

	if list.Head != nil {
		t.Error("Head should be nil after removing only node")
	}
	if list.Tail != nil {
		t.Error("Tail should be nil after removing only node")
	}
	if list.Size != 0 {
		t.Errorf("Size should be 0, got %d", list.Size)
	}
}

func TestRemoveLRU_MultipleNodes(t *testing.T) {
	list := &LRUList{}
	list.AddToFront("A")
	list.AddToFront("B")
	list.AddToFront("C")

	// List: C -> B -> A
	// Remove A (tail)
	node := list.RemoveLRU()

	if node == nil {
		t.Fatal("expected a node, got nil")
	}
	if node.Key != "A" {
		t.Errorf("Expected 'A', got '%s'", node.Key)
	}
	if list.Tail.Key != "B" {
		t.Errorf("New tail should be B, got %s", list.Tail.Key)
	}
	if list.Tail.Next != nil {
		t.Error("New tail's Next should be nil")
	}
	if list.Size != 2 {
		t.Errorf("Size should be 2, got %d", list.Size)
	}

	// Remove B (new tail)
	node = list.RemoveLRU()
	if node == nil {
		t.Fatal("expected a node, got nil")
	}
	if node.Key != "B" {
		t.Errorf("Expected 'B', got '%s'", node.Key)
	}
	if list.Tail.Key != "C" {
		t.Errorf("New tail should be C, got %s", list.Tail.Key)
	}
	if list.Head != list.Tail {
		t.Error("With one node, Head should equal Tail")
	}
}

// ========================================
// MoveToFront Tests
// ========================================

func TestMoveToFront_AlreadyAtHead(t *testing.T) {
	list := &LRUList{}
	tail := list.AddToFront("A")
	list.AddToFront("B")
	list.AddToFront("C")

	// List: C -> B -> A
	headBefore := list.Head

	list.MoveToFront(list.Head)

	// Should be no-op
	if list.Head != headBefore {
		t.Error("Head should not change")
	}
	if list.Head.Key != "C" {
		t.Errorf("Head should still be C, got %s", list.Head.Key)
	}
	if list.Tail != tail {
		t.Error("Tail should remain A")
	}
}

func TestMoveToFront_FromTail(t *testing.T) {
	list := &LRUList{}
	nodeA := list.AddToFront("A")
	list.AddToFront("B")
	list.AddToFront("C")

	// List: C -> B -> A
	// Move A (tail) to front
	list.MoveToFront(nodeA)

	// List should be: A -> C -> B
	if list.Head.Key != "A" {
		t.Errorf("Head should be A, got %s", list.Head.Key)
	}
	if list.Tail.Key != "B" {
		t.Errorf("Tail should be B, got %s", list.Tail.Key)
	}
	if list.Head.Next.Key != "C" {
		t.Errorf("A.Next should be C, got %s", list.Head.Next.Key)
	}
	if list.Tail.Prev.Key != "C" {
		t.Errorf("B.Prev should be C, got %s", list.Tail.Prev.Key)
	}
}

func TestMoveToFront_FromMiddle(t *testing.T) {
	list := &LRUList{}
	list.AddToFront("A")
	nodeB := list.AddToFront("B")
	list.AddToFront("C")

	// List: C -> B -> A
	// Move B to front
	list.MoveToFront(nodeB)

	// List should be: B -> C -> A
	if list.Head.Key != "B" {
		t.Errorf("Head should be B, got %s", list.Head.Key)
	}
	if list.Tail.Key != "A" {
		t.Errorf("Tail should be A, got %s", list.Tail.Key)
	}
	if list.Head.Next.Key != "C" {
		t.Errorf("B.Next should be C, got %s", list.Head.Next.Key)
	}
	if list.Tail.Prev.Key != "C" {
		t.Errorf("A.Prev should be C, got %s", list.Tail.Prev.Key)
	}
}

func TestMoveToFront_TwoNodes(t *testing.T) {
	list := &LRUList{}
	nodeA := list.AddToFront("A")
	list.AddToFront("B")

	// List: B -> A
	// Move A to front
	list.MoveToFront(nodeA)

	// List should be: A -> B
	if list.Head.Key != "A" {
		t.Errorf("Head should be A, got %s", list.Head.Key)
	}
	if list.Tail.Key != "B" {
		t.Errorf("Tail should be B, got %s", list.Tail.Key)
	}
	if list.Head.Next != list.Tail {
		t.Error("Head.Next should point to Tail")
	}
	if list.Tail.Prev != list.Head {
		t.Error("Tail.Prev should point to Head")
	}
}

// ========================================
// Remove Tests
// ========================================

func TestRemove_SingleNode(t *testing.T) {
	list := &LRUList{}
	node := list.AddToFront("only")

	list.Remove(node)

	if list.Head != nil {
		t.Error("Head should be nil")
	}
	if list.Tail != nil {
		t.Error("Tail should be nil")
	}
	if list.Size != 0 {
		t.Errorf("Size should be 0, got %d", list.Size)
	}
}

func TestRemove_Head(t *testing.T) {
	list := &LRUList{}
	list.AddToFront("A")
	list.AddToFront("B")
	nodeC := list.AddToFront("C")

	// List: C -> B -> A
	list.Remove(nodeC)

	// List should be: B -> A
	if list.Head.Key != "B" {
		t.Errorf("Head should be B, got %s", list.Head.Key)
	}
	if list.Head.Prev != nil {
		t.Error("New head's Prev should be nil")
	}
	if list.Size != 2 {
		t.Errorf("Size should be 2, got %d", list.Size)
	}
}

func TestRemove_Tail(t *testing.T) {
	list := &LRUList{}
	nodeA := list.AddToFront("A")
	list.AddToFront("B")
	list.AddToFront("C")

	// List: C -> B -> A
	list.Remove(nodeA)

	// List should be: C -> B
	if list.Tail.Key != "B" {
		t.Errorf("Tail should be B, got %s", list.Tail.Key)
	}
	if list.Tail.Next != nil {
		t.Error("New tail's Next should be nil")
	}
	if list.Size != 2 {
		t.Errorf("Size should be 2, got %d", list.Size)
	}
}

func TestRemove_Middle(t *testing.T) {
	list := &LRUList{}
	list.AddToFront("A")
	nodeB := list.AddToFront("B")
	list.AddToFront("C")

	// List: C -> B -> A
	list.Remove(nodeB)

	// List should be: C -> A
	if list.Head.Key != "C" {
		t.Errorf("Head should be C, got %s", list.Head.Key)
	}
	if list.Tail.Key != "A" {
		t.Errorf("Tail should be A, got %s", list.Tail.Key)
	}
	if list.Head.Next != list.Tail {
		t.Error("Head.Next should be Tail")
	}
	if list.Tail.Prev != list.Head {
		t.Error("Tail.Prev should be Head")
	}
	if list.Size != 2 {
		t.Errorf("Size should be 2, got %d", list.Size)
	}
}

// ========================================
// Complex Scenario Tests
// ========================================

func TestComplexScenario(t *testing.T) {
	list := &LRUList{}

	// Build the list by pushing nodes: C -> B -> A
	nodeA := list.AddToFront("A")
	nodeB := list.AddToFront("B")
	nodeC := list.AddToFront("C")

	// Move A to front: A -> C -> B
	list.MoveToFront(nodeA)
	if list.Head != nodeA {
		t.Error("Head should be A after move")
	}

	// Remove C: A -> B
	list.Remove(nodeC)
	if list.Size != 2 {
		t.Errorf("Size should be 2, got %d", list.Size)
	}

	// RemoveLRU (remove B): A
	node := list.RemoveLRU()
	if node != nodeB {
		t.Errorf("Expected to remove B, got %v", node)
	}
	if list.Head != list.Tail {
		t.Error("With one node, Head should equal Tail")
	}

	// RemoveLRU (remove A): empty
	node = list.RemoveLRU()
	if node != nodeA {
		t.Errorf("Expected to remove A, got %v", node)
	}
	if list.Head != nil || list.Tail != nil {
		t.Error("List should be empty")
	}

	// RemoveLRU on empty: should return nil
	node = list.RemoveLRU()
	if node != nil {
		t.Errorf("Expected nil, got %v", node)
	}
}
