package socket

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
)

// Server handles Unix domain socket communication
type Server struct {
	listener  net.Listener
	clients   map[net.Conn]bool
	mu        sync.RWMutex
	onCommand func(conn net.Conn, msg SocketMessage)
}

// NewServer creates a new socket server
func NewServer(path string, onCommand func(net.Conn, SocketMessage)) (*Server, error) {
	// Remove existing socket file
	os.Remove(path)

	listener, err := net.Listen("unix", path)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket: %w", err)
	}

	// Set socket permissions (read/write for owner)
	if err := os.Chmod(path, 0600); err != nil {
		listener.Close()
		return nil, fmt.Errorf("failed to set socket permissions: %w", err)
	}

	server := &Server{
		listener:  listener,
		clients:   make(map[net.Conn]bool),
		onCommand: onCommand,
	}

	// Start accepting connections
	go server.accept()

	return server, nil
}

// accept waits for new client connections
func (s *Server) accept() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.listener != nil {
				fmt.Printf("Error accepting connection: %v\n", err)
			}
			return
		}

		s.mu.Lock()
		s.clients[conn] = true
		s.mu.Unlock()

		go s.handleConnection(conn)
	}
}

// handleConnection reads messages from a client
func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Bytes()
		var msg SocketMessage
		if err := json.Unmarshal(line, &msg); err != nil {
			fmt.Printf("Error parsing message: %v\n", err)
			continue
		}

		// Handle command from client
		if s.onCommand != nil {
			s.onCommand(conn, msg)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Connection error: %v\n", err)
	}
}

// Broadcast sends a message to all connected clients
func (s *Server) Broadcast(msg SocketMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	data = append(data, '\n')

	// Phase 1: iterate under RLock, collect dead clients
	s.mu.RLock()
	var dead []net.Conn
	for client := range s.clients {
		_, err := client.Write(data)
		if err != nil {
			fmt.Printf("Error writing to client: %v\n", err)
			dead = append(dead, client)
		}
	}
	s.mu.RUnlock()

	// Phase 2: clean up dead clients under Lock
	if len(dead) > 0 {
		s.mu.Lock()
		for _, client := range dead {
			delete(s.clients, client)
			client.Close()
		}
		s.mu.Unlock()
	}

	return nil
}

// SendMessage sends a single message to a specific connection
func SendMessage(conn net.Conn, msg SocketMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	data = append(data, '\n')
	_, err = conn.Write(data)
	return err
}

// Close shuts down the server
func (s *Server) Close() error {
	if s.listener != nil {
		err := s.listener.Close()
		s.listener = nil
		return err
	}
	return nil
}

// ClientCount returns the number of connected clients
func (s *Server) ClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}
