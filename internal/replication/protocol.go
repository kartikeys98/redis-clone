package replication

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type OpType string

const (
	OpSet    OpType = "SET"
	OpDelete OpType = "DELETE"
	OpFlush  OpType = "FLUSH"
	OpPing   OpType = "PING"
	OpPong   OpType = "PONG"
)

type Operation struct {
	Type      OpType
	Key       string
	Value     string
	TTL       time.Duration
	Timestamp int64
}

// Serialize operation to wire format
// Format varies by operation type
func (op *Operation) String() string {
	switch op.Type {
	case OpSet:
		ttlMillis := op.TTL.Milliseconds()
		return fmt.Sprintf("%s %s %s %d %d\n", op.Type, op.Key, op.Value, ttlMillis, op.Timestamp)
	case OpDelete:
		return fmt.Sprintf("%s %s %d\n", op.Type, op.Key, op.Timestamp)
	case OpFlush:
		return fmt.Sprintf("%s %d\n", op.Type, op.Timestamp)
	case OpPing, OpPong:
		return fmt.Sprintf("%s %d\n", op.Type, op.Timestamp)
	default:
		return fmt.Sprintf("%s\n", op.Type)
	}
}

// Parse operation from wire format
func ParseOperation(line string) (*Operation, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid operation format")
	}

	op := &Operation{
		Type: OpType(parts[0]),
	}

	// Parse based on type
	switch op.Type {
	case OpSet:
		if len(parts) < 5 {
			return nil, fmt.Errorf("SET requires 5 parts")
		}
		op.Key = parts[1]
		op.Value = parts[2]

		ttlMillis, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return nil, err
		}
		op.TTL = time.Duration(ttlMillis) * time.Millisecond

		op.Timestamp, err = strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			return nil, err
		}

	case OpDelete:
		if len(parts) < 3 {
			return nil, fmt.Errorf("DELETE requires 3 parts")
		}
		op.Key = parts[1]

		var err error
		op.Timestamp, err = strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, err
		}

	case OpFlush:
		if len(parts) < 2 {
			return nil, fmt.Errorf("FLUSH requires 2 parts")
		}
		var err error
		op.Timestamp, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}

	case OpPing, OpPong:
		if len(parts) < 2 {
			return nil, fmt.Errorf("%s requires 2 parts", op.Type)
		}
		var err error
		op.Timestamp, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return op, nil
}
