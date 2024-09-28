package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	s, err := newServer(":8080")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	s.Start()

	// Wait for a SIGINT or SIGTERM signal to signal fo rshutdown

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting Down The Server ... ... .. . . .. . ")
	s.Stop()
	fmt.Println("Server Shut Down - / | / - | /")
}

/*

Create a new net.Listener object that listens for incoming connections on a specified address using the net.Listen function.
Start two goroutines to handle incoming connections concurrently: one to accept new connections and another to handle them.

In the acceptConnections goroutine, use a for loop and a select statement to listen for incoming connections on the listener and send them over a channel.
In the handleConnections goroutine, use a for loop and a select statement to receive connections from the channel and handle them in separate goroutines.
In the handleConnection function, handle the incoming connection by performing any necessary processing or sending data to the client.

Implement graceful shutdown by creating a shutdown channel that signals to the goroutines that they should stop processing connections.
Close the listener and wait for the goroutines to finish using a sync.WaitGroup.

*/

type server struct {
	wg         sync.WaitGroup
	listener   net.Listener
	shutdown   chan struct{}
	connection chan net.Conn
}

func newServer(addr string) (*server, error) {
	// Create a new net.Listener object that listens for incoming connections on a specified address
	// using the net.Listen function.

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Failed to listen on address %s : %w", addr, err)
	}
	return &server{
		listener:   listener,
		shutdown:   make(chan struct{}),
		connection: make(chan net.Conn),
	}, nil

}

func (s *server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				continue

			}
			s.connection <- conn
		}

	}

}
func (s *server) handleConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		case conn := <-s.connection:
			go s.handleConnection(conn)

		}

	}

}

func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Logic for handle func
	fmt.Fprintf(conn, "Welcome to my TCP server!\n The connection is getting handled by the server")
	time.Sleep(5 * time.Second)
	fmt.Fprintf(conn, "Connection is handled\n Goodbye!\n")

}
func (s *server) Start() {
	// Start two goroutines to handle incoming connections concurrently:
	// one to accept new connections and another to handle them.
	s.wg.Add(2)
	go s.acceptConnections()
	go s.handleConnections()

}

func (s *server) Stop() {
	close(s.shutdown)
	s.listener.Close()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		fmt.Println("Timed out ... !! Waiting for connections to finish .. !!")
		return

	}
}
