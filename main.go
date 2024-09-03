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

type server struct {
	wg         sync.WaitGroup
	listener   net.Listener
	shutdown   chan struct{}
	connection chan net.Conn
}

func newServer(addr string) (*server, error) {
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
