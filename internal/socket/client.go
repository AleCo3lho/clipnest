package socket

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// Client connects to the daemon via Unix socket
type Client struct {
	conn    net.Conn
	scanner *bufio.Scanner
}

// NewClient dials the daemon socket
func NewClient(socketPath string) (*Client, error) {
	conn, err := net.DialTimeout("unix", socketPath, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to daemon: %w", err)
	}

	return &Client{
		conn:    conn,
		scanner: bufio.NewScanner(conn),
	}, nil
}

// Send sends a SocketMessage to the daemon
func (c *Client) Send(msg SocketMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	data = append(data, '\n')
	_, err = c.conn.Write(data)
	return err
}

// Receive reads one JSON response from the daemon
func (c *Client) Receive() (SocketMessage, error) {
	_ = c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	if !c.scanner.Scan() {
		if err := c.scanner.Err(); err != nil {
			return SocketMessage{}, fmt.Errorf("read error: %w", err)
		}
		return SocketMessage{}, fmt.Errorf("connection closed")
	}

	var msg SocketMessage
	if err := json.Unmarshal(c.scanner.Bytes(), &msg); err != nil {
		return SocketMessage{}, fmt.Errorf("failed to parse response: %w", err)
	}

	return msg, nil
}

// Close closes the connection
func (c *Client) Close() error {
	return c.conn.Close()
}
