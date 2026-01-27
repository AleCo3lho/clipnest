package storage

import (
	"testing"
	"time"
)

func setupTestStorage() (*Storage, func()) {
	store, _ := NewStorage(5)
	cleanup := func() {
		store.Close()
	}
	return store, cleanup
}

func TestStorage_Add(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

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

	// Verify clip exists
	retrieved, err := store.Get(id)
	if err != nil {
		t.Fatalf("Failed to get clip: %v", err)
	}

	if retrieved.Content != clip.Content {
		t.Fatalf("Expected content %s, got %s", clip.Content, retrieved.Content)
	}
}

func TestStorage_Pin(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	// Add unpinned clip
	clip := Clip{
		Content:   "test content",
		Type:      "text",
		Timestamp: time.Now(),
		Pinned:    false,
	}

	id, _ := store.Add(clip)

	// Pin it
	err := store.Pin(id)
	if err != nil {
		t.Fatalf("Failed to pin clip: %v", err)
	}

	// Verify it's marked as pinned
	retrieved, _ := store.Get(id)
	if !retrieved.Pinned {
		t.Fatal("Expected clip to be pinned")
	}
}

func TestStorage_Unpin(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	// Add a pinned clip
	clip := Clip{
		Content:   "test content",
		Type:      "text",
		Timestamp: time.Now(),
		Pinned:    true,
	}

	id, _ := store.Add(clip)

	// Verify it's pinned
	retrieved, _ := store.Get(id)
	if !retrieved.Pinned {
		t.Fatal("Expected clip to be pinned")
	}

	// Unpin it
	err := store.Unpin(id)
	if err != nil {
		t.Fatalf("Failed to unpin clip: %v", err)
	}

	// Verify it's unpinned
	retrieved, _ = store.Get(id)
	if retrieved.Pinned {
		t.Fatal("Expected clip to be unpinned")
	}
}

func TestStorage_List(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	// Add multiple clips
	for i := 0; i < 3; i++ {
		clip := Clip{
			Content:   "content " + string(rune('a'+i)),
			Type:      "text",
			Timestamp: time.Now(),
			Pinned:    i == 0, // First one pinned
		}
		store.Add(clip)
	}

	clips, err := store.List(10)
	if err != nil {
		t.Fatalf("Failed to list clips: %v", err)
	}

	if len(clips) != 3 {
		t.Fatalf("Expected 3 clips, got %d", len(clips))
	}

	// First should be pinned
	if !clips[0].Pinned {
		t.Fatal("Expected first clip to be pinned")
	}
}

func TestStorage_Search(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	// Add clips with different content
	clips := []Clip{
		{Content: "api_key_123", Type: "text", Timestamp: time.Now(), Pinned: true},
		{Content: "database_url", Type: "text", Timestamp: time.Now(), Pinned: false},
		{Content: "another value", Type: "text", Timestamp: time.Now(), Pinned: false},
	}

	for _, clip := range clips {
		store.Add(clip)
	}

	// Search for "api"
	results, err := store.Search("api", 10)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("Expected search results")
	}

	// Should find "api_key_123"
	found := false
	for _, result := range results {
		if result.Content == "api_key_123" {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("Search didn't find 'api_key_123'")
	}
}

func TestStorage_Remove(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	clip := Clip{
		Content:   "test content",
		Type:      "text",
		Timestamp: time.Now(),
		Pinned:    true,
	}

	id, _ := store.Add(clip)

	err := store.Remove(id)
	if err != nil {
		t.Fatalf("Failed to remove clip: %v", err)
	}

	// Verify it's gone
	_, err = store.Get(id)
	if err == nil {
		t.Fatal("Clip still exists after removal")
	}
}

func TestStorage_Clear(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	// Add clips
	for i := 0; i < 5; i++ {
		clip := Clip{
			Content:   "content " + string(rune('a'+i)),
			Type:      "text",
			Timestamp: time.Now(),
			Pinned:    i < 2, // First 2 pinned
		}
		store.Add(clip)
	}

	// Clear all
	err := store.Clear()
	if err != nil {
		t.Fatalf("Failed to clear storage: %v", err)
	}

	// Verify everything is gone
	clips, err := store.List(10)
	if err != nil {
		t.Fatalf("Failed to list clips: %v", err)
	}

	if len(clips) != 0 {
		t.Fatalf("Expected 0 clips after clear, got %d", len(clips))
	}

	pinned, _ := store.GetPinned()
	if len(pinned) != 0 {
		t.Fatalf("Expected 0 pinned clips after clear, got %d", len(pinned))
	}
}

func TestStorage_MemoryEviction(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	// Add 5 clips (exceeds memory limit of 5)
	for i := 0; i < 5; i++ {
		clip := Clip{
			Content:   "content " + string(rune('a'+i)),
			Type:      "text",
			Timestamp: time.Now(),
			Pinned:    false,
		}
		store.Add(clip)
	}

	// Add one more clip
	clip := Clip{
		Content:   "extra content",
		Type:      "text",
		Timestamp: time.Now(),
		Pinned:    false,
	}
	store.Add(clip)

	// Only last 5 should remain
	clips, err := store.List(10)
	if err != nil {
		t.Fatalf("Failed to list clips: %v", err)
	}

	if len(clips) != 5 {
		t.Fatalf("Expected 5 clips (memory limit), got %d", len(clips))
	}
}
