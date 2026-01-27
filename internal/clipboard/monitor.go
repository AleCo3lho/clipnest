package clipboard

import (
	"time"

	"github.com/atotto/clipboard"
)

// Monitor watches clipboard for changes
type Monitor struct {
	lastContent string
	lastType    string
	onChange    func(content, clipType string)
	interval    time.Duration
	running     bool
}

// NewMonitor creates a new clipboard monitor
func NewMonitor(interval time.Duration, onChange func(content, clipType string)) *Monitor {
	return &Monitor{
		interval: interval,
		onChange: onChange,
	}
}

// Start begins monitoring clipboard changes
func (m *Monitor) Start() {
	m.running = true

	// Initialize with current clipboard
	content, clipType := m.readClipboard()
	if content != "" {
		m.lastContent = content
		m.lastType = clipType
	}

	go m.poll()
}

// Stop stops monitoring
func (m *Monitor) Stop() {
	m.running = false
}

// poll checks clipboard periodically
func (m *Monitor) poll() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for m.running {
		<-ticker.C

		content, clipType := m.readClipboard()
		if content == "" {
			continue
		}

		// Check if changed
		if content != m.lastContent || clipType != m.lastType {
			m.lastContent = content
			m.lastType = clipType

			// Notify callback
			if m.onChange != nil {
				m.onChange(content, clipType)
			}
		}
	}
}

// readClipboard reads the current clipboard content
func (m *Monitor) readClipboard() (string, string) {
	content, err := clipboard.ReadAll()
	if err != nil {
		return "", "error"
	}

	if content != "" {
		return content, "text"
	}

	return "", "unknown"
}

// Copy writes content to clipboard
func Copy(content string) error {
	return clipboard.WriteAll(content)
}

// Read reads current clipboard content
func Read() (string, string, error) {
	content, err := clipboard.ReadAll()
	if err != nil {
		return "", "error", err
	}
	return content, "text", nil
}
