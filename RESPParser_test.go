package parser

import (
	"errors"
	"strings"
	"testing"
)

func TestParseRESP(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "Simple String",
			input:    "+OK\r\n",
			expected: "OK",
		},
		{
			name:     "Error",
			input:    "-Error message\r\n",
			expected: errors.New("Error message"),
		},
		{
			name:     "Integer",
			input:    ":1000\r\n",
			expected: int64(1000),
		},
		{
			name:     "Bulk String",
			input:    "$6\r\nfoobar\r\n",
			expected: "foobar",
		},
		{
			name:     "Array",
			input:    "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			expected: []interface{}{"foo", "bar"},
		},
		{
			name:     "Null Bulk String",
			input:    "$-1\r\n",
			expected: nil,
		},
		{
			name:     "Null Array",
			input:    "*-1\r\n",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := ParseRESP(reader)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !equal(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func equal(a, b interface{}) bool {
	switch a := a.(type) {
	case *string:
		if a == nil {
			return b == nil
		}
		return *a == b
	case error:
		b, ok := b.(error)
		return ok && a.Error() == b.Error()
	case *[]interface{}:
		if a == nil {
			return b == nil
		}
		b, ok := b.([]interface{})
		if !ok || len(*a) != len(b) {
			return false
		}
		for i := range *a {
			if !equal((*a)[i], b[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}
