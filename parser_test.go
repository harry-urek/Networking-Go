package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

// MockLogger is a simplified logger for testing
type MockLogger struct{}

// LogCommand simulates logging a command (for testing purposes)
func (l *MockLogger) LogCommand(clientID, command string) {
	fmt.Printf("Logged command: %s\n", command)
}

// TestParserOneCommand tests the ParseOne method for parsing a single RESP command
func TestParserOneCommand(t *testing.T) {
	mockLogger := &MockLogger{}

	// Example RESP commands
	testCases := []struct {
		input    string
		expected interface{}
	}{
		{"+OK\r\n", "OK"},                                                 // Simple string
		{"-Error message\r\n", "Error message"},                           // Error string
		{":1000\r\n", int64(1000)},                                        // Integer
		{"$6\r\nfoobar\r\n", "foobar"},                                    // Bulk string
		{"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", []interface{}{"foo", "bar"}}, // Array
	}

	for _, tc := range testCases {
		conn := strings.NewReader(tc.input)
		parser := &RESPParser{
			logger: mockLogger,
			conn:   conn,
			buf:    &bytes.Buffer{},
			tbuf:   make([]byte, 1024), // Temporary buffer for testing
		}

		// Parsing the input
		result, err := parser.ParseOne()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Check the result
		if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
			t.Errorf("Expected %v, got %v", tc.expected, result)
		} else {
			fmt.Printf("Test passed for input: %s, result: %v\n", tc.input, result)
		}
	}
}

// TestParserMultipleCommands tests the ParseMultiple method for parsing multiple RESP commands
func TestParserMultipleCommands(t *testing.T) {
	mockLogger := &MockLogger{}

	// Example RESP commands
	input := "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	expected := []interface{}{
		[]interface{}{"foo", "bar"},
		[]interface{}{"hello", "world"},
	}

	conn := strings.NewReader(input)
	parser := &RESPParser{
		logger: mockLogger,
		conn:   conn,
		buf:    &bytes.Buffer{},
		tbuf:   make([]byte, 1024),
	}

	result, err := parser.ParseMultiple()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check the result
	if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	} else {
		fmt.Printf("Test passed for input: %s, result: %v\n", input, result)
	}
}

// Main function for running the tests
func main() {
	fmt.Println("Testing single command parsing:")
	TestParserOneCommand(nil)

	fmt.Println("\nTesting multiple command parsing:")
	TestParserMultipleCommands(nil)
}
