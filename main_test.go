package main

import (
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	s, err := newServer(":8080")
	if err != nil {
		t.Fatal(err)
	}
	s.Start()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	expected := "Welcome to my TCP server!\n The connection is getting handled by the server"
	actual := make([]byte, len(expected))

	if _, err := conn.Read(actual); err != nil {
		t.Fatal(err)
	}

	if string(actual) != expected {
		t.Error("expected %q, but got %q", expected, actual)

	}

	//  Stop the Server

	s.Stop()

}
