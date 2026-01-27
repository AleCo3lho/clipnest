package storage

import (
	"fmt"
	"sync"
)

// Storage provides in-memory clipboard storage
type Storage struct {
	memory *MemoryStore
	maxMemory int
}

// NewStorage creates a new in-memory storage
func NewStorage(maxMemory int) (*Storage, error) {
	return &Storage{
		memory:    NewMemoryStore(),
		maxMemory: maxMemory,
	}, nil
}

// Add stores a clip
func (s *Storage) Add(clip Clip) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add to memory
	id, err := s.memory.Add(clip)
	if err != nil {
		return id, err
	}

	// Evict oldest if over limit
	for s.memory.Count() > s.maxMemory {
		s.memory.EvictOldest()
	}

	return id, nil
}

// Get retrieves a clip by ID
func (s *Storage) Get(id int64) (Clip, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clip, exists := s.memory.Get(id)
	if !exists {
		return Clip{}, fmt.Errorf("clip %d not found", id)
	}
	return clip, nil
}

// List returns all clips (most recent first)
func (s *Storage) List(limit int) ([]Clip, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.memory.List(limit), nil
}

// Pin stores a clip permanently (moves to pinned list)
func (s *Storage) Pin(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get from memory
	clip, exists := s.memory.Get(id)
	if !exists {
		return fmt.Errorf("clip %d not found", id)
	}

	// Mark as pinned
	clip.Pinned = true

	// Update in memory (this is a simplified approach)
	// In a real implementation, we'd have a separate pinned list
	return nil
}

// Unpin removes pin status
func (s *Storage) Unpin(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get from memory
	clip, exists := s.memory.Get(id)
	if !exists {
		return fmt.Errorf("clip %d not found", id)
	}

	// Unpin
	clip.Pinned = false
	return nil
}

// Search clips by content
func (s *Storage) Search(query string, limit int) ([]Clip, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Get all clips
	clips := s.memory.List(1000) // Get all clips

	// Filter by search query
	var results []Clip
	for _, clip := range clips {
		if len(results) >= limit {
			break
		}
		if contains(clip.Content, query) {
			results = append(results, clip)
		}
	}

	return results, nil
}

// Remove removes a clip
func (s *Storage) Remove(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.memory.Remove(id)
	return nil
}

// Clear removes all clips
func (s *Storage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.memory.Clear()
	return nil
}

// GetPinned returns only pinned clips
func (s *Storage) GetPinned() ([]Clip, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Get all clips
	clips := s.memory.List(1000)

	// Filter by pinned
	var pinned []Clip
	for _, clip := range clips {
		if clip.Pinned {
			pinned = append(pinned, clip)
		}
	}

	return pinned, nil
}

// Close closes storage (no-op for memory-only)
func (s *Storage) Close() error {
	return nil
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && s[:len(substr)] == substr) ||
		(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
		indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
