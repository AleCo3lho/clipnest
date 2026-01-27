package storage

import (
	"fmt"
	"sync"

	"clipnest/internal/config"
)

// Storage combines memory and persistent storage
type Storage struct {
	memory    *MemoryStore
	db        *SQLiteStore
	mu        sync.RWMutex
	maxMemory int
}

// NewStorage creates a new hybrid storage
func NewStorage(cfg config.Config) (*Storage, error) {
	memory := NewMemoryStore()
	db, err := NewSQLiteStore(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlite store: %w", err)
	}

	return &Storage{
		memory:    memory,
		db:        db,
		maxMemory: cfg.MaxMemoryClips,
	}, nil
}

// Add stores a clip (memory only by default)
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

	// Only add to DB if pinned
	if clip.Pinned {
		if err := s.db.Add(clip); err != nil {
			return id, err
		}
	}

	return id, nil
}

// Get retrieves a clip (checks memory first, then DB)
func (s *Storage) Get(id int64) (Clip, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check memory first
	if clip, exists := s.memory.Get(id); exists {
		return clip, nil
	}

	// Check DB (pinned items)
	return s.db.Get(id)
}

// List returns all clips (pinned first, then recent)
func (s *Storage) List(limit int) ([]Clip, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Clip

	// Get pinned items from DB
	pinned, err := s.db.List()
	if err != nil {
		return nil, err
	}
	result = append(result, pinned...)

	// Get recent from memory
	remaining := limit - len(result)
	if remaining > 0 {
		recent := s.memory.List(remaining)
		result = append(result, recent...)
	}

	return result, nil
}

// Pin moves a clip from memory to persistent storage
func (s *Storage) Pin(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get from memory
	clip, exists := s.memory.Get(id)
	if !exists {
		return fmt.Errorf("clip %d not found in memory", id)
	}

	// Mark as pinned
	clip.Pinned = true

	// Add to DB
	if err := s.db.Add(clip); err != nil {
		return err
	}

	// Remove from memory (optional - we keep it for quick access)
	// s.memory.Remove(id)

	return nil
}

// Unpin removes a clip from persistent storage
func (s *Storage) Unpin(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove from DB
	return s.db.Remove(id)
}

// Search clips in both memory and DB
func (s *Storage) Search(query string, limit int) ([]Clip, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Clip

	// Search DB (pinned items)
	pinned, err := s.db.Search(query, limit)
	if err != nil {
		return nil, err
	}
	result = append(result, pinned...)

	// Simple search in memory
	memoryClips := s.memory.List(limit)
	for _, clip := range memoryClips {
		if len(result) >= limit {
			break
		}
		if contains(clip.Content, query) {
			result = append(result, clip)
		}
	}

	return result, nil
}

// Remove removes a clip from both stores
func (s *Storage) Remove(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.memory.Remove(id)
	return s.db.Remove(id)
}

// Clear clears all clips (both memory and DB)
func (s *Storage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.memory.Clear()

	// Clear DB
	_, err := s.db.db.Exec("DELETE FROM clips")
	return err
}

// Close closes the database connection
func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Close()
}

// GetPinned returns only pinned clips
func (s *Storage) GetPinned() ([]Clip, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.db.List()
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
