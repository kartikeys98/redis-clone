package cache

type Node struct {
	Key  string
	Prev *Node
	Next *Node
}

type LRUList struct {
	Head *Node
	Tail *Node
	Size int
}

func (l *LRUList) AddToFront(key string) *Node {
	node := &Node{
		Key:  key,
		Prev: nil,
		Next: nil,
	}
	if l.Head == nil {
		l.Head = node
		l.Tail = node
	} else {
		node.Next = l.Head
		l.Head.Prev = node
		l.Head = node
	}
	l.Size++
	return node
}

func (l *LRUList) MoveToFront(node *Node) {
	if node == l.Head {
		return
	}
	if node == l.Tail {
		l.Tail = node.Prev
		l.Tail.Next = nil
	} else {
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
	}
	node.Next = l.Head
	node.Prev = nil
	node.Next.Prev = node
	l.Head = node
}

func (l *LRUList) RemoveLRU() string {
	if l.Tail == nil {
		return ""
	}
	node := l.Tail
	l.Tail = node.Prev
	if l.Tail != nil {
		l.Tail.Next = nil
	} else {
		l.Head = nil
	}
	node.Prev = nil
	node.Next = nil
	l.Size--
	key := node.Key
	return key
}

func (l *LRUList) Remove(node *Node) {
	if node == l.Head {
		l.Head = node.Next
		if l.Head != nil {
			l.Head.Prev = nil
		} else {
			l.Tail = nil
		}
	} else if node == l.Tail {
		l.Tail = node.Prev
		l.Tail.Next = nil
	} else {
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
	}
	node.Prev = nil
	node.Next = nil
	l.Size--
}
