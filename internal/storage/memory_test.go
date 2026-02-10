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

func TestMemoryStore_Update(t *testing.T) {
	store := NewMemoryStore()

	clip := Clip{
		Content:   "test content",
		Type:      "text",
		Timestamp: time.Now(),
	}

	id, _ := store.Add(clip)

	// Retrieve, modify, and update
	retrieved, _ := store.Get(id)
	retrieved.Pinned = true
	updated := store.Update(retrieved)
	if !updated {
		t.Fatal("Expected Update to return true")
	}

	// Verify the change persisted
	after, _ := store.Get(id)
	if !after.Pinned {
		t.Fatal("Expected clip to be pinned after Update")
	}
}

func TestMemoryStore_Update_NotFound(t *testing.T) {
	store := NewMemoryStore()

	clip := Clip{ID: 999, Content: "nope"}
	if store.Update(clip) {
		t.Fatal("Expected Update to return false for non-existent clip")
	}
}

func TestMemoryStore_EvictOldest_SkipsPinned(t *testing.T) {
	store := NewMemoryStore()

	// Add 3 clips: first two are pinned
	for i := 0; i < 3; i++ {
		clip := Clip{
			Content:   "content " + string(rune('a'+i)),
			Type:      "text",
			Timestamp: time.Now(),
			Pinned:    i < 2, // first two pinned
		}
		store.Add(clip)
	}

	// Evict should skip pinned and remove the unpinned one (clip 3, content "c")
	evicted := store.EvictOldest()
	if !evicted {
		t.Fatal("Expected eviction to succeed")
	}

	if store.Count() != 2 {
		t.Fatalf("Expected 2 clips remaining, got %d", store.Count())
	}

	// Both remaining should be pinned
	clips := store.List(10)
	for _, c := range clips {
		if !c.Pinned {
			t.Fatalf("Expected all remaining clips to be pinned, found unpinned: %s", c.Content)
		}
	}

	// Another eviction should fail (all pinned)
	evicted = store.EvictOldest()
	if evicted {
		t.Fatal("Expected eviction to fail when all clips are pinned")
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
