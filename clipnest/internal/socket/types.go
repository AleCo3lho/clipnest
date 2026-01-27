package socket

// SocketMessage represents a message sent over the socket
type SocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// ClipData represents clip information in messages
type ClipData struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Pinned    bool   `json:"pinned"`
}

// CommandData represents a command from the client
type CommandData struct {
	ID int64 `json:"id"`
}

// CopyClipCommand sends a clip to clipboard
type CopyClipCommand struct {
	ID int64 `json:"id"`
}

// PinCommand pins a clip
type PinCommand struct {
	ID int64 `json:"id"`
}

// UnpinCommand unpins a clip
type UnpinCommand struct {
	ID int64 `json:"id"`
}

// SearchCommand searches clips
type SearchCommand struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

// ListCommand lists clips
type ListCommand struct {
	Limit int `json:"limit"`
}
