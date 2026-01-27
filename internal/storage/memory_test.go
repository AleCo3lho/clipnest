package storage

import (
	"testing"
	"time"
)

func TestMemoryStore_Add(t *testing.T) {
	store := NewMemoryStore()

	clip := Clip{
		Content:   "test content",
		Type:      "text",
		Timestamp: time.Now(),
	}

	id, err := store.Add(clip)
	if err != nil {
		t.Fatalf("Failed to add clip: %v", err)
	}

	if id == 0 {
		t.Fatal("Expected non-zero ID")
	}
}

func TestMemoryStore_Get(t *testing.T) {
	store := NewMemoryStore()

	clip := Clip{
		Content:   "test content",
		Type:      "text",
		Timestamp: time.Now(),
	}

	id, _ := store.Add(clip)

	retrieved, exists := store.Get(id)
	if !exists {
		t.Fatal("Clip not found")
	}

	if retrieved.Content != clip.Content {
		t.Fatalf("Expected content %s, got %s", clip.Content, retrieved.Content)
	}
}

func TestMemoryStore_List(t *testing.T) {
	store := NewMemoryStore()

	// Add multiple clips
	for i := 0; i < 5; i++ {
		clip := Clip{
			Content:   "test content " + string(rune('a'+i)),
			Type:      "text",
			Timestamp: time.Now(),
		}
		store.Add(clip)
	}

	clips := store.List(3)
	if len(clips) != 3 {
		t.Fatalf("Expected 3 clips, got %d", len(clips))
	}
}

func TestMemoryStore_Remove(t *testing.T) {
	store := NewMemoryStore()

	clip := Clip{
		Content:   "test content",
		Type:      "text",
		Timestamp: time.Now(),
	}

	id, _ := store.Add(clip)

	removed := store.Remove(id)
	if !removed {
		t.Fatal("Failed to remove clip")
	}

	_, exists := store.Get(id)
	if exists {
		t.Fatal("Clip still exists after removal")
	}
}

func TestMemoryStore_EvictOldest(t *testing.T) {
	store := NewMemoryStore()

	// Add clips
	for i := 0; i < 3; i++ {
		clip := Clip{
			Content:   "test content " + string(rune('a'+i)),
			Type:      "text",
			Timestamp: time.Now(),
		}
		store.Add(clip)
	}

	initialCount := store.Count()

	evicted := store.EvictOldest()
	if !evicted {
		t.Fatal("Failed to evict oldest clip")
	}

	if store.Count() != initialCount-1 {
		t.Fatalf("Expected count %d, got %d", initialCount-1, store.Count())
	}
}

func TestMemoryStore_Clear(t *testing.T) {
	store := NewMemoryStore()

	// Add clips
	for i := 0; i < 5; i++ {
		clip := Clip{
			Content:   "test content",
			Type:      "text",
			Timestamp: time.Now(),
		}
		store.Add(clip)
	}

	store.Clear()

	if store.Count() != 0 {
		t.Fatalf("Expected empty store after clear, got %d clips", store.Count())
	}
}

func TestMemoryStore_Deduplication(t *testing.T) {
	store := NewMemoryStore()

	clip := Clip{
		Content:   "duplicate content",
		Type:      "text",
		Timestamp: time.Now(),
	}

	id1, _ := store.Add(clip)
	id2, _ := store.Add(clip)

	if id1 != id2 {
		t.Fatalf("Expected same ID for duplicate content, got %d and %d", id1, id2)
	}
}
