package parser

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestParseRESP(t *testing.T) {

	nilByteArray := []byte(nil)
	nilInterfaceArray := []interface{}(nil)

	tests := []struct {
		name   string
		resp   string
		native interface{}
	}{
		{
			name:   "Simple String",
			resp:   "+OK\r\n",
			native: "OK",
		},
		{
			name:   "Error",
			resp:   "-Error message\r\n",
			native: errors.New("Error message"),
		},
		{
			name:   "Integer",
			resp:   ":1000\r\n",
			native: int64(1000),
		},
		{
			name:   "Bulk String",
			resp:   "$6\r\nfoobar\r\n",
			native: []byte("foobar"),
		},
		{
			name:   "Array",
			resp:   "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			native: []interface{}{[]byte("foo"), []byte("bar")},
		},
		{
			name:   "Empty Bulk String",
			resp:   "$0\r\n\r\n",
			native: []byte(""),
		},
		{
			name:   "Null Bulk String",
			resp:   "$-1\r\n",
			native: nilByteArray,
		},
		{
			name:   "Null Array",
			resp:   "*-1\r\n",
			native: nilInterfaceArray,
		},
		{
			name:   "Empty Array",
			resp:   "*0\r\n",
			native: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run("Parse "+tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.resp)
			result, err := ParseRESP(reader)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !equal(result, tt.native) {
				t.Errorf("expected %v, got %v", tt.native, result)
			}
		})

		t.Run("Generate "+tt.name, func(t *testing.T) {
			result, err := GenerateRESP(tt.native)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !equal(result, tt.resp) {
				t.Errorf("expected %v, got %v", tt.resp, result)
			}
		})
	}
}

func equal(a, b interface{}) bool {
	switch a := a.(type) {
	case []byte:
		if a == nil {
			return reflect.ValueOf(b).IsNil()
		}
		return string(a) == string(b.([]byte))
	case error:
		b, ok := b.(error)
		return ok && a.Error() == b.Error()
	case []interface{}:
		if a == nil {
			return reflect.ValueOf(b).IsNil()
		}
		b, ok := b.([]interface{})
		if !ok || len(a) != len(b) {
			return false
		}
		for i := range a {
			if !equal((a)[i], b[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}
