package storage

import (
	"container/list"
	"sync"
)

// MemoryStore stores clips in memory with LRU eviction
type MemoryStore struct {
	elements map[int64]*list.Element
	order    *list.List
	mu       sync.RWMutex
	nextID   int64
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		elements: make(map[int64]*list.Element),
		order:    list.New(),
		nextID:   1,
	}
}

// Add stores a clip in memory
func (m *MemoryStore) Add(clip Clip) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if this content already exists (deduplicate)
	for id, elem := range m.elements {
		stored := elem.Value.(Clip)
		if stored.Content == clip.Content && stored.Type == clip.Type {
			// Move to front (most recently used)
			m.order.MoveToFront(elem)
			return id, nil
		}
	}

	clip.ID = m.nextID
	m.nextID++

	// Store in list (front = most recent)
	elem := m.order.PushFront(clip)
	m.elements[clip.ID] = elem

	return clip.ID, nil
}

// Get retrieves a clip by ID
func (m *MemoryStore) Get(id int64) (Clip, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	elem, exists := m.elements[id]
	if !exists {
		return Clip{}, false
	}
	return elem.Value.(Clip), true
}

// List returns recent clips (most recent first)
func (m *MemoryStore) List(limit int) []Clip {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Clip, 0, limit)
	count := 0

	for elem := m.order.Front(); elem != nil && count < limit; elem = elem.Next() {
		result = append(result, elem.Value.(Clip))
		count++
	}

	return result
}

// Remove removes a clip from memory
func (m *MemoryStore) Remove(id int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	elem, exists := m.elements[id]
	if !exists {
		return false
	}

	m.order.Remove(elem)
	delete(m.elements, id)
	return true
}

// EvictOldest removes the oldest clip
func (m *MemoryStore) EvictOldest() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.order.Len() == 0 {
		return false
	}

	oldest := m.order.Back()
	if oldest != nil {
		clip := oldest.Value.(Clip)
		m.order.Remove(oldest)
		delete(m.elements, clip.ID)
		return true
	}

	return false
}

// Count returns the number of clips in memory
func (m *MemoryStore) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.order.Len()
}

// Clear removes all clips from memory
func (m *MemoryStore) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.elements = make(map[int64]*list.Element)
	m.order = list.New()
}
