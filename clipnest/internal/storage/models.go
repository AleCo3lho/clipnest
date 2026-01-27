package storage

import "time"

// Clip represents a clipboard entry
type Clip struct {
	ID        int64
	Content   string
	Type      string // "text", "image"
	Timestamp time.Time
	Pinned    bool
}

// Message types for socket communication
type SocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// ClipData for socket messages
type ClipData struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Pinned    bool      `json:"pinned"`
}

// CommandData for UI commands
type CommandData struct {
	ID int64 `json:"id"`
}
