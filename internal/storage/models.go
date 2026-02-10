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
