package services

import (
	"bytes"
	"errors"
	/* "io" */
	"testing"
)

// Mock connection to simulate a ReadWriter
type mockConn struct {
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
}

func newMockConn(input string) *mockConn {
	return &mockConn{
		readBuffer:  bytes.NewBufferString(input),
		writeBuffer: new(bytes.Buffer),
	}
}

func (m *mockConn) Read(p []byte) (int, error) {
	return m.readBuffer.Read(p)
}

func (m *mockConn) Write(p []byte) (int, error) {
	return m.writeBuffer.Write(p)
}

// Helper function for initializing the parser
func newTestParser(input string) *RESPParser {
	conn := newMockConn(input)
	logger := &Logger{file: nil} // No actual logging for the tests
	parser := new(RESPParser)
	return parser.NewRESP(conn, logger)
}

// Test simple RESP "+" message
func TestParseSimpleString(t *testing.T) {
	parser := newTestParser("+OK\r\n")

	result, err := parser.ParseOne()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != "OK" {
		t.Errorf("Expected 'OK', got '%v'", result)
	}
}

// Test RESP "-" error message
func TestParseError(t *testing.T) {
	parser := newTestParser("-Error message\r\n")

	result, err := parser.ParseOne()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != errors.New("Error message") {
		t.Errorf("Expected 'Error message', got '%v'", result)
	}
}

// Test RESP ":" integer message
func TestParseInteger(t *testing.T) {
	parser := newTestParser(":1000\r\n")

	result, err := parser.ParseOne()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != int64(1000) {
		t.Errorf("Expected '1000', got '%v'", result)
	}
}

// Test RESP "$" bulk string
func TestParseBulkString(t *testing.T) {
	parser := newTestParser("$6\r\nfoobar\r\n")

	result, err := parser.ParseOne()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != "foobar" {
		t.Errorf("Expected 'foobar', got '%v'", result)
	}
}

// Test RESP "*" array
func TestParseArray(t *testing.T) {
	parser := newTestParser("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")

	result, err := parser.ParseOne()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	array, ok := result.([]interface{})
	if !ok || len(array) != 2 {
		t.Fatalf("Expected array of length 2, got %v", result)
	}

	if array[0] != "foo" || array[1] != "bar" {
		t.Errorf("Expected ['foo', 'bar'], got %v", array)
	}
}

// Test parsing multiple RESP messages in a single input stream
func TestParseMultiple(t *testing.T) {
	parser := newTestParser("+OK\r\n:123\r\n$6\r\nfoobar\r\n")

	results, err := parser.ParseMultiple()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	if results[0] != "OK" || results[1] != int64(123) || results[2] != "foobar" {
		t.Errorf("Unexpected results: %v", results)
	}
}
