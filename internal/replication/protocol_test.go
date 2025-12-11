package replication

import (
	"testing"
	"time"
)

func TestOperationSerializeDeserialize(t *testing.T) {
	tests := []struct {
		name string
		op   *Operation
	}{
		{
			name: "SET with TTL",
			op: &Operation{
				Type:      OpSet,
				Key:       "mykey",
				Value:     "myvalue",
				TTL:       60 * time.Second,
				Timestamp: 1234567890,
			},
		},
		{
			name: "SET without TTL",
			op: &Operation{
				Type:      OpSet,
				Key:       "key2",
				Value:     "value2",
				TTL:       0,
				Timestamp: 1234567890,
			},
		},
		{
			name: "DELETE",
			op: &Operation{
				Type:      OpDelete,
				Key:       "oldkey",
				Timestamp: 1234567890,
			},
		},
		{
			name: "FLUSH",
			op: &Operation{
				Type:      OpFlush,
				Timestamp: 1234567890,
			},
		},
		{
			name: "PING",
			op: &Operation{
				Type:      OpPing,
				Timestamp: 1234567890,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			serialized := tt.op.String()

			// Deserialize
			parsed, err := ParseOperation(serialized)
			if err != nil {
				t.Fatalf("ParseOperation failed: %v", err)
			}

			// Compare
			if parsed.Type != tt.op.Type {
				t.Errorf("Type mismatch: got %v, want %v", parsed.Type, tt.op.Type)
			}
			if parsed.Key != tt.op.Key {
				t.Errorf("Key mismatch: got %v, want %v", parsed.Key, tt.op.Key)
			}
			if parsed.Value != tt.op.Value {
				t.Errorf("Value mismatch: got %v, want %v", parsed.Value, tt.op.Value)
			}
			if parsed.TTL != tt.op.TTL {
				t.Errorf("TTL mismatch: got %v, want %v", parsed.TTL, tt.op.TTL)
			}
			if parsed.Timestamp != tt.op.Timestamp {
				t.Errorf("Timestamp mismatch: got %v, want %v", parsed.Timestamp, tt.op.Timestamp)
			}
		})
	}
}

func TestParseOperationErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"too few parts", "SET key"},
		{"SET missing parts", "SET key value 60"}, // Missing timestamp
		{"DELETE missing timestamp", "DELETE key"},
		{"invalid TTL", "SET key value abc 123"},
		{"invalid timestamp", "SET key value 60 abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseOperation(tt.input)
			if err == nil {
				t.Errorf("Expected error for input %q, got nil", tt.input)
			}
		})
	}
}
